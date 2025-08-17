package queues

import (
	"errors"
	"syscall/js"

	"github.com/syumai/workers/cloudflare/env"
	jsclass "github.com/syumai/workers/internal/class"
)

// Consumer is a function that received a batch of messages from Cloudflare Queues.
// The function should be set using Consume or ConsumeNonBlock.
// A returned error will cause the batch to be retried (unless the batch or individual messages are acked).
// NOTE: to do long-running message processing task within the Consumer, use cloudflare.WaitUntil, this will postpone the message
// acknowledgment until the task is completed witout blocking the queue consumption.
type Consumer func(batch *MessageBatch) error

var consumer Consumer = func(batch *MessageBatch) error {
	return errors.New("no consumer implemented")
}

func init() {
	handleBatchCallback := js.FuncOf(func(this js.Value, args []js.Value) any {
		batch := args[0]
		envObj := args[1]
		ctxObj := args[2]
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]
			go func() {
				err := consumeBatch(batch, envObj, ctxObj)
				if err != nil {
					reject.Invoke(jsclass.ToJSError(err))
					return
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})
		return jsclass.Promise.New(cb)
	})
	js.Global().Get("cf").Set("queue", handleBatchCallback)
}

func consumeBatch(batch, envObj, ctxObj js.Value) error {
	jsclass.Env = envObj
	jsclass.ExcutionContext = ctxObj

	env.LoadEnvs()
	b, err := newMessageBatch(batch)
	if err != nil {
		return errors.New("failed to parse message batch: " + err.Error())
	}

	if err := consumer(b); err != nil {
		return err
	}
	return nil
}

// ConsumeNonBlock sets the Consumer function to receive batches of messages from Cloudflare Queues.
// This function is intented to be used when the worker has other purposes (e.g. handling HTTP requests).
// The worker will not block receiving messages and will continue to execute other tasks.
// ConsumeNonBlock should be called before setting other blocking handlers (e.g. workers.Serve).
func ConsumeNonBlock(f Consumer) {
	consumer = f
}
