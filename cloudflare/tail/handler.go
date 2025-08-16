//go:build js && wasm

package tail

import (
	"fmt"
	"syscall/js"

	jstail "github.com/syumai/workers/internal/tail"
	jsutil "github.com/syumai/workers/internal/utils"
)

type TailConsumer func(f *[]jstail.TailEvent) error

var consumer TailConsumer = func(_ *[]jstail.TailEvent) error {
	return fmt.Errorf("no consumer implemented")
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
					reject.Invoke(jsutil.Error(err.Error()))
				} else {
					resolve.Invoke(true)
				}
			}()

			return nil
		})

		return jsutil.NewPromise(cb)
	})

	js.Global().Get("cf").Set("tail", promise)
}

func handler(eventsObj, envObj, ctxObj js.Value) error {
	jsutil.RuntimeEnv = envObj
	jsutil.RuntimeExcutionContext = ctxObj

	events := jstail.NewEvents(eventsObj)

	return consumer(events)
}

func ConsumeNonBlock(c TailConsumer) {
	consumer = c
}
