//go:build js && wasm

package kv

import (
	"errors"
	"syscall/js"

	"github.com/Darckfast/workers-go/cloudflare/lifecycle"
)

type Namespace struct {
	js.Value
}

func NewNamespace(binding string) (*Namespace, error) {
	inst := lifecycle.Env.Get(binding)
	if !inst.Truthy() {
		return nil, errors.New(binding + " is undefined")
	}
	return &Namespace{inst}, nil
}
