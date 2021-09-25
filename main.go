package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/structtag"
	"github.com/go-openapi/spec"
	"github.com/swaggo/swag"
)

func reverseSliceString(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func main() {
	operation := swag.NewOperation(nil)

	sourceFileName := "./test/test.go"

	routesMethods := []string{
		"main.(*api).updateUser-fm",
	}

	compilerClousureSuffix := "-fm"

	requestBodyIdentifier := func() []string {
		return []string{"baseController", "ValidateBody"}
	}
	//
	//requestQueryIdentifier := func() {
	//
	//}
	//
	//requestErrorIdentifier := func() {
	//
	//}

	findMethod := routesMethods[0]
	findMethod = strings.TrimSuffix(findMethod, compilerClousureSuffix)

	fmt.Println(routesMethods, sourceFileName)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, sourceFileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: support methods and functions
	// TODO: support recursion search

	fmt.Println("########### Manual Iteration ###########")

	fmt.Println("Imports:")
	for _, i := range node.Imports {
		fmt.Println(i.Path.Value)
	}

	fmt.Println("Comments:")
	for _, c := range node.Comments {
		fmt.Print(c.Text())
	}

	fmt.Println("Functions:")

	findBody := func(exp ast.Expr) {
		call, ok := exp.(*ast.CallExpr)
		if !ok {
			return
		}

		fun, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		ids := requestBodyIdentifier()
		reverseSliceString(ids)

		var f *ast.SelectorExpr
		f = fun

		expect := len(ids)
		current := 0

		for _, id := range ids {
			if f.Sel.Name != id {
				break
			}

			current++

			switch t := f.X.(type) {
			case *ast.CallExpr:
				if s, ok := t.Fun.(*ast.SelectorExpr); ok {
					f = s
				}
			case *ast.SelectorExpr:
				f = t
			}
		}

		if expect == current {
			// TODO: support non pointer request struct also
			requestBodyPointer := call.Args[1]

			type abc struct {
				name   string
				fields []string
			}

			a := &abc{
				fields: make([]string, 0),
			}

			// parse body struct to swagger (comments, types, names, validation, json tags)
			if pointer, ok := requestBodyPointer.(*ast.UnaryExpr); ok {
				if ident, ok := pointer.X.(*ast.Ident); ok {
					if valueSpec, ok := ident.Obj.Decl.(*ast.ValueSpec); ok {
						if value, ok := valueSpec.Type.(*ast.Ident); ok {

							if typeSpec, ok := value.Obj.Decl.(*ast.TypeSpec); ok {
								name := typeSpec.Name.Name
								a.name = name

								if typeSpec, ok := value.Obj.Decl.(*ast.TypeSpec); ok {
									if structSpec, ok := typeSpec.Type.(*ast.StructType); ok {
										for _, field := range structSpec.Fields.List {

											tags, err := structtag.Parse(strings.ReplaceAll(field.Tag.Value, "`", ""))
											if err != nil {
												panic(err)
											}

											tag, err := tags.Get("json")
											if err != nil {
												panic(err)
											}

											a.fields = append(a.fields, tag.Value())
										}
									}
								}
							}
						}
					}
				}

				resp := &spec.Response{
					ResponseProps: spec.ResponseProps{
						Description: "",
						Schema: &spec.Schema{
							SchemaProps:        spec.SchemaProps{},
						},
						Headers:  nil,
						Examples: nil,
					},
				}

				operation.AddResponse(http.StatusOK, resp)
				fmt.Println("parse request body to swagger", pointer)
			}
			// parse body argument
		}
	}

root:
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}

		pkg := node.Name.Name
		methodName := fn.Name.Name

		// find router method

		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			for _, recv := range fn.Recv.List {
				// TODO: support pointers and non pointers
				if star, ok := recv.Type.(*ast.StarExpr); ok {
					if ident, ok := star.X.(*ast.Ident); ok {
						structName := ident.Name

						join := strings.Join([]string{pkg, fmt.Sprintf("(*%s)", structName), methodName}, ".")

						if join == findMethod {
							fmt.Println("ok")
							if fn.Body != nil && len(fn.Body.List) > 0 {
								for _, body := range fn.Body.List {
									switch t := body.(type) {
									case *ast.AssignStmt:
										if t.Rhs == nil {
											continue
										}

										for _, rhs := range t.Rhs {
											findBody(rhs)
										}
									case *ast.IfStmt:
										if t.Init == nil {
											continue
										}

										init, ok := t.Init.(*ast.AssignStmt)
										if !ok || init.Rhs == nil {
											continue
										}

										for _, rhs := range init.Rhs {
											findBody(rhs)
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

		fmt.Println(fn.Name.Name)
	}

	fmt.Println("########### Inspect ###########")
	ast.Inspect(node, func(n ast.Node) bool {
		// Find Return Statements
		ret, ok := n.(*ast.ReturnStmt)
		if ok {
			fmt.Printf("return statement found on line %d:\n\t", fset.Position(ret.Pos()).Line)
			printer.Fprint(os.Stdout, fset, ret)
			return true
		}
		// Find Functions
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			var exported string
			if fn.Name.IsExported() {
				exported = "exported "
			}
			fmt.Printf("%sfunction declaration found on line %d: \n\t%s\n", exported, fset.Position(fn.Pos()).Line, fn.Name.Name)
			return true
		}
		return true
	})
	fmt.Println()
}
