//go:build js && wasm

package jsconv

import (
	"strconv"
	"syscall/js"
	"time"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	"github.com/mailru/easyjson"
)

func ArrayFrom(v js.Value) js.Value {
	return jsclass.Array.Call("from", v)
}

func StrRecordToMap(v js.Value) jsclass.GenericStringMap {
	m := jsclass.GenericStringMap{}

	if v.IsUndefined() || v.IsNull() {
		return m
	}

	obj := jsclass.JSON.Stringify(v)
	_ = easyjson.Unmarshal([]byte(obj.String()), &m)

	return m
}

func MapToJSValue(v jsclass.GenericAnyMap) js.Value {
	b, _ := easyjson.Marshal(v)
	return jsclass.JSON.Call("parse", string(b))
}

func JSMapToMap(v js.Value) (map[string]any, error) {
	obj := map[string]any{}
	if !v.Truthy() {
		return obj, nil
	}

	vm := jsclass.Object.FromEntries(v)

	return JSValueToMap(vm)
}

func JSValueToMapString(v js.Value) (jsclass.GenericStringMap, error) {
	obj := jsclass.GenericStringMap{}
	if !v.Truthy() {
		return obj, nil
	}

	jsonStr := jsclass.JSON.Stringify(v).String()
	err := easyjson.Unmarshal([]byte(jsonStr), &obj)

	return obj, err
}

func JSValueToMap(v js.Value) (jsclass.GenericAnyMap, error) {
	obj := jsclass.GenericAnyMap{}
	if !v.Truthy() {
		return obj, nil
	}

	jsonStr := jsclass.JSON.Stringify(v).String()
	err := easyjson.Unmarshal([]byte(jsonStr), &obj)

	return obj, err
}

func MaybeStringList(v js.Value) []string {
	if v.Truthy() {
		list := []string{}
		for i := 0; i < v.Length(); i++ {
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
