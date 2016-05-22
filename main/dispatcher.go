package main

import (
	"github.com/mutaphore/docker-dispatch/docker"
)

func main() {
	dclient := docker.NewDockerClient()
	dclient.GetContainers()
}
