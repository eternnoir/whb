package main

import (
	"os"

	_ "github.com/eternnoir/whb/converters"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	aListenAddr = ""
	log         = logrus.WithField("module", "whb")
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:        "addr, a",
		Usage:       "Server listen address.",
		EnvVar:      "ADDR",
		Value:       ":8080",
		Destination: &aListenAddr,
	},
}

func start(c *cli.Context) error {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))
	e.Use(middleware.Logger())
	bindRoute(e)
	return e.Start(aListenAddr)
}

func main() {
	app := cli.NewApp()
	app.Name = "WHB"
	app.Usage = "Webhook Bridge"
	app.Flags = flags
	app.Action = start
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
