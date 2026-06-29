//go:build js && wasm

/*
Package rpc is the glue code for Cloudflare's Worker RPC stubs
*/
package rpc

import (
	"context"
	"io"
	"net/http"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/bind"
)

var jsProto js.Value

func initProto() {
	if !jsProto.Truthy() {
		jsWorkerapp := js.Global().Get("workerapp")

		if !jsWorkerapp.Truthy() {
			println("using RPC but globalThis.workerapp is undefined")
		}

		jsProto = jsclass.Object.GetPrototypeOf(jsWorkerapp)

		bind.Env.LoadEnvs(js.Value{})
		bind.Ctx.Init(js.Value{})
	}
}

type RPCStubStreamFunc func(c context.Context, w http.ResponseWriter, body io.ReadCloser, args [][]byte)
type RPCStubFunc func(c context.Context, args [][]byte) [][]byte

func RPCStub(name string, h RPCStubFunc) {
	initProto()

	var hrp = js.FuncOf(func(this js.Value, args []js.Value) any {
		bufs := make([][]byte, len(args))

		for i, a := range args {
			if a.Truthy() {
				bufs[i] = make([]byte, a.Length())
				js.CopyBytesToGo(bufs[i], a)
			} else {
				bufs[i] = nil
			}
		}

		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]

			go func() {
				ctx := context.Background()

				out := h(ctx, bufs)

				dstArr := jsclass.Array.New(len(out))
				for i, o := range out {
					dst := jsclass.Uint8Array.New(len(o))
					js.CopyBytesToJS(dst, o)
					dstArr.SetIndex(i, dst)
				}

				resolve.Invoke(dstArr)
			}()

			return nil
		})

		return jsclass.Promise.New(cb)
	})

	jsProto.Set(name, hrp)
}
