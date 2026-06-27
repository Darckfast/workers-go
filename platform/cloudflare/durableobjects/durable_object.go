//go:build js && wasm

package durableobjects

import (
	"context"
	"io"
	"iter"
	"net/http"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"codeberg.org/darckfast/workers-go/internal/jshelper"
	"codeberg.org/darckfast/workers-go/internal/jshttp"
	"codeberg.org/darckfast/workers-go/internal/jsruntime"
	"codeberg.org/darckfast/workers-go/internal/jstry"
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
	var handleRequestPromise = js.FuncOf(func(this js.Value, args []js.Value) any {
		reqObj := args[0]
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]

			go func(resolve js.Value, reject js.Value) {
				jsclass.Env.LoadEnvs(js.Value{})

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

				go func(w *jshttp.ResponseWriter, r *http.Request) {
					defer func() {
						w.Ready()
						err := writer.Close()

						if err != nil {
							println("error closing response body writer", err.Error())
						}
					}()

					fetch(w, req)
				}(w, req)

				<-w.ReadyCh
				resolve.Invoke(w.ToJSResponse())
			}(resolve, reject)

			return nil
		})

		return jsclass.Promise.New(cb)
	})

	d.__prototype__.Set("fetch", handleRequestPromise)
}
