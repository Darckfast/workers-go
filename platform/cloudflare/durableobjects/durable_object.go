//go:build js && wasm

package durableobjects

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"iter"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"codeberg.org/darckfast/workers-go/internal/jshelper"
	"codeberg.org/darckfast/workers-go/internal/jshttp"
	"codeberg.org/darckfast/workers-go/internal/jstry"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/bind"
	"github.com/mailru/easyjson"
)

type DurableObjectSQLApi struct {
	v jshelper.LazyJSVal
}

func (d *DurableObjectSQLApi) Exec(query string, bindings ...any) *SqlStorageCursor {
	var cursor js.Value
	if len(bindings) == 0 {
		cursor = d.v.Call("exec", query)
	} else {
		cursor = d.v.Call("exec", query, bindings)
	}

	return &SqlStorageCursor{v: cursor}
}

type DurableObjectStorage struct {
	v   jshelper.LazyJSVal
	sql DurableObjectSQLApi
}

type SqlStorageCursor struct {
	v js.Value
}

// func (s *SqlStorageCursor) Raw() {}

func (s *SqlStorageCursor) One() ([]byte, error) {
	r, err := jstry.And.Catch(s.v, "one")
	if err != nil {
		return nil, err
	}

	str := jsclass.JSON.Stringify(r).String()
	return []byte(str), nil
}

func (s *SqlStorageCursor) ToArray() []byte {
	str := jsclass.JSON.Stringify(s.v.Call("toArray")).String()
	return []byte(str)
}

func (s *SqlStorageCursor) Next() iter.Seq[map[string]any] {
	return func(yield func(map[string]any) bool) {
		for {
			var r SqlStorageCursorProp
			str := jsclass.JSON.Stringify(s.v.Call("next")).String()
			_ = easyjson.Unmarshal([]byte(str), &r)

			if r.Done || !yield(r.Value) {
				return
			}
		}
	}
}

func (d *DurableObjectStorage) Sql() *DurableObjectSQLApi {
	d.sql.v.Do(func() {
		d.sql.v.Value = d.v.Get("sql")
	})

	return &d.sql
}

func (d *DurableObjectStorage) Get(name string) {
	d.v.Call("get", name)
}

type CtxJS struct {
	v       jshelper.LazyJSVal
	storage DurableObjectStorage
}

func (c *CtxJS) Storage() *DurableObjectStorage {
	c.storage.v.Do(func() {
		c.storage.v.Value = c.v.Get("storage")
	})

	return &c.storage
}

type DurableObject struct {
	__prototype__ js.Value
	class         js.Value
	ctx           CtxJS
	env           jshelper.LazyJSVal
}

func (d *DurableObject) Env() js.Value {
	d.env.Do(func() {
		d.env.Value = d.class.Get("ctx")
	})

	return d.env.Value
}

func (d *DurableObject) Ctx() *CtxJS {
	d.ctx.v.Do(func() {
		d.ctx.v.Value = d.class.Get("ctx")
	})

	return &d.ctx
}

func (d *DurableObject) New(classname string) {
	d.class = js.Global().Get("_durableobjects").Call("get", classname)
	d.__prototype__ = jsclass.Object.GetPrototypeOf(d.class)
}

type rpcfunc = func(ctx context.Context, i [][]byte) [][]byte

func (d *DurableObject) AddRPC(name string, h rpcfunc) {
	var jsfun = js.FuncOf(func(this js.Value, args []js.Value) any {
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

			go func(resolve js.Value) {
				ctx := context.Background()
				out := h(ctx, bufs) // fun here

				dstArr := jsclass.Array.New(len(out))
				for i, o := range out {
					dst := jsclass.Uint8Array.New(len(o))
					js.CopyBytesToJS(dst, o)
					dstArr.SetIndex(i, dst)
				}

				resolve.Invoke(dstArr)
			}(resolve)

			return nil
		})

		return jsclass.Promise.New(cb)
	})

	d.__prototype__.Set(name, jsfun)
}

type FetchHandler = func(w http.ResponseWriter, r *http.Request)

func (d *DurableObject) AddFetch(fetch FetchHandler) {
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

		bind.Env.LoadEnvs(d.Env())
		bind.Ctx.Init(d.ctx.v.Value)
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
						w.Ready()
					}(w)

					fetch(w, req)
				}(w, req)

				// Once the writable is about to to written, we can
				// return the http.headers and wait for the readable to
				// start pulling
				<-w.ReadyCh

				buf := new(bytes.Buffer)
				buf.WriteString("HTTP/1.1 ")
				buf.WriteString(strconv.Itoa(w.StatusCode))
				buf.WriteString(" ")
				buf.WriteString(http.StatusText(w.StatusCode))
				buf.WriteString("\n")

				_ = w.HeaderValue.Write(buf)
				b := jsclass.Uint8Array.New(buf.Len())
				js.CopyBytesToJS(b, buf.Bytes())

				resolve.Invoke(b)
			}(resolve)

			return nil
		})

		return jsclass.Promise.New(cb)
	})

	d.__prototype__.Set("_fetch", promise)
}
