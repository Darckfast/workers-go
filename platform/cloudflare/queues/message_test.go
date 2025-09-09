//go:build js && wasm

package queues

import (
	"bytes"
	"syscall/js"
	"testing"
	"time"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumerMessage(t *testing.T) {
	ts := time.Now()
	jsTs := jsconv.TimeToDate(ts)
	id := "some-message-id"
	m := map[string]any{
		"body":      "hello",
		"timestamp": jsTs,
		"id":        id,
		"attempts":  1,
	}

	got, err := newMessage(js.ValueOf(m))
	assert.Nil(t, err)
	assert.Equal(t, "hello", got.Body.String())
	assert.Equal(t, id, got.ID)
	assert.Equal(t, 1, got.Attempts)
	assert.True(t, ts.Equal(got.Timestamp))
}

func TestConsumerMessage_Ack(t *testing.T) {
	ackCalled := false
	jsObj := jsclass.Object.New()
	jsObj.Set("ack", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ackCalled = true
		return nil
	}))
	m := &Message{
		instance: jsObj,
	}

	m.Ack()

	assert.True(t, ackCalled)
}

func TestConsumerMessage_Retry(t *testing.T) {
	retryCalled := false
	jsObj := jsclass.Object.New()
	jsObj.Set("retry", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		retryCalled = true
		return nil
	}))
	m := &Message{
		instance: jsObj,
	}

	m.Retry()

	assert.True(t, retryCalled)
}

func TestConsumerMessage_RetryWithDelay(t *testing.T) {
	retryCalled := false
	jsObj := jsclass.Object.New()
	jsObj.Set("retry", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		retryCalled = true
		if len(args) != 1 {
			t.Fatalf("retry() called with %d arguments, want 1", len(args))
		}

		opts := args[0]
		if opts.Type() != js.TypeObject {
			t.Fatalf("retry() called with argument of type %v, want object", opts.Type())
		}

		if delay := opts.Get("delaySeconds").Int(); delay != 10 {
			t.Fatalf("delaySeconds = %v, want %v", delay, 10)
		}

		return nil
	}))

	m := &Message{
		instance: jsObj,
	}

	m.Retry(WithRetryDelay(10 * time.Second))

	if !retryCalled {
		t.Fatalf("RetryAll() did not call retryAll")
	}
}

func TestNewConsumerMessage_StringBody(t *testing.T) {
	tests := []struct {
		name    string
		body    func() js.Value
		want    string
		wantErr bool
	}{
		{
			name: "string",
			body: func() js.Value {
				return js.ValueOf("hello")
			},
			want: "hello",
		},
		{
			name: "uint8 array",
			body: func() js.Value {
				v := jsclass.Uint8Array.New(3)
				js.CopyBytesToJS(v, []byte("foo"))
				return v
			},
			wantErr: true,
		},
		{
			name: "int",
			body: func() js.Value {
				return js.ValueOf(42)
			},
			wantErr: true,
		},
		{
			name: "undefined",
			body: func() js.Value {
				return js.Undefined()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Body: tt.body(),
			}

			got, err := m.StringBody()
			if (err != nil) != tt.wantErr {
				t.Fatalf("StringBody() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Fatalf("StringBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsumerMessage_BytesBody(t *testing.T) {
	tests := []struct {
		name    string
		body    func() js.Value
		want    []byte
		wantErr bool
	}{
		{
			name: "uint8 array",
			body: func() js.Value {
				v := jsclass.Uint8Array.New(3)
				js.CopyBytesToJS(v, []byte("foo"))
				return v
			},
			want: []byte("foo"),
		},
		{
			name: "uint8 clamped array",
			body: func() js.Value {
				v := jsclass.Uint8ClampedArray.New(3)
				js.CopyBytesToJS(v, []byte("bar"))
				return v
			},
			want: []byte("bar"),
		},
		{
			name: "incorrect type",
			body: func() js.Value {
				return js.ValueOf("hello")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Body: tt.body(),
			}

			got, err := m.BytesBody()
			if (err != nil) != tt.wantErr {
				t.Fatalf("BytesBody() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !bytes.Equal(got, tt.want) {
				t.Fatalf("BytesBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
