package main

import (
	"fmt"
)

type Dispatcher struct {
	client   *DockerClient
	inbound  <-chan Message // receive only channel for Messages
	outbound chan Result    // outbound channel of results
	attach   bool           // bind to container stdout
}

func NewDispatcher(hostAddr string, attach bool) *Dispatcher {
	return &Dispatcher{
		client: NewDockerClient(hostAddr),
		attach: attach,
	}
}

func (d *Dispatcher) Start(inbound <-chan Message) <-chan Result {
	d.inbound = inbound
	d.outbound = make(chan Result)
	go func() {
		for m := range d.inbound {
			if m.Dockercmd == "run" {
				go d.dispatchRun(m)
			} else {
				d.outbound <- Result{data: fmt.Sprintf("Error: Unsupported operation %s", m.Dockercmd)}
			}
		}
		close(d.outbound)
	}()
	return d.outbound
}

func (d *Dispatcher) dispatchRun(m Message) {
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
	// Return container id
	d.outbound <- Result{data: container.Id}
	// Attach to container
	if d.attach == true {
		stdout, err := d.client.AttachContainer(name)
		if err != nil {
			d.outbound <- Result{data: err.Error()}
			return
		}
		for s := range stdout {
			d.outbound <- Result{data: s}
		}
	}
}
