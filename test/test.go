package main

import (
	"errors"
	"fmt"

	"github.com/kataras/iris/v12"
)

type (
	request struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	response struct {
		ID      uint64 `json:"id"`
		Message string `json:"message"`
	}
)

type api struct {
	baseController *baseController
}

func main() {
	app := iris.New()

	a := &api{
		baseController: NewBaseController(),
	}

	app.Put("/users/{id:uint64}", a.updateUser)
	//routes := app.APIBuilder.GetRoutes()
	fmt.Println(app.APIBuilder.GetRoutes()[0].MainHandlerName)
	app.Listen(":8080")
}

func (a *api) updateUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint64("id")

	var req request

	if err := a.baseController.ValidateBody(ctx, &req); err != nil {
		a.baseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	resp := response{
		ID:      id,
		Message: req.Firstname + " updated successfully",
	}

	ctx.JSON(resp)
}
