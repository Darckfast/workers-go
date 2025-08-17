package jsclass

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAwaitOnResolve(t *testing.T) {
	r, e := Await(Promise.Call("resolve", true))

	assert.Nil(t, e)
	assert.True(t, r.Bool())
}

func TestAwaitOnReject(t *testing.T) {
	r, e := Await(Promise.Call("reject", "Error message"))

	assert.NotNil(t, e)
	assert.Equal(t, "failed on promise: Error: Error message", e.Error())
	assert.True(t, r.IsUndefined())
}
