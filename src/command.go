/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 11:36:12
 * @LastEditTime: 2021-09-21 20:18:06
 */
package main

import (
	"Scachet/src/cgroups/subsystems"
	"Scachet/src/container"
	"fmt"
	"os"

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
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
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
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
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

		// io redirect
		tty := c.Bool("ti")
		// backgroud
		detach := c.Bool("d")
		if tty && detach {
			return fmt.Errorf("ti and d paramter can not both provided")
		}
		log.Infof("createTty %v", tty)
		// volume
		volume := c.String("v")
		containerName := c.String("name")
		Run(tty, cmdArray, resConf, volume, containerName)
		return nil
	},
}

// commit container
var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit a container into inage",
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		imageName := c.Args().Get(0)
		commitContainer(imageName)

		return nil
	},
}

// list container
var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all the containers",
	Action: func(c *cli.Context) error {
		ListContainers()

		return nil
	},
}

var logCommand = cli.Command{
	Name:  "logs",
	Usage: "print logs of a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Please input your container name")
		}
		containerName := context.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}

// exec container
var execCommand = cli.Command{
	Name:  "exec",
	Usage: "exec a command into container",
	Action: func(context *cli.Context) error {
		//This is for callback
		if os.Getenv(ENV_EXEC_PID) != "" {
			log.Infof("pid callback pid %s", os.Getgid())
			return nil
		}
		log.Infof("first call")
		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container name or command")
		}

		containerName := context.Args().Get(0)
		var commandArray []string
		for _, arg := range context.Args().Tail() {
			commandArray = append(commandArray, arg)
		}

		ExecContainer(containerName, commandArray)
		return nil
	},
}
