package restflix

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kataras/iris/v12/core/router"
	gogitignore "github.com/sabhiram/go-gitignore"
)

type SearchIdentifier struct {
	MethodStatement  []string
	ArgumentPosition int
}

/*
	TODO: support errors

	```
	var req request

	if err := ctx.ReadJSON(req | &req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}
   ```
*/
// TODO: issues with "-": "",
// TODO: custom types
// TODO: embedded structs
// TODO: nested interface/struct method response (outside pkg where definition)
// TODO: support
/*
	support nested within function make
	```
	response := &objects.GetVisitorsSegmentResponse{
		Results: make([]*objects.VisitorSegmentStats, 0),
	}
	```

	support variable declaration point to outside pkg

	var (
		card         *objects.Card
		billingInfo  *objects.BillingInfo
		subscription *objects.Subscription
	)
	```
	response := iris.Map{
		"card": card,
		...
	}
	```
*/

func parseRouterMethod(
	openapi *openapi3.T,
	sourceFileName string,
	findMethodStatement string,
	route *router.Route,
	structsMapping map[string]map[string]*ast.TypeSpec,
	searchIdentifiers []*SearchIdentifier,
	goModName string,
) {
	if route.Method == http.MethodOptions {
		return
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, sourceFileName, nil, parser.ParseComments) // 1. parse AST for router source file
	if err != nil {
		log.Fatal(err)
	}

	bodyParsed := false
	responseParsed := false

	operationName := fmt.Sprintf("[%s]%s", route.Method, strings.ReplaceAll(route.Path, "/", "-"))

	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	// path
	var item *openapi3.PathItem

	if _, ok := openapi.Paths[route.Path]; ok {
		item = openapi.Paths[route.Path]
	} else {
		item = &openapi3.PathItem{}
	}

	pathOperation := &openapi3.Operation{
		RequestBody: &openapi3.RequestBodyRef{
			Ref: "#/components/requestBodies/" + operationName,
		},
		// TODO: 200, 400, 404, 500
		Responses: openapi3.NewResponses(),
	}

	switch route.Method {
	case http.MethodGet:
		item.Get = pathOperation
	case http.MethodPost:
		item.Post = pathOperation
	case http.MethodPut:
		item.Put = pathOperation
	case http.MethodPatch:
		item.Patch = pathOperation
	case http.MethodDelete:
		item.Delete = pathOperation
	}

	openapi.Paths[route.Path] = item

root:
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl) // 2. check functions/methods only  // TODO: check if works with router inside method not function
		if !ok {
			continue
		}

		methodName := fn.Name.Name

		if fn.Recv == nil || len(fn.Recv.List) <= 0 {
			continue
		}

		// 3. find router method e.g func (a *api) updateUser(ctx iris.Context)
		for _, recv := range fn.Recv.List { // 4. list pointers
			// TODO: support pointers and non pointers
			if star, ok := recv.Type.(*ast.StarExpr); ok { // 5. check if pointer e.g *api
				if _, ok := star.X.(*ast.Ident); ok { // 6. get pointer name e.g "api"
					routerMethodMatchWithASTMethod := methodName == findMethodStatement // 7. check if ast method match with iris router method

					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}

					if !routerMethodMatchWithASTMethod {
						continue
					}

					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}

					if fn.Body != nil && len(fn.Body.List) > 0 { // 8. check body - there should be request validation
						for _, body := range fn.Body.List {
							if bodyParsed && responseParsed {
								break root
							}
							switch t := body.(type) {
							case *ast.AssignStmt: // 9. maybe validation in assign statement
								if t.Rhs == nil {
									continue
								}

								for _, rhs := range t.Rhs {
									if operationName == debugOperationMethod {
										fmt.Sprintf("d")
									}
									if parseHandlerBodyAST(node, sourceFileName, structsMapping, operationName, searchIdentifiers, rhs, openapi, route, goModName) {
										bodyParsed = true
									}
								}
							case *ast.IfStmt: // 10. or maybe it's in if statement first
								if t.Init == nil {
									continue
								}

								if operationName == debugOperationMethod {
									fmt.Sprintf("d")
								}

								init, ok := t.Init.(*ast.AssignStmt) // 11. and then check assign (inside if)
								if !ok || init.Rhs == nil {
									continue
								}

								for _, rhs := range init.Rhs {
									if operationName == debugOperationMethod {
										fmt.Sprintf("d")
									}
									if parseHandlerBodyAST(node, sourceFileName, structsMapping, operationName, searchIdentifiers, rhs, openapi, route, goModName) {
										bodyParsed = true
									}
								}

								// parse response body
							case *ast.ExprStmt: // TODO: response body, support 400, 404, 500 etc.
								if operationName == debugOperationMethod {
									fmt.Sprintf("d")
								}

								if parseHandlerResponseAST(node, structsMapping, sourceFileName, t.X, operationName, openapi, route, goModName) {
									responseParsed = true
								}
							}
						}
					}
					continue root
				}
			}
		}
	}
}

func parseHandlerBodyAST(
	node *ast.File,
	sourceFileName string,
	structsMapping map[string]map[string]*ast.TypeSpec,
	operationName string,
	searchIdentifiers []*SearchIdentifier,
	exp ast.Expr,
	openapi *openapi3.T,
	route *router.Route,
	goModName string,
) bool {
	callExpresssions, argumentPosition, match := matchCallExpression(operationName, searchIdentifiers, exp) // 1. match call expression e.g a.baseController.ValidateBody(ctx, &req)
	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}
	if !match {
		return false
	}

	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	// TODO: support non pointer request struct also
	requestBodyArg := callExpresssions.Args[argumentPosition] // 2. get validation struct e.g a.baseController.ValidateBody(ctx, &req) || &req ctx.ValidateRequest(&req)

	parsed := false
	var bodyOpenApiRef *openapi3.SchemaRef

	sourceDir, _ := filepath.Split(strings.ReplaceAll(sourceFileName, "./", ""))
	astEncoder := &astSchemaEncoder{
		node:           node,
		structsMapping: structsMapping,
		sourceDir:      sourceDir,
		goModName:      goModName,
	}

	switch reqBody := requestBodyArg.(type) {
	// TODO: var req pkg.request
	case *ast.UnaryExpr: //  if var req request
		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}
		if ident, ok := reqBody.X.(*ast.Ident); ok {
			if valueSpec, ok := ident.Obj.Decl.(*ast.ValueSpec); ok {
				if ref, _ := astEncoder.astExprToSchemaRef(operationName, valueSpec.Type); ref != nil {
					parsed = true
					bodyOpenApiRef = ref
				}
			}
		}
	case *ast.Ident: //  if req := &request{}
		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}
		if assign, ok := reqBody.Obj.Decl.(*ast.AssignStmt); ok {
			if expr, ok := assign.Rhs[0].(*ast.UnaryExpr); ok {
				if composeList, ok := expr.X.(*ast.CompositeLit); ok {
					if ref, _ := astEncoder.astExprToSchemaRef(operationName, composeList.Type); ref != nil {
						parsed = true
						bodyOpenApiRef = ref
					}
				}
			}
		}
	}

	if !parsed {
		return false
	}

	requestBody := openapi3.NewRequestBody().
		WithJSONSchemaRef(bodyOpenApiRef)

	// TODO: someRequestBody
	// request body ref
	openapi.Components.RequestBodies[operationName] = &openapi3.RequestBodyRef{
		Value: requestBody,
	}

	operation := openAPIOperationByMethod(openapi.Paths[route.Path], route.Method)

	if operation == nil {
		log.Printf("operation not found, [path=%s, method=%s]", route.Path, route.Method)
		return false
	}

	operation.RequestBody.Value = requestBody

	// parse body argument

	return true
}

func parseHandlerResponseAST(
	node *ast.File,
	structsMapping map[string]map[string]*ast.TypeSpec,
	sourceFileName string,
	exp ast.Expr,
	operationName string,
	openapi *openapi3.T,
	route *router.Route,
	goModName string,
) bool {
	sourceDir, _ := filepath.Split(strings.ReplaceAll(sourceFileName, "./", ""))

	encoder := &astSchemaEncoder{
		node:           node,
		structsMapping: structsMapping,
		sourceDir:      sourceDir,
		goModName:      goModName,
	}

	// TODO: support 400, 404, 500 etc.
	defaultResponse := func(schemaRef *openapi3.SchemaRef) {
		responseJSON := openapi3.NewResponse().
			WithJSONSchemaRef(schemaRef)

		operation := openAPIOperationByMethod(openapi.Paths[route.Path], route.Method)
		operation.Responses.Default().Value = responseJSON
	}

	parseContextJSONStruct := func(structsMapping map[string]map[string]*ast.TypeSpec, compositeList *ast.CompositeLit) bool {
		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}

		schemaRef, err := encoder.astExprToSchemaRef(operationName, compositeList)
		if err != nil {
			log.Println(err)
			return false
		}

		defaultResponse(schemaRef)

		return true
	}

	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	call, ok := exp.(*ast.CallExpr)
	if !ok {
		return false
	}

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if selector.Sel.Name != "JSON" {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	f, ok := ident.Obj.Decl.(*ast.Field)
	if !ok {
		return false
	}

	tx, ok := f.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if c, ok := tx.X.(*ast.Ident); !ok {
		return false
	} else if c.Name != "iris" {
		return false
	}

	if tx.Sel.Name != "Context" {
		return false
	}

	// yes it's iris context.JSON()

	if len(call.Args) < 1 {
		return false
	}

	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	parseExpr := func(expr ast.Expr) bool {
		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}
		switch t := expr.(type) {
		case *ast.CompositeLit:
			return parseContextJSONStruct(structsMapping, t)
		case *ast.UnaryExpr:
			compositeList, ok := t.X.(*ast.CompositeLit) // if ctx.JSON(&Struct{})
			if !ok {
				return false
			}

			return parseContextJSONStruct(structsMapping, compositeList)
		case *ast.CallExpr:
			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}

			argIdent, ok := call.Args[0].(*ast.Ident)
			if !ok {
				return false
			}

			assignStmt, ok := argIdent.Obj.Decl.(*ast.AssignStmt)
			if !ok {
				return false
			}

			if schemaRef := encoder.astMethodToSchemaRef(operationName, assignStmt, argIdent.Name); schemaRef != nil {
				defaultResponse(schemaRef)
			}
		}

		return false
	}

	argIdent, ok := call.Args[0].(*ast.Ident)
	if !ok {
		return parseExpr(call.Args[0])
	}

	assignStmt, ok := argIdent.Obj.Decl.(*ast.AssignStmt)
	if !ok {
		return false
	}

	if len(assignStmt.Rhs) < 1 {
		return false
	}

	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	return parseExpr(assignStmt.Rhs[0])
}

func matchCallExpression(operationName string, searchIdentifiers []*SearchIdentifier, exp ast.Expr) (callExpr *ast.CallExpr, foundArgumentPosition int, found bool) {
	var selectorExpr *ast.SelectorExpr

	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	call, ok := exp.(*ast.CallExpr)

	if !ok {
		return
	}

	selectorExpr, ok = call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	for _, identifier := range searchIdentifiers {
		ids := make([]string, len(identifier.MethodStatement))
		copy(ids, identifier.MethodStatement)

		reverseSliceString(ids)

		expect := len(ids)
		current := 0

		for i, id := range ids {
			nextId := ""

			if len(ids)-1 >= i+1 {
				nextId = ids[i+1]
			}

			if selectorExpr.Sel.Name != id {
				break
			}

			current++

			switch t := selectorExpr.X.(type) {
			case *ast.CallExpr:
				if s, ok := t.Fun.(*ast.SelectorExpr); ok {
					selectorExpr = s
				}
			case *ast.SelectorExpr:
				selectorExpr = t
			case *ast.Ident:
				if t.Obj == nil {
					selectorExpr = &ast.SelectorExpr{ // hack for method expression e.g  ["iris", "Context", "ReadJSON"] => ctx.ReadJSON
						Sel: &ast.Ident{
							Name: t.Name,
						},
					}

					continue
				}

				field, ok := t.Obj.Decl.(*ast.Field)
				if !ok {
					continue
				}

				s, ok := field.Type.(*ast.SelectorExpr)
				if ok {
					selectorExpr = s
					continue
				}

				starExpr, ok := field.Type.(*ast.StarExpr)
				if !ok {
					continue
				}

				ident, ok := starExpr.X.(*ast.Ident)
				if !ok {
					continue
				}

				typeSpec, ok := ident.Obj.Decl.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				for _, field := range structType.Fields.List {
					switch t := field.Type.(type) {
					case *ast.StarExpr: // if pointer e.g *baseController
						ptr, ok := field.Type.(*ast.StarExpr)
						if !ok {
							continue
						}

						ident, ok := ptr.X.(*ast.Ident)
						if !ok {
							continue
						}

						if operationName == debugOperationMethod {
							fmt.Sprintf("d")
						}

						if ident.Name != nextId {
							continue
						}

						selectorExpr = &ast.SelectorExpr{ // hack for method shortcut e.g api.BaseController.Abc() and api.Abc()
							Sel: &ast.Ident{
								Name: ident.Name,
							},
						}

						break
					case *ast.Ident: // if interface e.g BaseController
						if operationName == debugOperationMethod {
							fmt.Sprintf("d")
						}

						if t.Name != nextId {
							continue
						}

						if operationName == debugOperationMethod {
							fmt.Sprintf("d")
						}

						selectorExpr = &ast.SelectorExpr{ // hack for method shortcut e.g api.BaseController.Abc() and api.Abc()
							Sel: &ast.Ident{
								Name: t.Name,
							},
						}

						break
					}

				}
			}
		}

		if expect == current {
			return call, identifier.ArgumentPosition, true
		}
	}

	return
}

func mapStructsFromFiles(root string) (map[string]map[string]*ast.TypeSpec, error) {
	walkRook := "."
	if root != "" {
		walkRook = root
	}

	ignored, err := gogitignore.CompileIgnoreFile("./gitignore")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// per file path, per struct name
	out := make(map[string]map[string]*ast.TypeSpec)

	// TODO: pass root folder to omit every folder checking
	return out, filepath.Walk(walkRook, func(path string, f os.FileInfo, err error) error {
		if ignored != nil {
			if ignored.MatchesPath(path) {
				return filepath.SkipDir
			}
		}

		// TODO: respect .gitignore
		switch path {
		case ".git":
			return filepath.SkipDir
		}

		// TODO: omit <name>_test.go
		if filepath.Ext(path) != ".go" {
			return nil
		}

		if root != "" {
			if is, _ := isSubPath(root, path); !is {
				return nil
			}
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments) // 1. parse AST for router source file
		if err != nil {
			return err
		}

		// TODO: check if absolute - if yes remove prefix  ~/Code/zdunecki -> prefix

		if _, ok := out[path]; !ok {
			out[path] = make(map[string]*ast.TypeSpec)
		}

		if node.Scope != nil {
			for _, v := range node.Scope.Objects {
				typeSpec, ok := v.Decl.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if _, ok := out[path][typeSpec.Name.Name]; !ok {
					out[path][typeSpec.Name.Name] = typeSpec
				}
			}
		}

		return nil
	})
}

type astSchemaEncoder struct {
	node           *ast.File
	structsMapping map[string]map[string]*ast.TypeSpec
	sourceDir      string
	goModName      string
}

func getImportPkgName(nodeImport *ast.ImportSpec) (importName, importPath string) {
	if nodeImport.Name != nil {
		importName = nodeImport.Name.Name
	}

	importPath, _ = strconv.Unquote(nodeImport.Path.Value)

	if importName == "" {
		_, importName = filepath.Split(importPath)
	}

	return
}

// TODO: inside pkg, outside pkg, interfaces, structs
func (encoder *astSchemaEncoder) astMethodToSchemaRef(operationName string, assignExpr *ast.AssignStmt, variableName string) (schemaRef *openapi3.SchemaRef) {
	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	if assignExpr == nil || len(assignExpr.Rhs) == 0 {
		return
	}

	call, ok := assignExpr.Rhs[0].(*ast.CallExpr)
	if !ok {
		return
	}

	methodName := ""
	if f, ok := call.Fun.(*ast.SelectorExpr); ok {
		methodName = f.Sel.Name
	}

	switch t := call.Fun.(type) {
	case *ast.SelectorExpr:
		funSelector, ok := t.X.(*ast.SelectorExpr)
		if !ok {
			return
		}

		structRef := funSelector.X
		refName := funSelector.Sel.Name

		ident, ok := structRef.(*ast.Ident)

		if !ok {
			return
		}

		if ident.Obj == nil {
			return
		}

		field, ok := ident.Obj.Decl.(*ast.Field)
		if !ok {
			return
		}

		starExpr, ok := field.Type.(*ast.StarExpr)
		if !ok {
			return
		}

		ident, ok = starExpr.X.(*ast.Ident)
		if !ok {
			return
		}

		typeSpec, ok := ident.Obj.Decl.(*ast.TypeSpec)
		if !ok {
			return
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return
		}

		var searchField *ast.Field

		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}

	fieldsLoop:
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				if refName == name.Name {
					searchField = field
					break fieldsLoop
				}
			}
		}

		if searchField == nil {
			return
		}

		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}

		switch tt := searchField.Type.(type) {
		case *ast.Ident: // inside pkg
			interfaceStructName := tt.Name // interface / struct name

			if ref := encoder.methodWithingPkgToSchemaRef(operationName, interfaceStructName, methodName, variableName, assignExpr); ref != nil {
				schemaRef = ref
			}
		case *ast.SelectorExpr: // outside pkg
			ident, ok := tt.X.(*ast.Ident)
			if !ok {
				return
			}

			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}

			for _, importNode := range encoder.node.Imports {
				importPkg, importPath := getImportPkgName(importNode)

				if ident.Name != importPkg {
					continue
				}

				interfaceStructName := tt.Sel.Name

				pkgDir := strings.TrimPrefix(importPath, encoder.goModName)

				// TODO: more optimized solution - cache, pass node as argument
				outsideEncoder := &astSchemaEncoder{
					node:           encoder.node, // TODO: i guess in the future it should be changed because what about nested node imports?
					structsMapping: encoder.structsMapping,
					sourceDir:      pkgDir,
					goModName:      encoder.goModName,
				}

				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}

				if ref := outsideEncoder.methodWithingPkgToSchemaRef(operationName, interfaceStructName, methodName, variableName, assignExpr); ref != nil {
					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}
					continue
				}
			}
		}
	case *ast.Ident:
		if t.Name == "make" {
			ref, err := encoder.astExprToSchemaRef(operationName, call.Args[0])
			if err != nil {
				fmt.Println(err)
				break
			}
			schemaRef = ref
		} else {
			// TODO: handle other function calls
		}
	}

	return
}

// TODO: search within package / another packages
func (encoder *astSchemaEncoder) astExprToSchemaRef(operationName string, expr ast.Expr) (*openapi3.SchemaRef, error) {
	var schemaRefFromStructField func(field *ast.Field) (*openapi3.SchemaRef, string, error)
	properties := make(openapi3.Schemas) // TODO: not always properties - sometimes primitive types or array of structs

	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	schemaRef := func(ident *ast.Ident) *openapi3.SchemaRef {
		return &openapi3.SchemaRef{
			Ref: "",
			Value: &openapi3.Schema{
				Type:    ident.Name,
				Example: getTypeDefaultValue(ident.Name),
			},
		}
	}

	schemaRefArray := func(ident *ast.Ident) *openapi3.SchemaRef {
		return &openapi3.SchemaRef{
			Ref: "",
			Value: &openapi3.Schema{
				Type:  "array",
				Items: schemaRef(ident),
			},
		}
	}

	schemaRefArrayOfObject := func(props openapi3.Schemas) *openapi3.SchemaRef {
		return &openapi3.SchemaRef{
			Ref: "",
			Value: &openapi3.Schema{
				Type: "array",
				Items: &openapi3.SchemaRef{
					Ref: "",
					Value: &openapi3.Schema{
						Type:       "object",
						Properties: props,
					},
				},
			},
		}
	}

	schemaRefFromStructField = func(field *ast.Field) (*openapi3.SchemaRef, string, error) {
		var parseExpr func(expr ast.Expr) (*openapi3.SchemaRef, string, error)

		if field.Tag == nil {
			return nil, "", nil
		}
		tags, err := structtag.Parse(strings.ReplaceAll(field.Tag.Value, "`", ""))
		if err != nil {
			return nil, "", err
		}

		tag, err := tags.Get("json")
		if err != nil {
			return nil, "", err
		}

		parseExpr = func(expr ast.Expr) (*openapi3.SchemaRef, string, error) {
			switch t := expr.(type) {
			case *ast.StarExpr:
				return parseExpr(t.X)
			case *ast.Ident: // &req{}
				if ref, _ := encoder.astToSchemaRefWithinPackage(operationName, t); ref != nil {
					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}
					return ref, tag.Name, err
				}

				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}

				return schemaRef(t), tag.Name, nil
			case *ast.SelectorExpr: // &pkg.req{}
				if ref, _ := encoder.astToSchemaOutsidePackage(operationName, t); ref != nil {
					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}
					return ref, tag.Name, err
				}

				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}

				ident := t.X.(*ast.Ident)

				if ref, _ := encoder.astToSchemaRefWithinPackage(operationName, ident); ref != nil {
					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}

					return ref, tag.Name, err
				}

				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}

				return schemaRef(ident), tag.Name, nil
			case *ast.ArrayType:
				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}

				schemaRefFromIdent := func(ident *ast.Ident) *openapi3.SchemaRef {
					if ident.Obj != nil { // it's struct
						typeSpec := ident.Obj.Decl.(*ast.TypeSpec)
						nestedSchemaRef, err := encoder.astExprToSchemaRef(operationName, typeSpec.Type)

						if err != nil {
							return nil
						}

						return schemaRefArrayOfObject(nestedSchemaRef.Value.Properties)
					} else {
						return schemaRefArray(ident)
					}
				}

				switch t := t.Elt.(type) {
				case *ast.Ident:
					if ref := schemaRefFromIdent(t); ref != nil {
						return ref, tag.Name, nil
					}

					return nil, "", nil
				case *ast.StarExpr:
					ident, ok := t.X.(*ast.Ident)
					if ok {
						if ref := schemaRefFromIdent(ident); ref != nil {
							return ref, tag.Name, nil
						}

						return nil, "", nil
					}
				}
			}

			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}

			return nil, "", nil
		}

		return parseExpr(field.Type)
	}

	parseStructType := func(structType *ast.StructType) error {
		for _, field := range structType.Fields.List {
			ref, name, err := schemaRefFromStructField(field)
			if err != nil {
				return err
			}

			if ref != nil {
				properties[name] = ref
			}
		}

		return nil
	}

	parseMap := func(t *ast.CompositeLit) error {
		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}
		for _, elt := range t.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			// TODO: string | int | bool keys
			key := ""

			if k, ok := kv.Key.(*ast.BasicLit); ok {
				if unquoted, err := strconv.Unquote(k.Value); err != nil {
					key = k.Value
				} else {
					key = unquoted
				}
			}

			if key == "" {
				continue
			}

			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}

		tttLoop:
			switch ttt := kv.Value.(type) {
			case *ast.BasicLit:
				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}
				valType := getTypeFromToken(ttt.Kind)
				properties[key] = &openapi3.SchemaRef{
					Ref: "",
					Value: &openapi3.Schema{
						Type:    valType,
						Example: getTypeDefaultValue(valType),
					},
				}
			case *ast.Ident:
				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}
				if ttt.Obj != nil {
					switch v := ttt.Obj.Decl.(type) {
					case *ast.ValueSpec:
						if operationName == debugOperationMethod {
							fmt.Sprintf("d")
						}

						schemaRef, err := encoder.astExprToSchemaRef(operationName, v.Type)
						if err != nil {
							log.Println(err)
							return err
						}

						if !encoder.emptySchemaProperties(schemaRef) {
							if operationName == debugOperationMethod {
								fmt.Sprintf("d")
							}

							properties[key] = schemaRef

							break tttLoop
						}

					case *ast.AssignStmt:
						if operationName == debugOperationMethod {
							fmt.Sprintf("d")
						}

						schemaRef := encoder.astMethodToSchemaRef(operationName, v, ttt.Name)

						if operationName == debugOperationMethod {
							fmt.Sprintf("d")
						}
						if !encoder.emptySchemaProperties(schemaRef) {
							if operationName == debugOperationMethod {
								fmt.Sprintf("d")
							}

							properties[key] = schemaRef

							break tttLoop
						}
					}
				}

				valType := getTypeFromIdent(ttt)
				properties[key] = &openapi3.SchemaRef{
					Ref: "",
					Value: &openapi3.Schema{
						Type:    valType,
						Example: getTypeDefaultValue(valType),
					},
				}
			case *ast.CompositeLit:
				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}

				schemaRef, err := encoder.astExprToSchemaRef(operationName, ttt.Type)
				if err != nil {
					log.Println(err)
					return err
				}

				properties[key] = schemaRef
			}
		}

		return nil
	}

	switch t := expr.(type) {
	case *ast.Ident: // &exampleStructInOtherFile{} | var req request
		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}
		if t.Obj == nil {
			if ref, _ := encoder.astToSchemaRefWithinPackage(operationName, t); ref != nil {
				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}

				properties = ref.Value.Properties
			}
		} else { // var req request
			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}
			if typeSpec, ok := t.Obj.Decl.(*ast.TypeSpec); ok {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					if err := parseStructType(structType); err != nil {
						return nil, err
					}
				}
			}
		}

	case *ast.SelectorExpr: // &objects.exampleStructInOtherFile{}
		if ref, _ := encoder.astToSchemaOutsidePackage(operationName, t); ref != nil {
			properties = ref.Value.Properties
		} else {
			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}
			switch tx := t.X.(type) {
			case *ast.Ident:
				if ref, _ := encoder.astToSchemaRefWithinPackage(operationName, tx); ref != nil {
					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}
					properties = ref.Value.Properties
				}

				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}
			}
		}
	case *ast.CompositeLit:
		switch tt := t.Type.(type) {
		case *ast.Ident:
			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}
			typeSpec, ok := tt.Obj.Decl.(*ast.TypeSpec)
			if !ok {
				return nil, nil
			}
			structSpec, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return nil, nil
			}

			schemaRef, err := encoder.astExprToSchemaRef(operationName, structSpec)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			properties = schemaRef.Value.Properties
		case *ast.SelectorExpr:
			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}
			schemaRef, err := encoder.astExprToSchemaRef(operationName, tt)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			if encoder.emptySchemaProperties(schemaRef) && t.Elts != nil && len(t.Elts) > 0 {
				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}
				if err := parseMap(t); err != nil {
					return nil, err
				}
			} else {
				properties = schemaRef.Value.Properties
			}
		default:
			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}
			if err := parseMap(t); err != nil {
				return nil, err
			}
		}
	case *ast.StructType:
		if err := parseStructType(t); err != nil {
			return nil, err
		}

		// TODO: array of another package structure
	case *ast.ArrayType:
		var arrExp ast.Expr

		switch t := t.Elt.(type) {
		case *ast.StarExpr:
			arrExp = t.X
		default:
			arrExp = t
		}

		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}

		nestedSechemaRef, err := encoder.astExprToSchemaRef(operationName, arrExp)
		if err != nil {
			return nil, err
		}

		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}
		ref := schemaRefArrayOfObject(nestedSechemaRef.Value.Properties)

		switch t := t.Elt.(type) {
		case *ast.Ident:
			properties[t.Name] = ref
		default:
			if ref.Value.Properties == nil {
				return ref, nil
			}

			properties = ref.Value.Properties
		}

		// TODO: check below
		//ident := t.Elt.(*ast.Ident)
		//if ident.Obj != nil { // it's struct
		//	typeSpec := ident.Obj.Decl.(*ast.TypeSpec)
		//	nestedSechemaRef, err := encoder.astExprToSchemaRef(typeSpec.Type)
		//
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	ref := schemaRefArrayOfObject(nestedSechemaRef.Value.Properties)
		//
		//	if ref != nil {
		//		properties[ident.Name] = ref
		//	}
		//} else {
		//	ref := schemaRefArray(ident)
		//
		//	if ref != nil {
		//		properties[ident.Name] = ref
		//	}
		//}
	}

	return &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Properties: properties,
		},
	}, nil
}

func (encoder astSchemaEncoder) astToSchemaOutsidePackage(operationName string, selector *ast.SelectorExpr) (*openapi3.SchemaRef, error) {
	for _, nodeImport := range encoder.node.Imports { // search struct inside imports
		importName, nodeImportPath := getImportPkgName(nodeImport)

		ident, ok := selector.X.(*ast.Ident)
		if !ok {
			continue
		}

		if importName != ident.Name {
			continue
		}

		// TODO: dir mapping instead of loop all files
		for mappingPath, structs := range encoder.structsMapping { // strategy for in another pkg
			mappingDir, _ := filepath.Split(mappingPath)

			if strings.HasSuffix(mappingDir, "/") {
				mappingDir = mappingDir[:len(mappingDir)-1]
			}
			if !strings.HasSuffix(nodeImportPath, mappingDir) {
				continue
			}

			if typeSpec, ok := structs[selector.Sel.Name]; ok && typeSpec != nil {
				return encoder.astExprToSchemaRef(operationName, typeSpec.Type)
			}
		}
	}

	return nil, nil
}

func (encoder astSchemaEncoder) astToSchemaRefWithinPackage(operationName string, ident *ast.Ident) (*openapi3.SchemaRef, error) {
	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}
	if ident.Obj == nil { // strategy for in another file
		return encoder.structNameWithinPkgToSchemaRef(operationName, encoder.sourceDir, ident.Name)
	}

	if typeSpec, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
		return encoder.astExprToSchemaRef(operationName, typeSpec.Type)
	}

	return nil, nil
}

func (encoder astSchemaEncoder) structNameWithinPkgToSchemaRef(operationName, pkgDir, structName string) (*openapi3.SchemaRef, error) {
	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	// TODO: dir mapping instead of loop all files
	for mappingPath, structs := range encoder.structsMapping {
		dir, _ := filepath.Split(mappingPath)

		withinInnerPackage := encoder.samePackageDir(pkgDir, dir)

		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}

		if withinInnerPackage {
			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}
			if typeSpec, ok := structs[structName]; ok {
				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}
				if typeSpec != nil {
					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}
					ref, err := encoder.astExprToSchemaRef(operationName, typeSpec.Type)

					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}

					return ref, err
				}
			}
		}
	}

	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}

	return nil, nil
}

// e.g interfaceStructName - "AccountsModel", variableName - "accounts", assignExpr - accounts, err := a.accountsModel.GetAccount
// TODO: check by method name also (interfaceStructMethodName)
func (encoder astSchemaEncoder) methodWithingPkgToSchemaRef(
	operationName,
	interfaceStructName,
	interfaceStructMethodName,
	variableName string,
	assignExpr *ast.AssignStmt,
) (schemaRef *openapi3.SchemaRef) {
	pkgDir := encoder.sourceDir
	if strings.HasPrefix(pkgDir, "/") {
		pkgDir = "." + pkgDir // make sure it's relative from root
	}

	if operationName == debugOperationMethod {
		fmt.Sprintf("d")
	}
	// TODO: cache results (use global variable)
	filepath.Walk(pkgDir, func(path string, info os.FileInfo, err error) error {
		if schemaRef != nil {
			return filepath.SkipDir
		}
		if info.IsDir() {
			if path == pkgDir {
				return nil
			}
			return filepath.SkipDir
		}

		fset := token.NewFileSet()
		searchNode, err := parser.ParseFile(fset, path, nil, parser.ParseComments) // 1. parse AST for router source file
		if err != nil {
			return err
		}

		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}

		if obj, ok := searchNode.Scope.Objects[interfaceStructName]; ok {
			typeSpec, ok := obj.Decl.(*ast.TypeSpec)
			if !ok {
				return nil
			}

			// TODO: struct methods
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return nil
			}

			if operationName == debugOperationMethod {
				fmt.Sprintf("d")
			}

			for _, method := range interfaceType.Methods.List {
				foundName := false
				for _, n := range method.Names {
					if n.Name == interfaceStructMethodName {
						foundName = true
					}
				}

				if !foundName {
					continue
				}

				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}

				methodFun, ok := method.Type.(*ast.FuncType)
				if !ok {
					continue
				}

				var index *int

				for i, l := range assignExpr.Lhs { // resp, err :=
					ident, ok := l.(*ast.Ident)
					if !ok {
						continue
					}

					if ident.Name == variableName {
						index = &i
						break
					}
				}

				if index == nil {
					continue
				}

				if methodFun.Results == nil {
					continue
				}

				result := methodFun.Results.List[*index]

				if operationName == debugOperationMethod {
					fmt.Sprintf("d")
				}

				switch t := result.Type.(type) {
				case *ast.StarExpr:
					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}

					ref, _ := encoder.astExprToSchemaRef(operationName, t.X)

					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}

					if !encoder.emptySchemaProperties(ref) {
						schemaRef = ref
						break
					}
				default:
					if operationName == debugOperationMethod {
						fmt.Sprintf("d")
					}

					ref, _ := encoder.astExprToSchemaRef(operationName, t)

					if !encoder.emptySchemaProperties(ref) {
						schemaRef = ref
						break
					}
				}
			}
		}

		return nil
	})

	return
}

func (encoder astSchemaEncoder) emptySchemaProperties(schemaRef *openapi3.SchemaRef) bool {
	emptyValue := schemaRef == nil || schemaRef.Value == nil

	if emptyValue {
		return true
	}

	if schemaRef.Value.Properties == nil || len(schemaRef.Value.Properties) == 0 {
		if schemaRef.Value.Items == nil || schemaRef.Value.Items.Value == nil || len(schemaRef.Value.Items.Value.Properties) == 0 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (encoder astSchemaEncoder) samePackageDir(p1, p2 string) bool {
	p1 = strings.TrimSuffix(strings.TrimPrefix(p1, "/"), "/")
	p2 = strings.TrimSuffix(strings.TrimPrefix(p2, "/"), "/")

	return p1 == p2
}
