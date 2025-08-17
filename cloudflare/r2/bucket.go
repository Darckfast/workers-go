package r2

import (
	"errors"
	"io"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
	jsstream "github.com/syumai/workers/internal/stream"
)

// Bucket represents interface of Cloudflare Worker's R2 Bucket instance.
//   - https://developers.cloudflare.com/r2/runtime-apis/#bucket-method-definitions
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1006
type Bucket struct {
	instance js.Value
}

// NewBucket returns Bucket for given variable name.
//   - variable name must be defined in wrangler.toml.
//   - see example: https://github.com/syumai/workers/tree/main/_examples/r2-image-viewer
//   - if the given variable name doesn't exist on runtime context, returns error.
//   - This function panics when a runtime context is not found.
func NewBucket(varName string) (*Bucket, error) {
	inst := jsclass.Env.Get(varName)
	if inst.IsUndefined() {
		return nil, errors.New("%s is undefined" + varName)
	}
	return &Bucket{instance: inst}, nil
}

// Head returns the result of `head` call to Bucket.
//   - Body field of *Object is always nil for Head call.
//   - if the object for given key doesn't exist, returns nil.
//   - if a network error happens, returns error.
func (r *Bucket) Head(key string) (*Object, error) {
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
func (r *Bucket) Get(key string) (*Object, error) {
	p := r.instance.Call("get", key)
	v, err := jsclass.Await(p)
	if err != nil {
		return nil, err
	}
	if v.IsNull() {
		return nil, nil
	}
	return toObject(v)
}

// PutOptions represents Cloudflare R2 put options.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1128
type PutOptions struct {
	HTTPMetadata   HTTPMetadata
	CustomMetadata map[string]string
	MD5            string
}

func (opts *PutOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}
	obj := jsclass.Object.New()
	if opts.HTTPMetadata != (HTTPMetadata{}) {
		obj.Set("httpMetadata", opts.HTTPMetadata.toJS())
	}
	if opts.CustomMetadata != nil {
		// convert map[string]string to map[string]any.
		// This makes the map convertible to JS.
		// see: https://pkg.go.dev/syscall/js#ValueOf
		customMeta := make(map[string]any, len(opts.CustomMetadata))
		for k, v := range opts.CustomMetadata {
			customMeta[k] = v
		}
		obj.Set("customMetadata", customMeta)
	}
	if opts.MD5 != "" {
		obj.Set("md5", opts.MD5)
	}
	return obj
}

// Put returns the result of `put` call to Bucket.
//   - This method copies all bytes into memory for implementation restriction.
//   - Body field of *Object is always nil for Put call.
//   - if a network error happens, returns error.
func (r *Bucket) Put(key string, value io.ReadCloser, size int64, opts *PutOptions) (*Object, error) {
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
func (r *Bucket) List() (*Objects, error) {
	p := r.instance.Call("list")
	v, err := jsclass.Await(p)
	if err != nil {
		return nil, err
	}
	return toObjects(v)
}
