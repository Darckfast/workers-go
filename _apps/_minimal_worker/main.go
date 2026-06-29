//go:build js && wasm

package main

import (
	"net/http"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/fetch"
	"github.com/julienschmidt/httprouter"
)

func main() {
	mux := httprouter.New()

	mux.HandlerFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Header().Add("content-type", "text/plain")
		_, _ = w.Write([]byte("hello from workers-go"))
	})

	fetch.ServeNonBlock(mux)

	<-make(chan struct{})
}
