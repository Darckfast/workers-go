//go:build js && wasm

package kv

import (
	"errors"
	"io"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
	"github.com/mailru/easyjson"
)

func (o *GetOptions) ToJS() js.Value {
	b, _ := easyjson.Marshal(o)
	j, _ := jsclass.JSON.Parse(string(b))

	return j
}

func (ns *Namespace) GetWithMetadata(key string, cacheTtl int) (*StringWithMetadata, error) {
	opts := GetOptions{CacheTTL: cacheTtl, Type: "text"}
	p := ns.Call("getWithMetadata", key, opts.ToJS())
	r, err := jsclass.Await(p)

	if err != nil {
		return nil, err
	}

	if r.IsNull() || r.IsUndefined() {
		return nil, errors.New("key has no value")
	}

	s := jsclass.JSON.Stringify(r)

	var sm StringWithMetadata

	err = easyjson.Unmarshal([]byte(s.String()), &sm)
	return &sm, err
}

func (ns *Namespace) Get(keysRaw []string, cacheTtl int) (map[string]any, error) {
	keys := make([]any, len(keysRaw))
	for i, v := range keysRaw {
		keys[i] = v
	}

	opts := GetOptions{CacheTTL: cacheTtl, Type: "text"}
	p := ns.Call("get", keys, opts.ToJS())
	v, err := jsclass.Await(p)

	if err != nil {
		return nil, err
	}

	return jsconv.JSMapToMap(v)
}

func (ns *Namespace) GetAsReader(key string, cacheTtl int) (io.ReadCloser, error) {
	opts := GetOptions{CacheTTL: cacheTtl, Type: "stream"}
	p := ns.Call("get", key, opts.ToJS())
	v, err := jsclass.Await(p)
	if err != nil {
		return nil, err
	}

	if v.IsNull() || v.IsUndefined() {
		return nil, errors.New("key has no value")
	}

	return jsstream.ReadableStreamToReadCloser(v), nil
}
