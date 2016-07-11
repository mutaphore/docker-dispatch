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
	// Create container
	param := CreateContainerParam{
		Image: m.Image,
		Cmd:   m.Cmd,
	}
	c, err := d.client.CreateContainer(m.Container, param)
	// If status is 404, pull, create again
	if 
	// Start container
	// If in detached mode, display container's id
	// else not in detached mode, attach to container
	r := Result{}
	d.outbound <- r
}
