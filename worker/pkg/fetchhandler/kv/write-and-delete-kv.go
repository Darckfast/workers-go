//go:build js && wasm

package httpkv

import (
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/Darckfast/workers-go/cloudflare/d1" // register driver

	"github.com/Darckfast/workers-go/cloudflare/kv"
)

var DELETE_KV = func(w http.ResponseWriter, r *http.Request) {
	namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
	err := namespace.Delete("count")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"has_error": err != nil})
}

var POST_KV = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	namespace, _ := kv.NewNamespace("TEST_NAMESPACE")

	countStr, _ := namespace.GetString("count", nil)
	count, _ := strconv.Atoi(countStr)

	nextCountStr := strconv.Itoa(count + 1)

	err := namespace.PutString("count", nextCountStr, nil)
	json.NewEncoder(w).Encode(map[string]any{"has_error": err != nil, "count": nextCountStr})
}
