//go:build js && wasm

package main

import (
	"net/http"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/durableobjects"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/fetch"
	"github.com/julienschmidt/httprouter"
)

func main() {
	mux := httprouter.New()

	mux.HandlerFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		n, _ := durableobjects.NewDurableObjectNamespace("MY_DURABLE_OBJECT")
		stub, _ := n.Get(n.IdFromName("id"))
		v, _ := stub.Call("SayHello")
		_, _ = w.Write([]byte(v.(string)))
	})

	fetch.ServeNonBlock(mux)

	select {}
}
