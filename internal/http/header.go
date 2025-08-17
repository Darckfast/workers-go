package jshttp

import (
	"encoding/json"
	"net/http"
	"strings"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
	jsutil "github.com/syumai/workers/internal/utils"
)

func ToHeaderV2(headers js.Value) http.Header {
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

func ToHeader(headers js.Value) http.Header {
	if !headers.Truthy() {
		return http.Header{}
	}

	entries := jsutil.ArrayFrom(headers.Call("entries"))
	headerLen := entries.Length()
	h := http.Header{}
	for i := range headerLen {
		entry := entries.Index(i)
		key := entry.Index(0).String()
		values := entry.Index(1).String()
		for value := range strings.SplitSeq(values, ",") {
			h.Add(key, value)
		}
	}
	return h
}

func ToJSHeaderV2(header http.Header) js.Value {
	hBytes, _ := json.Marshal(header)
	hObj, _ := jsclass.JSON.Parse(string(hBytes))
	// Returning as an object is faster, but it has problems with Get(key)
	// on some headers keys
	h := jsutil.HeadersClass.New(hObj)

	return h
}

func ToJSHeader(header http.Header) js.Value {
	h := jsutil.HeadersClass.New()
	for key, values := range header {
		for _, value := range values {
			h.Call("append", key, value)
		}
	}
	return h
}
