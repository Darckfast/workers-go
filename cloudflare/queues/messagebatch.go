//go:build js && wasm

package queues

import (
	"errors"
	"syscall/js"
)

// MessageBatch represents a batch of messages received by the consumer. The size of the batch is determined by the
// worker configuration.
//   - https://developers.cloudflare.com/queues/configuration/configure-queues/#consumer
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#messagebatch
type MessageBatch struct {
	// instance - The underlying instance of the JS message object passed by the cloudflare
	instance js.Value

	// Queue - The name of the queue from which the messages were received
	Queue string

	// Messages - The messages in the batch
	Messages []*Message
}

func newMessageBatch(obj js.Value) (*MessageBatch, error) {
	msgArr := obj.Get("messages")
	messages := make([]*Message, msgArr.Length())
	for i := 0; i < msgArr.Length(); i++ {
		m, err := newMessage(msgArr.Index(i))
		if err != nil {
			return nil, errors.New("failed to parse message " + err.Error())
		}
		messages[i] = m
	}

	return &MessageBatch{
		instance: obj,
		Queue:    obj.Get("queue").String(),
		Messages: messages,
	}, nil
}

// AckAll acknowledges all messages in the batch as successfully delivered despite the result returned from the consuming function.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#messagebatch
func (b *MessageBatch) AckAll() {
	b.instance.Call("ackAll")
}

// RetryAll marks all messages in the batch to be re-delivered.
// The messages will be retried after the optional delay configured with RetryOption.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#messagebatch
func (b *MessageBatch) RetryAll(opts ...RetryOption) {
	var o *retryOptions
	if len(opts) > 0 {
		o = &retryOptions{}
		for _, opt := range opts {
			opt(o)
		}
	}

	b.instance.Call("retryAll", o.toJS())
}
