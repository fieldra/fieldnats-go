// Package fieldnats provides a NATS client implementation with support for protobuf message handling.
// It simplifies the process of publishing and subscribing to messages using Protocol Buffers
// for message serialization and deserialization.
package fieldnats

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

// ProtoHandler is a function type that defines the signature for handling protobuf messages.
// It processes a decoded protobuf message along with its NATS message context and returns an error if processing fails.
type ProtoHandler func(ctx context.Context, msg proto.Message, m jetstream.Msg) error

// Subscriber interface defines the contract for types that can handle protobuf message subscriptions.
// It requires implementations to provide both the protobuf message type and a handler function.
type Subscriber interface {
	// ProtoMessage returns a new instance of the protobuf message type that this consumer handles
	ProtoMessage() proto.Message
	// Handler processes the received protobuf message and its associated NATS message
	Handler(ctx context.Context, msg proto.Message, m jetstream.Msg) error
}

// Subscribe sets up a subscription to a NATS subject with protobuf message handling.
// It automatically decodes incoming messages using the provided protobuf message factory
// and processes them with the specified handler.
//
// Parameters:
//   - subject: The NATS subject to subscribe to
//   - stream: The stream name for the subscription (for JetStream persistence)
//   - durable: The durable name for the subscription (for JetStream persistence)
//   - factory: A function that creates new instances of the protobuf message type
//   - handler: A function that processes decoded protobuf messages
//   - opts: Optional subscription options that override the default settings
//
// Default behavior:
//   - Uses durable subscriptions for message persistence
//   - Requires manual message acknowledgment
//   - Sets a 30-second acknowledgment timeout
//
// Returns:
//   - error: Returns an error if the subscription setup fails
func (n *fieldNats) Subscribe(
	ctx context.Context,
	subject string,
	stream string,
	durable string,
	factory func() proto.Message,
	handler ProtoHandler,
	opts ...jetstream.PullConsumeOpt,
) error {
	jetStream, err := n.js.Stream(ctx, stream)
	if err != nil {
		return err
	}

	consumer, err := jetStream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          durable,
		Durable:       durable,
		AckWait:       30 * time.Second,
		FilterSubject: subject,
	})
	if err != nil {
		return err
	}

	// Subscribe to the subject with the provided options
	_, err = consumer.Consume(func(m jetstream.Msg) {
		// Create a new instance of the protobuf message
		msg := factory()
		if err := proto.Unmarshal(m.Data(), msg); err != nil {
			if n.cfg.Debug {
				log.Printf("🚨 fieldnats: failed to decode protobuf: %v", err)
			}

			_ = m.Nak() // NACK to let NATS know we couldn't process the message
			return
		}

		// Call the handler to process the message
		if err := handler(ctx, msg, m); err != nil {
			if n.cfg.Debug {
				log.Printf("🚨 fieldnats: handler error: %v", err)
			}

			_ = m.Nak() // NACK if the handler fails
			return
		}
	}, opts...)

	if err != nil {
		if n.cfg.Debug {
			log.Printf("❌ fieldnats: failed to subscribe to subject: %s: %v", subject, err)
		}
		return err
	}

	if n.cfg.Debug {
		log.Printf("🚀 fieldnats: successfully subscribed to subject: %s", subject)
	}

	return err
}
