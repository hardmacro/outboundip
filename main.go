package main

import (
	"github.com/hardmacro/outboundip/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func main() {
	logger.InitLogger(false)

	logger.Info("hello")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		ip := strings.Split(c.Request().RemoteAddr, ":")[0] // TODO: probably actually the real IP is a header?
		logger.Info(ip)
		return c.String(http.StatusOK, ip)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
