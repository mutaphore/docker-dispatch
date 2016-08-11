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
			switch m.Dockercmd {
			case "run":
				go d.dispatchRun(m)
			case "stop":
				go d.dispatchStop(m)
			case "remove":
				go d.dispatchRemove(m)
			default:
				d.outbound <- Result{data: fmt.Sprintf("Error: Unsupported operation %s", m.Dockercmd)}
			}
		}
		close(d.outbound)
	}()
	return d.outbound
}

// Dispatch a run container command
func (d *Dispatcher) dispatchRun(m Message) {
	// Generate random string if container name option doesn't exist
	name := m.Options.Name
	if name == nil {
		_, name = genRandStr(32)
	}
	// 1. Create a container
	var attach []string
	if m.Options.Attach {
		attach = m.Options.Attach
	}
	param := CreateContainerParam{
		AttachStdin:  itemInList("STDIN", attach),
		AttachStderr: itemInList("STDERR", attach),
		AttachStdout: itemInList("STDOUT", attach),
		Image:        m.Image,
		Cmd:          m.Cmd,
	}
	container, err := d.client.CreateContainer(name, param)
	// TODO: If err status is 404, pull image first, create again
	if err != nil {
		d.outbound <- Result{data: fmt.Sprintf("Error: %s", err.Error())}
		return
	}
	// Return container id
	d.outbound <- Result{data: fmt.Sprintf("Container id: %s", container.Id)}
	// 2. Attach to container
	if d.attach == true {
		stdout, err := d.client.AttachContainer(name)
		if err != nil {
			d.outbound <- Result{data: err.Error()}
			return
		}
		go func() {
			for s := range stdout {
				d.outbound <- Result{data: s}
			}
		}()
	}
	// 3. Start container
	err = d.client.StartContainer(name)
	if err != nil {
		d.outbound <- Result{data: fmt.Sprintf("Error: %s", err.Error())}
		// TODO: remove attached loop
		return
	}
}

// Dispatch a stop container command
func (d *Dispatcher) dispatchStop(m Message) {

}

// Dispatch a remove container command
func (d *Dispatcher) dispatchRemove(m Message) {

}
