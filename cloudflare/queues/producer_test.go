//go:build js && wasm

package queues

import (
	"errors"
	"fmt"
	"syscall/js"
	"testing"
	"time"

	jsclass "github.com/Darckfast/workers-go/internal/class"
)

func validatingProducer(t *testing.T, validateFn func(message js.Value, options js.Value) error) *Producer {
	sendFn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		sendArg := args[0] // this should be batch (in case of SendBatch) or a single message (in case of Send)
		var options js.Value
		if len(args) > 1 {
			options = args[1]
		}
		return jsclass.Promise.New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			go func() {
				if err := validateFn(sendArg, options); err != nil {
					// must be non-fatal to avoid a deadlock
					t.Errorf("validation failed: %v", err)
				}
				resolve.Invoke(js.Undefined())
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
		err := producer.SendText("hello")
		if err != nil {
			t.Fatalf("Send failed: %v", err)
		}
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
		err := producer.SendJSON("hello")
		if err != nil {
			t.Fatalf("Send failed: %v", err)
		}
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
	err := producer.SendBatch(batch)
	if err != nil {
		t.Fatalf("SendBatch failed: %v", err)
	}
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
	err := producer.SendBatch(batch, WithBatchDelaySeconds(5*time.Second))
	if err != nil {
		t.Fatalf("SendBatch failed: %v", err)
	}
}
