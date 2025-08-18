//go:build js && wasm

package jshttp

import (
	"net/http"
	"strconv"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
)

func ToResponse(res js.Value) *http.Response {
	body := jsstream.ReadableStreamToReadCloser(res.Get("body"))
	status := res.Get("status").Int()
	header := ToHeader(res.Get("headers"))
	contentLength, _ := strconv.ParseInt(header.Get("Content-Length"), 10, 64)

	return &http.Response{
		Status:        http.StatusText(status),
		StatusCode:    status,
		Header:        header,
		Body:          body,
		ContentLength: contentLength,
	}
}

func ToJSResponse(res *http.Response) js.Value {
	status := res.StatusCode
	if status == 0 {
		status = http.StatusOK
	}

	respInit := jsclass.Object.New()
	respInit.Set("status", status)
	respInit.Set("statusText", http.StatusText(status))
	respInit.Set("headers", ToJSHeader(res.Header))
	readableStream := jsclass.Null

	if status == http.StatusSwitchingProtocols ||
		status == http.StatusNoContent ||
		status == http.StatusResetContent ||
		status == http.StatusNotModified {
		return jsclass.Response.New(readableStream, respInit)
	}

	contentLength, _ := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
	if jsclass.MaybeFixedLengthStream.Truthy() && contentLength > 0 {
		readableStream = jsstream.ReadCloserToFixedLengthStream(res.Body, contentLength)
	} else {
		readableStream = jsstream.ReadCloserToReadableStream(res.Body)
	}

	return jsclass.Response.New(readableStream, respInit)

}
