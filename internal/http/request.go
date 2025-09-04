//go:build js && wasm

package jshttp

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
	"github.com/mailru/easyjson"
)

func ToBody(body js.Value) io.ReadCloser {
	if !body.Truthy() {
		return io.NopCloser(bytes.NewReader([]byte{}))
	}

	return jsstream.ReadableStreamToReadCloser(body)
}

func ToRequest(req js.Value) *http.Request {
	reqStr := jsclass.JSON.Stringify(req, []any{"method", "url", "headers"})
	var reqMap JSRequest

	_ = easyjson.Unmarshal([]byte(reqStr.String()), &reqMap)
	reqUrl, _ := url.Parse(reqMap.Url)
	header, _ := MapToHeader(reqMap.Headers)

	contentLength, _ := strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	return &http.Request{
		Method:           reqMap.Method,
		URL:              reqUrl,
		Header:           header,
		Body:             ToBody(req.Get("body")),
		ContentLength:    contentLength,
		TransferEncoding: header["transfer-encoding"],
		Host:             header.Get("Host"),
	}
}

func ToJSRequest(req *http.Request) js.Value {
	jsReqOptions := jsclass.Object.New()
	jsReqOptions.Set("method", req.Method)
	jsReqOptions.Set("headers", ToJSHeader(req.Header))
	jsReqBody := js.Undefined()
	if req.Body != nil && req.Method != http.MethodGet {
		jsReqBody = jsstream.ReadCloserToReadableStream(req.Body)
		jsReqOptions.Set("duplex", "half")
	}
	jsReqOptions.Set("body", jsReqBody)
	jsReq := jsclass.Request.New(req.URL.String(), jsReqOptions)
	return jsReq
}
