//go:build js && wasm

package httpenv

import (
	"encoding/json"
	"net/http"
	"os"
)

var GET_ENV = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(map[string]any{
		"result": os.Environ(),
	})
}
