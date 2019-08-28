package controllers

import (
	"fmt"

	"github.com/pangpanglabs/goutils/behaviorlog"

	"nomni/utils/api"

	"github.com/labstack/echo"
)

func ReturnApiFail(c echo.Context, status int, apiError api.Error, detail ...interface{}) error {
	behaviorlog.FromCtx(c.Request().Context()).WithError(apiError)
	for _, d := range detail {
		if d != nil {
			apiError.Details = fmt.Sprint(detail...)
		}
	}
	return c.JSON(status, api.Result{
		Success: false,
		Error:   apiError,
	})
}

func ReturnApiSucc(c echo.Context, status int, result interface{}) error {
	behaviorlog.FromCtx(c.Request().Context()).WithBizAttrs(map[string]interface{}{"resp": result})
	return c.JSON(status, api.Result{
		Success: true,
		Result:  result,
	})
}
