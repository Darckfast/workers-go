package queuehandler

import (
	"strings"

	"github.com/Darckfast/workers-go/cloudflare/kv"
	"github.com/Darckfast/workers-go/cloudflare/queues"
)

func New() {
	/*
	 * This functions must be called to instantiate the queue consumer, and make
	 * globalThis.cf.queue() defined on JS global scope
	 */
	queues.ConsumeNonBlock(func(batch *queues.MessageBatch) error {
		for _, msg := range batch.Messages {

			v := strings.ToUpper(msg.Body.String())

			namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
			namespace.PutString("queue:result", v, nil)

			msg.Ack()
		}

		return nil
	})
}
