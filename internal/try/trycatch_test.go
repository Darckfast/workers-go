//go:build js && wasm

package jstry

import (
	"syscall/js"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTryCatchShouldReturnErrAsVal(t *testing.T) {
	j := js.Global().Get("JSON")
	v, e := TryCatch(j, "parse", `{test:1}`)

	assert.Error(t, e, "JavaScript error: Expected property name or '}' in JSON at position 1 (line 1 column 2)")
	assert.Equal(t, js.Undefined(), v)
}

func TestTryCatchShouldReturnVal(t *testing.T) {
	j := js.Global().Get("JSON")
	v, e := TryCatch(j, "parse", `{"test":1}`)

	assert.Nil(t, e)
	assert.Equal(t, v.Get("test").Int(), 1)
}
