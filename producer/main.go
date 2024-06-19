package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Name string
}

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

	//create new queue...
	q, err := ch.QueueDeclare("events", true, false, false, true, amqp.Table{
		"x-queue-type":                    "stream",
		"x-stream-max-segment-size-bytes": 30000,
		"x-max-length-bytes":              150000,
		// "x-max-age": "1h",
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// publish 1001 messages to the stream...
	for i := 0; i < 1001; i++ {
		event := Event{
			Name: "test",
		}

		data, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}

		err = ch.PublishWithContext(ctx, "", "events", false, false, amqp.Publishing{
			Body:          data,
			CorrelationId: uuid.NewString(),
		})

		if err != nil {
			panic(err)
		}
	}

	fmt.Printf(q.Name)
}
