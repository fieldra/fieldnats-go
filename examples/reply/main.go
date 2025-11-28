package main

import (
	"context"
	"fmt"

	"github.com/fieldra/fieldnats-go"
	v1 "github.com/fieldra/fieldnats-go/gen/fieldra/fieldnats/v1"
	"google.golang.org/protobuf/proto"
)

var (
	client = fieldnats.New("nats://localhost:4222")
)

func init() {
	client.Connect()
}

func main() {
	defer client.Close()

	_ = client.Reply("example.say.hello",
		func() proto.Message { return &v1.SayHelloRequest{} },
		func(ctx context.Context, req proto.Message) (proto.Message, error) {
			request := req.(*v1.SayHelloRequest)
			return &v1.SayHelloResponse{Message: fmt.Sprintf("Hello Reply %s", request.GetName())}, nil
		},
	)

	select {}
}
