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
		d.outbound <- Result{data: fmt.Sprintf("Error: %s", err.Error())}
		return
	}

	// 1.1 Return container id
	d.outbound <- Result{data: fmt.Sprintf("Container id: %s", container.Id)}

	// 2. Attach to container
	if stdin || stderr || stdout {
		// default logs and stream to true
		output, err := d.client.AttachContainer(name, true, true, stdin, stdout, stderr)
		if err != nil {
			d.outbound <- Result{data: err.Error()}
			return
		}
		go func() {
			for s := range output {
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
