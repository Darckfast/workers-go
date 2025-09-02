//go:build js && wasm

package jsclass

import (
	"syscall/js"

	jstry "github.com/Darckfast/workers-go/internal/try"
)

// JSON.stringify and JSON.parse
type JSONWrap struct {
	js.Value
}

func (j *JSONWrap) Stringify(args ...any) js.Value {
	return j.Call("stringify", args...)
}

func (j *JSONWrap) Parse(args ...any) (js.Value, error) {
	cb := js.FuncOf(func(_ js.Value, _ []js.Value) any {
		return j.Call("parse", args...)
	})

	defer cb.Release()

	return jstry.TryCatch(cb)
}

// Object.fromEntries
type ObjectWrap struct {
	js.Value
}

func (o *ObjectWrap) FromEntries(args ...any) js.Value {
	return o.Call("fromEntries", args...)
}

type ExecutionContextWrap struct {
	Ctx js.Value
}

func (e *ExecutionContextWrap) WaitUntil(task func() error) {
	var cb js.Func

	cb = js.FuncOf(func(this js.Value, pArgs []js.Value) any {
		resolve := pArgs[0]
		reject := pArgs[1]

		go func() {
			defer cb.Release()

			err := task()
			if err != nil {
				reject.Invoke(ToJSError(err))
			} else {
				resolve.Invoke(true)
			}
		}()

		return nil
	})

	e.Ctx.Call("waitUntil", Promise.New(cb))
}

func (e *ExecutionContextWrap) PassThroughOnException() {
	e.Ctx.Call("passThroughOnException")
}

type ConsoleWrap struct {
	js.Value
}

func (c *ConsoleWrap) Log(v js.Value) {
	c.Call("log", v)
}
