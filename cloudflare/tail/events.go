//go:build js && wasm

package tail

import (
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
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
