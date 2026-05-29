//go:build js && wasm

package rpchandler

import (
	"context"
	"io"
	"net/http"

	"github.com/Darckfast/workers-go/platform/cloudflare/rpc"
)

func New() {
	/*
	 * RPCStub must be called to instantiate the RPC function, and make
	 * globalThis.cf.<stub-name>() defined on JS global scope
	 */
	rpc.RPCStub("echo", func(c context.Context, args [][]byte) [][]byte {
		return args
	})

	/*
	* RPCStubStream works similar to RPCStub, using ReadableStream and Writers
	* to handle better large payloads, but its slower
	 */
	rpc.RPCStubStream("echoStream", func(c context.Context, w http.ResponseWriter, body io.ReadCloser, args [][]byte) {
		defer body.Close()

		b, _ := io.ReadAll(body)

		w.Write(b)
	})
}
