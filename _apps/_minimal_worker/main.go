//go:build js && wasm

package main

import (
	"net/http"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/fetch"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello from workers-go"))
	})

	fetch.ServeNonBlock(mux)

	<-make(chan struct{})
}
