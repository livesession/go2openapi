package main

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kataras/iris/v12"
)

func irisRouterStrategy(app *iris.Application, openapi *openapi3.T, validationIdentifier func() []string) {
	// 1. get iris routers
	for _, route := range app.APIBuilder.GetRoutes() {
		sourceFileName := route.SourceFileName

		p := strings.Split(route.MainHandlerName, "/")
		findMethod := p[len(p)-1] // 2. get router method name => <package>.(*<struct>).<method> e.g app.(*api).updateUser
		findMethod = strings.TrimSuffix(findMethod, compilerClousureSuffix)

		parseRouterMethod(openapi, sourceFileName, findMethod, route, validationIdentifier)
	}
}
