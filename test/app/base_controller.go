package app

import (
	"net/http"

	"github.com/kataras/iris/v12"
)

type baseController struct{}

func NewBaseController() *baseController {
	return &baseController{}
}

func (c *baseController) InternalError(ctx iris.Context, err error) {
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.JSON(err)
}

func (c *baseController) ValidateBody(ctx iris.Context, data interface{}) error {
	return nil
}

func (c *baseController) ValidateQuery(ctx iris.Context, data interface{}) error {
	return nil
}
