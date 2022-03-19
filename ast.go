package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/structtag"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/kataras/iris/v12/core/router"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
)

func parseRouterMethod(openapi *openapi3.T, sourceFileName string, findMethodStatement string, route *router.Route, validationIdentifier func() []string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, sourceFileName, nil, parser.ParseComments) // 1. parse AST for router source file
	if err != nil {
		log.Fatal(err)
	}

	bodyParsed := false
	responseParsed := false

	operationName := fmt.Sprintf("[%s]%s", route.Method, strings.ReplaceAll(route.Path, "/", "-"))

	// path
	{
		var item *openapi3.PathItem

		if _, ok := openapi.Paths[route.Path]; ok {
			item = openapi.Paths[route.Path]
		} else {
			item = &openapi3.PathItem{}
		}

		// TODO: response
		pathOperation := &openapi3.Operation{
			RequestBody: &openapi3.RequestBodyRef{
				Ref: "#/components/requestBodies/" + operationName, // TODO: request someRequestBody
				//Value: requestBody, // TODO: request body
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
	}

root:
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl) // 2. check functions/methods only  // TODO: check if works with router inside method not function
		if !ok {
			continue
		}

		pkg := node.Name.Name
		methodName := fn.Name.Name

		if fn.Recv == nil || len(fn.Recv.List) <= 0 {
			continue
		}

		// 3. find router method e.g func (a *api) updateUser(ctx iris.Context)
		for _, recv := range fn.Recv.List { // 4. list pointers
			// TODO: support pointers and non pointers
			if star, ok := recv.Type.(*ast.StarExpr); ok { // 5. check if pointer e.g *api
				if ident, ok := star.X.(*ast.Ident); ok { // 6. get pointer name e.g "api"
					structName := ident.Name

					join := strings.Join([]string{pkg, fmt.Sprintf("(*%s)", structName), methodName}, ".") //
					routerMethodMatchWithASTMethod := join == findMethodStatement                          // 7. check if ast method match with iris router method

					if !routerMethodMatchWithASTMethod {
						continue
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
									if parseHandlerBodyAST(operationName, validationIdentifier, rhs, openapi, route) {
										bodyParsed = true
									}
								}
							case *ast.IfStmt: // 10. or maybe it's in if statement first
								if t.Init == nil {
									continue
								}

								init, ok := t.Init.(*ast.AssignStmt) // 11. and then check assign (inside if)
								if !ok || init.Rhs == nil {
									continue
								}

								for _, rhs := range init.Rhs {
									if parseHandlerBodyAST(operationName, validationIdentifier, rhs, openapi, route) {
										bodyParsed = true
									}
								}

							case *ast.ExprStmt: // TODO: response body, support 400, 404, 500 etc.
								parseContextJSONStruct := func(compositeList *ast.CompositeLit) {
									compositeTypIdent, ok := compositeList.Type.(*ast.Ident)
									if !ok {
										return
									}

									typeSpec, ok := compositeTypIdent.Obj.Decl.(*ast.TypeSpec)
									if !ok {
										return
									}

									responseBodyStruct := dynamicstruct.NewStruct()

									if err := parseStruct(typeSpec, responseBodyStruct); err != nil {
										panic(err)
									}

									requestBodyStructInstance := responseBodyStruct.Build().New()

									schemaRef, _, err := openapi3gen.NewSchemaRefForValue(requestBodyStructInstance)
									if err != nil {
										panic(err)
									}
									responseJSON := openapi3.NewResponse().
										WithJSONSchemaRef(schemaRef)

									// TOOD: support 400, 404, 500 etc.
									//openapi.Components.Responses[operationName] = &openapi3.ResponseRef{
									//	Value: responseJSON,
									//} // TODO:
									operation := openAPIOperationByMethod(openapi.Paths[route.Path], route.Method)
									operation.Responses.Default().Value = responseJSON

									responseParsed = true
									return
								}

								call, ok := t.X.(*ast.CallExpr)
								if !ok {
									continue
								}

								selector, ok := call.Fun.(*ast.SelectorExpr)
								if !ok {
									continue
								}

								if selector.Sel.Name != "JSON" {
									continue
								}

								ident, ok := selector.X.(*ast.Ident)
								if !ok {
									continue
								}

								f, ok := ident.Obj.Decl.(*ast.Field)
								if !ok {
									continue
								}

								tx, ok := f.Type.(*ast.SelectorExpr)
								if !ok {
									continue
								}

								if c, ok := tx.X.(*ast.Ident); !ok {
									continue
								} else if c.Name != "iris" {
									continue
								}

								if tx.Sel.Name != "Context" {
									continue
								}

								// yes it's iris context.JSON()

								if len(call.Args) < 1 {
									continue
								}

								argIdent, ok := call.Args[0].(*ast.Ident) // ctx.JSON(resp)
								if !ok {
									if compositeList, ok := call.Args[0].(*ast.CompositeLit); ok { // if ctx.JSON(&Struct{})
										parseContextJSONStruct(compositeList)
									}
									continue
								}

								// else ctx.JSON(&variable)

								assignStmt, ok := argIdent.Obj.Decl.(*ast.AssignStmt)
								if !ok {
									continue
								}

								if len(assignStmt.Rhs) < 1 {
									continue
								}

								compositeList, ok := assignStmt.Rhs[0].(*ast.CompositeLit)
								if !ok {
									continue
								}

								parseContextJSONStruct(compositeList)
							}
						}
					}
					continue root
				}
			}
		}
	}
}

func parseHandlerBodyAST(operationName string, validationIdentifier func() []string, exp ast.Expr, openapi *openapi3.T, route *router.Route) bool {
	callExpresssions, match := matchCallExpression(validationIdentifier, exp) // 1. match call expression e.g a.baseController.ValidateBody(ctx, &req)
	if !match {
		return false
	}

	// TODO: support non pointer request struct also
	requestBodyPointer := callExpresssions.Args[1] // 2. get validation struct e.g a.baseController.ValidateBody(ctx, &req) -> &req

	requestBodyStruct := dynamicstruct.NewStruct()

	// parse body struct to swagger (comments, types, names, validation, json tags)
	if pointer, ok := requestBodyPointer.(*ast.UnaryExpr); ok {
		if ident, ok := pointer.X.(*ast.Ident); ok {
			if valueSpec, ok := ident.Obj.Decl.(*ast.ValueSpec); ok {
				if value, ok := valueSpec.Type.(*ast.Ident); ok {

					if typeSpec, ok := value.Obj.Decl.(*ast.TypeSpec); ok {
						if err := parseStruct(typeSpec, requestBodyStruct); err != nil {
							panic(err)
						}
					}
				}
			}
		}

		requestBodyStructInstance := requestBodyStruct.Build().New()

		schemaRef, _, err := openapi3gen.NewSchemaRefForValue(requestBodyStructInstance)
		if err != nil {
			panic(err)
		}

		requestBody := openapi3.NewRequestBody().
			WithJSONSchemaRef(schemaRef)

		// TODO: someRequestBody
		// request body ref
		openapi.Components.RequestBodies[operationName] = &openapi3.RequestBodyRef{
			Value: requestBody,
		}

		operation := openAPIOperationByMethod(openapi.Paths[route.Path], route.Method)

		if operation == nil {
			log.Fatalf("operation not found, [path=%s, method=%s]", route.Path, route.Method)
			return false
		}

		operation.RequestBody.Value = requestBody
	}

	// parse body argument

	return true
}

func parseStruct(typeSpec *ast.TypeSpec, structBuilder dynamicstruct.Builder) error {
	structSpec, ok := typeSpec.Type.(*ast.StructType)

	if !ok {
		return errors.New("invalid type")
	}

	for _, field := range structSpec.Fields.List {

		tags, err := structtag.Parse(strings.ReplaceAll(field.Tag.Value, "`", ""))
		if err != nil {
			return err
		}

		tag, err := tags.Get("json")
		if err != nil {
			return err
		}

		if f, ok := field.Type.(*ast.Ident); ok {
			name := field.Names[0].Name // TODO: [0] is ok? why its an array
			structBuilder.AddField(name, getType(f.Name), fmt.Sprintf(`json:"%s"`, tag.Name))
		} else if f, ok := field.Type.(*ast.SelectorExpr); ok {
			x := f.X.(*ast.Ident)

			switch t := getType(x.Name).(type) {
			case *time.Time:
				name := field.Names[0].Name
				structBuilder.AddField(name, t, fmt.Sprintf(`json:"%s"`, tag.Name))
			}
		}
	}

	return nil
}
