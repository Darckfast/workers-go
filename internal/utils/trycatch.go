package jsutil

import (
	"errors"
	"syscall/js"
)

func TryCatch(fn js.Func) (js.Value, error) {
	fnResultVal := js.Global().Call("tryCatch", fn)
	resultVal := fnResultVal.Get("result")
	errorVal := fnResultVal.Get("error")
	if !errorVal.IsUndefined() {
		return js.Value{}, errors.New(errorVal.String())
	}
	return resultVal, nil
}
