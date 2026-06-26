//go:build !js && !wasm

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInline_Empty(t *testing.T) {
	assert.Equal(t, "", Inline("a", []string{}))
}

func TestInline_Single(t *testing.T) {
	assert.Equal(t, "a0", Inline("a", []string{"x"}))
}

func TestInline_Multiple(t *testing.T) {
	assert.Equal(t, "p0,p1,p2", Inline("p", []string{"x", "y", "z"}))
}

func TestInlineType_Empty(t *testing.T) {
	assert.Equal(t, "", InlineType("a", []string{}, []string{}))
}

func TestInlineType_Single(t *testing.T) {
	assert.Equal(t, "a0:string", InlineType("a", []string{"string"}, []string{"x"}))
}

func TestInlineType_Multiple(t *testing.T) {
	assert.Equal(t, "p0:string,p1:number", InlineType("p", []string{"string", "number"}, []string{"x", "y"}))
}

func TestFuncReturn_Void(t *testing.T) {
	assert.Equal(t, "void", funcReturn([]string{""}))
}

func TestFuncReturn_Single(t *testing.T) {
	assert.Equal(t, "string", funcReturn([]string{"string"}))
}

func TestFuncReturn_Multiple(t *testing.T) {
	assert.Equal(t, "[string,number]", funcReturn([]string{"string", "number"}))
}

func makeDom(t *testing.T, key string, funcs []DurableObjectFunc) *DOM {
	t.Helper()
	dom := &DOM{}
	dom.Set(key, funcs)
	return dom
}

func TestGenDurableObjects_EmptyDOM(t *testing.T) {
	indir := t.TempDir()
	outdir := t.TempDir()
	dom := &DOM{}
	comp, files, err := genDurableObjects(indir, outdir, dom)
	assert.NoError(t, err)
	assert.Empty(t, comp)
	assert.Empty(t, files)
}

func TestGenDurableObjects_GenGoAndTSFiles(t *testing.T) {
	indir := t.TempDir()
	outdir := t.TempDir()
	srcDir := t.TempDir()

	dom := makeDom(t, "MyDO", []DurableObjectFunc{
		{Path: filepath.Join(srcDir, "do.go"), FuncName: "Greet", Args: []string{"string"}, TSArgs: []string{"string"}},
	})

	comp, files, err := genDurableObjects(indir, outdir, dom)
	assert.NoError(t, err)

	_, err = os.Stat(filepath.Join(outdir, "durable_objects"))
	assert.NoError(t, err)

	require.Len(t, comp, 1)
	assert.Equal(t, srcDir, comp[0].In)
	assert.Equal(t, "durableobject", comp[0].Tag)
	assert.Contains(t, comp[0].Out, "MyDO.wasm")

	require.Len(t, files, 1)
	assert.Contains(t, files[0], "mydo_generated.go")

	_, err = os.Stat(filepath.Join(srcDir, "mydo_generated.go"))
	assert.NoError(t, err)

	_, err = os.Stat(filepath.Join(outdir, "durable_objects", "MyDO.ts"))
	assert.NoError(t, err)
}
