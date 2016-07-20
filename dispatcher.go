package main

import (
	"fmt"
)

type Dispatcher struct {
	client   *DockerClient
	inbound  <-chan Message // receive only channel for Messages
	outbound chan Result
}

func NewDispatcher(hostAddr string) *Dispatcher {
	return &Dispatcher{
		client: NewDockerClient(hostAddr),
	}
}

func (d *Dispatcher) Start(inbound <-chan Message) <-chan Result {
	d.inbound = inbound
	d.outbound = make(chan Result)
	go func() {
		for m := range d.inbound {
			if m.Dockercmd == "run" {
				go d.dispatchRunner(m)
			} else {
				d.outbound <- Result{data: fmt.Sprintf("Error: Unsupported operation %s", m.Dockercmd)}
			}
		}
	}()
	return d.outbound
}

func (d *Dispatcher) dispatchRunner(m Message) {
	// Create a container
	name := m.Container
	param := CreateContainerParam{
		Image: m.Image,
		Cmd:   m.Cmd,
	}
	container, err := d.client.CreateContainer(name, param)
	if err != nil {
		d.outbound <- Result{data: fmt.Sprintf("Error: %s", err.Error())}
		return
	}
	// TODO: If err status is 404, pull, create again
	// Start container
	err = d.client.StartContainer(name)
	if err != nil {
		d.outbound <- Result{data: fmt.Sprintf("Error: %s", err.Error())}
		return
	}
	// In detached mode, display container's id
	r := Result{
		data: container.Id,
	}
	d.outbound <- r
	// Not in detached mode, attach to container
}
