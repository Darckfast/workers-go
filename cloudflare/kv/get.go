//go:build js && wasm

package kv

import (
	"encoding/json"
	"errors"
	"io"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
)

type GetOptions struct {
	Type     string `json:"type"`
	CacheTTL int    `json:"cacheTtl,omitempty"`
}

func (o *GetOptions) ToJS() js.Value {
	b, _ := json.Marshal(o)
	j, _ := jsclass.JSON.Parse(string(b))

	return j
}

type StringWithMetadata struct {
	Value    string         `json:"value"`
	Metadata map[string]any `json:"metadata"`
}

func (ns *Namespace) GetWithMetadata(key string, cacheTtl int) (*StringWithMetadata, error) {
	opts := GetOptions{CacheTTL: cacheTtl, Type: "text"}
	p := ns.instance.Call("getWithMetadata", key, opts.ToJS())
	r, err := jsclass.Await(p)

	if err != nil {
		return nil, err
	}

	if r.IsNull() || r.IsUndefined() {
		return nil, errors.New("key has no value")
	}

	s := jsclass.JSON.Stringify(r)

	var sm StringWithMetadata

	err = json.Unmarshal([]byte(s.String()), &sm)
	return &sm, err
}

func (ns *Namespace) Get(keysRaw []string, cacheTtl int) (map[string]any, error) {
	keys := make([]any, len(keysRaw))
	for i, v := range keysRaw {
		keys[i] = v
	}

	opts := GetOptions{CacheTTL: cacheTtl, Type: "text"}
	p := ns.instance.Call("get", keys, opts.ToJS())
	v, err := jsclass.Await(p)

	if err != nil {
		return nil, err
	}

	jsclass.Console.Log(v)
	return jsconv.JSMapToMap(v)
}

func (ns *Namespace) GetAsReader(key string, cacheTtl int) (io.ReadCloser, error) {
	opts := GetOptions{CacheTTL: cacheTtl, Type: "stream"}
	p := ns.instance.Call("get", key, opts.ToJS())
	v, err := jsclass.Await(p)
	if err != nil {
		return nil, err
	}

	if v.IsNull() || v.IsUndefined() {
		return nil, errors.New("key has no value")
	}

	return jsstream.ReadableStreamToReadCloser(v), nil
}
