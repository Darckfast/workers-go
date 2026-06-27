//go:build js && wasm

package jstry

import (
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jshelper"
)

var And LazyJSVal

type LazyJSVal struct {
	jshelper.LazyJSVal
}

func (l *LazyJSVal) init() {
	l.Do(func() {
		l.Value = js.Global().Get("tryCatch")
		if !l.Truthy() {
			// Due Clouflare Workers limitation, this function cannot be
			// instantiated within Go runtime
			l.Value = js.Global().Get("Function").New("o", "fn", "args", `{
        try {
        if (fn) {
        return { data: o[fn](...args) };
        }

        return { data: o(...args) };
        } catch (err) {
        if (!(err instanceof Error)) {
        if (err instanceof Object) {
        err = JSON.stringify(err);
        }
        err = new Error(err || "no error message");
        }
        return { error: err };
        }
        }`)
		}
	})
}

func (l *LazyJSVal) Catch(o js.Value, fn string, args ...any) (js.Value, error) {
	l.init()
	fnResult := l.Invoke(o, fn, args)
	resultVal := fnResult.Get("data")
	errorVal := fnResult.Get("error")

	if errorVal.Truthy() {
		return js.Undefined(), js.Error{Value: errorVal}
	}

	return resultVal, nil
}
