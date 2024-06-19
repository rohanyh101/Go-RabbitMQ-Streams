package main

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

const (
	EVENTSTREAM = "events"
)

type Event struct {
	Message string
}

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

	// declare the stream, set segment size and max bytes on stream..
	err = env.DeclareStream(EVENTSTREAM, stream.NewStreamOptions().
		SetMaxSegmentSizeBytes(stream.ByteCapacity{}.MB(1)).
		SetMaxLengthBytes(stream.ByteCapacity{}.MB(2)),
	)

	if err != nil {
		panic(err)
	}

	// crete a new producer...
	producerOptions := stream.NewProducerOptions()
	producerOptions.SetProducerName("producer")

	// batch the 100 events in the same frame, the the SDK will handle the batching...
	// it means how many messages will be batched in the same frame...
	producerOptions.SetSubEntrySize(100)
	producerOptions.SetCompression(stream.Compression{}.Gzip())

	producer, err := env.NewProducer(EVENTSTREAM, producerOptions)
	if err != nil {
		panic(err)
	}

	// ctx := context.Background()

	// publish 6001 messages to the stream...
	for i := 0; i < 6001; i++ {
		event := Event{
			Message: "HI MOM!",
		}

		data, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}

		message := amqp.NewMessage(data)

		// apply the prorperties to the message...
		props := &amqp.MessageProperties{
			CorrelationID: uuid.NewString(),
		}

		message.Properties = props

		// if the message queue is doesn't have "subentry" and compression then,
		// assigning the publishing id will not allow the rabbitmq to store duplicate messages...
		message.SetPublishingId(int64(i))

		// send the message to the stream...
		if err := producer.Send(message); err != nil {
			panic(err)
		}
	}

	producer.Close()
}

// but subentry and compression will indeed speed the data processing but
// it will come with the cost of extra computation & decompression... thing...
// and even after using "subentry" the messages will be ducplicated in the stream...

// so, it depends on the use case and business logic to use the "subentry" and "compression"...

// the banchmarks by rabbitmq, u to 1M messages per second can be achieved with the "subentry" and "compression"...
// and without "subentry" and "compression" the rate will be 100K messages per second...
