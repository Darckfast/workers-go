//go:build js && wasm

package httpcontainer

import (
	"encoding/json"
	"net/http"

	"github.com/Darckfast/workers-go/cloudflare/durableobjects"
)

var GET_CONTAINER = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	c, err := durableobjects.GetContainer("GO_CONTAINER", "test")

	rs, _ := c.ContainerFetch(r)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"has_error": err != nil,
		"result":    rs.Status,
	})
}
