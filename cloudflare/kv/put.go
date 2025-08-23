//go:build js && wasm

package kv

import (
	"encoding/json"
	"io"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
)

type PutOptions struct {
	Expiration    int            `json:"expiration,omitempty"`
	ExpirationTTL int            `json:"expirationTtl,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

func (o *PutOptions) ToJS() js.Value {
	b, _ := json.Marshal(o)
	js, _ := jsclass.JSON.Parse(string(b))

	return js
}

func (ns *Namespace) Put(key string, value string, opts *PutOptions) error {
	p := ns.Value.Call("put", key, value, opts.ToJS())
	_, err := jsclass.Await(p)

	return err
}

func (ns *Namespace) PutReader(key string, value io.ReadCloser, opts *PutOptions) error {
	stream := jsstream.ReadCloserToReadableStream(value)
	p := ns.Value.Call("put", key, stream, opts.ToJS())
	_, err := jsclass.Await(p)

	return err
}
