package tailhandler

import (
	"encoding/json"
	"fmt"

	"github.com/syumai/workers/cloudflare/kv"
	"github.com/syumai/workers/cloudflare/tail"
	jstail "github.com/syumai/workers/internal/tail"
)

func New() {
	tail.ConsumeNonBlock(func(f *[]jstail.TailEvent) error {
		fmt.Println(f)
		namespace, _ := kv.NewNamespace("TEST_NAMESPACE")
		bjson, _ := json.Marshal(f)
		namespace.PutString("tail:result", string(bjson), nil)

		return nil
	})
}
