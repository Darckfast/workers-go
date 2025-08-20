//go:build js && wasm

package httpkv

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Darckfast/workers-go/cloudflare/kv"
)

var GET_KV_LIST = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	kvStore, _ := kv.NewNamespace("TEST_NAMESPACE")
	data, _ := kvStore.List(nil)

	_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
}

var GET_KV = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	key := r.URL.Query().Get("key")

	if key == "" {
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "missing key"})
		w.WriteHeader(400)
		return
	}

	kvStore, _ := kv.NewNamespace("TEST_NAMESPACE")
	data, err := kvStore.GetReader(key, nil)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	_, err = io.Copy(w, data)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
}

var DELETE_KV = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key := r.URL.Query().Get("key")

	if key == "" {
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "missing key"})
		w.WriteHeader(400)
		return
	}

	kvStore, _ := kv.NewNamespace("TEST_NAMESPACE")
	err := kvStore.Delete(key)

	_ = json.NewEncoder(w).Encode(map[string]any{"error": err})
}

var POST_KV = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key := r.URL.Query().Get("key")

	if key == "" {
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "missing key"})
		w.WriteHeader(400)
		return
	}

	kvStore, _ := kv.NewNamespace("TEST_NAMESPACE")

	defer func() {
		_ = r.Body.Close()
	}()

	err := kvStore.PutReader(key, r.Body, nil)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": err,
	})
}
