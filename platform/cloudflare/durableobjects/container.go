//go:build js && wasm

package durableobjects

import (
	"net/http"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"codeberg.org/darckfast/workers-go/internal/jshttp"
)

type Container struct {
	*DurableObjectStub
}

func (s *Container) ContainerFetch(req *http.Request) (*http.Response, error) {
	jsReq := jshttp.ToJSRequest(req)

	promise := s.val.Call("containerFetch", jsReq)
	jsRes, err := jsclass.Await(promise)
	if err != nil {
		return nil, err
	}

	return jshttp.ToResponse(jsRes), nil
}

func GetContainer(binding string, id string) (*Container, error) {
	inst := jsclass.Env.Get(binding)
	donamespace := &DurableObjectNamespace{instance: inst}
	objID := donamespace.IdFromName(id)
	obj, err := donamespace.Get(objID)

	if err != nil {
		return nil, err
	}

	return &Container{obj}, nil
}
