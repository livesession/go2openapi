package main

import (
	restlix "github.com/zdunecki/restflix"
	"github.com/zdunecki/restflix/test/app"
)

// TODO: support methods and functions
// TODO: support recursion search
// TODO: support query
func main() {
	restlix.Iris(app.App(), []*restlix.SearchIdentifier{
		{
			MethodStatement:  []string{"BaseController", "ValidateBody"},
			ArgumentPosition: 1,
		},
		{
			MethodStatement:  []string{"iris", "Context", "ReadJSON"},
			ArgumentPosition: 0,
		},
	}, "./test", "")

	app.Init() // TODO:

	return
}
