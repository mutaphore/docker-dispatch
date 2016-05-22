package docker

import (
	"fmt"
	"github.com/mutaphore/docker-dispatch/utils"
	"io/ioutil"
	"log"
	"net/http"
)

type DockerClient struct {
	httpClient *http.Client
}

func NewDockerClient() *DockerClient {
	return &DockerClient{
		httpClient: &http.Client{},
	}
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
	fmt.Printf("Returned %s", body)
}
