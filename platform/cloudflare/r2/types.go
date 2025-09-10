package r2

import (
	"io"
	"net/http"
	"time"
)

//easyjson:json
type R2Conditional struct {
	EtagMatches      string    `json:"etagMatches"`
	EtagDoesNotMatch string    `json:"etagDoesNotMatch"`
	UploadedBefore   time.Time `json:"uploadedBefore"`
	UploadedAfter    time.Time `json:"uploadedAfter"`
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
	Limit     int64    `json:"limit,omitempty"`
	Prefix    string   `json:"prefix,omitempty"`
	Cursor    string   `json:"cursor,omitempty"`
	Delimiter string   `json:"delimiter,omitempty"`
	Include   []string `json:"include,omitempty"`
}

//easyjson:json
type UploadedPart struct {
	ETag       string `json:"etag"`
	PartNumber int64  `json:"partNumber"`
}

//easyjson:json
type R2Object struct {
	Key            string            `json:"key"`
	Version        string            `json:"version"`
	Size           int               `json:"size"`
	ETag           string            `json:"etag"`
	HTTPETag       string            `json:"httpEtag"`
	Uploaded       time.Time         `json:"uploaded"`
	HTTPMetadata   map[string]string `json:"httpMetadata"`
	CustomMetadata map[string]string `json:"customMetadata"`
	Body           io.Reader         `json:"-"`
}

//easyjson:json
type R2Objects struct {
	Objects           []*R2Object `json:"objects,omitempty"`
	Truncated         bool        `json:"truncated"`
	Cursor            string      `json:"cursor,omitempty"`
	DelimitedPrefixes []string    `json:"delimitedPrefixes,omitempty"`
}
