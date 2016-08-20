package main

import (
	"fmt"
)

type Dispatcher struct {
	client   *DockerClient  // internal Docker client
	inbound  <-chan Message // receive only channel for Messages
	outbound chan Result    // outbound channel of results
}

func NewDispatcher(hostAddr string, inbound <-chan Message) *Dispatcher {
	return &Dispatcher{
		client:  NewDockerClient(hostAddr),
		inbound: inbound,
	}
}

func (d *Dispatcher) Start() <-chan Result {
	d.outbound = make(chan Result)
	go func() {
		for m := range d.inbound {
			switch m.Dockercmd {
			case "run":
				go d.run(m)
			case "stop":
				go d.stop(m)
			case "remove":
				go d.remove(m)
			default:
				d.outbound <- Result{Data: fmt.Sprintf("Error: Unsupported operation %s", m.Dockercmd)}
			}
		}
		close(d.outbound)
	}()
	return d.outbound
}

// Dispatch a run container command
func (d *Dispatcher) run(m Message) {
	// Generate random string for container name if container name option doesn't exist
	name := m.Options.Name
	if name == "" {
		_, name = genRandStr(32)
	}

	// 1. Create a container
	stdin := itemInList("STDIN", m.Options.Attach)
	stderr := itemInList("STDERR", m.Options.Attach)
	stdout := itemInList("STDOUT", m.Options.Attach)
	param := CreateContainerParam{
		AttachStdin:  stdin,
		AttachStderr: stderr,
		AttachStdout: stdout,
		Image:        m.Image,
		Cmd:          m.Cmd,
	}
	container, err := d.client.CreateContainer(name, param)
	// TODO: If err status is 404, pull image first, create again
	if err != nil {
		d.outbound <- Result{Data: fmt.Sprintf("Error: %s", err.Error())}
		return
	}

	// 2. Attach to container
	if stdin || stderr || stdout {
		// default logs and stream to true
		output, err := d.client.AttachContainer(name, true, true, stdin, stdout, stderr)
		if err != nil {
			d.outbound <- Result{Id: container.Id, Data: err.Error()}
			return
		}
		go func() {
			for s := range output {
				d.outbound <- Result{Id: container.Id, Data: s}
			}
		}()
	} else {
		// if not attaching, just return container id
		d.outbound <- Result{Id: container.Id, Data: fmt.Sprintf("Container id: %s", container.Id)}
	}

	// 3. Start container
	err = d.client.StartContainer(name)
	if err != nil {
		d.outbound <- Result{Data: fmt.Sprintf("Error: %s", err.Error())}
		// TODO: remove attached loop
		return
	}
}

// Dispatch a stop container command
func (d *Dispatcher) stop(m Message) {

}

// Dispatch a remove container command
func (d *Dispatcher) remove(m Message) {

}
