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

func args() (*string, *string, *bool, *bool, *arglist, error) {
	exports := arglist{}
	entry := flag.String("i", "", "Root directory of your Go worker")
	out := flag.String("o", "./bin", "Output directory")
	silent := flag.Bool("s", false, "Hide info logs")
	tiny := flag.Bool("tiny", false, "Use tinygo to compile the project")
	flag.Var(&exports, "ex", "Include a exports * from directory - the directory must contain a index.js(ts) file with the desired exports")

	flag.Parse()

	if *entry == "" {
		erro("Root directory is {Bold}required")
		return nil, nil, nil, nil, nil, errors.New("root is required")
	}

	fp, _ := filepath.Abs(*entry)
	if strings.HasSuffix(fp, ".go") {
		fp = filepath.Join(fp, "..")
	}

	fo, _ := filepath.Abs(*out)
	_, err := os.Stat(fp)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		erro("Root directory {Bold}%s{Reset} does not exist", fp)
		os.Exit(1)
	} else if err != nil {
		erro("Error probing directory {Bold}%s{Reset}\n%s", fp, err)
		return nil, nil, nil, nil, nil, err
	}

	return &fp, &fo, silent, tiny, &exports, err
}
