package restflix

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"testing"
)

func TestAstToSchemaRef(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "./test/ast/objects.go", nil, parser.ParseComments) // 1. parse AST for router source file
	if err != nil {
		log.Fatal(err)
	}

	structsMapping, err := mapStructsFromFiles("./test/ast")
	if err != nil {
		log.Fatal(err)
	}

	enc := astSchemaEncoder{
		node:           node,
		structsMapping: structsMapping,
		sourceFileName: "./test/ast/objects.go",
	}

	for _, obj := range node.Scope.Objects {
		if typeSpec, ok := obj.Decl.(*ast.TypeSpec); ok {
			schemaRef, err := enc.astExprToSchemaRef(typeSpec.Type)
			if err != nil {
				t.Fatal(err)
			}

			t.Log(schemaRef)
		}
	}
}
