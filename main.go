package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/hardmacro/outboundip/logger"
	"github.com/labstack/echo/v4"
	flag "github.com/spf13/pflag"
)

var flagDebug bool
var flagIpFromXFF bool
var flagListenPort string

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "log in debug mode")
	flag.BoolVar(&flagIpFromXFF, "ip-from-xff", false, "get all remote IPs from 'X-Forwarded-For' header")
	flag.StringVar(&flagListenPort, "listen-port", "8080", "port to listen to for web traffic")
}
func main() {
	flag.Parse()
	logger.InitLogger(flagDebug)

	logger.Infow("startup")

	app := echo.New()
	app.HideBanner = true
	app.HidePort = true

	if flagIpFromXFF {
		_, ipV4, _ := net.ParseCIDR("0.0.0.0/0")
		_, ipV6, _ := net.ParseCIDR("0:0:0:0:0:0:0:0/0")
		app.IPExtractor = echo.ExtractIPFromXFFHeader(echo.TrustIPRange(ipV4), echo.TrustIPRange(ipV6))
	}

	app.GET("/", func(c echo.Context) error {
		logger.Infow("get",
			"path", "/",
			"remote", c.RealIP(),
			"userAgent", c.Request().UserAgent(),
		)
		return c.String(http.StatusOK, c.RealIP())
	})
	app.Logger.Fatal(
		app.Start(fmt.Sprintf(":%s", flagListenPort)),
	)
}
