//go:build js && wasm

package rpchandler

import (
	"context"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/rpc"
)

func New() {
	/*
	 * RPCStub must be called to instantiate the RPC function, and make
	 * globalThis.cf.<stub-name>() defined on JS global scope
	 */
	rpc.RPCStub("echo", func(c context.Context, args [][]byte) [][]byte {
		return args
	})
}
