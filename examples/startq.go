package tests

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Invalid number of arguments: startq qname")
		os.Exit(1)
	}
	qName := os.Args[1]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

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
