package main

import (
	"fmt"

	"github.com/docker/docker/pkg/namesgenerator"
)

type Dispatcher struct {
	client   *DockerClient  // internal Docker client
	inbound  <-chan Message // receive only channel for Messages
	outbound chan Result    // outbound channel of results
	exited   chan string    // channel sending back exited containers
}

func NewDispatcher(hostAddr string, inbound <-chan Message) *Dispatcher {
	return &Dispatcher{
		client:  NewDockerClient(hostAddr),
		inbound: inbound,
		exited:  make(chan string, 200),
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
		name = namesgenerator.GetRandomName(0)
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
			d.outbound <- Result{Id: container.Id, Name: name, Data: err.Error()}
			return
		}
		go func() {
			for s := range output {
				d.outbound <- Result{Id: container.Id, Name: name, Data: s}
			}
		}()
	} else {
		// if we're not attaching to container output, just return container id
		d.outbound <- Result{Id: container.Id, Name: name, Data: fmt.Sprintf("Container id: %s", container.Id)}
	}

	// 3. Start container
	err = d.client.StartContainer(name)
	if err != nil {
		d.outbound <- Result{Id: container.Id, Name: name, Data: fmt.Sprintf("Error: %s", err.Error())}
		// TODO: remove attached loop
		return
	}

	// 4. Wait for container to finish and remove it
	if m.Options.Remove {
		err = d.client.WaitContainer(container.Id)
		if err != nil {
			d.outbound <- Result{Id: container.Id, Name: name, Data: fmt.Sprintf("Error: %s", err.Error())}
			return
		}
		err = d.client.RemoveContainer(container.Id, m.Options.Volumes, m.Options.Force)
		if err != nil {
			d.outbound <- Result{Id: container.Id, Name: name, Data: fmt.Sprintf("Error: %s", err.Error())}
			return
		}
		d.exited <- container.Id
	}
}

// Dispatch a stop container command
func (d *Dispatcher) stop(m Message) {
	// get number of seconds to wait before killing the container
	time := m.Options.Time
	if time == 0 {
		time = 10
	}
	err := d.client.StopContainer(m.Container, time)
	if err != nil {
		d.outbound <- Result{Data: err.Error()}
		return
	}
}

// Dispatch a remove container command
func (d *Dispatcher) remove(m Message) {
	err := d.client.RemoveContainer(m.Container, m.Options.Volumes, m.Options.Force)
	if err != nil {
		d.outbound <- Result{Data: err.Error()}
		return
	}
}
