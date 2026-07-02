//go:build !js && !wasm

package main

import (
	_ "embed"
	"errors"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type arglist []string

func (a *arglist) String() string {
	return strings.Join(*a, ", ")
}

func (a *arglist) Set(value string) error {
	*a = append(*a, value)
	return nil
}

type Args struct {
	EntryDir      string
	OutDir        string
	Tiny          bool
	Exports       arglist
	NoCompression bool
}

func args() (*Args, error) {
	a := Args{}

	flag.StringVar(&a.EntryDir, "i", "", "Root directory of your Go worker")
	flag.StringVar(&a.OutDir, "o", "./bin", "Output directory")
	flag.BoolVar(&a.Tiny, "tiny", false, "Use tinygo to compile the project")
	flag.BoolVar(&a.NoCompression, "nocompression", false, "Skips wams-opt compression step")
	flag.Var(&a.Exports, "ex", "Include a exports * from directory - the directory must contain a index.js(ts) file with the desired exports")

	flag.Parse()

	if a.EntryDir == "" {
		erro("Root directory is {Bold}required")
		return nil, errors.New("root is required")
	}

	fp, _ := filepath.Abs(a.EntryDir)
	if strings.HasSuffix(fp, ".go") {
		fp = filepath.Join(fp, "..")
	}

	fo, _ := filepath.Abs(a.OutDir)
	_, err := os.Stat(fp)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		erro("Root directory {Bold}%s{Reset} does not exist", fp)
		os.Exit(1)
	} else if err != nil {
		erro("Error probing directory {Bold}%s{Reset}\n%s", fp, err)
		return nil, err
	}

	a.EntryDir = fp
	a.OutDir = fo

	return &a, nil
}
