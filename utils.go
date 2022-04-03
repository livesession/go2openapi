package restflix

import (
	"go/ast"
	"go/token"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

var compilerClousureSuffix = "-fm"

func isSubPath(parent, sub string) (bool, error) {
	up := ".." + string(os.PathSeparator)

	// path-comparisons using filepath.Abs don't work reliably according to docs (no unique representation).
	rel, err := filepath.Rel(parent, sub)
	if err != nil {
		return false, err
	}
	if !strings.HasPrefix(rel, up) && rel != ".." {
		return true, nil
	}
	return false, nil
}

func reverseSliceString(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

var primitiveTypeZeroValue = map[string]interface{}{
	"string":  "",
	"int":     0,
	"int8":    0,
	"int16":   0,
	"int32":   0,
	"int64":   0,
	"uint":    0,
	"uint8":   0,
	"uint16":  0,
	"uint32":  0,
	"uint64":  0,
	"float":   0.0,
	"float32": 0.0,
	"float64": 0.0,
	"bool":    false,
}

// TODO: pointers
func getPrimitiveTypeDefaultValue(typ string) interface{} {
	if v, ok := primitiveTypeZeroValue[typ]; ok {
		return v
	}

	switch typ {
	case "time":
		return &time.Time{}
	}

	return nil
}

func getTypeFromToken(t token.Token) string {
	switch t {
	case token.INT:
		return "int"
	case token.FLOAT:
		return "float"
	case token.CHAR:
		return "string"
	case token.STRING:
		return "string"
	}

	return ""
}

func getTypeFromIdent(t *ast.Ident) string {
	switch t.Name {
	case "true":
		return "bool"
	case "false":
		return "bool"
	}

	return ""
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
