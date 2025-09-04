//go:build js && wasm

package main

import (
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	httpd1 "worker/pkg/fetchhandler/d1"

	"github.com/Darckfast/workers-go/cloudflare/fetch"
	"github.com/julienschmidt/httprouter"
)

func main() {
	mux := httprouter.New()
	mux.HandlerFunc("GET", "/d1", httpd1.GET_D1)
	mux.HandlerFunc("GET", "/", httpd1.GET_D1_TOTAL)
	mux.HandlerFunc("GET", "/debug/pprof", pprof.Index)
	mux.HandlerFunc("GET", "/debug/pprof/:name", pprof.Index)
	mux.HandlerFunc("GET", "/debug/profile", pprof.Profile)
	mux.HandlerFunc("POST", "/d1", httpd1.POST_D1)
	mux.HandlerFunc("GET", "/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	fetch.ServeNonBlock(mux)
	<-make(chan struct{})
}
