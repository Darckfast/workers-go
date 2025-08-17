package cache

import (
	"errors"
	"net/http"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
	jshttp "github.com/syumai/workers/internal/http"
)

type Cache struct {
	instance js.Value
}

func (c *Cache) Open(namespace string) error {
	v, err := jsclass.Await(jsclass.Caches.Call("open", namespace))
	if err != nil {
		return err
	}

	c.instance = v
	return nil
}

func New() *Cache {
	return &Cache{
		instance: jsclass.Caches.Get("default"),
	}
}

// Put attempts to add a response to the cache, using the given request as the key.
// Returns an error for the following conditions
// - the request passed is a method other than GET.
// - the response passed has a status of 206 Partial Content.
// - Cache-Control instructs not to cache or if the response is too large.
// docs: https://developers.cloudflare.com/workers/runtime-apis/cache/#put
func (c *Cache) Put(req *http.Request, res *http.Response) error {
	_, err := jsclass.Await(c.instance.Call("put", jshttp.ToJSRequest(req), jshttp.ToJSResponse(res)))
	if err != nil {
		return err
	}
	return nil
}

var ErrCacheNotFound = errors.New("cache not found")

type MatchOptions struct {
	IgnoreMethod bool
}

func (opts *MatchOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}
	obj := jsclass.Object.New()
	obj.Set("ignoreMethod", opts.IgnoreMethod)
	return obj
}

// Match returns the response object keyed to that request.
// docs: https://developers.cloudflare.com/workers/runtime-apis/cache/#match
func (c *Cache) Match(req *http.Request, opts *MatchOptions) (*http.Response, error) {
	res, err := jsclass.Await(c.instance.Call("match", jshttp.ToJSRequest(req), opts.toJS()))
	if err != nil {
		return nil, err
	}
	if res.IsUndefined() {
		return nil, ErrCacheNotFound
	}
	return jshttp.ToResponse(res), nil
}

type DeleteOptions struct {
	IgnoreMethod bool
}

func (opts *DeleteOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}
	obj := jsclass.Object.New()
	obj.Set("ignoreMethod", opts.IgnoreMethod)
	return obj
}

// Delete removes the Response object from the cache.
// This method only purges content of the cache in the data center that the Worker was invoked.
// Returns ErrCacheNotFount if the response was not cached.
func (c *Cache) Delete(req *http.Request, opts *DeleteOptions) error {
	res, err := jsclass.Await(c.instance.Call("delete", jshttp.ToJSRequest(req), opts.toJS()))
	if err != nil {
		return err
	}
	if !res.Bool() {
		return ErrCacheNotFound
	}
	return nil
}
