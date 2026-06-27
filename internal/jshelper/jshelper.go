//go:build js && wasm

package jshelper

import (
	"sync"
	"syscall/js"
)

// LazyJSVal Generic interface for lazy JS Value
type LazyJSVal struct {
	js.Value
	sync.Once
}

func (l *LazyJSVal) Init(id string) {
	if id == "" {
		panic("no id passed during Lazy js.Value initialition")
	}

	l.Do(func() {
		l.Value = js.Global().Get(id)
	})
}

// NewClass generic interface for JS classes
type NewClass struct {
	ID string
	V  LazyJSVal
}

func (r *NewClass) New(args ...any) js.Value {
	r.V.Init(r.ID)
	return r.V.New(args...)
}

func (r *NewClass) Class() js.Value {
	r.V.Init(r.ID)
	return r.V.Value
}

func (r *NewClass) Truthy() bool {
	r.V.Init(r.ID)
	return r.V.Truthy()
}

// PromiseClass Interface for JS Promise
type PromiseClass struct {
	v LazyJSVal
}

func (p *PromiseClass) New(args ...any) js.Value {
	p.v.Init("Promise")
	return p.v.New(args...)
}
func (p *PromiseClass) Reject(args ...any) js.Value {
	p.v.Init("Promise")
	return p.v.Call("reject", args...)
}

func (p *PromiseClass) Resolve(args ...any) js.Value {
	p.v.Init("Promise")
	return p.v.Call("resolve", args...)
}

// ArrayClass interface for JS Array class
type ArrayClass struct {
	v LazyJSVal
}

func (a *ArrayClass) New(args ...any) js.Value {
	a.v.Init("Array")
	return a.v.New(args...)
}
func (a *ArrayClass) From(args ...any) js.Value {
	a.v.Init("Array")
	return a.v.Call("from", args...)
}

// CacheInterface interface for JS Caches API
type CacheInterface struct {
	v LazyJSVal
}

func (c *CacheInterface) Open(ns string) js.Value {
	c.v.Init("caches")
	return c.v.Call("open", ns)
}
func (c *CacheInterface) Default() js.Value {
	c.v.Init("caches")
	return c.v.Get("default")
}

// AbortSignal interface for JS AbortSignal
type AbortSignal struct {
	v NewClass
}

func (a *AbortSignal) Timeout(t int64) js.Value {
	a.v.V.Init("AbortSignal")
	return a.v.V.Call("timeout", t)
}

func (a *AbortSignal) Any(args ...any) js.Value {
	a.v.V.Init("AbortSignal")
	return a.v.V.Call("any", args...)
}
