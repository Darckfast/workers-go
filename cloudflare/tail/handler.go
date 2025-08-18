//go:build js && wasm

package tail

import (
	"errors"
	"syscall/js"

	"github.com/Darckfast/workers-go/cloudflare/env"
	jsclass "github.com/Darckfast/workers-go/internal/class"
	jstail "github.com/Darckfast/workers-go/internal/tail"
)

type TailConsumer func(f *[]jstail.TailEvent) error

var consumer TailConsumer = func(_ *[]jstail.TailEvent) error {
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
	jsclass.Env = envObj
	jsclass.ExcutionContext = ctxObj

	env.LoadEnvs()
	events := jstail.NewEvents(eventsObj)

	return consumer(events)
}

func ConsumeNonBlock(c TailConsumer) {
	consumer = c
}
