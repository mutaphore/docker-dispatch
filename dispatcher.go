package dockerdispatch

type Dispatcher struct {
	client   *DockerClient
	inbound  <-chan Message // receive only channel for Messages
	outbound chan<- Result  // send only channgel for Results
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
		for m := range disp.inbound {
			if m.Dockercmd == "run" {
				go disp.DispatchRun(m)
			} else {
				d.outbound <- Result{code: -1, message: "Unsupported operation"}
			}
		}
	}()
	return d.outbound
}

func (d *Dispatcher) DispatchRun(m Message) {
	// Create a container
	name := m.Container
	param := CreateContainerParam{
		Image: m.Image,
		Cmd:   m.Cmd,
	}
	container, err := d.client.CreateContainer(name, param)
	// TODO: If status is 404, pull, create again
	// Start container
	err = d.client.StartContainer(name)
	// If in detached mode, display container's id
	// else not in detached mode, attach to container
	r := Result{
		data: container.Id,
	}
	d.outbound <- r
}
