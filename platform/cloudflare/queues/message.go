//go:build js && wasm

package queues

import (
	"errors"
	"syscall/js"
	"time"

	jsclass "codeberg.org/darckfast/workers-go/internal/class"
	jsconv "codeberg.org/darckfast/workers-go/internal/conv"
)

// Message represents a message of the batch received by the consumer.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#message
type Message struct {
	Timestamp time.Time
	instance  js.Value
	Body      js.Value
	ID        string
	Attempts  int
}

func newMessage(obj js.Value) (*Message, error) {
	timestamp := jsconv.DateToTime(obj.Get("timestamp"))

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
	if m.Body.InstanceOf(jsclass.ArrayBuffer) {
		b := make([]byte, m.Body.Get("byteLength").Int())
		js.CopyBytesToGo(b, jsclass.Uint8Array.New(m.Body))
		return b, nil
	}

	if m.Body.Type() == js.TypeObject &&
		(m.Body.InstanceOf(jsclass.Uint8Array) ||
			m.Body.InstanceOf(jsclass.Uint8ClampedArray)) {
		b := make([]byte, m.Body.Get("byteLength").Int())
		js.CopyBytesToGo(b, m.Body)
		return b, nil
	}
	return nil, errors.New("message body is not a byte array: " + m.Body.Type().String())
}
