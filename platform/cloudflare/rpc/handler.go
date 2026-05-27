//go:build js && wasm

package rpc

import (
	"context"
	"log"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	"github.com/Darckfast/workers-go/platform/cloudflare/env"
	"github.com/Darckfast/workers-go/platform/cloudflare/lifecycle"
)

func init() {
	js.Global().Get("cf").Set("rpc", jsclass.Object)
	lifecycle.Env = js.Global().Get("cf").Get("env")
	err := env.LoadEnvs()

	if err != nil {
		log.Println("error loading envs: " + err.Error())
	}
}

type RPCHandler func(c context.Context, args ...[]byte) ([]byte, error)

func AddRPCHandler(name string, h RPCHandler) {
	var hrp = js.FuncOf(func(this js.Value, args []js.Value) any {
		bufs := make([][]byte, len(args))

		for i, a := range args {
			bufs[i] = make([]byte, a.Get("length").Int())
			js.CopyBytesToGo(bufs[i], a)
		}

		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]

			go func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				out, err := h(ctx, bufs...)

				if err != nil {
					reject.Invoke(jsclass.ToJSError(err))
				} else {
					res := jsclass.Uint8Array.New(len(out))
					js.CopyBytesToJS(res, out)
					resolve.Invoke(res)
				}
			}()

			return nil
		})

		return jsclass.Promise.New(cb)
	})

	js.Global().Get("cf").Get("rpc").Set(name, hrp)
}
