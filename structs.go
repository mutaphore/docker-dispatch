package main

// Common data structures

// Docker API Parameters

type CreateContainerParam struct {
	Hostname        string
	Domainname      string
	User            string
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	Tty             bool
	OpenStdin       bool
	StdinOnce       bool
	Env             []string
	Labels          map[string]string
	Cmd             []string
	Entrypoint      []string
	Image           string
	Volumes         map[string]map[string]interface{}
	WorkingDir      string
	NetworkDisabled bool
	ExposedPorts    map[string]map[string]interface{}
	StopSignal      string
	HostConfig      map[string]interface{}
}

// Docker API Responses

type DockerImage struct {
	Id          string
	ParentId    string
	RepoTags    []string
	RepoDigests interface{}
	Created     int64
	Size        int64
	VirtualSize int64
	Labels      interface{}
}

type DockerContainer struct {
	Id         string
	Names      []string
	Image      string
	ImageID    string
	Command    string
	Created    int64
	State      string
	Status     string
	Ports      []string
	Labels     map[string]string
	SizeRw     int64
	SizeRootFs int64
	HostConfig map[string]interface{}
}

type DockerInfo map[string]interface{}

// Messages

// docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
// An example message object:
// Message {
//   "run"
//   RunOptions {
//   	 "rm"
//   }
//   "hello-world"
//   ""
//   "echo hello"
//   ""
// }
type Message struct {
	Dockercmd string
	Options   Options
	Image     string
	Container string
	Cmd       []string
}

type Options struct {
	Attach     []string
	Entrypoint string
	Name       string
	Time       int
	Volumes    bool
	Force      bool
}

type Result struct {
	Id   string      // container unique id
	Name string      // container name
	Data interface{} // data
}
