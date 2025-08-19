//go:build js && wasm

package httpr2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Darckfast/workers-go/cloudflare/r2"
)

var GET_R2 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bucket, _ := r2.NewBucket("TEST_BUCKET")

	result, err := bucket.Get("count")
	rawBody, _ := io.ReadAll(result.Body)
	b64 := base64.StdEncoding.EncodeToString(rawBody)

	json.NewEncoder(w).Encode(map[string]any{
		"has_error": err != nil,
		"result":    result,
		"body":      b64,
	})
}

var POST_R2 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bucket, _ := r2.NewBucket("TEST_BUCKET")

	b64 := r.FormValue("b64")
	data, err := base64.StdEncoding.DecodeString(b64)

	reader := io.NopCloser(bytes.NewReader(data))
	result, err := bucket.Put("count", reader, int64(len(data)), nil)

	json.NewEncoder(w).Encode(map[string]any{
		"has_error": err != nil,
		"result":    result,
	})
}
