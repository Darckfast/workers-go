//go:build js && wasm

/*
Package tail is the glue code for Cloudflare's Worker Tail handler
*/
package tail

import (
	"syscall/js"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"github.com/mailru/easyjson"
)

func NewEvents(eventsJs js.Value) (*Traces, error) {
	traces := Traces{}

	if !eventsJs.Truthy() {
		return &traces, nil
	}

	str := jsclass.JSON.Stringify(eventsJs)
	err := easyjson.Unmarshal([]byte(str.String()), &traces)

	return &traces, err
}
