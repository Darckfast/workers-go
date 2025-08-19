//go:build js && wasm

package jstry

import (
	"syscall/js"
)

var catchThis = js.Global().Get("tryCatch")

func init() {
	if !catchThis.Truthy() {
		fn := js.Global().Get("Function")

		// Sync fn with JS error normalization
		// This is only to go test run without crashing
		catchThis = fn.New("fn", `{
      try {
      return {
      data: fn(),
      };
      } catch (e) {
      if (!(e instanceof Error)) {
      if (e instanceof Object) {
      e = JSON.stringify(e)
      }
      e = new Error(e || 'no error message')
      }
      return {
      error: e,
      };
      }
      }`)
	}
}

func TryCatch(fn js.Func) (js.Value, error) {
	fnResult := catchThis.Invoke(fn)
	resultVal := fnResult.Get("data")
	errorVal := fnResult.Get("error")

	if !errorVal.IsUndefined() {
		return js.Value{}, js.Error{Value: errorVal}
	}

	return resultVal, nil
}
