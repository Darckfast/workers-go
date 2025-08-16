package queuehandler

import (
	"strings"

	"github.com/syumai/workers/cloudflare/kv"
	"github.com/syumai/workers/cloudflare/queues"
)

func New() {
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
