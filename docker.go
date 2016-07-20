package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"strconv"

	"golang.org/x/net/websocket"
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

func (d *DockerClient) makeRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return d.httpClient.Do(req)
}

// Get list of images
func (d *DockerClient) GetImages() ([]DockerImage, error) {
	resp, err := d.makeRequest("GET", d.pathPrefix+"/images/json", nil)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: %s", resp.StatusCode)
	}
	defer resp.Body.Close()
	var images []DockerImage
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &images)
	return images, err
}

// Get list of containers
func (d *DockerClient) GetContainers(all bool, filters map[string][]string) ([]DockerContainer, error) {
	resp, err := d.makeRequest("GET", d.pathPrefix+"/containers/json?all="+strconv.FormatBool(all), nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: %s", resp.StatusCode)
	}
	var containers []DockerContainer
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &containers)
	return containers, err
}

// Get Docker info
func (d *DockerClient) GetInfo() (*DockerInfo, error) {
	resp, err := d.makeRequest("GET", d.pathPrefix+"/info", nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: %s", resp.StatusCode)
	}
	info := DockerInfo{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &info)
	return &info, err
}

// Create a container
func (d *DockerClient) CreateContainer(name string, param CreateContainerParam) (*DockerContainer, error) {
	b, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	resp, err := d.makeRequest("POST", d.pathPrefix+"/containers/create?name="+name, bytes.NewReader(b))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 201 {
		return nil, fmt.Errorf("%s", resp.StatusCode)
	}
	container := DockerContainer{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &container)
	return &container, err
}

// Start a container
func (d *DockerClient) StartContainer(idOrName string) error {
	resp, err := d.makeRequest("POST", d.pathPrefix+"/containers/"+idOrName+"/start", nil)
	if err != nil {
		return err
	} else if resp.StatusCode != 204 {
		return fmt.Errorf("%s", resp.StatusCode)
	}
	return nil
}

// Attach to a container
func (d *DockerClient) AttachContainer(idOrName string) (chan string, error) {
	ws, err := websocket.Dial(d.pathPrefix+"/containers/"+idOrName+"/attach/ws?stream=1&stdout=1", "", "http://localhost/")
	outbound := make(chan string)
	if err != nil {
		return nil, err
	}
	go func() {
		var msg = make([]bytes, 512)
		var n int
		for {
			n, err := ws.Read(msg)
			if err != nil {
				close(out)
				return
			}
			outbound <- string(msg[:n])
		}
	}()
	return outbound, nil
}
