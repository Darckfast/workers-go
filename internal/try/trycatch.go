//go:build js && wasm

package jstry

import (
	"syscall/js"
)

func init() {
	if !js.Global().Get("tryCatch").Truthy() {

		fn := js.Global().Get("Function")

		// Sync fn with JS error normalization: https://github.com/Darckfast/catch-this/blob/f96ceb6b7b5060281a152494353e07f9d6cbf0ca/index.js#L1
		tryCatchFn := fn.New("fn", `{
      try {
      return {
      result: fn(),
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

		js.Global().Set("tryCatch", tryCatchFn)
	}
}

func TryCatch(fn js.Func) (js.Value, error) {
	fnResultVal := js.Global().Call("tryCatch", fn)
	resultVal := fnResultVal.Get("result")
	errorVal := fnResultVal.Get("error")
	if !errorVal.IsUndefined() {
		return js.Value{}, js.Error{Value: errorVal}
	}
	return resultVal, nil
}
