package fieldnats

import (
	"context"
	"log"

	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

// Publish publishes a protobuf message to the specified NATS subject.
// It marshals the protobuf message into bytes and publishes it using JetStream.
//
// Parameters:
//   - ctx: Context for the operation (currently unused)
//   - subject: The NATS subject to publish the message to
//   - msg: The protobuf message to be published
//   - opts: Optional publishing options for NATS
//
// Returns:
//   - error: Returns an error if marshaling fails or if publishing fails
func (n *fieldNats) Publish(ctx context.Context, subject string, msg proto.Message, opts ...jetstream.PublishOpt) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		if n.cfg.Debug {
			log.Printf("❌ fieldnats: failed to encode protobuf: %v", err)
		}

		return err
	}

	ack, err := n.js.Publish(ctx, subject, data, opts...)
	if err != nil {
		if n.cfg.Debug {
			log.Printf("❌ fieldnats: failed to publish message: %v", err)
		}

		return err
	}

	if n.cfg.Debug {
		log.Printf("🚀 fieldnats: published message on domain: %s", ack.Domain)
		log.Printf("🚀 fieldnats: published message on sequence: %d", ack.Sequence)
		log.Printf("🚀 fieldnats: published message on duplicate: %v", ack.Duplicate)
		log.Printf("🚀 fieldnats: published message on stream: %s", ack.Stream)
	}

	return err
}
