//go:build js && wasm

package fetchhandler

import (
	httpcache "worker/pkg/fetchhandler/cache-api"
	httpcontainer "worker/pkg/fetchhandler/container"
	httpd1 "worker/pkg/fetchhandler/d1"
	httpdurableobject "worker/pkg/fetchhandler/durable-objects"
	httpenv "worker/pkg/fetchhandler/env"
	errorshandler "worker/pkg/fetchhandler/errors"
	httpsimple "worker/pkg/fetchhandler/http"
	httpkv "worker/pkg/fetchhandler/kv"
	httpqueue "worker/pkg/fetchhandler/queue"
	httpr2 "worker/pkg/fetchhandler/r2"
	httpsocket "worker/pkg/fetchhandler/sockets"
	httptail "worker/pkg/fetchhandler/tail"

	"github.com/julienschmidt/httprouter"

	"github.com/Darckfast/workers-go/cloudflare/fetch"
)

func New() {
	router := httprouter.New()
	// HTTP
	router.HandlerFunc("GET", "/hello", httpsimple.GET_HELLO)
	router.HandlerFunc("GET", "/application/json", httpsimple.GET_JSON)
	router.HandlerFunc("POST", "/application/json", httpsimple.POST_JSON)
	router.HandlerFunc("POST", "/application/x-www-form-urlencoded", httpsimple.POST_FORM_URLENCODED)
	router.HandlerFunc("POST", "/multipart/form-data", httpsimple.POST_MULTIPART_FORM_DATA)

	//", "KV
	router.HandlerFunc("DELETE", "/kv", httpkv.DELETE_KV)
	router.HandlerFunc("POST", "/kv", httpkv.POST_KV)
	router.HandlerFunc("POST", "/kv/meta", httpkv.POST_KV_META)
	router.HandlerFunc("GET", "/kv/meta", httpkv.GET_KV_META)
	router.HandlerFunc("GET", "/kv", httpkv.GET_KV)
	router.HandlerFunc("GET", "/kvs", httpkv.GET_KVS)
	router.HandlerFunc("GET", "/kv/list", httpkv.GET_KV_LIST)

	//", "R2
	router.HandlerFunc("GET", "/r2", httpr2.GET_R2)
	router.HandlerFunc("POST", "/r2", httpr2.POST_R2)
	//", "Cache api
	router.HandlerFunc("GET", "/cache", httpcache.GET_CACHE)
	//", "D1
	router.HandlerFunc("GET", "/d1", httpd1.GET_D1)
	router.HandlerFunc("POST", "/d1", httpd1.POST_D1)
	router.HandlerFunc("PUT", "/d1", httpd1.PUT_D1)
	router.HandlerFunc("DELETE", "/d1", httpd1.DELETE_D1)
	router.HandlerFunc("POST", "/d1/batch", httpd1.POST_D1_BATCH)
	//", "Queue
	router.HandlerFunc("GET", "/queue", httpqueue.GET_QUEUE)
	router.HandlerFunc("POST", "/queue", httpqueue.POST_QUEUE)
	//", "Socket
	router.HandlerFunc("GET", "/socket", httpsocket.GET_SOCKET_TCPBIN)
	//", "Tail
	router.HandlerFunc("GET", "/tail", httptail.GET_TAIL)
	//", "Container
	router.HandlerFunc("GET", "/container", httpcontainer.GET_CONTAINER)
	//", "Durable Object
	router.HandlerFunc("GET", "/do", httpdurableobject.GET_DO)
	//", "Env
	router.HandlerFunc("GET", "/env", httpenv.GET_ENV)

	//Error
	router.HandlerFunc("GET", "/error", errorshandler.GET_ERROR)

	/*
	 * Fetch handler uses http.DefaulServeMux as default, calling this function is
	 * optional, unless a different handler must be used
	 */
	fetch.ServeNonBlock(router)
}
