package docker

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
