//go:build js && wasm

package httpenv

import (
	"encoding/json"
	"net/http"
	"os"

	_ "github.com/Darckfast/workers-go/cloudflare/d1" // register driver
)

var GET_ENV = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(map[string]any{
		"result": os.Environ(),
	})
}
