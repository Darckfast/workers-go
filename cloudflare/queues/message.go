//go:build js && wasm

package queues

import (
	"errors"
	"syscall/js"
	"time"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
)

// Message represents a message of the batch received by the consumer.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#message
type Message struct {
	// instance - The underlying instance of the JS message object passed by the cloudflare
	instance js.Value

	// ID - The unique Cloudflare-generated identifier of the message
	ID string
	// Timestamp - The time when the message was enqueued
	Timestamp time.Time
	// Body - The message body. Could be accessed directly or using converting helpers as StringBody, BytesBody, IntBody, FloatBody.
	Body js.Value
	// Attempts - The number of times the message delivery has been retried.
	Attempts int
}

func newMessage(obj js.Value) (*Message, error) {
	timestamp, err := jsconv.DateToTime(obj.Get("timestamp"))
	if err != nil {
		return nil, errors.New("failed to parse message timestamp: " + err.Error())
	}

	return &Message{
		instance:  obj,
		ID:        obj.Get("id").String(),
		Body:      obj.Get("body"),
		Attempts:  obj.Get("attempts").Int(),
		Timestamp: timestamp,
	}, nil
}

// Ack acknowledges the message as successfully delivered despite the result returned from the consuming function.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#message
func (m *Message) Ack() {
	m.instance.Call("ack")
}

// Retry marks the message to be re-delivered.
// The message will be retried after the optional delay configured with RetryOption.
func (m *Message) Retry(opts ...RetryOption) {
	var o *retryOptions
	if len(opts) > 0 {
		o = &retryOptions{}
		for _, opt := range opts {
			opt(o)
		}
	}

	m.instance.Call("retry", o.toJS())
}

func (m *Message) StringBody() (string, error) {
	if m.Body.Type() != js.TypeString {
		return "", errors.New("message body is not a string: " + m.Body.Type().String())
	}
	return m.Body.String(), nil
}

func (m *Message) BytesBody() ([]byte, error) {
	if m.Body.Type() != js.TypeObject ||
		!(m.Body.InstanceOf(jsclass.Uint8Array) || m.Body.InstanceOf(jsclass.Uint8ClampedArray)) {
		return nil, errors.New("message body is not a byte array: " + m.Body.Type().String())
	}
	b := make([]byte, m.Body.Get("byteLength").Int())
	js.CopyBytesToGo(b, m.Body)
	return b, nil
}
