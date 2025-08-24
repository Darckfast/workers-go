//go:build js && wasm

package httpd1

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	_ "github.com/Darckfast/workers-go/cloudflare/d1" // register driver
)

type TestingRow struct {
	Data      string `json:"data"`
	Id        int64  `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

var GET_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, _ := sql.Open("d1", "DB")

	id := r.URL.Query().Get("id")
	idi, _ := strconv.Atoi(id)

	result := db.QueryRow("SELECT id, data, created_at, updated_at FROM Testing WHERE id = ?", idi)

	var a TestingRow
	err := result.Scan(&a.Id, &a.Data, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": err.Error(),
		})
	} else {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": a,
		})
	}
}

var PUT_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, _ := sql.Open("d1", "DB")

	id := r.URL.Query().Get("id")
	idi, _ := strconv.Atoi(id)

	defer func() {
		_ = r.Body.Close()
	}()
	data, _ := io.ReadAll(r.Body)
	result, err := db.Exec("UPDATE Testing SET data = ?, updated_at = ( unixepoch('subsec') * 1000 ) WHERE id = ?", string(data), idi)
	rows, _ := result.RowsAffected()

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": err,
		"data":  rows,
	})
}

var POST_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, _ := sql.Open("d1", "DB")

	defer func() {
		_ = r.Body.Close()
	}()
	data, _ := io.ReadAll(r.Body)
	// D1 seems to not work with Vite, at least it does not find the table, even though it exists
	// testing directly with wrangler dev works fine
	result, err := db.Exec("INSERT INTO Testing (data) VALUES (?)", string(data))

	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": err.Error(),
		})
		return
	}
	id, _ := result.LastInsertId()
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": nil,
		"data": map[string]int64{
			"id": id,
		},
	})
}

var DELETE_D1 = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, _ := sql.Open("d1", "DB")

	id := r.URL.Query().Get("id")
	idi, _ := strconv.Atoi(id)
	result, err := db.Exec("DELETE FROM Testing  WHERE id = ?", idi)
	rows, _ := result.RowsAffected()

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": err,
		"data":  rows,
	})
}
