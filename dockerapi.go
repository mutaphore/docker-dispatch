package dockerdispatch

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
)

type DockerClient struct {
	hostAddr   string
	httpClient *http.Client
	pathPrefix string
	logger     *log.Logger
}

func NewDockerClient(hostAddr string) *DockerClient {
	dockerClient := &DockerClient{
		hostAddr: hostAddr,
	}
	// Determine if this is a valid TCP or unix socket address
	_, err := net.ResolveTCPAddr("tcp", hostAddr)
	if err == nil {
		// http client transport uses TCP by default
		dockerClient.httpClient = &http.Client{}
		dockerClient.pathPrefix = "http://" + hostAddr
	} else if path.IsAbs(hostAddr) {
		// custom transport for unix socket
		dockerClient.httpClient = &http.Client{
			Transport: NewUnixTransport(hostAddr),
		}
		dockerClient.pathPrefix = "unix://"
	} else {
		log.Fatal("Unsupported transport method: must be unix socket or tcp")
	}
	return dockerClient
}

func (d *DockerClient) makeRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (d *DockerClient) GetImages() ([]DockerImage, error) {
	body, err := d.makeRequest("GET", d.pathPrefix+"/images/json", nil)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, err
	}
	var images []DockerImage
	err = json.Unmarshal(body, &images)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	return images, err
}

func (d *DockerClient) GetContainers(all bool, filters map[string][]string) ([]DockerContainer, error) {
	body, err := d.makeRequest("GET", d.pathPrefix+"/containers/json", nil)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, err
	}
	var containers []DockerContainer
	err = json.Unmarshal(body, &containers)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	return containers, err
}

func (d *DockerClient) GetInfo() (DockerInfo, error) {
	body, err := d.makeRequest("GET", d.pathPrefix+"/info", nil)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, err
	}
	var info DockerInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	return info, err
}

func (d *DockerClient) CreateContainer(name string, param CreateContainerParam) (DockerContainer, error) {
	b, err := json.Marshal(param)
	if err != nil {
		log.Printf("Error marshalling parameter: %v", param)
		return nil, err
	}
	body, err := d.makeRequest("POST", d.pathPrefix+"/containers/create?name="+name, bytes.NewReader(b))
}
