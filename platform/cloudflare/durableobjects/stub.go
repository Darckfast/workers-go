//go:build js && wasm

/*
Package durableobjects is the glue code for Cloudflare's DurableObjects bindings
*/
package durableobjects

import (
	"errors"
	"net/http"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"codeberg.org/darckfast/workers-go/internal/jshttp"
	"codeberg.org/darckfast/workers-go/internal/jstry"
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

//nolint:staticcheck // drops error warning for ST1003
func (ns *DurableObjectNamespace) IdFromName(name string) *DurableObjectID {
	id := ns.instance.Call("idFromName", name)
	return &DurableObjectID{val: id}
}

//nolint:staticcheck // drops error warning for ST1003
func (ns *DurableObjectNamespace) IdFromString(id string) (*DurableObjectID, error) {
	idStr, err := jstry.And.Catch(ns.instance, "idFromString", id)
	if err != nil {
		return nil, err
	}

	return &DurableObjectID{val: idStr}, nil
}

//nolint:staticcheck // drops error warning for ST1003
func (ns *DurableObjectNamespace) NewUniqueId() *DurableObjectID {
	id := ns.instance.Call("newUniqueId")
	return &DurableObjectID{val: id}
}

func (ns *DurableObjectNamespace) Jurisdiction(jur string) *DurableObjectNamespace {
	inst := ns.instance.Call("jurisdiction")
	return &DurableObjectNamespace{instance: inst}
}

func (ns *DurableObjectNamespace) GetByName(id string) *DurableObjectStub {
	stub := ns.instance.Call("getByName", id)
	return &DurableObjectStub{val: stub}
}

func (ns *DurableObjectNamespace) Get(id *DurableObjectID) (*DurableObjectStub, error) {
	if id == nil || id.val.IsUndefined() {
		return nil, errors.New("invalid UniqueGlobalId")
	}
	stub := ns.instance.Call("get", id.val)
	return &DurableObjectStub{val: stub}, nil
}

type DurableObjectID struct {
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

func (s *DurableObjectStub) Call(funcName string, args ...any) (any, error) {
	promise := s.val.Call(funcName, args...)
	//TODO: stringify the return
	jsRes, err := jsclass.Await(promise)
	if err != nil {
		return nil, err
	}

	return jsRes.String(), nil
}
