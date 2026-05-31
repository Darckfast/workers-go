//go:build js && wasm

package queuehandler

import (
	"context"
	"strings"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/kv"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/queues"
)

func New() {
	/*
	 * ConsumeNonBlock must be called to instantiate the queue consumer, and make
	 * globalThis.cf.queue() defined on JS global scope
	 */
	queues.ConsumeNonBlock(func(c context.Context, batch *queues.MessageBatch) error {
		for _, msg := range batch.Messages {

			v := strings.ToUpper(msg.Body.String())

			namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
			err := namespace.Put("queue:result", v, nil)

			if err != nil {
				return err
			}

			msg.Ack()
		}

		return nil
	})
}
