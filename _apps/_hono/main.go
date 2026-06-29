//go:build js && wasm

package main

import (
	httpsimple "worker/pkg/fetchhandler/http"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/fetch"
	"github.com/julienschmidt/httprouter"
)

func main() {
	mux := httprouter.New()

	mux.HandlerFunc("GET", "/", httpsimple.GET_HELLO)
	mux.HandlerFunc("GET", "/application/json", httpsimple.GET_JSON)
	mux.HandlerFunc("POST", "/application/json", httpsimple.POST_JSON)
	mux.HandlerFunc("POST", "/application/x-www-form-urlencoded", httpsimple.POST_FORM_URLENCODED)
	mux.HandlerFunc("POST", "/multipart/form-data", httpsimple.POST_MULTIPART_FORM_DATA)

	fetch.ServeNonBlock(mux)
	<-make(chan struct{})
}
