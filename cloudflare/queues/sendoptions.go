//go:build js && wasm

package queues

import (
	"syscall/js"
	"time"

	jsclass "github.com/Darckfast/workers-go/internal/class"
)

type sendOptions struct {
	// ContentType - Content type of the message
	// Default is "json"
	ContentType contentType

	// DelaySeconds - The number of seconds to delay the message.
	// Default is 0
	DelaySeconds int
}

func (o *sendOptions) toJS() js.Value {
	obj := jsclass.Object.New()
	obj.Set("contentType", string(o.ContentType))

	if o.DelaySeconds != 0 {
		obj.Set("delaySeconds", o.DelaySeconds)
	}

	return obj
}

type SendOption func(*sendOptions)

// WithDelaySeconds changes the number of seconds to delay the message.
func WithDelaySeconds(d time.Duration) SendOption {
	return func(o *sendOptions) {
		o.DelaySeconds = int(d.Seconds())
	}
}
