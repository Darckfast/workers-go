//go:build js && wasm

package r2

import (
	"errors"
	"io"
	"syscall/js"
	"time"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
	jsstream "github.com/Darckfast/workers-go/internal/stream"
)

// Object represents Cloudflare R2 object.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1094
type Object struct {
	instance       js.Value
	Key            string
	Version        string
	Size           int
	ETag           string
	HTTPETag       string
	Uploaded       time.Time
	HTTPMetadata   HTTPMetadata
	CustomMetadata map[string]string
	// Body is a body of Object.
	// This value is nil for the result of the `Head` or `Put` method.
	Body io.Reader
}

// TODO: implement
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1106
// func (o *Object) WriteHTTPMetadata(headers http.Header) {
// }

func (o *Object) BodyUsed() (bool, error) {
	v := o.instance.Get("bodyUsed")
	if v.IsUndefined() {
		return false, errors.New("bodyUsed doesn't exist for this Object")
	}
	return v.Bool(), nil
}

// toObject converts JavaScript side's Object to *Object.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1094
func toObject(v js.Value) (*Object, error) {
	uploaded, err := jsconv.DateToTime(v.Get("uploaded"))
	if err != nil {
		return nil, errors.New("error converting uploaded: " + err.Error())
	}
	r2Meta, err := toHTTPMetadata(v.Get("httpMetadata"))
	if err != nil {
		return nil, errors.New("error converting httpMetadata: " + err.Error())
	}
	bodyVal := v.Get("body")
	var body io.Reader
	if !bodyVal.IsUndefined() {
		body = jsstream.ReadableStreamToReadCloser(v.Get("body"))
	}
	return &Object{
		instance:       v,
		Key:            v.Get("key").String(),
		Version:        v.Get("version").String(),
		Size:           v.Get("size").Int(),
		ETag:           v.Get("etag").String(),
		HTTPETag:       v.Get("httpEtag").String(),
		Uploaded:       uploaded,
		HTTPMetadata:   r2Meta,
		CustomMetadata: jsconv.StrRecordToMap(v.Get("customMetadata")),
		Body:           body,
	}, nil
}

// HTTPMetadata represents metadata of Object.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1053
type HTTPMetadata struct {
	ContentType        string
	ContentLanguage    string
	ContentDisposition string
	ContentEncoding    string
	ContentLength      int64
	CacheControl       string
	CacheExpiry        time.Time
}

func toHTTPMetadata(v js.Value) (HTTPMetadata, error) {
	if v.IsUndefined() || v.IsNull() {
		return HTTPMetadata{}, nil
	}
	cacheExpiry, err := jsconv.MaybeDate(v.Get("cacheExpiry"))
	if err != nil {
		return HTTPMetadata{}, errors.New("error converting cacheExpiry: " + err.Error())
	}
	return HTTPMetadata{
		ContentType:        jsconv.MaybeString(v.Get("contentType")),
		ContentLanguage:    jsconv.MaybeString(v.Get("contentLanguage")),
		ContentDisposition: jsconv.MaybeString(v.Get("contentDisposition")),
		ContentEncoding:    jsconv.MaybeString(v.Get("contentEncoding")),
		ContentLength:      jsconv.MaybeInt64(v.Get("contentLength")),
		CacheControl:       jsconv.MaybeString(v.Get("cacheControl")),
		CacheExpiry:        cacheExpiry,
	}, nil
}

func (md *HTTPMetadata) toJS() js.Value {
	obj := jsclass.Object.New()
	kv := map[string]any{
		"contentType":        md.ContentType,
		"contentLanguage":    md.ContentLanguage,
		"contentLength":      md.ContentLength,
		"contentDisposition": md.ContentDisposition,
		"contentEncoding":    md.ContentEncoding,
		"cacheControl":       md.CacheControl,
	}
	for k, v := range kv {
		if v != "" {
			obj.Set(k, v)
		}
	}
	if !md.CacheExpiry.IsZero() {
		obj.Set("cacheExpiry", jsconv.TimeToDate(md.CacheExpiry))
	}
	return obj
}
