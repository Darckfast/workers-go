//go:build js && wasm

package fetchhandler

import (
	"net/http"

	httpcache "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/cache-api"
	httpcontainer "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/container"
	httpd1 "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/d1"
	httpdurableobject "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/durable-objects"
	httpenv "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/env"
	errorshandler "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/errors"
	httpsimple "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/http"
	httpkv "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/kv"
	httpqueue "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/queue"
	httpr2 "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/r2"
	httpsocket "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/sockets"
	httptail "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/tail"

	"github.com/Darckfast/workers-go/cloudflare/fetch"
)

func New() {
	// HTTP
	http.HandleFunc("GET /hello", httpsimple.GET_HELLO)
	http.HandleFunc("GET /application/json", httpsimple.GET_JSON)
	http.HandleFunc("POST /application/json", httpsimple.POST_JSON)
	http.HandleFunc("POST /application/x-www-form-urlencoded", httpsimple.POST_FORM_URLENCODED)
	http.HandleFunc("POST /multipart/form-data", httpsimple.POST_MULTIPART_FORM_DATA)

	// KV
	http.HandleFunc("DELETE /kv", httpkv.DELETE_KV)
	http.HandleFunc("POST /kv", httpkv.POST_KV)
	http.HandleFunc("POST /kv/meta", httpkv.POST_KV_META)
	http.HandleFunc("GET /kv/meta", httpkv.GET_KV_META)
	http.HandleFunc("GET /kv", httpkv.GET_KV)
	http.HandleFunc("GET /kvs", httpkv.GET_KVS)
	http.HandleFunc("GET /kv/list", httpkv.GET_KV_LIST)

	// R2
	http.HandleFunc("GET /r2", httpr2.GET_R2)
	http.HandleFunc("POST /r2", httpr2.POST_R2)
	// Cache api
	http.HandleFunc("GET /cache", httpcache.GET_CACHE)
	// D1
	http.HandleFunc("GET /d1", httpd1.GET_D1)
	http.HandleFunc("POST /d1", httpd1.POST_D1)
	http.HandleFunc("PUT /d1", httpd1.PUT_D1)
	http.HandleFunc("DELETE /d1", httpd1.DELETE_D1)
	// Queue
	http.HandleFunc("GET /queue", httpqueue.GET_QUEUE)
	http.HandleFunc("POST /queue", httpqueue.POST_QUEUE)
	// Socket
	http.HandleFunc("GET /socket", httpsocket.GET_SOCKET_TCPBIN)
	// Tail
	http.HandleFunc("GET /tail", httptail.GET_TAIL)
	// Container
	http.HandleFunc("GET /container", httpcontainer.GET_CONTAINER)
	// Durable Object
	http.HandleFunc("GET /do", httpdurableobject.GET_DO)
	// Env
	http.HandleFunc("GET /env", httpenv.GET_ENV)

	//Error
	http.HandleFunc("GET /error", errorshandler.GET_ERROR)

	/*
	 * Fetch handler uses http.DefaulServeMux as default, calling this function is
	 * optional, unless a different handler must be used
	 */
	fetch.ServeNonBlock(nil)
}
