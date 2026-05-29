//go:build js && wasm

package rpc

import (
	"context"
	"io"
	"log"
	"net/http"
	"syscall/js"

	jsclass "codeberg.org/darckfast/workers-go/internal/class"
	jshttp "codeberg.org/darckfast/workers-go/internal/http"
	jsruntime "codeberg.org/darckfast/workers-go/internal/runtime"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/env"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/lifecycle"
)

func init() {
	js.Global().Get("cf").Set("rpc", jsclass.Object.Value)
	lifecycle.Env = js.Global().Get("workerapp").Get("env")
	lifecycle.Ctx = jsclass.ExecutionContextWrap{Ctx: js.Global().Get("workerapp").Get("ctx")}

	err := env.LoadEnvs()

	if err != nil {
		log.Println("error loading envs: " + err.Error())
	}
}

type RPCStubStreamFunc func(c context.Context, w http.ResponseWriter, body io.ReadCloser, args [][]byte)
type RPCStubFunc func(c context.Context, args [][]byte) [][]byte

func RPCStubStream(name string, h RPCStubStreamFunc) {
	var hrp = js.FuncOf(func(this js.Value, args []js.Value) any {
		bufs := make([][]byte, len(args[1:]))

		for i, a := range args[1:] {
			bufs[i] = make([]byte, a.Length())
			js.CopyBytesToGo(bufs[i], a)
		}

		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]

			go func() {
				body := jshttp.ToBody(args[0])
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				ctx = jsruntime.New(ctx, args[0])
				reader, writer := io.Pipe()

				w := &jshttp.ResponseWriter{
					Reader:  reader,
					Writer:  writer,
					ReadyCh: make(chan struct{}),
				}

				go func() {
					defer func() {
						w.Ready()
						err := writer.Close()

						if err != nil {
							log.Println("error closing response body writer", err.Error())
						}
					}()

					h(ctx, w, body, bufs)
				}()

				<-w.ReadyCh
				resolve.Invoke(w.ToReadableStream())
			}()

			return nil
		})

		return jsclass.Promise.New(cb)
	})

	js.Global().Get("cf").Get("rpc").Set(name, hrp)
}

func RPCStub(name string, h RPCStubFunc) {
	var hrp = js.FuncOf(func(this js.Value, args []js.Value) any {
		bufs := make([][]byte, len(args))

		for i, a := range args {
			bufs[i] = make([]byte, a.Length())
			js.CopyBytesToGo(bufs[i], a)
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

	js.Global().Get("cf").Get("rpc").Set(name, hrp)
}
