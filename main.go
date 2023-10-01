package main

import (
	"net/http"

	"github.com/hardmacro/outboundip/logger"
	"github.com/labstack/echo/v4"
)

func main() {
	logger.InitLogger(true)

	logger.Info("hello")

	e := echo.New()

	e.GET("/", func(c echo.Context) error {

		ip := "unknown"
		for key, values := range c.Request().Header {
			logger.Infow("headers",
				"key", key,
				"values", values,
			)
			if key == "Do-Connecting-Ip" {
				if len(values) > 0 {
					ip = values[0]
				}
			}
		}

		// logger.Infow(c.Request().RemoteAddr)
		// ip := strings.Split(c.Request().RemoteAddr, ":")[0] // TODO: probably actually the real IP is a header?
		// logger.Info(ip)
		return c.String(http.StatusOK, ip)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
