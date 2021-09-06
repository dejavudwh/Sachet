package main

import (
	"Scachet/src/container"
	"os"

	log "github.com/Sirupsen/logrus"
)

/*
* Create a new process through the created command
 */
func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	parent.Wait()
	os.Exit(-1)
}
