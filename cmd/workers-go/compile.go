//go:build !js && !wasm

package main

import (
	"context"
	_ "embed"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var execCommand = exec.CommandContext

func copyWasmExecJs(ctx context.Context, tiny bool, out []compile) error {
	var wjs *os.File
	if tiny {
		tinyroot, err := execCommand(ctx, "tinygo", "env", "TINYGOROOT").Output()
		if err != nil {
			erro("Error getting TinyGo root path: {Bold}%s", err)
			return err
		}

		in, err := os.ReadFile(filepath.Join(strings.TrimSpace(string(tinyroot)), "targets/wasm_exec.js"))
		if err != nil {
			erro("Error reading wasm_exec.js file: %s", err)
			return err
		}

		fixedIn := strings.Replace(string(in), "global.require = require;", "// Polyfill removed due compatibility issues with Cloudflare Workers", 1)
		wjs, err = os.CreateTemp("", "*.js")
		if err != nil {
			erro("Error creating temporary file: {Bold}%s", err)
			return err
		}
		_, err = wjs.WriteString(fixedIn)
		if err != nil {
			erro("Error writing wasm_exec.js file: {Bold}%s", err)
			return err
		}

		err = wjs.Sync()
		if err != nil {
			erro("Error syncing wasm_exec.js file: {Bold}%s", err)
			return err
		}
		wjs.Seek(0, 0)
	} else {
		goroot, err := execCommand(ctx, "go", "env", "GOROOT").Output()
		if err != nil {
			erro("Error getting go root path: {Bold}%s", err)
			return err
		}

		wjs, err = os.Open(filepath.Join(strings.TrimSpace(string(goroot)), "/lib/wasm/wasm_exec.js"))
		if err != nil {
			erro("Error reading wasm_exec file: {Bold}%s", err)
			return err
		}
	}

	var wout []io.Writer
	var wdirs []string
	for _, c := range out {
		trimpath := filepath.Dir(c.Out)
		if !slices.Contains(wdirs, trimpath) {
			wdirs = append(wdirs, trimpath)

			filename := filepath.Join(trimpath, "wasm_exec.js")
			df, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			defer df.Close()
			wout = append(wout, df)
		}
	}

	mw := io.MultiWriter(wout...)
	_, err := io.Copy(mw, wjs)
	if err != nil {
		erro("Error writing wasm_exec file: {Bold}%s", err)
		return err
	}

	return nil
}

func compileGo(ctx context.Context, a *Args, comp []compile) ([]string, error) {
	outputfiles := []string{}
	var cmd *exec.Cmd

	err := copyWasmExecJs(ctx, a.Tiny, comp)
	if err != nil {
		return nil, err
	}

	if a.Tiny {
		warn("⚠  Using TinyGo might result in some unexpected bugs due compatibility issues ⚠")
	}

	for i, c := range comp {
		if a.Tiny {
			cmd = execCommand(ctx, "tinygo", "build", "-no-debug", "-o", c.Out)
		} else {
			cmd = execCommand(ctx, "go", "build", "-trimpath", "-ldflags", "-s -w -buildid=", "-o", c.Out)
		}

		if c.Tag != "" {
			cmd.Args = append(cmd.Args, "-tags", c.Tag)
		}
		cmd.Dir = c.In
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOOS=js", "GOARCH=wasm")

		out, err := cmd.CombinedOutput()
		if err != nil {
			erro("Error compiling %s: {Bold}%s\n%s", c.Out, err, string(out))
			return nil, err
		}

		if !a.NoCompression {
			err = compress(ctx, c.Out)
			if err != nil {
				return nil, err
			}
		}

		s, _ := os.Stat(c.Out)
		ch := "├─"
		if i == len(comp)-1 {
			ch = "└─"
		}
		rl, _ := filepath.Rel(a.EntryDir, c.Out)
		outputfiles = append(outputfiles, "  "+ch+" "+rl+" ("+fmtBytes(s.Size())+")")
	}

	return outputfiles, nil
}

func compress(ctx context.Context, p string) error {
	_, err := exec.LookPath("wasm-opt")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			warn("{Bold}wasm-opt{Reset} not found - skipping extra compression step")
		} else {
			erro("Error looking for cmd wasm-opt: {Bold}%s", err)
		}
	} else {
		cmd := execCommand(ctx, "wasm-opt", "--all-features", "-Os", p, "-o", p)
		out, err := cmd.CombinedOutput()
		if err != nil {
			erro("Error compressing app.wasm: %s\n%s", err, string(out))
			return err
		}
	}
	return nil
}
