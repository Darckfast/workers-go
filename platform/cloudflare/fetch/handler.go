//go:build js && wasm

package fetch

import (
	"context"
	"io"
	"log"
	"net/http"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jshttp "github.com/Darckfast/workers-go/internal/http"
	jsruntime "github.com/Darckfast/workers-go/internal/runtime"
	"github.com/Darckfast/workers-go/platform/cloudflare/env"
	"github.com/Darckfast/workers-go/platform/cloudflare/lifecycle"
)

var httpHandler http.Handler = http.DefaultServeMux

func init() {
	var handleRequestPromise = js.FuncOf(func(this js.Value, args []js.Value) any {
		reqObj := args[0]
		envObj := jsclass.Null
		ctxObj := jsclass.Null

		if len(args) >= 2 {
			envObj = args[1]
		}
		if len(args) >= 3 {
			ctxObj = args[2]
		}
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]

			go func() {
				res, err := handler(reqObj, envObj, ctxObj)
				if err != nil {
					reject.Invoke(jsclass.ToJSError(err))
				} else {
					resolve.Invoke(res)
				}
			}()

			return nil
		})

		return jsclass.Promise.New(cb)
	})

	js.Global().Get("cf").Set("fetch", handleRequestPromise)
}

func handler(reqObj js.Value, envObj js.Value, ctxObj js.Value) (js.Value, error) {
	lifecycle.Env = envObj
	lifecycle.Ctx = jsclass.ExecutionContextWrap{Ctx: ctxObj}

	err := env.LoadEnvs()
	if err != nil {
		return js.Value{}, err
	}

	req := jshttp.ToRequest(reqObj)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signal := reqObj.Get("signal")
	var cbCancel js.Func
	defer cbCancel.Release()

	cbCancel = js.FuncOf(func(this js.Value, args []js.Value) any {
		cancel()
		return nil
	})

	if signal.Truthy() {
		signal.Call("addEventListener", "abort", cbCancel)
	} else {
		reqObj.Call("on", "close", cbCancel)
	}

	ctx = context.WithValue(ctx, jsruntime.CtxSignal{}, signal)

	ctx = jsruntime.New(ctx, reqObj)
	req = req.WithContext(ctx)
	reader, writer := io.Pipe()

	w := &jshttp.ResponseWriter{
		HeaderValue: http.Header{},
		StatusCode:  http.StatusOK,
		Reader:      reader,
		Writer:      writer,
		ReadyCh:     make(chan struct{}),
	}

	go func() {
		defer func() {
			w.Ready()
			err := writer.Close()

			if err != nil {
				log.Println("error closing response body writer", err.Error())
			}
		}()

		httpHandler.ServeHTTP(w, req)
	}()

	<-w.ReadyCh

	return w.ToJSResponse(), nil
}

func ServeNonBlock(handler http.Handler) {
	if handler != nil {
		httpHandler = handler
	}
}
