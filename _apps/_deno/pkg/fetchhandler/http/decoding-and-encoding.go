//go:build js && wasm

package httpsimple

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

var GET_HELLO = func(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("hello"))
}

var GET_JSON = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"vitest": true})
}

var POST_JSON = func(w http.ResponseWriter, r *http.Request) {
	payload := map[string]any{}

	defer func() {
		_ = r.Body.Close()
	}()
	_ = json.NewDecoder(r.Body).Decode(&payload)

	b, _ := json.Marshal(payload)
	h := r.Header.Get("X-Test-Id")
	size := len(strconv.Quote(string(b)))
	w.Header().Set("Content-Type", "application/json")
	result := map[string]any{"raw": string(b), "size": size, "header": h, "query": r.URL.Query().Encode()}
	_ = json.NewEncoder(w).Encode(result)
}

var POST_FORM_URLENCODED = func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	alpha := r.FormValue("alpha")
	url := r.FormValue("url")
	name := r.FormValue("fullname")
	num := r.FormValue("number")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":       id,
		"alpha":    alpha,
		"url":      url,
		"fullname": name,
		"number":   num,
	})
}

var POST_MULTIPART_FORM_DATA = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = r.ParseForm()

	f, fh, err := r.FormFile("img")
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"has_error": true,
			"error":     err.Error(),
		})
		return
	}
	defer func() {
		_ = f.Close()
	}()

	buf := bytes.NewBuffer(make([]byte, 0))
	_, _ = io.Copy(buf, f)

	jsonStr := r.FormValue("json")
	var j map[string]any
	_ = json.Unmarshal([]byte(jsonStr), &j)
	jb, _ := json.Marshal(j)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"has_error":    err != nil,
		"size":         fh.Size,
		"filename":     fh.Filename,
		"actual-size":  buf.Len(),
		"content-type": fh.Header.Get("content-type"),
		"json":         string(jb),
	})
}
