//go:build js && wasm

package kv

import (
	"io"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
	"github.com/mailru/easyjson"
)

func (o *PutOptions) ToJS() js.Value {
	b, _ := easyjson.Marshal(o)
	js, _ := jsclass.JSON.Parse(string(b))

	return js
}

func (ns *Namespace) Put(key string, value string, opts *PutOptions) error {
	p := ns.Call("put", key, value, opts.ToJS())
	_, err := jsclass.Await(p)

	return err
}

func (ns *Namespace) PutReader(key string, value io.ReadCloser, opts *PutOptions) error {
	stream := jsstream.ReadCloserToReadableStream(value)
	p := ns.Call("put", key, stream, opts.ToJS())
	_, err := jsclass.Await(p)

	return err
}
