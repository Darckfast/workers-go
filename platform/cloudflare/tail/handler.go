//go:build js && wasm

package tail

import (
	"context"
	"errors"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
)

type TailConsumer func(ctx context.Context, traces *Traces) error

var consumer TailConsumer = func(_ context.Context, _ *Traces) error {
	return errors.New("no consumer implemented")
}

func init() {
	var promise = js.FuncOf(func(this js.Value, args []js.Value) any {
		events := args[0]
		jsenv := args[1]
		jsctx := args[2]
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]

			go func() {
				err := handler(events, jsenv, jsctx)
				if err != nil {
					reject.Invoke(jsclass.ToJSError(err))
				} else {
					resolve.Invoke(true)
				}
			}()

			return nil
		})

		return jsclass.Promise.New(cb)
	})

	jsclass.CF.Set("tail", promise)
}

func handler(eventsObj, envObj, ctxObj js.Value) error {
	jsclass.Ctx.Init(ctxObj)
	jsclass.Env.LoadEnvs(envObj)

	events, err := NewEvents(eventsObj)

	if err != nil {
		return err
	}

	ctx := context.Background()
	return consumer(ctx, events)
}

func ConsumeNonBlock(c TailConsumer) {
	consumer = c
}
