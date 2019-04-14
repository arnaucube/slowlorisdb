package main

import (
	"os"

	"github.com/arnaucube/slowlorisdb/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "slowlorisdb"
	app.Version = "0.0.1-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, cmd.Commands...)
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
