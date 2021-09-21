/*
 * @Author: dejavudwh
 * @Date: 2021-09-05 23:09:31
 * @LastEditTime: 2021-09-21 21:20:56
 */
package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `
	Sachet is a simple container runtime implementation.
	Just for fun.
`

func main() {
	app := cli.NewApp()
	app.Name = "Sachet"
	app.Usage = usage

	// init command for container
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		commitCommand,
		listCommand,
		logCommand,
		execCommand,
		stopCommand,
	}

	// init logrus
	app.Before = func(c *cli.Context) error {
		// set JSON formatter
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
