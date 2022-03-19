package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
)

var app *iris.Application

func init() {
	app = iris.New()
}

func App() *iris.Application {
	return app
}

type (
	request struct {
		Firstname string    `json:"firstname"`
		Lastname  string    `json:"lastname"`
		Age       uint8     `json:"age"`
		Points    float64   `json:"points"`
		Awards    int       `json:"awards"`
		IsAdmin   bool      `json:"is_admin"`
		CreatedAt time.Time `json:"created_at"`
	}

	response struct {
		ID      uint64 `json:"id"`
		Message string `json:"message"`
	}

	response2 struct {
		Success bool `json:"success"`
	}

	response3 struct {
		ID string `json:"id"`
	}
)

type api struct {
	baseController *baseController
}

func Init() {
	a := &api{
		baseController: NewBaseController(),
	}

	app.Put("/users/{id:uint64}", a.updateUser2)
	app.Post("/users/{id:uint64}", a.updateUser3)
	app.Get("/users/{id:uint64}", a.updateUser)
	app.Delete("/users/{id:uint64}", a.updateUser)

	//routes := app.APIBuilder.GetRoutes()
	fmt.Println(app.APIBuilder.GetRoutes()[0].MainHandlerName)
	app.Listen(":8080")
}

func (a *api) updateUser2(ctx iris.Context) {
	var req request

	if err := a.baseController.ValidateBody(ctx, &req); err != nil {
		a.baseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	resp := response2{
		Success: true,
	}

	ctx.JSON(resp)
}

func (a *api) updateUser3(ctx iris.Context) {
	var req request

	if err := a.baseController.ValidateBody(ctx, &req); err != nil {
		a.baseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	ctx.JSON(response3{
		ID: "1",
	})
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
