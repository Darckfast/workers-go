//go:build js && wasm

/*
Package kv is the glue code for Cloudflare's KV Worker
*/
package kv

import (
	jsclass "codeberg.org/darckfast/workers-go/internal/class"
)

func (ns *Namespace) Delete(key string) error {
	p := ns.Call("delete", key)
	_, err := jsclass.Await(p)
	return err
}
