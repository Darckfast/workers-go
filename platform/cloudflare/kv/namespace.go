//go:build js && wasm

package kv

import (
	"errors"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
)

type Namespace struct {
	js.Value
}

func NewNamespace(binding string) (*Namespace, error) {
	inst := jsclass.Env.Get(binding)
	if !inst.Truthy() {
		return nil, errors.New(binding + " is undefined")
	}
	return &Namespace{inst}, nil
}
