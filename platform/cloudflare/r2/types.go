/*
Package r2 is the glue code for Cloudflare's R2 bindings
*/
package r2

import (
	"io"
	"net/http"
	"time"
)

//easyjson:json
type R2Conditional struct {
	UploadedBefore   time.Time `json:"uploadedBefore"`
	UploadedAfter    time.Time `json:"uploadedAfter"`
	EtagMatches      string    `json:"etagMatches"`
	EtagDoesNotMatch string    `json:"etagDoesNotMatch"`
}

type R2Range struct {
	Offset int64 `json:"offset,omitempty"`
	Length int64 `json:"length,omitempty"`
	Suffix int64 `json:"suffix,omitempty"`
}

//easyjson:json
type GetOptions struct {
	OnlyIf  *R2Conditional `json:"onlyIf,omitempty"`
	Range   *R2Range       `json:"range"`
	SSecKey string         `json:"ssecKey,omitempty"`
}

//easyjson:json
type PutOptions struct {
	OnlyIf         *R2Conditional    `json:"onlyIf,omitempty"`
	HTTPMetadata   http.Header       `json:"httpMetadata,omitempty"`
	CustomMetadata map[string]string `json:"customMetadata,omitempty"`
	MD5            string            `json:"md5,omitempty"`
	SHA1           string            `json:"sha1,omitempty"`
	SHA256         string            `json:"sha256,omitempty"`
	SHA384         string            `json:"sha384,omitempty"`
	SHA512         string            `json:"sha512,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	SSecKey        string            `json:"ssecKey,omitempty"`
}

//easyjson:json
type MultipartOptions struct {
	HTTPMetadata   http.Header       `json:"httpMetadata,omitempty"`
	CustomMetadata map[string]string `json:"customMetadata,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	SSecKey        string            `json:"ssecKey,omitempty"`
}

//easyjson:json
type ListOptions struct {
	Prefix    string   `json:"prefix,omitempty"`
	Cursor    string   `json:"cursor,omitempty"`
	Delimiter string   `json:"delimiter,omitempty"`
	Include   []string `json:"include,omitempty"`
	Limit     int64    `json:"limit,omitempty"`
}

//easyjson:json
type UploadedPart struct {
	ETag       string `json:"etag"`
	PartNumber int64  `json:"partNumber"`
}

//easyjson:json
type R2Object struct {
	Uploaded       time.Time         `json:"uploaded"`
	Body           io.Reader         `json:"-"`
	HTTPMetadata   map[string]string `json:"httpMetadata"`
	CustomMetadata map[string]string `json:"customMetadata"`
	Key            string            `json:"key"`
	Version        string            `json:"version"`
	ETag           string            `json:"etag"`
	HTTPETag       string            `json:"httpEtag"`
	Size           int               `json:"size"`
}

//easyjson:json
type R2Objects struct {
	Cursor            string      `json:"cursor,omitempty"`
	Objects           []*R2Object `json:"objects,omitempty"`
	DelimitedPrefixes []string    `json:"delimitedPrefixes,omitempty"`
	Truncated         bool        `json:"truncated"`
}
