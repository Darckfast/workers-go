package jshttp

import (
	"io"
	"net/http"
	"sync"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
	jsstream "github.com/syumai/workers/internal/stream"
)

type ResponseWriter struct {
	HeaderValue http.Header
	StatusCode  int
	Reader      io.ReadCloser
	Writer      *io.PipeWriter
	ReadyCh     chan struct{}
	Once        sync.Once
	RawJSBody   *js.Value
	Length      int64
}

var (
	_ http.ResponseWriter = (*ResponseWriter)(nil)
	_ http.Flusher        = (*ResponseWriter)(nil)
)

// Ready indicates that ResponseWriter is ready to be converted to Response.
func (w *ResponseWriter) Ready() {
	w.Once.Do(func() {
		close(w.ReadyCh)
	})
}

func (w *ResponseWriter) Write(data []byte) (n int, err error) {
	w.Ready()
	n, e := w.Writer.Write(data)
	w.Length += int64(n)
	return n, e
}

func (w *ResponseWriter) Header() http.Header {
	return w.HeaderValue
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

// Flush is a no-op implementation of http.Flusher.
//
// * PipeWriter does not have buffer, and JS-side Response does not have flush method.
// * But some libraries like `mcp-go` requires this method.
// * So implement this method as a workaround.
func (w *ResponseWriter) Flush() {
	// no-op
}

func (w *ResponseWriter) ToJSResponse() js.Value {
	respInit := jsclass.Object.New()
	respInit.Set("status", w.StatusCode)
	respInit.Set("statusText", http.StatusText(w.StatusCode))
	respInit.Set("headers", ToJSHeader(w.HeaderValue))
	readableStream := jsclass.Null

	if jsclass.MaybeFixedLengthStream.Truthy() && w.Length > 0 {
		readableStream = jsstream.ReadCloserToFixedLengthStream(w.Reader, w.Length)
	} else {
		readableStream = jsstream.ReadCloserToReadableStream(w.Reader)
	}

	return jsclass.Response.New(readableStream, respInit)
}
