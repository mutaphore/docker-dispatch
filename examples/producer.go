package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func FailOnError(err error, msg string) {
	if err != nil {
		s := fmt.Sprintf("%s - %s", msg, err)
		log.Fatal(s)
		panic(s)
	}
}

type Message struct {
	Dockercmd string
	Options   string
	Image     string
	Container string
	Cmd       []string
	Args      string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Invalid number of arguments: producer <qname>")
	}
	qName := os.Args[1]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	defer ch.Close()
	FailOnError(err, "Failed to open a channel")

	body, err := json.Marshal(Message{
		Dockercmd: "run",
		Image:     "debian:jessie",
		Container: "sayhello2",
	})
	FailOnError(err, "Failed to marshal body")

	err = ch.Publish(
		"",
		qName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	FailOnError(err, "Failed to publish message")
}
