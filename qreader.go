package dockerdispatch

import (
	"github.com/streadway/amqp"
	"log"
)

type QueueReader struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	inbound <-chan amqp.Delivery // receive only channel
}

func NewQueueReader(url string) (*QueueReader, error) {
	qreader := QueueReader{}
	var err error

	qreader.conn, err = amqp.Dial(url)
	if err != nil {
		log.Printf("Failed to connect to queue address url: %s\n", url)
		return nil, err
	}

	qreader.channel, err = qreader.conn.Channel()
	if err != nil {
		log.Printf("Failed to create channel\n")
		return nil, err
	}

	return &qreader, nil
}

func (q *QueueReader) Consume(queue string) (<-chan []byte, error) {
	var err error
	outbound := make(chan []byte)
	q.inbound, err = q.channel.Consume(
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
