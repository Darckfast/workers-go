//go:build js && wasm

package main

import (
	"net/http/pprof"
	_ "net/http/pprof"
	httpsimple "worker/pkg/fetchhandler/http"

	"github.com/Darckfast/workers-go/cloudflare/fetch"
	"github.com/julienschmidt/httprouter"
)

func main() {
	mux := httprouter.New()

	mux.HandlerFunc("GET", "/hello", httpsimple.GET_HELLO)
	mux.HandlerFunc("GET", "/application/json", httpsimple.GET_JSON)
	mux.HandlerFunc("POST", "/application/json", httpsimple.POST_JSON)
	mux.HandlerFunc("POST", "/application/x-www-form-urlencoded", httpsimple.POST_FORM_URLENCODED)
	mux.HandlerFunc("POST", "/multipart/form-data", httpsimple.POST_MULTIPART_FORM_DATA)

	mux.HandlerFunc("GET", "/debug/pprof", pprof.Index)
	mux.HandlerFunc("GET", "/debug/pprof/:name", pprof.Index)
	mux.HandlerFunc("GET", "/debug/profile", pprof.Profile)

	fetch.ServeNonBlock(mux)
	<-make(chan struct{})
}
