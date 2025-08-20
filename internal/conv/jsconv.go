//go:build js && wasm

package jsconv

import (
	"encoding/json"
	"strconv"
	"syscall/js"
	"time"

	jsclass "github.com/Darckfast/workers-go/internal/class"
)

func ArrayFrom(v js.Value) js.Value {
	return jsclass.Array.Call("from", v)
}

// TODO: redo this for map[string]any
func StrRecordToMap(v js.Value) map[string]string {
	if v.IsUndefined() || v.IsNull() {
		return map[string]string{}
	}
	entries := jsclass.Object.Call("entries", v)
	entriesLen := entries.Get("length").Int()
	result := make(map[string]string, entriesLen)
	for i := range entriesLen {
		entry := entries.Index(i)
		key := entry.Index(0).String()
		value := entry.Index(1).String()
		result[key] = value
	}
	return result
}

func MapToJSValue(v map[string]any) js.Value {
	b, _ := json.Marshal(v)
	return jsclass.JSON.Call("parse", string(b))
}

func JSValueToMapString(v js.Value) map[string]string {
	obj := map[string]string{}
	if !v.Truthy() {
		return obj
	}

	jsonStr := jsclass.JSON.Stringify(v).String()
	json.Unmarshal([]byte(jsonStr), &obj)

	return obj
}

func JSValueToMap(v js.Value) map[string]any {
	obj := map[string]any{}
	if !v.Truthy() {
		return obj
	}

	jsonStr := jsclass.JSON.Stringify(v).String()
	json.Unmarshal([]byte(jsonStr), &obj)

	return obj
}

func MaybeStringList(v js.Value) []string {
	if v.Truthy() {
		list := []string{}
		for i := range v.Length() {

			list = append(list, v.Index(i).String())
		}

		return list
	}

	return []string{}
}

func MaybeString(v js.Value) string {
	if v.IsUndefined() {
		return ""
	}
	return v.String()
}

func MaybeBool(v js.Value) bool {
	return v.Truthy()
}

func MaybeInt(v js.Value) int {
	if v.IsUndefined() {
		return 0
	}
	return v.Int()
}

func MaybeInt64(v js.Value) int64 {
	if v.Truthy() {
		vs := jsclass.String.Invoke(v)
		vi, _ := strconv.ParseInt(vs.String(), 10, 64)

		return vi
	}

	return 0
}

func MaybeDate(v js.Value) time.Time {
	if v.IsUndefined() {
		return time.Time{}
	}
	return DateToTime(v)
}

func DateToTimestamp(v js.Value) int64 {
	if v.InstanceOf(jsclass.Date) {
		ms := MaybeInt64(v.Call("getTime"))

		return ms
	}

	return -1
}

func DateToTime(v js.Value) time.Time {
	milli := MaybeInt64(v.Call("getTime"))
	return time.UnixMilli(milli)
}

func TimeToDate(t time.Time) js.Value {
	return jsclass.Date.New(t.UnixMilli())
}
