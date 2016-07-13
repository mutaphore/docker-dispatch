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
		log.Fatal("Invalid number of arguments: startq <qname>")
	}
	qName := os.Args[1]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	defer ch.Close()
	FailOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		qName,
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	fmt.Printf("Queue declared: %s\n", q.Name)
}
