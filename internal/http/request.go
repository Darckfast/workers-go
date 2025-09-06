//go:build js && wasm

package jshttp

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	reqStr := jsclass.JSON.Stringify(req, []any{"method", "url"})
	var reqMap JSRequest

	_ = easyjson.Unmarshal([]byte(reqStr.String()), &reqMap)

	reqUrl, _ := url.Parse(reqMap.Url)
	headers, _ := ToHeader(req.Get("headers"))

	contentLength, _ := strconv.ParseInt(headers.Get("Content-Length"), 10, 64)
	return &http.Request{
		Method:           reqMap.Method,
		URL:              reqUrl,
		Header:           headers,
		Body:             ToBody(req.Get("body")),
		ContentLength:    contentLength,
		TransferEncoding: strings.Split(headers.Get("Transfer-Encoding"), ","),
		Host:             headers.Get("Host"),
	}
}

func ToJSRequest(req *http.Request) js.Value {
	jsReq := JSRequest{
		Url:     req.URL.String(),
		Method:  req.Method,
		Headers: HeaderToMap(req.Header),
	}

	jsReqB, _ := easyjson.Marshal(jsReq)
	jsReqOptions, _ := jsclass.JSON.Parse(jsReqB)
	jsReqBody := js.Undefined()

	if req.Body != nil && req.Method != http.MethodGet {
		jsReqBody = jsstream.ReadCloserToReadableStream(req.Body)
		jsReqOptions.Set("duplex", "half")
	}
	jsReqOptions.Set("body", jsReqBody)
	return jsclass.Request.New(req.URL.String(), jsReqOptions)
}
