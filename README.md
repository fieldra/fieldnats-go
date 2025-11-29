# FIELDNats
A fieldra package for handling events with nats jetstream

### Publishing events
An example of publishing an event is shown below:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rimdesk/fieldcore-go"
	v1 "github.com/fieldra/fieldnats-go/gen/shooters/nexor/v1"
)

var (
	client = nexor.New("nats://localhost:4222")
)

func init() {
	client.Connect()
}

func main() {
	defer client.Close()

	ctx := context.Background()

	// Initialize the event bus
	if err := client.CreateStream(ctx, jetstream.StreamConfig{
		Name:        "product_stream",
		Description: "Sample events",
		Subjects:    []string{"sample.>"},
		MaxBytes:    1024 * 1024 * 1024,
	}); err != nil {
		log.Println("🚨 [FIELDNats]: Failed to initialize event bus:", err)
		os.Exit(1)
	}

	for {
		// Get the current time
		currentTime := time.Now()
		// List of sample words
		words := []string{"Apple", "Banana", "Orange", "Mango", "Grape", "Peach", "Plum", "Cherry", "Lemon", "Lime"}
		subjects := []string{"sample.created", "sample.updated"}
		// Generate a random word
		randomWord := words[currentTime.UnixNano()%int64(len(words))]
		subject := subjects[currentTime.UnixNano()%int64(len(subjects))]

		// Create a ProductCreated event
		event := &v1.Event{
			Name: subject,
			Product: &v1.ProductCreated{
				Id:         uuid.NewString(),
				Name:       randomWord,
				SupplierId: uuid.NewString(),
				CreatedAt:  currentTime.UnixMilli(),
			},
		}

		// Publish the event
		if err := client.Publish(ctx, subject, event); err != nil {
			log.Fatalf("🚨 Failed to publish event: %v", err)
		}

		fmt.Printf("🚀 Event published to subject: %s successfully! 🚀\n", subject)

		time.Sleep(time.Duration(3) * time.Second)
	}
}

```

### Environment variables:
The default parameters can be overridden by setting the following environment variables:

```text
FIELDNATS.CLIENT=Nexor
FIELDNATS.DEBUG=false
FIELDNATS.URL=nats://localhost:4222
FIELDNATS.MAX_CONNECTIONS=5
FIELDNATS.MAX_RECONNECT_WAIT=5
```