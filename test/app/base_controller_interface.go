package app

import (
	"net/http"

	"github.com/kataras/iris/v12"
)

type BaseController interface {
	InternalError(ctx iris.Context, err error)
	ValidateBody(ctx iris.Context, data interface{}) error
	ValidateQuery(ctx iris.Context, data interface{}) error
}

type baseControllerInterface struct{}

func NewBaseController2() *baseControllerInterface {
	return &baseControllerInterface{}
}

func (c *baseControllerInterface) InternalError(ctx iris.Context, err error) {
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.JSON(err)
}

func (c *baseControllerInterface) ValidateBody(ctx iris.Context, data interface{}) error {
	return nil
}

func (c *baseControllerInterface) ValidateQuery(ctx iris.Context, data interface{}) error {
	return nil
}
