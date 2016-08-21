package main

import (
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

	containerName := "sayhello"

	// remove container
	rmStr := fmt.Sprintf(`{
		"Dockercmd": "remove",
		"Container": "%s"
	}`, containerName)
	err = ch.Publish(
		"",
		qName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(rmStr),
		},
	)
	FailOnError(err, "Failed to publish message")

	// run container
	runStr := fmt.Sprintf(`{
		"Dockercmd": "run",
		"Options": {
			"Attach": ["STDERR", "STDOUT"],
			"Name": "%s"
		},
		"Image": "debian:jessie",
		"Cmd": ["echo", "hello"]
	}`, containerName)
	err = ch.Publish(
		"",
		qName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(runStr),
		},
	)
	FailOnError(err, "Failed to publish message")

}
