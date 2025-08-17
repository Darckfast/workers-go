package kv

import (
	jsclass "github.com/syumai/workers/internal/class"
)

// Delete deletes key-value pair specified by the key.
//   - if a network error happens, returns error.
func (ns *Namespace) Delete(key string) error {
	p := ns.instance.Call("delete", key)
	_, err := jsclass.Await(p)
	if err != nil {
		return err
	}
	return nil
}
