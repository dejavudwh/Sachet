/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 14:08:11
 * @LastEditTime: 2021-09-06 17:46:53
 */
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
