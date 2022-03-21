package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	objects2 "github.com/zdunecki/restflix/test/app/objects"
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
	//ab string
	//*baseController
	//xd string
	BaseController
}

func middleware(ctx context.Context) {

}

// TODO: array parser issues

func Init() {
	a := &api{
		//baseController: NewBaseController(),
		BaseController: NewBaseController2(),
	}

	app.Put("/test-basecontroller/{id:uint64}", a.testBaseController)
	app.Post("/test-basecontroller-shortcut/{id:uint64}", a.testBaseControllerShortcut)
	app.Post("/test-basecontroller-ctx-read-json/{id:uint64}", a.testBaseControllerCtxReadJson)
	app.Get("/test-ctx-json-struct/{id:uint64}", a.testCtxJsonStruct)
	app.Delete("/test-ctx-json-variable/{id:uint64}", a.testCtxJsonVariable)
	app.Delete("/test-basecontroller-request-struct-in-variable/{id:uint64}", a.testBaseControllerRequestStructInVariable)
	app.Delete("/test-basecontroller-request-struct-in-different-file/{id:uint64}", a.testBaseControllerRequestStructInDifferentFile)
	app.Delete("/test-basecontroller-request-struct-in-different-file-and-package/{id:uint64}", a.testBaseControllerRequestStructInDifferentFileAndPackage) // TODO:

	app.Delete(
		"/test-basecontroller-middleware/{id:uint64}",
		middleware,
		a.testBaseControllerMiddleware,
	) // TODO:

	//routes := app.APIBuilder.GetRoutes()
	fmt.Println(app.APIBuilder.GetRoutes()[0].MainHandlerName)
	app.Listen(":8080")
}

func (a *api) testBaseController(ctx iris.Context) {
	var req request

	if err := a.BaseController.ValidateBody(ctx, &req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	resp := response2{
		Success: true,
	}

	ctx.JSON(resp)
}

func (a *api) testBaseControllerShortcut(ctx iris.Context) {
	var req request

	if err := a.ValidateBody(ctx, &req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	ctx.JSON(response3{
		ID: "1",
	})
}

func (a *api) testBaseControllerCtxReadJson(ctx iris.Context) {
	var req request

	if err := ctx.ReadJSON(&req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	ctx.JSON(response3{
		ID: "1",
	})
}

func (a *api) testCtxJsonStruct(ctx iris.Context) {
	var req request

	if err := a.BaseController.ValidateBody(ctx, &req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	ctx.JSON(response3{
		ID: "1",
	})
}

func (a *api) testCtxJsonVariable(ctx iris.Context) {
	id, _ := ctx.Params().GetUint64("id")

	var req request

	if err := a.BaseController.ValidateBody(ctx, &req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	resp := response{
		ID:      id,
		Message: req.Firstname + " updated successfully",
	}

	ctx.JSON(resp)
}

func (a *api) testBaseControllerRequestStructInVariable(ctx iris.Context) {
	req := &request{}

	if err := ctx.ReadJSON(req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	ctx.JSON(response3{
		ID: "1",
	})
}

func (a *api) testBaseControllerRequestStructInDifferentFile(ctx iris.Context) {
	req := &exampleStructInOtherFile{}

	if err := ctx.ReadJSON(req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	ctx.JSON(response3{
		ID: "1",
	})
}

func (a *api) testBaseControllerRequestStructInDifferentFileAndPackage(ctx iris.Context) {
	req := &objects2.ExampleStructInOtherFileAndPackage{}

	if err := ctx.ReadJSON(req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	ctx.JSON(response3{
		ID: "1",
	})
}

func (a *api) testBaseControllerMiddleware(ctx iris.Context) {
	req := &objects2.ExampleStructInOtherFileAndPackage{}

	if err := ctx.ReadJSON(req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	ctx.JSON(response3{
		ID: "1",
	})
}
