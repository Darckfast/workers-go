//go:build js && wasm

package jshttp

import (
	"encoding/json"
	"net/http"
	"strings"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
)

func ToHeader(headers js.Value) (http.Header, error) {
	if !headers.Truthy() {
		return http.Header{}, nil
	}

	hObj := jsclass.Object.FromEntries(headers.Call("entries"))
	hStr := jsclass.JSON.Stringify(hObj).String()
	var hMap map[string]string

	err := json.Unmarshal([]byte(hStr), &hMap)

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
	hBytes, _ := json.Marshal(header)
	hObj, _ := jsclass.JSON.Parse(string(hBytes))
	// Returning as an object is faster, but it has problems with Get(key)
	// on some headers keys
	h := jsclass.Headers.New(hObj)

	return h
}
