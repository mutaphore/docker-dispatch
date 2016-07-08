package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Invalid number of arguments: producer qname")
		os.Exit(1)
	}
	qName := os.Args[1]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	body := "Hello world"
	err = ch.Publish(
		"",
		qName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	FailOnError(err, "Failed to publish message")
}
