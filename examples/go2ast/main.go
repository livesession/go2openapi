package main

import (
	"os"

	"github.com/livesession/restflix/go2ast"
)

//	TODO: check this conecpt

func main() {
	go2ast.Generate(`
		req := &exampleStructInOtherFile{}

		if err := ctx.ReadJSON(req); err != nil {
			a.BaseController.InternalError(ctx, errors.New("validation error"))
			return
		}
	`, os.Stdout)

}
