//go:build js && wasm

package lifecycle

import (
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
)

var Ctx = jsclass.ExecutionContextWrap{Ctx: js.Global().Get("ctx")}
var Env = js.Global().Get("env")
