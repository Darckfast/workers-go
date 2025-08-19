//go:build js && wasm

package httptail

import (
	"encoding/json"
	"net/http"

	"github.com/Darckfast/workers-go/cloudflare/kv"
)

var GET_TAIL = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
	result, _ := namespace.GetString("tail:result", nil)

	if result == "<null>" {
		w.WriteHeader(404)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"result": result,
	})
}
