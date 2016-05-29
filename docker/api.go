package docker

import (
	"encoding/json"
	"fmt"
	"github.com/mutaphore/docker-dispatch/utils"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
)

type DockerClient struct {
	httpClient *http.Client
	hostAddr   string
	pathPrefix string
}

func NewDockerClient(hostAddr string) *DockerClient {
	var httpClient *http.Client
	var pathPrefix string
	// Determine if this is a valid TCP or unix socket address
	_, err := net.ResolveTCPAddr("tcp", hostAddr)
	if err == nil {
		// http client transport uses TCP by default
		httpClient = &http.Client{}
		pathPrefix = "http://" + hostAddr
	} else if path.IsAbs(hostAddr) {
		httpClient = &http.Client{
			Transport: NewUnixTransport(hostAddr),
		}
		pathPrefix = "unix://"
	} else {
		log.Fatal("Unsupported transport method: must be unix socket or tcp")
	}
	return &DockerClient{
		httpClient: httpClient,
		hostAddr:   hostAddr,
		pathPrefix: pathPrefix,
	}
}

func (d *DockerClient) GetImages() []DockerImage {
	url := d.pathPrefix + "/images/json"
	req, err := http.NewRequest("GET", url, nil)
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
	var images []DockerImage
	json.Unmarshal(body, &images)
	return images
}

func (d *DockerClient) GetContainers() {
	url := utils.DockerHost() + "/containers/json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := d.httpClient.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var containers []interface{}
	json.Unmarshal(body, &containers)
	fmt.Printf("Number of containers %d\n", len(containers))
}
