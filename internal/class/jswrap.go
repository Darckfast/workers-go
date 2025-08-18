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
	return jstry.TryCatch(js.FuncOf(func(_ js.Value, _ []js.Value) any {
		return j.Call("parse", args...)
	}))
}

// Object.fromEntries
type ObjectWrap struct {
	js.Value
}

func (o *ObjectWrap) FromEntries(args ...any) js.Value {
	return o.Call("fromEntries", args...)
}
