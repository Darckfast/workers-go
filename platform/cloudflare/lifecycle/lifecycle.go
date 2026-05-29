//go:build js && wasm

package lifecycle

import (
	"syscall/js"

	jsclass "codeberg.org/darckfast/workers-go/internal/class"
)

var Ctx = jsclass.ExecutionContextWrap{Ctx: js.Global().Get("ctx")}
var Env = js.Global().Get("env")
