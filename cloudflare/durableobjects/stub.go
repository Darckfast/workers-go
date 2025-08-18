//go:build js && wasm

package durableobjects

import (
	"errors"
	"net/http"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jshttp "github.com/Darckfast/workers-go/internal/http"
	jstry "github.com/Darckfast/workers-go/internal/try"
)

type DurableObjectNamespace struct {
	instance js.Value
}

func NewDurableObjectNamespace(varName string) (*DurableObjectNamespace, error) {
	inst := jsclass.Env.Get(varName)
	if inst.IsUndefined() {
		return nil, errors.New("%s is undefined" + varName)
	}
	return &DurableObjectNamespace{instance: inst}, nil
}

func (ns *DurableObjectNamespace) IdFromName(name string) *DurableObjectId {
	id := ns.instance.Call("idFromName", name)
	return &DurableObjectId{val: id}
}

func (ns *DurableObjectNamespace) IdFromString(id string) (*DurableObjectId, error) {
	idStr, err := jstry.TryCatch(js.FuncOf(func(_ js.Value, args []js.Value) any {
		return ns.instance.Call("idFromString", id)
	}))

	if err != nil {
		return nil, err
	}

	return &DurableObjectId{val: idStr}, nil
}

func (ns *DurableObjectNamespace) NewUniqueId() *DurableObjectId {
	id := ns.instance.Call("newUniqueId")
	return &DurableObjectId{val: id}
}

func (ns *DurableObjectNamespace) Jurisdiction(jur string) *DurableObjectNamespace {
	inst := ns.instance.Call("jurisdiction")
	return &DurableObjectNamespace{instance: inst}
}

func (ns *DurableObjectNamespace) Get(id *DurableObjectId) (*DurableObjectStub, error) {
	if id == nil || id.val.IsUndefined() {
		return nil, errors.New("invalid UniqueGlobalId")
	}
	stub := ns.instance.Call("get", id.val)
	return &DurableObjectStub{val: stub}, nil
}

type DurableObjectId struct {
	val js.Value
}

type DurableObjectStub struct {
	val js.Value
}

func (s *DurableObjectStub) Fetch(req *http.Request) (*http.Response, error) {
	jsReq := jshttp.ToJSRequest(req)

	promise := s.val.Call("fetch", jsReq)
	jsRes, err := jsclass.Await(promise)
	if err != nil {
		return nil, err
	}

	return jshttp.ToResponse(jsRes), nil
}

func (s *DurableObjectStub) Call(funcName string) (any, error) {
	promise := s.val.Call(funcName)
	jsRes, err := jsclass.Await(promise)
	if err != nil {
		return nil, err
	}

	return jsRes.String(), nil
}
