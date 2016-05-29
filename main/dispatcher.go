package main

import (
	"fmt"
	"github.com/mutaphore/docker-dispatch/docker"
	// "log"
	// "net"
)

func main() {
	// dclient := docker.NewDockerClient("172.17.0.1:2375")
	dclient := docker.NewDockerClient("/var/run/docker.sock")
	images := dclient.GetImages()
	fmt.Printf("Number of images %d\n", len(images))
	fmt.Printf("Images %v\n", images[0].Id)
}
