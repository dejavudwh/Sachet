package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

// Callable internally
var initCommand = cli.Command{
	Name: "init",
	Usage: `init container process run user's process in container.
			-- Do not call it outside.`,
	/*
		init container
	*/
	Action: func(c *cli.Context) error {
		log.Infof("Init contanier ...")
		cmd := c.Args().Get(0)
		log.Infof("command %s", cmd)
		// TODO: run init process

		return nil
	},
}

// run command for container
var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups
			limit Sachet run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},
	/*
		run container
	*/
	Action: func(c *cli.Context) error {
		// determine the length of the command
		if len(c.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := c.Args().Get(0)
		tty := c.Bool("ti")
		Run(tty, cmd)
		return nil
	},
}
