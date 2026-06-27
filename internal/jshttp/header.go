//go:build js && wasm

package jshttp

import (
	"net/http"
	"strings"
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
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

	hObj := headers
	if headers.InstanceOf(jsclass.Headers.Class()) {
		hObj = jsclass.Object.FromEntries(headers)
	}

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

func HeaderToMap(header http.Header) map[string]string {
	hMap := jsclass.GenericStringMap{}
	for k, v := range header {
		if len(v) > 0 {
			hMap[k] = strings.Join(v, ",")
		}
	}
	return hMap
}

func ToJSHeader(header http.Header) js.Value {
	hMap := jsclass.GenericStringMap{}
	cookies := []string{}

	for k, v := range header {
		if len(v) > 0 {
			if strings.ToLower(k) == "set-cookie" {
				cookies = append(cookies, v...)
			} else {
				hMap[k] = strings.Join(v, ", ")
			}
		}
	}
	hBytes, _ := easyjson.Marshal(hMap)
	hObj, _ := jsclass.JSON.Parse(string(hBytes))
	// Returning as an object is faster, but it has problems with Get(key)
	// on some headers keys
	h := jsclass.Headers.New(hObj)

	if len(cookies) > 0 {
		for _, v := range cookies {
			h.Call("append", "set-cookie", v)
		}
	}

	return h
}
