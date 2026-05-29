//go:build js && wasm

package main

import (
	"log"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/queues"
)

func main() {
	queues.ConsumeNonBlock(func(batch *queues.MessageBatch) error {
		for _, msg := range batch.Messages {

			b, _ := msg.StringBody()
			log.Println("message body:", b)

			msg.Ack()
		}

		return nil
	})

	<-make(chan struct{})
}
