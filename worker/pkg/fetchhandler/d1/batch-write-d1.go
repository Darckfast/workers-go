//go:build js && wasm

package httpd1

import (
	"encoding/json"
	"net/http"

	"github.com/Darckfast/workers-go/cloudflare/d1/v2"
)

var POST_D1_BATCH = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, _ := d1.GetDB("DB")

	stmt1 := db.Prepare("INSERT INTO Testing (data) VALUES (?)").Bind("test")
	stmt2 := db.Prepare("INSERT INTO Testing (data) VALUES (?)").Bind("test 2")

	res, err := db.Batch([]d1.D1PreparedStatment{*stmt1, *stmt2})

	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": err.Error(),
		})
	} else {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": res,
		})
	}
}
