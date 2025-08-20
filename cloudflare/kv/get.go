//go:build js && wasm

package kv

import (
	"errors"
	"io"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
)

type GetOptions struct {
	CacheTTL int
}

func (opts *GetOptions) toJS(type_ string) js.Value {
	obj := jsclass.Object.New()
	obj.Set("type", type_)
	if opts == nil {
		return obj
	}
	if opts.CacheTTL != 0 {
		obj.Set("cacheTtl", opts.CacheTTL)
	}
	return obj
}

func (ns *Namespace) GetStringWithMetadata(key string, opts *GetOptions) (string, string, error) {
	p := ns.instance.Call("get", key, opts.toJS("text"))
	r, err := jsclass.Await(p)

	if err != nil {
		return "", "", err
	}

	if r.IsNull() || r.IsUndefined() {
		return "", "", errors.New("key has no value")
	}

	vm := r.Get("metadata")
	v := r.Get("value")

	if vm.IsNull() || vm.IsUndefined() {
		return v.String(), "", nil
	}

	return v.String(), vm.String(), nil
}

func (ns *Namespace) GetStrings(keys []string, opts *GetOptions) (map[string]string, error) {
	p := ns.instance.Call("get", keys, opts.toJS("text"))
	v, err := jsclass.Await(p)

	if err != nil {
		return nil, err
	}

	if v.IsNull() || v.IsUndefined() {
		return nil, errors.New("key has no value")
	}

	return jsconv.JSValueToMapString(v), nil
}

func (ns *Namespace) GetString(key string, opts *GetOptions) (string, error) {
	p := ns.instance.Call("get", key, opts.toJS("text"))
	v, err := jsclass.Await(p)

	if err != nil {
		return "", err
	}

	if v.IsNull() || v.IsUndefined() {
		return "", errors.New("key has no value")
	}

	return v.String(), nil
}

func (ns *Namespace) GetReader(key string, opts *GetOptions) (io.ReadCloser, error) {
	p := ns.instance.Call("get", key, opts.toJS("stream"))
	v, err := jsclass.Await(p)
	if err != nil {
		return nil, err
	}

	if v.IsNull() || v.IsUndefined() {
		return nil, errors.New("key has no value")
	}

	return jsstream.ReadableStreamToReadCloser(v), nil
}
