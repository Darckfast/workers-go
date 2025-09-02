//go:build js && wasm

package r2

import (
	"errors"
	"syscall/js"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	"github.com/mailru/easyjson"
)

func toObjects(v js.Value) (*R2Objects, error) {
	var objects R2Objects
	str := jsclass.JSON.Stringify(v)
	err := easyjson.Unmarshal([]byte(str.String()), &objects)

	if err != nil {
		return nil, err
	}

	objectsVal := v.Get("objects")

	for i := 0; i < len(objects.Objects); i++ {
		obj, err := toObject(objectsVal.Index(i))
		if err != nil {
			return nil, errors.New("error converting to Object: " + err.Error())
		}
		objects.Objects[i] = obj
	}

	return &objects, nil
}
