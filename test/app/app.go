package app

import (
	"errors"
	"fmt"
	"net/http"
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

	requestNested struct {
		Nested request `json:"nested_request"`
		Wuwu   string  `json:"wuwu"`
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

	response4 struct {
		Request *request `json:"request"`
	}
)

type api struct {
	//ab string
	//*baseController
	//xd string
	BaseController
	model        Model
	modelOutside objects2.Model
}

func middleware(ctx context.Context) {

}

// TODO: array parser issues

func Init() {
	a := &api{
		//baseController: NewBaseController(),
		BaseController: NewBaseController2(),
		model:          newModel(),
		modelOutside:   objects2.NewModel(),
	}

	// simple
	//app.Put("/test-basecontroller/{id:uint64}", a.testBaseController)
	//app.Get("/test-ctx-json-struct/{id:uint64}", a.testCtxJsonStruct)

	// shortcut properties declaration
	//app.Post("/test-basecontroller-shortcut/{id:uint64}", a.testBaseControllerShortcut)

	// read json
	//app.Post("/test-basecontroller-ctx-read-json/{id:uint64}", a.testBaseControllerCtxReadJson)

	// variable assign
	//app.Delete("/test-ctx-json-variable/{id:uint64}", a.testCtxJsonVariable)
	//app.Delete("/test-basecontroller-request-struct-in-variable/{id:uint64}", a.testBaseControllerRequestStructInVariable)

	// within pkg
	//app.Delete("/test-basecontroller-request-struct-in-different-file/{id:uint64}", a.testBaseControllerRequestStructInDifferentFile)

	// outside pkg
	//app.Delete("/test-basecontroller-request-struct-in-different-file-and-package/{id:uint64}", a.testBaseControllerRequestStructInDifferentFileAndPackage)

	// nested
	//app.Delete("/test-basecontroller-request-struct-nested/{id:uint64}", a.testBaseControllerRequestStructNested)
	//
	// map
	//app.Delete("/test-basecontroller-request-struct-map-response/{id:uint64}", a.testBaseControllerRequestStructMapResponse)

	// method response within pkg
	//app.Delete("/test-basecontroller-request-struct-method-response/{id:uint64}", a.testBaseControllerRequestStructMetodResponse)
	//app.Delete("/test-basecontroller-request-struct-method-response-inside/{id:uint64}", a.testBaseControllerRequestStructMetodResponseInside)
	//app.Delete("/test-basecontroller-request-struct-method-response-map/{id:uint64}", a.testBaseControllerRequestStructMetodResponseMap)

	// method response outside pkg
	//app.Delete("/test-basecontroller-request-struct-method-response-outside-pkg/{id:uint64}", a.testBaseControllerRequestStructMetodResponseOutSidePkg)
	//app.Delete("/test-basecontroller-request-struct-method-response-outside-pkg-array/{id:uint64}", a.testBaseControllerRequestStructMetodResponseOutSidePkgArray)
	//app.Delete("/test-basecontroller-request-struct-method-response-outside-pkg-array-make/{id:uint64}", a.testBaseControllerRequestStructMetodResponseOutSidePkgArrayMake)
	//app.Delete("/test-basecontroller-request-struct-method-response-outside-pkg-array-make-nested/{id:uint64}", a.testBaseControllerRequestStructMetodResponseOutSidePkgArrayMakeNested)

	app.Delete("/test-basecontroller-request-struct-test", a.testBaseControllerRequestStructTEST)

	// var declaration
	//app.Delete("/test-basecontroller-struct-method-response-var-declaration-mix", a.testBaseControllerStructMetodResponseVarDeclarationMix)

	// middleware
	//app.Delete("/test-basecontroller-middleware/{id:uint64}", middleware, a.testBaseControllerMiddleware)

	//routes := app.APIBuilder.GetRoutes()
	fmt.Println(app.APIBuilder.GetRoutes()[0].MainHandlerName)
	app.Listen(":9090")
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

func (a *api) testBaseControllerRequestStructNested(ctx iris.Context) {
	req := &requestNested{}

	if err := ctx.ReadJSON(req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	ctx.JSON(response3{
		ID: "1",
	})
}

func (a *api) testBaseControllerRequestStructMapResponse(ctx iris.Context) {
	var req request

	if err := ctx.ReadJSON(&req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	ctx.JSON(map[string]interface{}{
		"a":                          "a",
		"b":                          1,
		"c":                          true,
		"above_struct":               request{},
		"within_package_struct":      exampleStructInOtherFile{},
		"outside_package_struct":     objects2.ExampleStructInOtherFileAndPackage{},
		"outside_package_struct_arr": []*objects2.ExampleStructInOtherFileAndPackage{},
	})
}

func (a *api) testBaseControllerRequestStructMetodResponse(ctx iris.Context) {
	resp, err := a.model.GetSomething()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	x := 10
	if x != 0 {
		fmt.Println(x)
	}

	example()

	ctx.JSON(&response4{
		Request: resp,
	})
}

func (a *api) testBaseControllerRequestStructMetodResponseInside(ctx iris.Context) {
	resp, err := a.model.GetSomething()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	x := 10
	if x != 0 {
		fmt.Println(x)
	}

	example()

	ctx.JSON(resp)
}

func (a *api) testBaseControllerRequestStructMetodResponseMap(ctx iris.Context) {
	resp, err := a.model.GetSomething()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	x := 10
	if x != 0 {
		fmt.Println(x)
	}

	example()

	ctx.JSON(iris.Map{
		"response": resp,
	})
}

func example() {

}

func (a *api) testBaseControllerRequestStructMetodResponseOutSidePkg(ctx iris.Context) {
	resp, err := a.modelOutside.GetSomething()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	x := 10
	if x != 0 {
		fmt.Println(x)
	}

	example()

	ctx.JSON(resp)
}

func (a *api) testBaseControllerRequestStructMetodResponseOutSidePkgArray(ctx iris.Context) {
	resp, err := a.modelOutside.GetSomethingArray()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	x := 10
	if x != 0 {
		fmt.Println(x)
	}

	example()

	ctx.JSON(iris.Map{
		"response": resp,
	})
}

func (a *api) testBaseControllerRequestStructMetodResponseOutSidePkgArrayMake(ctx iris.Context) {
	out := make([]*objects2.ExampleStructInOtherFileAndPackage, 0)

	ctx.JSON(iris.Map{
		"response": out,
	})
}

func (a *api) testBaseControllerRequestStructMetodResponseOutSidePkgArrayMakeNested(ctx iris.Context) {
	out := &objects2.ExampleStructInOtherFileAndPackageNested{
		X: make([]*objects2.ExampleStructInOtherFileAndPackage, 0),
	}

	ctx.JSON(iris.Map{
		"response": out,
	})
}

func (a *api) testBaseControllerRequestStructTEST(ctx iris.Context) {
	req := &objects2.CreateAgentsBatchRequest{}
	if err := a.ValidateBody(ctx, req); err != nil {
		a.BaseController.InternalError(ctx, errors.New("validation error"))
		return
	}

	/*
		TODO:
		1. ctx.JSON(iris.Map{})

		2. ctx.JSON(iris.Map{
			"abc": &objects2.CreateAgentRequest{
				Email: "email",
				Role:  "role",
			},
		})

		3.
	*/

	ctx.JSON(iris.Map{
		"abc": &objects2.CreateAgentRequest{
			Email: "email",
			Role:  "role",
		},
	})
}

func (a *api) testBaseControllerStructMetodResponseVarDeclarationMix(ctx iris.Context) {
	var (
		x  *objects2.ExampleStructInOtherFileAndPackageNested
		x2 exampleStructInOtherFile
	)

	ctx.JSON(iris.Map{
		"response":  x,
		"response2": x2,
	})
}
