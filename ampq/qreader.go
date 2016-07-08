package main

import (
	"github.com/mutaphore/docker-dispatch/utils"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type QueueReader struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	inbound chan amqp.Delivery
}

func NewQueueReader(url, queue string) (*QueueReader, error) {
	qreader := QueueReader{}
	var err error

	qreader.conn, err = amqp.Dial(url)
	if err != nil {
		log.Printf("Failed to connect to queue address %s\n", addr)
		return nil, err
	}

	qreader.channel, err = conn.Channel()
	if err != nil {
		log.Printf("Failed to create channel\n")
		return nil, err
	}

	return &qreader
}

func (q *QueueReader) Consume(queue) {
	var err error
	q.inbound, err = ch.Consume(
		queue,
		"qreader",
		true,
		false,
		false,
		false,
		nil,
	)
	go func() {

	}()
}

func main() {
	if len(os.Args) != 2 {
		log.Println("Invalid number of arguments: consumer qname")
		os.Exit(1)
	}
	qName := os.Args[1]

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
