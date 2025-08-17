package r2

import (
	"errors"
	"syscall/js"

	jsconv "github.com/syumai/workers/internal/conv"
)

// Objects represents Cloudflare R2 objects.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1121
type Objects struct {
	Objects   []*Object
	Truncated bool
	// Cursor indicates next cursor of Objects.
	//   - This becomes empty string if cursor doesn't exist.
	Cursor            string
	DelimitedPrefixes []string
}

// toObjects converts JavaScript side's Objects to *Objects.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1121
func toObjects(v js.Value) (*Objects, error) {
	objectsVal := v.Get("objects")
	objects := make([]*Object, objectsVal.Length())
	for i := 0; i < len(objects); i++ {
		obj, err := toObject(objectsVal.Index(i))
		if err != nil {
			return nil, errors.New("error converting to Object: " + err.Error())
		}
		objects[i] = obj
	}
	prefixesVal := v.Get("delimitedPrefixes")
	prefixes := make([]string, prefixesVal.Length())
	for i := 0; i < len(prefixes); i++ {
		prefixes[i] = prefixesVal.Index(i).String()
	}
	return &Objects{
		Objects:           objects,
		Truncated:         v.Get("truncated").Bool(),
		Cursor:            jsconv.MaybeString(v.Get("cursor")),
		DelimitedPrefixes: prefixes,
	}, nil
}
