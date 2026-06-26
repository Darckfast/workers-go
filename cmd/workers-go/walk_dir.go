//go:build !js && !wasm

package main

import (
	_ "embed"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type handlerName string

const (
	kindFetch     handlerName = "fetch"
	kindScheduled handlerName = "cron"
	kindQueue     handlerName = "queues"
	kindTail      handlerName = "tail"
	kindEmail     handlerName = "email"
	kindSocket    handlerName = "sockets"
	kindRPC       handlerName = "rpc"
)

var trgFuncs = map[string]bool{
	"ServeNonBlock":        true,
	"ScheduleTaskNonBlock": true,
	"ConsumeNonBlock":      true,
	"Connect":              true,
	"RPCStub":              true,
	"RPCStubStream":        true,
}

type DurableObjectFunc struct {
	Path     string
	FuncName string
	Args     []string
	Rargs    []string
	TSArgs   []string
	TSRarg   []string
}

func scandir(e *string) (*WM, *DOM, error) {
	fset := token.NewFileSet()
	var workersMap = WM{}
	var durableObjects = DOM{}

	filesChan := make(chan string, 100)
	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			for path := range filesChan {
				file, err := parser.ParseFile(fset, path, nil, 0)

				if err != nil {
					erro("[thread %d] Error parsing file: {Bold}%", i, err)
					continue
				}

				a2k := A2K{}
				for _, imp := range file.Imports {
					path := strings.Trim(imp.Path.Value, `"`)
					parts := strings.Split(path, "/")
					seg := handlerName(parts[len(parts)-1])

					switch seg {
					case kindFetch:
						a2k[kindFetch] = struct{}{}
					case kindTail:
						a2k[kindTail] = struct{}{}
					case kindEmail:
						a2k[kindEmail] = struct{}{}
					case kindQueue:
						a2k[kindQueue] = struct{}{}
					case kindScheduled:
						a2k[kindScheduled] = struct{}{}
					case kindSocket:
						a2k[kindSocket] = struct{}{}
					case kindRPC:
						a2k[kindRPC] = struct{}{}
					}
				}

				ast.Inspect(file, func(n ast.Node) bool {
					LookupHandlers(n, &a2k, &workersMap)
					return true
				})

				var fields []string
				ast.Inspect(file, func(n ast.Node) bool {
					fields = append(fields, LookupStructs(n)...)
					return true
				})

				ast.Inspect(file, func(n ast.Node) bool {
					LookupDurableObjects(n, fields, &durableObjects, &path)
					return true
				})
			}
		}(i)
	}

	err := filepath.WalkDir(*e, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() &&
			strings.HasSuffix(d.Name(), ".go") &&
			!strings.HasSuffix(d.Name(), "_test.go") &&
			!strings.HasSuffix(path, "_generated.go") {
			filesChan <- path
		}

		return nil
	})

	if err != nil {
		erro("Error walking on directories %s", err)
		return nil, nil, err
	}

	close(filesChan)
	wg.Wait()

	return &workersMap, &durableObjects, nil
}
