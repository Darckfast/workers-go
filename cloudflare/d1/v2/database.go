//go:build js && wasm

package d1

import (
	"errors"
	"syscall/js"

	"github.com/Darckfast/workers-go/cloudflare/lifecycle"
	jsclass "github.com/Darckfast/workers-go/internal/class"
	"github.com/mailru/easyjson"
)

type D1Db struct {
	v js.Value
}

func GetDB(binding string) (*D1Db, error) {
	v := lifecycle.Env.Get(binding)

	if !v.Truthy() {
		return nil, errors.New("d1 binding not found " + binding)
	}

	return &D1Db{v: v}, nil
}

type D1PreparedStatment struct {
	v js.Value
}

func (s *D1PreparedStatment) Bind(variable ...any) *D1PreparedStatment {
	s.v = s.v.Call("bind", variable...)

	return s
}

func (s *D1PreparedStatment) Run() (*D1Result, error) {
	r, err := jsclass.Await(s.v.Call("run"))

	if err != nil {
		return nil, err
	}

	var result D1Result
	str := jsclass.JSON.Stringify(r)
	err = easyjson.Unmarshal([]byte(str.String()), &result)

	return &result, err
}

func (s *D1PreparedStatment) Raw(columnNames bool) (D1RawResults, error) {
	arg := jsclass.Object.New()
	arg.Set("columnNames", columnNames)
	r, err := jsclass.Await(s.v.Call("raw", arg))

	if err != nil {
		return nil, err
	}

	var result D1RawResults
	str := jsclass.JSON.Stringify(r)
	err = easyjson.Unmarshal([]byte(str.String()), &result)

	return result, err
}

func (s *D1PreparedStatment) First(columnName string) (*D1FirstResult, error) {
	r, err := jsclass.Await(s.v.Call("first", columnName))

	if err != nil {
		return nil, err
	}

	var result D1FirstResult
	str := jsclass.JSON.Stringify(r)
	err = easyjson.Unmarshal([]byte(str.String()), &result)

	return &result, err
}

func (d *D1Db) Prepare(query string) *D1PreparedStatment {
	stmtObj := d.v.Call("prepare", query)
	return &D1PreparedStatment{stmtObj}
}

func (d *D1Db) Batch(stmts []D1PreparedStatment) (D1BatchResults, error) {
	jsList := jsclass.Array.New()

	for _, st := range stmts {
		jsList.Call("push", st.v)
	}

	batchResult, err := jsclass.Await(d.v.Call("batch", jsList))

	if err != nil {
		return nil, err
	}

	var results D1BatchResults
	str := jsclass.JSON.Stringify(batchResult)
	err = easyjson.Unmarshal([]byte(str.String()), &results)

	return results, err
}

func (d *D1Db) Exec(query string) (*D1ExecResult, error) {
	result, err := jsclass.Await(d.v.Call("exec", query))

	if err != nil {
		return nil, err
	}

	var d1Result D1ExecResult
	str := jsclass.JSON.Stringify(result)
	err = easyjson.Unmarshal([]byte(str.String()), &d1Result)

	return &d1Result, err
}

type D1DatabaseSession struct {
	v js.Value
}

func (d *D1DatabaseSession) Prepare(query string) *D1PreparedStatment {
	stmtObj := d.v.Call("prepare", query)
	return &D1PreparedStatment{v: stmtObj}
}

func (d *D1DatabaseSession) Batch(stmts ...D1PreparedStatment) ([]D1Result, error) {
	jsList := jsclass.Array.New()

	for _, st := range stmts {
		jsList.Call("push", st.v)
	}

	batchResult, err := jsclass.Await(d.v.Call("batch", jsList))

	if err != nil {
		return nil, err
	}

	var results D1BatchResults
	str := jsclass.JSON.Stringify(batchResult)
	err = easyjson.Unmarshal([]byte(str.String()), &results)

	return results, err
}

func (d *D1Db) WithSession(param string) *D1DatabaseSession {
	s := d.v.Call("withSession", param)
	return &D1DatabaseSession{s}
}

func (d *D1DatabaseSession) GetBookmark() string {
	str := d.v.Call("getBookmark")

	if str.IsNull() {
		return ""
	}

	return str.String()
}
