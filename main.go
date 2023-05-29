package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/hardmacro/outboundip/logger"
	"github.com/labstack/echo-contrib/echoprometheus"
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

	// metrics
	// TODO: potential for high entropy as it records paths for 404's (and other) codes leaving it open to abuse
	//       use a custom skipper and maybe just count 404's indivudually
	app.Use(echoprometheus.NewMiddleware("outboundip"))
	go func() {
		metrics := echo.New()
		metrics.HideBanner = true
		metrics.HidePort = true
		metrics.GET("/metrics", echoprometheus.NewHandler())
		if err := metrics.Start(":9090"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalw("startup failure", "error", err)
		}
	}()

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
