//go:build js && wasm

/*
Package rpc is the glue code for Cloudflare's Worker RPC stubs
*/
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
)

var jsProto js.Value

func initProto() {
	if !jsProto.Truthy() {
		jsWorkerapp := js.Global().Get("workerapp")

		if !jsWorkerapp.Truthy() {
			log.Panicln("using RPC but globalThis.workerapp is undefined")
		}

		jsProto = jsclass.Object.GetPrototypeOf(jsWorkerapp)

		err := env.LoadEnvs()
		if err != nil {
			log.Println("error loading envs: " + err.Error())
		}
	}
}

type RPCStubStreamFunc func(c context.Context, w http.ResponseWriter, body io.ReadCloser, args [][]byte)
type RPCStubFunc func(c context.Context, args [][]byte) [][]byte

func RPCStubStream(name string, h RPCStubStreamFunc) {
	initProto()

	var hrp = js.FuncOf(func(this js.Value, args []js.Value) any {
		bufs := make([][]byte, len(args[1:]))

		for i, a := range args[1:] {
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

	jsProto.Set(name, hrp)
}

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
