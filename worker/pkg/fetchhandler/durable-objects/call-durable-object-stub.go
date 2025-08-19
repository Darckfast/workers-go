//go:build js && wasm

package httpdurableobject

import (
	"encoding/json"
	"net/http"

	"github.com/Darckfast/workers-go/cloudflare/durableobjects"
)

var GET_DO = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	n, _ := durableobjects.NewDurableObjectNamespace("TEST_DO")
	objId := n.IdFromName("id")
	stub, _ := n.Get(objId)

	rs, err := stub.Call("sayHello")

	json.NewEncoder(w).Encode(map[string]any{
		"has_error": err != nil,
		"result":    rs,
	})
}
