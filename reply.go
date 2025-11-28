package fieldnats

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

// Reply sets up a handler that receives protobuf request messages and responds with protobuf replies.
// - subject: Subject to listen for requests on
// - reqFactory: Function that returns a new instance of the request message type
// - handler: Function to handle the request and return a response
func (n *fieldNats) Reply(subject string, reqFactory func() proto.Message, handler func(context.Context, proto.Message) (proto.Message, error)) error {
	_, err := n.conn.Subscribe(subject, func(m *nats.Msg) {
		req := reqFactory()
		if err := proto.Unmarshal(m.Data, req); err != nil {
			if n.cfg.Debug {
				log.Printf("❌ fieldnats: failed to unmarshal request: %v", err)
			}
			return
		}

		resp, err := handler(context.Background(), req)
		if err != nil {
			if n.cfg.Debug {
				log.Printf("❌ fieldnats: request handler failed: %v", err)
			}
			// Optionally send an error message (could serialize error into protobuf)
			_ = m.Respond([]byte{})
			return
		}

		data, err := proto.Marshal(resp)
		if err != nil {
			if n.cfg.Debug {
				log.Printf("❌ fieldnats: failed to marshal response: %v", err)
			}
			return
		}

		_ = m.Respond(data)
	})

	if err != nil && n.cfg.Debug {
		log.Printf("❌ fieldnats: failed to subscribe for reply on %s: %v", subject, err)
	}

	return err
}
