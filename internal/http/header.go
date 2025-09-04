//go:build js && wasm

package jshttp

import (
	"net/http"
	"strings"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	"github.com/mailru/easyjson"
)

func MapToHeader(headers jsclass.GenericStringMap) (http.Header, error) {
	h := http.Header{}
	for i := range headers {
		values := headers[i]
		for _, value := range strings.Split(values, ",") {
			h.Add(i, value)
		}
	}

	return h, nil
}
func ToHeader(headers js.Value) (http.Header, error) {
	if !headers.Truthy() {
		return http.Header{}, nil
	}

	hObj := jsclass.Object.FromEntries(headers.Call("entries"))
	hStr := jsclass.JSON.Stringify(hObj).String()
	var hMap jsclass.GenericStringMap

	err := easyjson.Unmarshal([]byte(hStr), &hMap)

	if err != nil {
		return http.Header{}, err
	}

	h := http.Header{}
	for i := range hMap {
		values := hMap[i]
		for _, value := range strings.Split(values, ",") {
			h.Add(i, value)
		}
	}

	return h, nil
}

func ToJSHeader(header http.Header) js.Value {
	hMap := jsclass.GenericStringMap{}
	for k, v := range header {
		if len(v) > 0 {
			hMap[k] = strings.Join(v, ",")
		}
	}
	hBytes, _ := easyjson.Marshal(hMap)
	hObj, _ := jsclass.JSON.Parse(string(hBytes))
	// Returning as an object is faster, but it has problems with Get(key)
	// on some headers keys
	h := jsclass.Headers.New(hObj)

	return h
}
