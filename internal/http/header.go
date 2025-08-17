package jshttp

import (
	"encoding/json"
	"net/http"
	"strings"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
)

func ToHeader(headers js.Value) http.Header {
	if !headers.Truthy() {
		return http.Header{}
	}

	hObj := jsclass.Object.FromEntries(headers.Call("entries"))
	hStr := jsclass.JSON.Stringify(hObj).String()
	var hMap map[string]string

	json.Unmarshal([]byte(hStr), &hMap)

	h := http.Header{}
	for i := range hMap {
		values := hMap[i]
		for value := range strings.SplitSeq(values, ",") {
			h.Add(i, value)
		}
	}
	return h
}

func ToJSHeader(header http.Header) js.Value {
	hBytes, _ := json.Marshal(header)
	hObj, _ := jsclass.JSON.Parse(string(hBytes))
	// Returning as an object is faster, but it has problems with Get(key)
	// on some headers keys
	h := jsclass.Headers.New(hObj)

	return h
}
