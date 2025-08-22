//go:build js && wasm

package kv

import (
	"encoding/json"
	"io"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
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
	p := ns.instance.Call("put", key, value, opts.ToJS())
	_, err := jsclass.Await(p)

	return err
}

// PutReader puts stream value into KV with key.
//   - This method copies all bytes into memory for implementation restriction.
//   - if a network error happens, returns error.
func (ns *Namespace) PutReader(key string, value io.Reader, opts *PutOptions) error {
	// fetch body cannot be ReadableStream. see: https://github.com/whatwg/fetch/issues/1438
	b, err := io.ReadAll(value)
	if err != nil {
		return err
	}
	ua := jsclass.Uint8Array.New(len(b))
	js.CopyBytesToJS(ua, b)
	p := ns.instance.Call("put", key, ua.Get("buffer"), opts.ToJS())
	_, err = jsclass.Await(p)
	if err != nil {
		return err
	}
	return nil
}
