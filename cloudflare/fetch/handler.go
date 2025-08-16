//go:build js && wasm

package fetch

import (
	"context"
	"io"
	"net/http"
	"syscall/js"

	jshttp "github.com/syumai/workers/internal/http"
	jsruntime "github.com/syumai/workers/internal/runtime"
	jsutil "github.com/syumai/workers/internal/utils"
)

var httpHandler http.Handler = http.DefaultServeMux

func init() {
	var handleRequestPromise = js.FuncOf(func(this js.Value, args []js.Value) any {
		reqObj := args[0]
		envObj := args[1]
		ctxObj := args[2]
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]

			go func() {
				res := handler(reqObj, envObj, ctxObj)
				resolve.Invoke(res)
			}()

			return nil
		})

		return jsutil.NewPromise(cb)
	})

	js.Global().Get("cf").Set("fetch", handleRequestPromise)
}

func handler(reqObj js.Value, envObj js.Value, ctxObj js.Value) js.Value {
	jsutil.RuntimeEnv = envObj
	jsutil.RuntimeExcutionContext = ctxObj

	req := jshttp.ToRequest(reqObj)
	ctx := jsruntime.New(context.Background(), reqObj)
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
			writer.Close()
		}()

		httpHandler.ServeHTTP(w, req)
	}()
	<-w.ReadyCh
	return w.ToJSResponse()
}

func ServeNonBlock(handler http.Handler) {
	if handler != nil {
		httpHandler = handler
	}
}
