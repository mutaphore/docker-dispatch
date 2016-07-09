package dockerdispatch

type Dispatcher struct {
	dclient  *DockerClient
	inbound  <-chan Message // receive only channel for Messages
	outbound chan<- Result  // send only channgel for Results
}

func NewDispatcher(hostAddr string) *Dispatcher {
	return &Dispatcher{
		dclient: NewDockerClient(hostAddr),
	}
}

func (d *Dispatcher) Start(inbound <-chan Message) <-chan Result {
	d.inbound = inbound
	d.outbound = make(chan Result)
	go func() {
		for m := range disp.inbound {
			if m.dockercmd == "run" {
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
	// If status is 404, pull, create again
	// Start container
	// If in detached mode, display container's id
	// else not in detached mode, attach to container
	r := Result{}
	d.outbound <- r
}
