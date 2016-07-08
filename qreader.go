package main

import (
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

func (q *QueueReader) Consume(queue) (chan []byte, error) {
	var err error
	outbound := make(chan []byte)
	q.inbound, err = ch.Consume(
		queue,
		"qreader",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to create consumer on channel\n")
		return nil, err
	}
	go func() {
		for msg := range q.inbound {
			outbound <- msg.Body
		}
	}()
	return outbound, nil
}
