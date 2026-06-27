//go:build js && wasm

package jsclass

import (
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jshelper"
)

type global struct {
	cf jshelper.LazyJSVal
}

func (g *global) Set(p string, a any) {
	g.cf.Do(func() {
		g.cf.Value = js.Global().Get("cf")

		if !g.cf.Truthy() {
			cfObj, _ := JSON.Parse(`{}`)
			js.Global().Set("cf", cfObj)
			g.cf.Value = cfObj
		}
	})

	g.cf.Set(p, a)
}

var (
	// MaybeFixedLengthStream is a class for FixedLengthStream.
	// * This class is only available in Cloudflare Workers.
	// * If this class is not available, the value will be undefined.
	MaybeFixedLengthStream = jshelper.NewClass{ID: "FixedLengthStream"}
	CF                     = global{}
	Connect                = ConnectWrap{}
	JSON                   = JSONJS{}
	Object                 = ObjectWrap{}
	Console                = ConsoleWrap{}
	Ctx                    = ExecutionContextWrap{}
	Env                    = EnvBinding{}
	Caches                 = jshelper.CacheInterface{}
	Promise                = jshelper.PromiseClass{}
	AbortSignal            = jshelper.AbortSignal{}
	Array                  = jshelper.ArrayClass{}
	AbortController        = jshelper.NewClass{ID: "AbortController"}
	Request                = jshelper.NewClass{ID: "Request"}
	Response               = jshelper.NewClass{ID: "Response"}
	Headers                = jshelper.NewClass{ID: "Headers"}
	ArrayBuffer            = jshelper.NewClass{ID: "ArrayBuffer"}
	Uint8Array             = jshelper.NewClass{ID: "Uint8Array"}
	Uint8ClampedArray      = jshelper.NewClass{ID: "Uint8ClampedArray"}
	Error                  = jshelper.NewClass{ID: "Error"}
	ReadableStream         = jshelper.NewClass{ID: "ReadableStream"}
	Date                   = jshelper.NewClass{ID: "Date"}
	Null                   = js.ValueOf(nil)
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
		if !jsErr.InstanceOf(Error.Class()) {
			if jsErr.InstanceOf(Object.Class()) {
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
