//go:build js && wasm

package httpd1

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/Darckfast/workers-go/cloudflare/d1" // register driver
)

var GET_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, _ := sql.Open("d1", "DB")

	result := db.QueryRow("SELECT current_timestamp")

	var a any
	result.Scan(&a)
	json.NewEncoder(w).Encode(map[string]any{
		"result": a,
	})
}
