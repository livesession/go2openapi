package restflix

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kataras/iris/v12"
)

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// TODO: reject middlewares - only last method from router
// TODO: support multiple method declarations in different files
func irisRouterStrategy(app *iris.Application, openapi *openapi3.T, searchIdentifiers []*SearchIdentifier, structsMappingRootPath string, goModName string, ignoreRoutes []string) error {
	structsMapping, err := mapStructsFromFiles(structsMappingRootPath)
	if err != nil {
		return err
	}

root:
	for _, route := range app.APIBuilder.GetRoutes() {
		if ignoreRoutes != nil {
			for _, pattern := range ignoreRoutes {
				found, _ := regexp.MatchString(pattern, route.Path)
				if found {
					continue root
				}
			}
		}

		sourceFileName := route.SourceFileName
		findMethod := route.MainHandlerName // github.com/livesession/restflix/test/app.(*api).testBaseController-fm

		lastHandler := route.Handlers[len(route.Handlers)-1]
		handlerName := getFunctionName(lastHandler)

		p := strings.Split(handlerName, "/")
		findMethod = p[len(p)-1]
		findMethod = strings.TrimSuffix(findMethod, compilerClousureSuffix) // app.(*api).testBaseController

		fullReferenceSplitter := strings.Split(findMethod, ".")
		findMethod = fullReferenceSplitter[len(fullReferenceSplitter)-1] // testBaseController

		operationName := fmt.Sprintf("[%s]%s", route.Method, strings.ReplaceAll(route.Path, "/", "-"))
		if operationName == debugOperationMethod {
			fmt.Sprintf("d")
		}

		parseRouterMethod(openapi, sourceFileName, findMethod, route, structsMapping, searchIdentifiers, goModName)
	}

	return nil
}
