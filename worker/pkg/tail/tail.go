//go:build js && wasm

package tailhandler

import (
	"encoding/json"

	"github.com/Darckfast/workers-go/cloudflare/kv"
	"github.com/Darckfast/workers-go/cloudflare/tail"
	jstail "github.com/Darckfast/workers-go/internal/tail"
)

func New() {
	/*
	 * This function must be called to instantiate the tail handler consumer, and
	 * make globalThis.cf.tail() defined in the global scope
	 */
	tail.ConsumeNonBlock(func(f *[]jstail.TailEvent) error {
		namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
		bjson, _ := json.Marshal(f)
		namespace.PutString("tail:result", string(bjson), nil)

		return nil
	})
}
