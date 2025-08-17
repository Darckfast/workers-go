package jsclass

import (
	"fmt"
	"syscall/js"
)

var (
	Object            = js.Global().Get("Object")
	Promise           = js.Global().Get("Promise")
	JSON              = js.Global().Get("JSON")
	Request           = js.Global().Get("Request")
	Response          = js.Global().Get("Response")
	Headers           = js.Global().Get("Headers")
	String            = js.Global().Get("String")
	Array             = js.Global().Get("Array")
	Number            = js.Global().Get("Number")
	Uint8Array        = js.Global().Get("Uint8Array")
	Uint8ClampedArray = js.Global().Get("Uint8ClampedArray")
	ErrorJS           = js.Global().Get("Error")
	ReadableStream    = js.Global().Get("ReadableStream")
	Date              = js.Global().Get("Date")
	Null              = js.ValueOf(nil)
	// MaybeFixedLengthStream is a class for FixedLengthStream.
	// * This class is only available in Cloudflare Workers.
	// * If this class is not available, the value will be undefined.
	MaybeFixedLengthStream = js.Global().Get("FixedLengthStream")
)

func Error(err error) js.Value {
	return ErrorJS.New(err.Error())
}

func Await(promise js.Value) (js.Value, error) {
	resultCh := make(chan js.Value)
	defer close(resultCh)

	then := js.FuncOf(func(_ js.Value, args []js.Value) any {
		resultCh <- args[0]
		return nil
	})
	defer then.Release()

	errCh := make(chan error)
	defer close(errCh)

	catch := js.FuncOf(func(_ js.Value, args []js.Value) any {
		errCh <- fmt.Errorf("failed on promise: %s", args[0].Call("toString").String())
		return nil
	})
	defer catch.Release()

	promise.Call("then", then).Call("catch", catch)

	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errCh:
		return js.Value{}, err
	}
}
