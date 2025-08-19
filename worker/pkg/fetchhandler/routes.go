//go:build js && wasm

package fetchhandler

import (
	"net/http"

	httpcache "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/cache-api"
	httpcontainer "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/container"
	httpd1 "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/d1"
	httpdurableobject "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/durable-objects"
	httpenv "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/env"
	httpsimple "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/http"
	httpkv "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/kv"
	httpqueue "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/queue"
	httpr2 "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/r2"
	httpsocket "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/sockets"
	httptail "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/tail"

	"github.com/Darckfast/workers-go/cloudflare/fetch"
)

func New() {
	http.HandleFunc("GET /hello", httpsimple.GET_HELLO)
	http.HandleFunc("GET /application/json", httpsimple.GET_JSON)
	http.HandleFunc("POST /application/json", httpsimple.POST_JSON)
	http.HandleFunc("POST /application/x-www-form-urlencoded", httpsimple.POST_FORM_URLENCODED)
	http.HandleFunc("POST /multipart/form-data", httpsimple.POST_MULTIPART_FORM_DATA)
	http.HandleFunc("DELETE /kv", httpkv.DELETE_KV)
	http.HandleFunc("POST /kv", httpkv.POST_KV)
	http.HandleFunc("GET /r2", httpr2.GET_R2)
	http.HandleFunc("POST /r2", httpr2.POST_R2)
	http.HandleFunc("GET /cache", httpcache.GET_CACHE)
	http.HandleFunc("GET /d1", httpd1.GET_D1)
	http.HandleFunc("GET /queue", httpqueue.GET_QUEUE)
	http.HandleFunc("POST /queue", httpqueue.POST_QUEUE)
	http.HandleFunc("GET /socket", httpsocket.GET_SOCKET_TCPBIN)
	http.HandleFunc("GET /tail", httptail.GET_TAIL)
	http.HandleFunc("GET /container", httpcontainer.GET_CONTAINER)
	http.HandleFunc("GET /do", httpdurableobject.GET_DO)
	http.HandleFunc("GET /env", httpenv.GET_ENV)

	/*
	 * Fetch handler uses http.DefaulServeMux as default, calling this function is
	 * optional, unless a different handler must be used
	 */
	fetch.ServeNonBlock(nil)
}
