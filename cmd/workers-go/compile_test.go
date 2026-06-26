//go:build !js && !wasm

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fakeCmd(stdout string, err error) func(context.Context, string, ...string) *exec.Cmd {
	return func(ctx context.Context, name string, args ...string) *exec.Cmd {
		if err != nil {
			cmd := exec.CommandContext(ctx, "/bin/false")
			return cmd
		}
		cmd := exec.CommandContext(ctx, "/bin/echo", "-n", stdout)
		return cmd
	}
}

func goroot(t *testing.T) string {
	t.Helper()
	dir, _ := os.MkdirTemp("", "*")
	os.MkdirAll(dir+"/lib/wasm/", os.ModePerm)
	os.MkdirAll(dir+"/target/", os.ModePerm)
	f, _ := os.Create(dir + "/lib/wasm/wasm_exec.js")
	f.WriteString("empty")
	f, _ = os.Create(dir + "/target/wasm_exec.js")
	f.WriteString("empty")

	return dir
}

func TestCopyWasmExecJs_GoEnvFails(t *testing.T) {
	old := execCommand
	execCommand = fakeCmd("", fmt.Errorf("go not found"))
	defer func() { execCommand = old }()

	outdir := t.TempDir()
	out := []compile{{In: outdir, Out: filepath.Join(outdir, "app.wasm")}}

	err := copyWasmExecJs(context.Background(), false, out)
	assert.Error(t, err)
}

func TestCopyWasmExecJs_EmptyOut(t *testing.T) {
	old := execCommand
	execCommand = fakeCmd(goroot(t), nil)
	defer func() { execCommand = old }()

	err := copyWasmExecJs(context.Background(), false, []compile{})
	assert.NoError(t, err)
}

func TestCopyWasmExecJs_DeduplicatesDirs(t *testing.T) {
	old := execCommand
	execCommand = fakeCmd(goroot(t), nil)
	defer func() { execCommand = old }()

	outdir := t.TempDir()
	out := []compile{
		{In: outdir, Out: filepath.Join(outdir, "a.wasm")},
		{In: outdir, Out: filepath.Join(outdir, "b.wasm")},
	}

	err := copyWasmExecJs(context.Background(), false, out)
	assert.NoError(t, err)

	_, statErr := os.Stat(filepath.Join(outdir, "wasm_exec.js"))
	assert.NoError(t, statErr)
}

func TestCopyWasmExecJs_MultipleOutputDirs(t *testing.T) {
	old := execCommand
	execCommand = fakeCmd(goroot(t), nil)
	defer func() { execCommand = old }()

	dir1 := t.TempDir()
	dir2 := t.TempDir()
	out := []compile{
		{In: dir1, Out: filepath.Join(dir1, "a.wasm")},
		{In: dir2, Out: filepath.Join(dir2, "b.wasm")},
	}

	err := copyWasmExecJs(context.Background(), false, out)
	assert.NoError(t, err)

	_, err1 := os.Stat(filepath.Join(dir1, "wasm_exec.js"))
	_, err2 := os.Stat(filepath.Join(dir2, "wasm_exec.js"))
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func setupCompileGo(t *testing.T, succeed bool) (srcDir, outDir string, restore func()) {
	t.Helper()
	srcDir = t.TempDir()
	outDir = t.TempDir()

	wasmOut := filepath.Join(outDir, "app.wasm")
	if succeed {
		assert.NoError(t, os.WriteFile(wasmOut, []byte("fake"), 0644))
	}

	old := execCommand
	if succeed {
		execCommand = fakeCmd(goroot(t), nil)
	} else {
		execCommand = fakeCmd("", fmt.Errorf("compile failed"))
	}
	restore = func() { execCommand = old }
	return
}

func TestCompileGo_OutputEntries(t *testing.T) {
	srcDir, outDir, restore := setupCompileGo(t, true)
	defer restore()

	assert.NoError(t, os.WriteFile(filepath.Join(outDir, "b.wasm"), []byte("fake"), 0644))

	comps := []compile{
		{In: srcDir, Out: filepath.Join(outDir, "app.wasm")},
		{In: srcDir, Out: filepath.Join(outDir, "b.wasm")},
	}

	entries, err := compileGo(context.Background(), outDir, outDir, comps, false)
	assert.NoError(t, err)

	assert.Len(t, entries, 2)
}

func TestCompileGo_CompileFailureReturnsError(t *testing.T) {
	srcDir, outDir, restore := setupCompileGo(t, false)
	defer restore()

	comps := []compile{{In: srcDir, Out: filepath.Join(outDir, "app.wasm")}}

	_, err := compileGo(context.Background(), outDir, outDir, comps, false)
	assert.Error(t, err)
}
