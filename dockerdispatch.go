package dockerdispatch

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var (
	verbose bool
)

func setupFlags() {
	flag.BoolVar(&verbose, "v", false, "Turn on debugging messages")
	flag.Parse()
}

// Print command line usage
// Docker host address should have format like 172.17.0.1:2375 for tcp or /var/run/docker.sock for socket
// Rabbit queue address should be like amqp://guest:guest@localhost:5672/
func usage() {
	flag.Usage()
	fmt.Println("docker-dispatch [options] DockerHostAddr AMPQAddr")
}

func checkURL(urlStr string) (*url.URL, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if strings.Compare(u.Scheme, "unix") != 0 && strings.Compare(u.Scheme, "http") != 0 && strings.Compare(u.Scheme, "https") != 0 {
		fmt.Printf("Scheme is %v\n", u.Scheme)
		return nil, fmt.Errorf("Scheme must be unix ie. unix:///var/run/daemon/sock:/path or http")
	}
	return u, nil
}

func checkArgs(args []string) {
	if len(args) != 2 {
		usage()
		os.Exit(1)
	}
}

func main() {
	setupFlags()
	args := flag.Args()
	checkArgs(args)

	hostAddr := args[0]
	dclient := NewDockerClient(hostAddr)

	images := dclient.GetImages()
	fmt.Printf("Number of images %d\n", len(images))
	fmt.Printf("Images %v\n", images[0])

	containers := dclient.GetContainers()
	fmt.Printf("Number of containers %d\n", len(containers))

	dclient.GetInfo()
}
