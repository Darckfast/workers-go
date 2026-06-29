//go:build js && wasm

/*
Package fetch is the glue code for Cloudflare's Worker fetch handler
*/
package fetch

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"codeberg.org/darckfast/workers-go/internal/jshttp"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/bind"
)

var httpHandler http.Handler = http.DefaultServeMux

func init() {
	var promise = js.FuncOf(func(this js.Value, args []js.Value) any {
		reqBody := args[0]
		reqMethod := args[1]
		reqUrl := args[2]
		reqHeaders := args[3]
		writable := args[4]

		signal := js.Value{}
		nodeSignal := js.Value{}

		if len(args) >= 6 {
			signal = args[5]
		}

		if len(args) >= 7 {
			nodeSignal = args[6]
		}

		bind.Env.LoadEnvs(js.Value{})
		bind.Ctx.Init(js.Value{})
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]

			go func(resolve js.Value) {
				var urlBytes = make([]byte, reqUrl.Length())
				js.CopyBytesToGo(urlBytes, reqUrl)

				var headersBytes = make([]byte, reqHeaders.Length())
				js.CopyBytesToGo(headersBytes, reqHeaders)

				var methodBytes = make([]byte, reqMethod.Length())
				js.CopyBytesToGo(methodBytes, reqMethod)

				reader := textproto.NewReader(bufio.NewReader(strings.NewReader(string(headersBytes))))
				headers, err := reader.ReadMIMEHeader()
				if err != nil && !errors.Is(err, io.EOF) {
					println("error decoding headers: ", err.Error())
				}

				url, _ := url.Parse(string(urlBytes))
				contentLength, _ := strconv.ParseInt(headers.Get("Content-Length"), 10, 64)
				req := &http.Request{
					Method:           string(methodBytes),
					URL:              url,
					Header:           http.Header(headers),
					Body:             jshttp.ToBody(reqBody),
					ContentLength:    contentLength,
					TransferEncoding: strings.Split(headers.Get("Transfer-Encoding"), ","),
					Host:             headers.Get("Host"),
				}

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				var cbCancel js.Func
				defer cbCancel.Release()
				cbCancel = js.FuncOf(func(this js.Value, args []js.Value) any {
					cancel()
					return nil
				})

				if signal.Truthy() {
					signal.Call("addEventListener", "abort", cbCancel)
				} else if nodeSignal.Truthy() {
					nodeSignal.Call("on", "close", cbCancel)
				}
				req = req.WithContext(ctx)

				w := &jshttp.ResponseWriter{
					HeaderValue: http.Header{},
					StatusCode:  http.StatusOK,
					ReadyCh:     make(chan struct{}),
					V:           writable,
				}

				go func(w *jshttp.ResponseWriter, req *http.Request) {
					defer func(w *jshttp.ResponseWriter) {
						select {
						case <-w.ReadyCh:
						default:
							w.ReadyCh <- struct{}{}
						}
					}(w)

					httpHandler.ServeHTTP(w, req)
				}(w, req)

				// Once the writable is about to to written, we can
				// return the http.headers and wait for the readable to
				// start pulling
				_, ok := <-w.ReadyCh
				buf := new(bytes.Buffer)
				buf.WriteString("HTTP/1.1 ")
				buf.WriteString(strconv.Itoa(w.StatusCode))
				buf.WriteString(" ")
				buf.WriteString(http.StatusText(w.StatusCode))
				buf.WriteString("\n")

				_ = w.HeaderValue.Write(buf)
				b := jsclass.Uint8Array.New(buf.Len())
				js.CopyBytesToJS(b, buf.Bytes())
				if ok {
					w.Close()
				}
				resolve.Invoke(b)
			}(resolve)

			return nil
		})

		return jsclass.Promise.New(cb)
	})

	jsclass.CF.Set("fetch", promise)
}

func ServeNonBlock(handler http.Handler) {
	if handler != nil {
		httpHandler = handler
	}
}
