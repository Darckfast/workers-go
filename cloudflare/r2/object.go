//go:build js && wasm

package r2

import (
	"errors"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
	"github.com/mailru/easyjson"
)

// TODO: implement awsfetch for url signed => due was limitations with syscall, using the aws-sdk-go is not viable
// afaik during init() it tries to access the home dir to load the local config, and this makes the process exit
type Object struct {
	instance js.Value
	R2Object
}

func (o *Object) BodyUsed() (bool, error) {
	v := o.instance.Get("bodyUsed")
	if v.IsUndefined() {
		return false, errors.New("bodyUsed doesn't exist for this Object")
	}
	return v.Bool(), nil
}

func toObject(v js.Value) (*R2Object, error) {
	str := jsclass.JSON.Stringify(v)
	var o R2Object
	err := easyjson.Unmarshal([]byte(str.String()), &o)

	if err != nil {
		return nil, err
	}

	bodyVal := v.Get("body")
	if bodyVal.Truthy() {
		o.Body = jsstream.ReadableStreamToReadCloser(bodyVal)
	}

	return &o, nil
}
