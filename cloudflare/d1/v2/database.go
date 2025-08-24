//go:build js && wasm

package d1

import (
	"encoding/json"
	"errors"
	"syscall/js"

	"github.com/Darckfast/workers-go/cloudflare/lifecycle"
	jsclass "github.com/Darckfast/workers-go/internal/class"
)

type D1Db struct {
	js.Value
}

func GetDB(binding string) (*D1Db, error) {
	v := lifecycle.Env.Get(binding)

	if !v.Truthy() {
		return nil, errors.New("d1 binding not found " + binding)
	}

	return &D1Db{v}, nil
}

type D1PreparedStatment struct {
	js.Value
}

func (s *D1PreparedStatment) Bind(variable ...any) *D1PreparedStatment {
	s.Value = s.Call("bind", variable...)

	return s
}

func (s *D1PreparedStatment) Run() (*D1Result, error) {
	r, err := jsclass.Await(s.Call("run"))

	if err != nil {
		return nil, err
	}

	var result D1Result
	str := jsclass.JSON.Stringify(r)
	err = json.Unmarshal([]byte(str.String()), &result)

	return &result, err
}

func (s *D1PreparedStatment) Raw(columnNames bool) ([]any, error) {
	arg := jsclass.Object.New()
	arg.Set("columnNames", columnNames)
	r, err := jsclass.Await(s.Call("raw", arg))

	if err != nil {
		return nil, err
	}

	var result []any
	str := jsclass.JSON.Stringify(r)
	err = json.Unmarshal([]byte(str.String()), &result)

	return result, err
}

func (s *D1PreparedStatment) First(columnName string) (*any, error) {
	r, err := jsclass.Await(s.Call("first", columnName))

	if err != nil {
		return nil, err
	}

	var result any
	str := jsclass.JSON.Stringify(r)
	err = json.Unmarshal([]byte(str.String()), &result)

	return &result, err
}

func (d *D1Db) Prepare(query string) *D1PreparedStatment {
	stmtObj := d.Call("prepare", query)
	return &D1PreparedStatment{stmtObj}
}

type D1ExecResult struct {
	Count    int `json:"count"`
	Duration int `json:"duration"`
}

type D1Result struct {
	Success bool  `json:"success"`
	Results []any `json:"results"`
	Meta    struct {
		ServedBy        string `json:"served_by"`
		ServedByRegion  string `json:"served_by_region"`
		ServedByPrimary bool   `json:"served_by_primary"`
		Timings         struct {
			SqlDurationMs int64 `json:"sql_duration_ms"`
		} `json:"timings"`
		Duration    int64 `json:"duration"`
		Changes     int64 `json:"changes"`
		LastRowId   int64 `json:"last_row_id"`
		ChangedDb   bool  `json:"changed_db"`
		SizeAfter   int64 `json:"size_after"`
		RowsRead    int64 `json:"rows_read"`
		RowsWritten int64 `json:"rows_written"`
	} `json:"meta"`
}

func (d *D1Db) Batch(stmts []D1PreparedStatment) ([]D1Result, error) {
	jsList := jsclass.Array.New()

	for _, st := range stmts {
		jsList.Call("push", st.Value)
	}

	batchResult, err := jsclass.Await(d.Call("batch", jsList))

	if err != nil {
		return nil, err
	}

	var results []D1Result
	str := jsclass.JSON.Stringify(batchResult)
	err = json.Unmarshal([]byte(str.String()), &results)

	return results, err
}

func (d *D1Db) Exec(query string) (*D1ExecResult, error) {
	result, err := jsclass.Await(d.Call("exec", query))

	if err != nil {
		return nil, err
	}

	var d1Result D1ExecResult
	str := jsclass.JSON.Stringify(result)
	err = json.Unmarshal([]byte(str.String()), &d1Result)

	return &d1Result, err
}

type D1DatabaseSession struct {
	js.Value
}

func (d *D1DatabaseSession) Prepare(query string) *D1PreparedStatment {
	stmtObj := d.Call("prepare", query)
	return &D1PreparedStatment{stmtObj}
}

func (d *D1DatabaseSession) Batch(stmts ...D1PreparedStatment) ([]D1Result, error) {
	jsList := jsclass.Array.New()

	for _, st := range stmts {
		jsList.Call("push", st.Value)
	}

	batchResult, err := jsclass.Await(d.Call("batch", jsList))

	if err != nil {
		return nil, err
	}

	var results []D1Result
	str := jsclass.JSON.Stringify(batchResult)
	err = json.Unmarshal([]byte(str.String()), &results)

	return results, err
}

func (d *D1Db) WithSession(param string) *D1DatabaseSession {
	s := d.Call("withSession", param)
	return &D1DatabaseSession{s}
}

func (d *D1DatabaseSession) GetBookmark() string {
	str := d.Call("getBookmark")

	if str.IsNull() {
		return ""
	}

	return str.String()
}
