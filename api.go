package dockerdispatch

import (
	"encoding/json"
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

func (d *DockerClient) makeRequest(method, url string) []byte {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := d.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func (d *DockerClient) GetImages() []DockerImage {
	url := d.pathPrefix + "/images/json"
	body := d.makeRequest("GET", url)
	var images []DockerImage
	err := json.Unmarshal(body, &images)
	if err != nil {
		log.Printf("Error parsing JSON: %s", err.Error())
	}
	return images
}

func (d *DockerClient) GetContainers() []DockerContainer {
	url := d.pathPrefix + "/containers/json"
	body := d.makeRequest("GET", url)
	var containers []DockerContainer
	err := json.Unmarshal(body, &containers)
	if err != nil {
		log.Printf("Error parsing JSON: %s", err.Error())
	}
	return containers
}

func (d *DockerClient) GetInfo() DockerInfo {
	url := d.pathPrefix + "/info"
	body := d.makeRequest("GET", url)
	var info DockerInfo
	err := json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("Error parsing JSON: %s", err.Error())
	}
	log.Printf("%v", info)
	return info
}
