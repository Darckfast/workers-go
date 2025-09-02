//go:build js && wasm

package r2

import (
	"errors"
	"io"
	"syscall/js"

	"github.com/Darckfast/workers-go/cloudflare/lifecycle"
	jsclass "github.com/Darckfast/workers-go/internal/class"
	jshttp "github.com/Darckfast/workers-go/internal/http"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
	"github.com/mailru/easyjson"
)

// Bucket represents interface of Cloudflare Worker's R2 Bucket instance.
//   - https://developers.cloudflare.com/r2/runtime-apis/#bucket-method-definitions
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1006
type Bucket struct {
	instance js.Value
}

// NewBucket returns Bucket for given variable name.
//   - variable name must be defined in wrangler.toml.
//   - see example: https://github.com/Darckfast/workers-go/tree/main/_examples/r2-image-viewer
//   - if the given variable name doesn't exist on runtime context, returns error.
//   - This function panics when a runtime context is not found.
func NewBucket(varName string) (*Bucket, error) {
	inst := lifecycle.Env.Get(varName)
	if inst.IsUndefined() {
		return nil, errors.New("%s is undefined" + varName)
	}
	return &Bucket{instance: inst}, nil
}

// Head returns the result of `head` call to Bucket.
//   - Body field of *Object is always nil for Head call.
//   - if the object for given key doesn't exist, returns nil.
//   - if a network error happens, returns error.
func (r *Bucket) Head(key string) (*R2Object, error) {
	p := r.instance.Call("head", key)
	v, err := jsclass.Await(p)
	if err != nil {
		return nil, err
	}
	if v.IsNull() {
		return nil, nil
	}
	return toObject(v)
}

// Get returns the result of `get` call to Bucket.
//   - if the object for given key doesn't exist, returns nil.
//   - if a network error happens, returns error.
func (r *Bucket) Get(key string, opts *GetOptions) (*R2Object, error) {
	b, _ := easyjson.Marshal(opts)
	optsJs := jsclass.JSON.Call("parse", string(b))

	p := r.instance.Call("get", key, optsJs)
	v, err := jsclass.Await(p)
	if err != nil {
		return nil, err
	}
	if v.IsNull() {
		return nil, nil
	}
	return toObject(v)
}

func (opts *PutOptions) toJS() js.Value {
	b, _ := easyjson.Marshal(opts)
	v := jsclass.JSON.Call("parse", string(b))
	v.Set("httpMetadata", jshttp.ToJSHeader(opts.HTTPMetadata))

	jsclass.Console.Log(v)
	return v
}

// Put returns the result of `put` call to Bucket.
//   - This method copies all bytes into memory for implementation restriction.
//   - Body field of *Object is always nil for Put call.
//   - if a network error happens, returns error.
func (r *Bucket) Put(key string, value io.ReadCloser, size int64, opts *PutOptions) (*R2Object, error) {
	readable := jsstream.ReadCloserToFixedLengthStream(value, size)
	p := r.instance.Call("put", key, readable, opts.toJS())
	v, err := jsclass.Await(p)
	if err != nil {
		return nil, err
	}
	return toObject(v)
}

// Delete returns the result of `delete` call to Bucket.
//   - if a network error happens, returns error.
func (r *Bucket) Delete(key string) error {
	p := r.instance.Call("delete", key)
	if _, err := jsclass.Await(p); err != nil {
		return err
	}
	return nil
}

// List returns the result of `list` call to Bucket.
//   - if a network error happens, returns error.
func (r *Bucket) List(opts *ListOptions) (*R2Objects, error) {
	b, _ := easyjson.Marshal(opts)
	optsJs := jsclass.JSON.Call("parse", string(b))

	p := r.instance.Call("list", optsJs)
	v, err := jsclass.Await(p)
	if err != nil {
		return nil, err
	}
	return toObjects(v)
}
