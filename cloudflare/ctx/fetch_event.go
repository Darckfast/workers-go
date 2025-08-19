//go:build js && wasm

package contx

import (
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
)

// WaitUntil extends the lifetime of the "fetch" event.
// It accepts an asynchronous task which the Workers runtime will execute before the handler terminates but without blocking the response.
// see: https://developers.cloudflare.com/workers/runtime-apis/fetch-event/#waituntil
func WaitUntil(task func()) {
	jsclass.ExcutionContext.Call("waitUntil", jsclass.Promise.New(js.FuncOf(func(this js.Value, pArgs []js.Value) any {
		resolve := pArgs[0]
		go func() {
			task()
			resolve.Invoke(js.Undefined())
		}()
		return js.Undefined()
	})))
}

// PassThroughOnException prevents a runtime error response when the Worker script throws an unhandled exception.
// Instead, the request forwards to the origin server as if it had not gone through the worker.
// see: https://developers.cloudflare.com/workers/runtime-apis/fetch-event/#passthroughonexception
func PassThroughOnException() {
	jsclass.Await(jsclass.Promise.New(js.FuncOf(func(this js.Value, pArgs []js.Value) any {
		resolve := pArgs[0]
		go func() {
			jsclass.ExcutionContext.Call("passThroughOnException")
			resolve.Invoke(js.Undefined())
		}()
		return js.Undefined()
	})))
}
