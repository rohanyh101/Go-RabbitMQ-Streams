package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", "guest", "guest", "localhost:5672", ""))
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer ch.Close()

	if err := ch.Qos(50, 0, false); err != nil {
		panic(err)
	}

	// auto-ack should be false, while consuming streams... or the server will break down...
	stream, err := ch.Consume("events", "events_consumer", false, false, false, false, amqp.Table{
		// offset-id from which the stream should start consuming...
		// by default it will be set to "next", which will only consumen the latest messages...
		"x-stream-offset": "10m",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting to consume stream")
	for event := range stream {
		fmt.Printf("Event: %s\n", event.CorrelationId)
		fmt.Printf("Headers: %v\n", event.Headers)

		// the payload will present in the body....
		fmt.Printf("Data: %v\n", string(event.Body))
	}
}
