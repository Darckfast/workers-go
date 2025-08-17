package jshttp

import (
	"io"
	"net/http"
	"strconv"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
	jsstream "github.com/syumai/workers/internal/stream"
	jsutil "github.com/syumai/workers/internal/utils"
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

func newJSResponse(statusCode int, headers http.Header, contentLength int64, body io.ReadCloser, rawBody *js.Value) js.Value {
	status := statusCode
	if status == 0 {
		status = http.StatusOK
	}
	respInit := jsutil.NewObject()
	respInit.Set("status", status)
	respInit.Set("statusText", http.StatusText(status))
	respInit.Set("headers", ToJSHeader(headers))
	if status == http.StatusSwitchingProtocols ||
		status == http.StatusNoContent ||
		status == http.StatusResetContent ||
		status == http.StatusNotModified {
		return jsutil.ResponseClass.New(jsutil.Null, respInit)
	}
	readableStream := func() js.Value {
		if rawBody != nil {
			return *rawBody
		}
		if !jsutil.MaybeFixedLengthStreamClass.IsUndefined() && contentLength > 0 {
			return jsstream.ReadCloserToFixedLengthStream(body, contentLength)
		}
		return jsstream.ReadCloserToReadableStream(body)
	}()
	return jsutil.ResponseClass.New(readableStream, respInit)
}

func ToJSResponse(res *http.Response) js.Value {
	status := res.StatusCode
	if status == 0 {
		status = http.StatusOK
	}

	respInit := jsutil.NewObject()
	respInit.Set("status", status)
	respInit.Set("statusText", http.StatusText(status))
	respInit.Set("headers", ToJSHeader(res.Header))
	readableStream := jsclass.Null

	if status == http.StatusSwitchingProtocols ||
		status == http.StatusNoContent ||
		status == http.StatusResetContent ||
		status == http.StatusNotModified {
		return jsutil.ResponseClass.New(readableStream, respInit)
	}

	contentLength, _ := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
	if jsclass.MaybeFixedLengthStream.Truthy() && contentLength > 0 {
		readableStream = jsstream.ReadCloserToFixedLengthStream(res.Body, contentLength)
	} else {
		readableStream = jsstream.ReadCloserToReadableStream(res.Body)
	}

	return jsutil.ResponseClass.New(readableStream, respInit)

}
