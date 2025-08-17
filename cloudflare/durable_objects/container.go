package durableobjects

import (
	"net/http"

	jsclass "github.com/syumai/workers/internal/class"
	jshttp "github.com/syumai/workers/internal/http"
	jsutil "github.com/syumai/workers/internal/utils"
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
	inst := jsutil.RuntimeEnv.Get(binding)
	donamespace := &DurableObjectNamespace{instance: inst}
	objId := donamespace.IdFromName(id)
	obj, err := donamespace.Get(objId)

	if err != nil {
		return nil, err
	}

	return &Container{obj}, nil
}
