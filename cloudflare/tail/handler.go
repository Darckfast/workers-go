//go:build js && wasm

package tail

import (
	"errors"
	"syscall/js"

	"github.com/Darckfast/workers-go/cloudflare/env"
	"github.com/Darckfast/workers-go/cloudflare/lifecycle"
	jsclass "github.com/Darckfast/workers-go/internal/class"
)

type TailConsumer func(f *[]TraceItem) error

var consumer TailConsumer = func(_ *[]TraceItem) error {
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

	js.Global().Get("cf").Set("tail", promise)
}

func handler(eventsObj, envObj, ctxObj js.Value) error {
	lifecycle.Env = envObj
	lifecycle.Ctx = jsclass.ExecutionContextWrap{Ctx: ctxObj}

	err := env.LoadEnvs()
	if err != nil {
		return err
	}

	events := NewEvents(eventsObj)

	return consumer(events)
}

func ConsumeNonBlock(c TailConsumer) {
	consumer = c
}
