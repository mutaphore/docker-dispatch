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
	wsPrefix   string
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
		dockerClient.wsPrefix = "ws://" + hostAddr
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
	resp, err := d.makeRequest("GET", fmt.Sprintf("%s/images/json", d.pathPrefix), nil)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GetImages: error status code %s", resp.StatusCode)
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
	resp, err := d.makeRequest("GET", fmt.Sprintf("%s/containers/json?all=%s", d.pathPrefix, strconv.FormatBool(all)), nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GetContainers: error status code %d", resp.StatusCode)
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
	resp, err := d.makeRequest("GET", fmt.Sprintf("%s/info", d.pathPrefix), nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GetInfo: error status code %d", resp.StatusCode)
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
	resp, err := d.makeRequest("POST", fmt.Sprintf("%s/containers/create?name=%s", d.pathPrefix, name), bytes.NewReader(b))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 201 {
		return nil, fmt.Errorf("CreateContainer: error status code %d", resp.StatusCode)
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
	resp, err := d.makeRequest("POST", fmt.Sprintf("%s/containers/%s/start", d.pathPrefix, idOrName), nil)
	if err != nil {
		return err
	} else if resp.StatusCode != 204 {
		return fmt.Errorf("StartContainier: error status code %d", resp.StatusCode)
	}
	return nil
}

// Stop a container after t seconds
func (d *DockerClient) StopContainer(idOrName string, t int) error {
	resp, err := d.makeRequest("POST", fmt.Sprintf("%s/containers/%s/stop?t=%d", d.pathPrefix, idOrName, t), nil)
	if err != nil {
		return err
	} else if resp.StatusCode != 204 {
		return fmt.Errorf("StopContainier: error status code %d", resp.StatusCode)
	}
	return nil
}

// Remove a container
func (d *DockerClient) RemoveContainer(idOrName string, volume, force bool) error {
	resp, err := d.makeRequest("DELETE", fmt.Sprintf("%s/containers/%s?v=%v&force=%v", d.pathPrefix, idOrName, volume, force), nil)
	if err != nil {
		return err
	} else if resp.StatusCode != 204 {
		return fmt.Errorf("RemoveContainier: error status code %d", resp.StatusCode)
	}
	return nil
}

// Attach to a container
func (d *DockerClient) AttachContainer(idOrName string, logs, stream, stdin, stdout, stderr bool) (chan string, error) {
	uri := fmt.Sprintf("%s/containers/%s/attach/ws?logs=%v&stream=%v&stdin=%v&stdout=%v", d.wsPrefix, idOrName, logs, stream, stdin, stdout, stderr)
	ws, err := websocket.Dial(uri, "", "http://127.0.0.1/")
	outbound := make(chan string)
	if err != nil {
		return nil, err
	}
	go func() {
		var msg = make([]byte, 4096)
		for {
			n, err := ws.Read(msg)
			if err != nil {
				close(outbound)
				return
			}
			outbound <- string(msg[:n])
		}
	}()
	return outbound, nil
}
