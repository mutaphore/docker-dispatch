package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var (
	verbose bool
	queue   string
)

func setupFlags() {
	flag.BoolVar(&verbose, "v", false, "Turn on debugging messages")
	flag.StringVar(&queue, "q", "", "Queue name")
	flag.Parse()
}

// Print command line usage
// Docker host address should have format like 172.17.0.1:2375 for tcp or /var/run/docker.sock for socket
// Rabbit queue address should be like amqp://guest:guest@localhost:5672/
func usage() {
	fmt.Println("Usage: docker-dispatch [OPTIONS] DOCKERHOST AMQPADDR\n")
	fmt.Println("Dispatch docker containers from commands passed through queue\n")
	fmt.Println("DOCKERHOST - ip:port tcp address or a unix socket path. Examples: 172.17.0.1:2375 or /var/run/docker.sock")
	fmt.Println("AMQPADDR   - full RabbitMQ uri. Example: amqp://guest:guest@localhost:5672\n")
	fmt.Println("Options:")
	fmt.Println("  -q, queue name")
	fmt.Println("  -v, verbose")
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
	// Parse command line arguments
	setupFlags()
	args := flag.Args()
	checkArgs(args)
	dockerAddr := args[0]
	queueAddr := args[1]

	// Create queue reader
	qreader, err := NewQueueReader(queueAddr)
	FailOnError(err, fmt.Sprintf("Error connecting to queue at %s\n", queueAddr))
	queueOutput, err := qreader.Consume(queue)
	FailOnError(err, fmt.Sprintf("Error reading from queue '%s'\n", queue))
	fmt.Printf("Connected to queue \"%s\" at %s\n", queue, queueAddr)

	// Create message parser
	parserOutput := NewMessageParser(queueOutput)

	// Create dispatcher
	dispatcher := NewDispatcher(dockerAddr, parserOutput)
	dispOutput := dispatcher.Start()
	fmt.Printf("Connected to docker at %s\n", dockerAddr)

	// Output results data
	for r := range dispOutput {
		fmt.Printf("%-20s | %v", r.Name, r.Data)
	}
}
