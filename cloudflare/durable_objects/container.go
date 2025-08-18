//go:build js && wasm

package durableobjects

import (
	"net/http"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jshttp "github.com/Darckfast/workers-go/internal/http"
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
	objId := donamespace.IdFromName(id)
	obj, err := donamespace.Get(objId)

	if err != nil {
		return nil, err
	}

	return &Container{obj}, nil
}
