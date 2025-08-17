package jsclass

import (
	"syscall/js"
)

func init() {
	cf := js.Global().Get("cf")
	if cf.IsUndefined() {
		// is faster to JSON.parse than create an object
		cfObj, _ := JSON.Parse(`{"ctx":{},"env":{},"handlers":{},"connect":{}}`)
		js.Global().Set("cf", cfObj)

		return
	}

	Env = cf.Get("env")
	ExcutionContext = cf.Get("ctx")
	Connect = cf.Get("connect")
}

var (
	Env               js.Value
	ExcutionContext   js.Value
	Connect           js.Value // replace with js.Global().Get("import")
	JSON              = JSONWrap{js.Global().Get("JSON")}
	Object            = ObjectWrap{js.Global().Get("Object")}
	Caches            = js.Global().Get("caches")
	Promise           = js.Global().Get("Promise")
	Request           = js.Global().Get("Request")
	Boolean           = js.Global().Get("Boolean")
	Response          = js.Global().Get("Response")
	Headers           = js.Global().Get("Headers")
	String            = js.Global().Get("String")
	Array             = js.Global().Get("Array")
	Number            = js.Global().Get("Number")
	Uint8Array        = js.Global().Get("Uint8Array")
	Uint8ClampedArray = js.Global().Get("Uint8ClampedArray")
	Error             = js.Global().Get("Error")
	ReadableStream    = js.Global().Get("ReadableStream")
	Date              = js.Global().Get("Date")
	Null              = js.ValueOf(nil)
	// MaybeFixedLengthStream is a class for FixedLengthStream.
	// * This class is only available in Cloudflare Workers.
	// * If this class is not available, the value will be undefined.
	MaybeFixedLengthStream = js.Global().Get("FixedLengthStream")
)

func ToJSError(err error) js.Value {
	return Error.New(err.Error())
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
		jsErr := args[0]
		if !jsErr.InstanceOf(Error) {
			if jsErr.InstanceOf(Object.Value) {
				jsErr = JSON.Stringify(jsErr)
			}

			jsErr = Error.New(jsErr)
		}
		errCh <- js.Error{Value: jsErr}
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
