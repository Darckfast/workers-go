//go:build js && wasm

package main

import (
	"context"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/queues"
)

func main() {
	queues.ConsumeNonBlock(func(ctx context.Context, batch *queues.MessageBatch) error {
		for _, msg := range batch.Messages {

			b, _ := msg.StringBody()
			println("message body:" + b)

			msg.Ack()
		}

		return nil
	})

	<-make(chan struct{})
}
