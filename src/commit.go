/*
 * @Author: dejavudwh
 * @Date: 2021-09-20 17:28:46
 * @LastEditTime: 2021-09-20 17:28:46
 */
package main

import (
	"fmt"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

func commitContainer(imageName string) {
	// TODO: hard code
	mntURL := "/home/mnt"
	imageTar := "/home/" + imageName + ".tar"
	fmt.Printf("%s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", mntURL, err)
	}
}
