//go:build js && wasm

package jstry

import (
	"syscall/js"
)

var catchThis js.Value

func init() {
	catchThis = js.Global().Get("Function").New("o", "fn", "args", `{
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

func TryCatch(o js.Value, fn string, args ...any) (js.Value, error) {
	fnResult := catchThis.Invoke(o, fn, args)
	resultVal := fnResult.Get("data")
	errorVal := fnResult.Get("error")

	if errorVal.Truthy() {
		return js.Undefined(), js.Error{Value: errorVal}
	}

	return resultVal, nil
}
