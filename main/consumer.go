package main

import (
	"github.com/mutaphore/docker-dispatch/utils"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Invalid number of arguments: consumer qname")
		os.Exit(1)
	}
	qName := os.Args[1]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")

	msgs, err := ch.Consume(
		qName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to open a channel")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}
