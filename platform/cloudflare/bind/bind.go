//go:build js && wasm

package bind

import "codeberg.org/darckfast/workers-go/internal/jsclass"

var (
	Ctx = jsclass.ExecutionContextWrap{}
	Env = jsclass.EnvBinding{}
)
