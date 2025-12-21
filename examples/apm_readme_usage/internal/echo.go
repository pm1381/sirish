package internal

import (
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm/module/apmechov4"
)

func StartEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(apmechov4.Middleware())
	// for gin, grpc and other modules Elastic provide middlewares. or you can use it pure with transactions although package has built-in support for it
	return e
}

func SetUpRoutes(e *echo.Echo, h *TestModuleSirishWrapperImpl) {
	e.GET("/apm/test1/:id", func(c echo.Context) error {
		id := c.Param("id")
		res, err := h.DoTest1(c.Request().Context(), DoTest1Request{id: id}, 0)
		if err != nil {
			return err
		}
		return c.JSON(200, map[string]interface{}{
			"res": res,
		})
	})
	e.POST("/apm/test2", func(c echo.Context) error {
		var w *DoTest2Request
		if err := c.Bind(w); err != nil {
			return err
		}
		res, err := h.DoTest2(c.Request().Context(), w)
		if err != nil {
			return err
		}
		return c.JSON(200, map[string]interface{}{
			"response": res,
		})
	})
}
