package jshttp

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"syscall/js"

	jsclass "github.com/syumai/workers/internal/class"
	jsstream "github.com/syumai/workers/internal/stream"
)

func ToBody(body js.Value) io.ReadCloser {
	if !body.Truthy() {
		return io.NopCloser(bytes.NewReader([]byte{}))
	}

	return jsstream.ReadableStreamToReadCloser(body)
}

func ToRequest(req js.Value) *http.Request {
	reqUrl, _ := url.Parse(req.Get("url").String())
	header := ToHeader(req.Get("headers"))

	contentLength, _ := strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	return &http.Request{
		Method:           req.Get("method").String(),
		URL:              reqUrl,
		Header:           header,
		Body:             ToBody(req.Get("body")),
		ContentLength:    contentLength,
		TransferEncoding: strings.Split(header.Get("Transfer-Encoding"), ","),
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
