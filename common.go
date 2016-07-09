package dockerdispatch

// Common data structures

// Docker

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
}

type DockerInfo map[string]interface{}

// Messages

// docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
// Example:
// Message {
//   "run"
//   "--rm"
//   "hello-world"
//   ""
//   "echo hello"
//   ""
// }
type Message struct {
	dockercmd string
	options   string
	image     string
	container string
	cmd       string
	args      string
}

type Result struct {
	code    int
	message string
}
