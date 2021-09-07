/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 11:36:12
 * @LastEditTime: 2021-09-07 18:42:04
 */
package main

import (
	"Scachet/src/cgroups/subsystems"
	"Scachet/src/container"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

// Callable internally
var initCommand = cli.Command{
	Name: "init",
	Usage: `init container process run user's process in container.
			-- Do not call it outside.`,
	/**
	* @description: Init container
	 */
	Action: func(c *cli.Context) error {
		log.Infof("Init contanier ...")
		cmd := c.Args().Get(0)
		log.Infof("command %s", cmd)
		// run init process
		container.RunContainerInitProcess()
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
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	/**
	* @description: Launch Container
	 */
	Action: func(c *cli.Context) error {
		// determine the length of the command
		if len(c.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range c.Args() {
			cmdArray = append(cmdArray, arg)
		}
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: c.String("m"),
			CpuSet:      c.String("cpuset"),
			CpuShare:    c.String("cpushare"),
		}

		tty := c.Bool("ti")
		Run(tty, cmdArray, resConf)
		return nil
	},
}
