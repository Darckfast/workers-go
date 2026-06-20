//go:build !js && !wasm

package main

import (
	"os"
	"path/filepath"
	"testing"
)

const fetchImpl = `
//go:build js && wasm

package main

import (
	"io"
	"net/http"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/fetch"
)

func main() {
	mux := http.NewServeMux()

  http.NewServeMux().HandleFunc()
	mux.HandlerFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
	  _, _ = w.Write([]byte("hello"))
  })

	fetch.ServeNonBlock(router)
}`

func TestFindFetch(t *testing.T) {
	tmpdir := t.TempDir()
	file, _ := os.Create(filepath.Join(tmpdir, "main.go"))

	file.Write([]byte(fetchImpl))

	main()
}
