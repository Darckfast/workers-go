//go:build js && wasm

package tailhandler

import (
	"encoding/json"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/kv"
	"codeberg.org/darckfast/workers-go/platform/cloudflare/tail"
)

func New() {
	/*
	 * ConsumeNonBlock must be called to instantiate the tail handler consumer, and
	 * make globalThis.cf.tail() defined in the global scope
	 */
	tail.ConsumeNonBlock(func(f *tail.Traces) error {
		namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
		bjson, _ := json.Marshal(f)
		return namespace.Put("tail:result", string(bjson), nil)
	})
}
