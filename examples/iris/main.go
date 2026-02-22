package main

import (
	"github.com/livesession/go2openapi"
	"github.com/livesession/go2openapi/test/app"
)

// TODO: support methods and functions
// TODO: support recursion search
// TODO: support query
func main() {
	go2openapi.Init((&go2openapi.Options{
		SearchIdentifiers: []*go2openapi.SearchIdentifier{
			{
				MethodStatement:  []string{"BaseController", "ValidateBody"},
				ArgumentPosition: 1,
			},
			{
				MethodStatement:  []string{"iris", "Context", "ReadJSON"},
				ArgumentPosition: 0,
			},
		},
		StructsMappingRootPath: "./test",
		SavePath:               "",
		GoModName:              "github.com/livesession/go2openapi",
	}).
		WithIris(app.App()),
	)

	app.Init() // TODO:

	return
}
