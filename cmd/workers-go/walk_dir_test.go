//go:build !js && !wasm

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0644))
}

func TestScandir_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	wm, dom, err := scandir(&dir)
	assert.NoError(t, err)
	assert.NotNil(t, wm)
	assert.NotNil(t, dom)
}

func TestScandir_IgnoresTestFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "handler_test.go", `package main
import "fetch"
func init() { fetch.ServeNonBlock() }
`)
	wm, _, err := scandir(&dir)
	assert.NoError(t, err)
	assert.Empty(t, wm.m)
}

func TestScandir_IgnoresGeneratedFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "foo_generated.go", `package main
import "fetch"
func init() { fetch.ServeNonBlock() }
`)
	wm, _, err := scandir(&dir)
	assert.NoError(t, err)
	assert.Empty(t, wm.m)
}

func TestScandir_DetectsFetchWorker(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "worker.go", `package main
import "fetch"
func init() { fetch.ServeNonBlock() }
`)
	wm, _, err := scandir(&dir)
	v, ok := wm.Get("fetch")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []string{"fetch"}, v)
}

func TestScandir_DetectsQueueWorker(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "worker.go", `package main
import "queues"
func init() { queues.ConsumeNonBlock() }
`)
	wm, _, err := scandir(&dir)
	v, ok := wm.Get("queues")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []string{"queues"}, v)
}

func TestScandir_DetectsRPCMethods(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "worker.go", `package main
import "rpc"
func init() {
	rpc.RPCStub("MethodA")
	rpc.RPCStub("MethodB")
}
`)
	wm, _, err := scandir(&dir)
	v, ok := wm.Get("rpc")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.ElementsMatch(t, []string{"MethodA", "MethodB"}, v)
}

func TestScandir_DetectsDurableObject(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "do.go", `package main
import "durableobjects"
type MyDO struct {
	durableobjects.DurableObject
}
func (d *MyDO) Greet(name string) string { return name }
`)
	_, dom, err := scandir(&dir)
	v, ok := dom.Get("MyDO")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Len(t, v, 1)
	assert.Equal(t, "Greet", v[0].FuncName)
}

func TestScandir_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "fetch_worker.go", `package main
import "fetch"
func init() { fetch.ServeNonBlock() }
`)
	writeFile(t, dir, "queue_worker.go", `package main
import "queues"
func init() { queues.ConsumeNonBlock() }
`)
	wm, _, err := scandir(&dir)
	_, hasFetch := wm.Get("fetch")
	_, hasQueue := wm.Get("queues")
	assert.NoError(t, err)
	assert.True(t, hasFetch)
	assert.True(t, hasQueue)
}
