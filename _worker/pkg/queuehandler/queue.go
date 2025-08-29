//go:build js && wasm

package queuehandler

import (
	"strings"

	"github.com/Darckfast/workers-go/cloudflare/kv"
	"github.com/Darckfast/workers-go/cloudflare/queues"
)

func New() {
	/*
	 * ConsumeNonBlock must be called to instantiate the queue consumer, and make
	 * globalThis.cf.queue() defined on JS global scope
	 */
	queues.ConsumeNonBlock(func(batch *queues.MessageBatch) error {
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
