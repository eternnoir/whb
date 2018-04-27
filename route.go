package main

import (
	"net/http"

	"github.com/eternnoir/whb/conmgr"
	"github.com/labstack/echo"
)

func bindRoute(ec *echo.Echo) {
	ec.GET("/", welcome)
	ec.POST("/:source/:target", bridge)
}

func welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Hi,here is WHB")
}

func bridge(c echo.Context) error {
	source := c.Param("source")
	target := c.Param("target")
	log.WithField("source", source).WithField("target", target).Info("Get Request")
	converter, err := conmgr.DefuaultConverterMgr.Get(source, target)
	if err != nil {
		log.WithError(err).Error(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"reason": err.Error()})
	}
	return converter.Process(c)
}
