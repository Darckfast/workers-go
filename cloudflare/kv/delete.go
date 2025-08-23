//go:build js && wasm

package kv

import (
	jsclass "github.com/Darckfast/workers-go/internal/class"
)

func (ns *Namespace) Delete(key string) error {
	p := ns.Call("delete", key)
	_, err := jsclass.Await(p)
	return err
}
