//go:build js && wasm

package httpqueue

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Darckfast/workers-go/cloudflare/kv"
	"github.com/Darckfast/workers-go/cloudflare/queues"
)

var GET_QUEUE = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
	result, _ := namespace.GetString("queue:result", nil)

	if result == "<null>" {
		w.WriteHeader(404)
	}
	json.NewEncoder(w).Encode(map[string]any{
		"result": result,
	})
}

var POST_QUEUE = func(w http.ResponseWriter, r *http.Request) {
	q, err := queues.NewProducer("TEST_QUEUE")
	content, _ := io.ReadAll(r.Body)
	err = q.SendText(string(content))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(map[string]any{
		"has_error": err != nil,
	})
}
