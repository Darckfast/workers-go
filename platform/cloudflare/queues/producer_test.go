//go:build js && wasm

package queues

import (
	"errors"
	"fmt"
	"syscall/js"
	"testing"
	"time"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
)

var startTime = time.Now().Unix()

func validatingProducer(t *testing.T, validateFn func(message js.Value, options js.Value) error) *Producer {
	sendFn := js.FuncOf(func(this js.Value, sargs []js.Value) interface{} {
		sendArg := sargs[0] // this should be batch (in case of SendBatch) or a single message (in case of Send)
		var options js.Value
		if len(sargs) > 1 {
			options = sargs[1]
		}
		return jsclass.Promise.New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]

			go func() {
				uint := js.Global().Get("TextEncoder").New().Call("encode", jsclass.JSON.Stringify(sendArg))
				sr := QueueSendResult{
					Metadata: Metadata{
						Metrics: QueueMetrics{
							BacklogCount:           int64(len(sargs) - 1),
							BacklogBytes:           int64(uint.Length()),
							OldestMessageTimestamp: startTime,
						},
					},
				}

				bj, _ := easyjson.Marshal(sr)
				jsV, _ := jsclass.JSON.Parse(string(bj))
				if err := validateFn(sendArg, options); err != nil {
					// must be non-fatal to avoid a deadlock
					t.Errorf("validation failed: %v", err)
				}
				resolve.Invoke(jsV)
			}()
			return js.Undefined()
		}))
	})

	queue := jsclass.Object.New()
	queue.Set("send", sendFn)
	queue.Set("sendBatch", sendFn)

	return &Producer{queue: queue}
}

func TestSend(t *testing.T) {
	t.Run("text content type", func(t *testing.T) {
		validation := func(message js.Value, options js.Value) error {
			if message.Type() != js.TypeString {
				return errors.New("message body must be a string")
			}
			if message.String() != "hello" {
				return errors.New("message body must be 'hello'")
			}
			if options.Get("contentType").String() != "text" {
				return errors.New("content type must be text")
			}
			return nil
		}

		producer := validatingProducer(t, validation)
		evt, err := producer.SendText("hello")
		if err != nil {
			t.Fatalf("Send failed: %v", err)
		}

		assert.Equal(t, int64(1), evt.Metadata.Metrics.BacklogCount)
		assert.Equal(t, int64(7), evt.Metadata.Metrics.BacklogBytes)
		assert.GreaterOrEqual(t, startTime, evt.Metadata.Metrics.OldestMessageTimestamp)
	})

	t.Run("json content type", func(t *testing.T) {
		validation := func(message js.Value, options js.Value) error {
			if message.Type() != js.TypeString {
				return errors.New("message body must be a string")
			}
			if message.String() != "hello" {
				return errors.New("message body must be 'hello'")
			}
			if options.Get("contentType").String() != "json" {
				return errors.New("content type must be json")
			}
			return nil
		}

		producer := validatingProducer(t, validation)
		evt, err := producer.SendJSON("hello")
		if err != nil {
			t.Fatalf("Send failed: %v", err)
		}

		assert.Equal(t, int64(1), evt.Metadata.Metrics.BacklogCount)
		assert.Equal(t, int64(7), evt.Metadata.Metrics.BacklogBytes)
		assert.GreaterOrEqual(t, startTime, evt.Metadata.Metrics.OldestMessageTimestamp)
	})
}

func TestSendBatch(t *testing.T) {
	validation := func(batch js.Value, options js.Value) error {
		if batch.Type() != js.TypeObject {
			return errors.New("message batch must be an object (array)")
		}
		if batch.Length() != 2 {
			return fmt.Errorf("expected 2 messages, got %d", batch.Length())
		}
		first := batch.Index(0)
		if first.Get("body").String() != "hello" {
			return fmt.Errorf("first message body must be 'hello', was %s", first.Get("body"))
		}
		if first.Get("options").Get("contentType").String() != "json" {
			return fmt.Errorf("first message content type must be json, was %s", first.Get("options").Get("contentType"))
		}

		second := batch.Index(1)
		if second.Get("body").String() != "world" {
			return fmt.Errorf("second message body must be 'world', was %s", second.Get("body"))
		}
		if second.Get("options").Get("contentType").String() != "text" {
			return fmt.Errorf("second message content type must be text, was %s", second.Get("options").Get("contentType"))
		}

		return nil
	}

	batch := []*MessageSendRequest{
		NewJSONMessageSendRequest("hello"),
		NewTextMessageSendRequest("world"),
	}

	producer := validatingProducer(t, validation)
	evt, err := producer.SendBatch(batch)
	if err != nil {
		t.Fatalf("SendBatch failed: %v", err)
	}

	assert.Equal(t, int64(1), evt.Metadata.Metrics.BacklogCount)
	assert.Equal(t, int64(101), evt.Metadata.Metrics.BacklogBytes)
	assert.GreaterOrEqual(t, startTime, evt.Metadata.Metrics.OldestMessageTimestamp)
}

func TestSendBatch_Options(t *testing.T) {
	validation := func(_ js.Value, options js.Value) error {
		if options.Get("delaySeconds").Int() != 5 {
			return fmt.Errorf("expected delay 5, got %d", options.Get("delaySeconds").Int())
		}
		return nil
	}

	batch := []*MessageSendRequest{
		NewTextMessageSendRequest("hello"),
	}

	producer := validatingProducer(t, validation)
	evt, err := producer.SendBatch(batch, WithBatchDelaySeconds(5*time.Second))
	if err != nil {
		t.Fatalf("SendBatch failed: %v", err)
	}

	assert.Equal(t, int64(1), evt.Metadata.Metrics.BacklogCount)
	assert.Equal(t, int64(51), evt.Metadata.Metrics.BacklogBytes)
	assert.GreaterOrEqual(t, startTime, evt.Metadata.Metrics.OldestMessageTimestamp)
}
