package main

import (
	"github.com/livesession/restflix"
	"github.com/livesession/restflix/test/app"
)

// TODO: support methods and functions
// TODO: support recursion search
// TODO: support query
func main() {
	restflix.Init((&restflix.Options{
		SearchIdentifiers: []*restflix.SearchIdentifier{
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
		GoModName:              "github.com/livesession/restflix",
	}).
		WithIris(app.App()),
	)

	app.Init() // TODO:

	return
}
