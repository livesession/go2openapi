package main

import (
	"go/ast"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

var compilerClousureSuffix = "-fm"

func reverseSliceString(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// TODO: pointers
func getType(typ string) interface{} {
	switch typ {
	case "string":
		return ""
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return 0
	case "float32", "float64":
		return 0.0
	case "bool":
		return false
	case "time":
		return &time.Time{}
	}

	return nil
}

func matchCallExpression(matchIdentifier func() []string, exp ast.Expr) (*ast.CallExpr, bool) {
	call, ok := exp.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	ids := matchIdentifier()
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

	return call, expect == current
}

func openAPIOperationByMethod(pathItem *openapi3.PathItem, method string) *openapi3.Operation {
	switch method {
	case http.MethodGet:
		return pathItem.Get
	case http.MethodPost:
		return pathItem.Post
	case http.MethodPut:
		return pathItem.Put
	case http.MethodPatch:
		return pathItem.Patch
	case http.MethodDelete:
		return pathItem.Delete
	}

	return nil
}
