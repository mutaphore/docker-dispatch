package main

import (
	"encoding/json"
	"log"
)

// TODO: check validity of message fields

// Returns a channel that parses bytes from an inbound channel into Message structs
func NewMessageParser(inbound <-chan []byte) <-chan Message {
	outbound := make(chan Message)
	go func() {
		var m Message
		for msg := range inbound {
			// parse json string
			err := json.Unmarshal(msg, &m)
			if err != nil {
				log.Printf("Error in decoding message: %s", msg)
				continue
			}
			outbound <- m
		}
	}()
	return outbound
}
