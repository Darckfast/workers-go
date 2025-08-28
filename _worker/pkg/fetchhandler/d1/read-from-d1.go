//go:build js && wasm

package httpd1

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Darckfast/workers-go/cloudflare/d1/v2"
)

type TestingRow struct {
	Data      string `json:"data"`
	Id        int64  `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

var GET_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var db, _ = d1.GetDB("DB")

	id := r.URL.Query().Get("id")
	idi, _ := strconv.Atoi(id)

	result, err := db.Prepare("SELECT id, data, created_at, updated_at FROM Testing WHERE id = ?").
		Bind(idi).
		Run()

	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": err.Error(),
		})
	} else {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": result,
		})
	}
}

var PUT_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var db, _ = d1.GetDB("DB")
	id := r.URL.Query().Get("id")
	idi, _ := strconv.Atoi(id)

	defer func() {
		_ = r.Body.Close()
	}()

	data, _ := io.ReadAll(r.Body)
	result, err := db.Prepare("UPDATE Testing SET data = ?, updated_at = ( unixepoch('subsec') * 1000 ) WHERE id = ?").
		Bind(string(data), idi).
		Run()

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": err,
		"data":  result,
	})
}

var POST_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var db, _ = d1.GetDB("DB")
	defer func() {
		_ = r.Body.Close()
	}()
	data, _ := io.ReadAll(r.Body)
	// D1 seems to not work with Vite, at least it does not find the table, even though it exists
	// testing directly with wrangler dev works fine
	result, err := db.Prepare("INSERT INTO Testing (data) VALUES (?)").
		Bind(string(data)).
		Run()

	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": err.Error(),
		})
		return
	}
	id := result.Meta.LastRowId
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": nil,
		"data": map[string]int64{
			"id": id,
		},
	})
}

var DELETE_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var db, _ = d1.GetDB("DB")
	id := r.URL.Query().Get("id")
	idi, _ := strconv.Atoi(id)
	result, err := db.Prepare("DELETE FROM Testing  WHERE id = ?").
		Bind(idi).
		Run()

	rows := result.Meta.Changes

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": err,
		"data":  rows,
	})
}
