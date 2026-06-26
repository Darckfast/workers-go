//go:build !js && !wasm

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeMapGet(t *testing.T) {
	var m SafeMap[string, []string]

	v, ok := m.Get("nil")

	assert.False(t, ok)
	assert.Empty(t, v)
}

func TestSafeMapSet(t *testing.T) {
	var m SafeMap[string, []string]

	m.Set("nil", []string{"-1"})
	v, ok := m.Get("nil")

	assert.True(t, ok)
	assert.Equal(t, []string{"-1"}, v)
}
