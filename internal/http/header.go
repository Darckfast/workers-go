package jshttp

import (
	"net/http"
	"strings"
	"syscall/js"

	jsutil "github.com/syumai/workers/internal/utils"
)

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

func ToJSHeader(header http.Header) js.Value {
	h := jsutil.HeadersClass.New()
	for key, values := range header {
		for _, value := range values {
			h.Call("append", key, value)
		}
	}
	return h
}
