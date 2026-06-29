//go:build js && wasm

package jshttp

import (
	"io"
	"net/http"
	"sync"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
)

type ResponseWriter struct {
	Writer      io.Writer
	V           js.Value
	w           js.Value
	HeaderValue http.Header
	ReadyCh     chan struct{}
	StatusCode  int
	Once        sync.Once
}

var (
	_ http.ResponseWriter = (*ResponseWriter)(nil)
)

// Ready indicates that ResponseWriter is ready to be converted to Response.
// If we start writing to the TransformStream() while the readable is not pulling,
// this will cause a deadlock and the JS loop will hang
func (rw *ResponseWriter) Ready() {
	rw.Once.Do(func() {
		close(rw.ReadyCh)
		rw.w = rw.V.Call("getWriter")
	})
}

func (rw *ResponseWriter) Write(data []byte) (n int, err error) {
	rw.Ready()
	_, err = jsclass.Await(rw.w.Get("ready"))
	if err != nil {
		rw.w.Call("abort", "writable ready promise returned error: "+err.Error())
		return 0, err
	}

	n = len(data)
	if n > 0 {
		b := jsclass.Uint8Array.New(n)
		js.CopyBytesToJS(b, data)
		rw.w.Call("write", b)
	}

	rw.Close()

	return n, nil
}

func (rw *ResponseWriter) Close() {
	rw.Ready()
	_, _ = jsclass.Await(rw.w.Get("ready"))
	rw.w.Call("close")
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.HeaderValue
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
}
