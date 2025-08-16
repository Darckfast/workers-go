package jshttp

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"syscall/js"

	jsutil "github.com/syumai/workers/internal/utils"
)

func ToBody(stream js.Value) io.ReadCloser {
	if stream.IsNull() || stream.IsUndefined() {
		return io.NopCloser(bytes.NewReader([]byte{}))
	}
	return jsutil.ReadableStreamToReadCloser(stream)
}

// ToRequest converts JavaScript sides Request to *http.Request.
//   - Request: https://developer.mozilla.org/docs/Web/API/Request
func ToRequest(req js.Value) *http.Request {
	reqUrl, _ := url.Parse(req.Get("url").String())
	header := ToHeader(req.Get("headers"))

	// ignore err
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

// ToJSRequest converts *http.Request to JavaScript sides Request.
//   - Request: https://developer.mozilla.org/docs/Web/API/Request
func ToJSRequest(req *http.Request) js.Value {
	jsReqOptions := jsutil.NewObject()
	jsReqOptions.Set("method", req.Method)
	jsReqOptions.Set("headers", ToJSHeader(req.Header))
	jsReqBody := js.Undefined()
	if req.Body != nil && req.Method != http.MethodGet {
		jsReqBody = jsutil.ReadCloserToReadableStream(req.Body)
	}
	jsReqOptions.Set("body", jsReqBody)
	jsReq := jsutil.RequestClass.New(req.URL.String(), jsReqOptions)
	return jsReq
}
