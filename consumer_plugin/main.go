package main

import (
	"fmt"
	"time"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

const (
	EVENTSTREAM = "events"
)

func main() {

	//connect to the stream plugin on rabbitmq...
	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost("localhost").
			SetPort(5552).
			SetUser("guest").
			SetPassword("guest"),
	)

	if err != nil {
		panic(err)
	}

	// create consumer options..
	consumerOptions := stream.NewConsumerOptions().
		SetConsumerName("consumer_1").
		SetOffset(stream.OffsetSpecification{}.First())

	// start consumer...
	consumer, err := env.NewConsumer(EVENTSTREAM, messageHandler, consumerOptions)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	consumer.Close()
}

func messageHandler(consumerContext stream.ConsumerContext, message *amqp.Message) {
	fmt.Printf("Message: %s\n", message.Properties.CorrelationID)
	fmt.Printf("Data: %v\n", string(message.GetData()))

	// we could unmashal the data into a struct... (depend on our business logic)
}
