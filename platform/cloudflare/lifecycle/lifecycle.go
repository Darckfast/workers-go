//go:build js && wasm

/*
Package lifecycle is the glue code for Cloudflare's Env and Ctx
*/
package lifecycle

import (
	"syscall/js"

	jsclass "codeberg.org/darckfast/workers-go/internal/class"
)

var Ctx jsclass.ExecutionContextWrap
var Env js.Value

func init() {
	// Auto init for class workers
	var workerapp = js.Global().Get("workerapp")

	if workerapp.Truthy() {
		Ctx = jsclass.ExecutionContextWrap{Ctx: workerapp.Get("ctx")}
		Env = workerapp.Get("env")
	}
}
