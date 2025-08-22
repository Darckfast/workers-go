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
	result, err := namespace.Get([]string{"tail:result"}, 0)

	if err != nil {
		w.WriteHeader(404)
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"result": result,
	})
}
