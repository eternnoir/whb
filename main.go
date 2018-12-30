package main

import (
	"os"

	logrusmiddleware "github.com/eternnoir/echo-logrusmiddleware"
	_ "github.com/eternnoir/whb/hangoutschat"
	_ "github.com/eternnoir/whb/msteams"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	aListenPort = ""
	aListenAddr = ""
	log         = logrus.WithField("module", "whb")
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:        "port, p",
		Usage:       "Server listen port.",
		EnvVar:      "PORT",
		Value:       "8080",
		Destination: &aListenPort,
	},
	cli.StringFlag{
		Name:        "addr, a",
		Usage:       "Server listen addr.",
		EnvVar:      "ADDR",
		Value:       "0.0.0.0",
		Destination: &aListenAddr,
	},
}

func start(c *cli.Context) error {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))
	e.Logger = logrusmiddleware.Logger{logrus.StandardLogger()}
	e.Use(logrusmiddleware.Hook())
	bindRoute(e)
	return e.Start(aListenAddr + ":" + aListenPort)
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
