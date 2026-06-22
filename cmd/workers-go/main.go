package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

type handlerName string

//go:embed main.ts.tmpl
var mainTSTmpl string

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

type unit struct {
	name string
	dur  time.Duration
}

var byteUnits = []string{"B", "KiB", "MiB"}

const byteSize = 1 << 10

var durationUnits = []unit{
	{"m", time.Minute},
	{"s", time.Second},
	{"ms", time.Millisecond},
	{"µs", time.Microsecond},
	{"ns", time.Nanosecond},
}

const (
	Bold   = "\033[1m"
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

type argList []string

func (a *argList) String() string {
	return strings.Join(*a, ", ")
}

func (a *argList) Set(value string) error {
	*a = append(*a, value)
	return nil
}

type SafeMap[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func (sm *SafeMap[K, V]) Set(key K, val V) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = val
}

func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	val, ok := sm.m[key]

	return val, ok
}

func FmtDuration(d time.Duration) string {
	for _, u := range durationUnits {
		if d >= u.dur {
			return strconv.FormatFloat(float64(d)/float64(u.dur), 'f', 2, 64) + " " + u.name
		}
	}

	return "0 ns"
}

func FmtBytes(b int64) string {
	v := float64(b)
	i := 0

	for v >= byteSize && i < len(byteUnits)-1 {
		v /= byteSize
		i++
	}

	return strconv.FormatFloat(v, 'f', 2, 64) + " " + byteUnits[i]
}

func scandir(e *string) *SafeMap[string, any] {
	fset := token.NewFileSet()
	safem := SafeMap[string, any]{m: make(map[string]any)}

	filesChan := make(chan string, 100)
	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			for path := range filesChan {
				file, err := parser.ParseFile(fset, path, nil, 0)

				if err != nil {
					log.Printf("%s[thread %d] Error parsing file: %s%s%s\n", Red, i, Bold, err, Reset)
					continue
				}

				aliasToKind := make(map[handlerName]bool)
				for _, imp := range file.Imports {
					path := strings.Trim(imp.Path.Value, `"`)
					parts := strings.Split(path, "/")
					seg := handlerName(parts[len(parts)-1])

					switch seg {
					case kindFetch:
						aliasToKind[kindFetch] = true
					case kindTail:
						aliasToKind[kindTail] = true
					case kindEmail:
						aliasToKind[kindEmail] = true
					case kindQueue:
						aliasToKind[kindQueue] = true
					case kindScheduled:
						aliasToKind[kindScheduled] = true
					case kindSocket:
						aliasToKind[kindSocket] = true
					case kindRPC:
						aliasToKind[kindRPC] = true
					}
				}

				ast.Inspect(file, func(n ast.Node) bool {
					call, ok := n.(*ast.CallExpr)

					if !ok {
						return true
					}

					sel, ok := call.Fun.(*ast.SelectorExpr)

					if !ok {
						return true
					}

					ident, ok := sel.X.(*ast.Ident)
					if !ok {
						return true
					}

					if !trgFuncs[sel.Sel.Name] {
						return true
					}

					if kind, ok := aliasToKind[handlerName(ident.Name)]; ok {
						if ident.Name == "rpc" {
							if arg, ok := call.Args[0].(*ast.BasicLit); ok && arg.Kind == token.STRING {
								cstr, err := strconv.Unquote(arg.Value)
								if err != nil {
									log.Printf("%s[thread %d] Error unquoting value: %s%s%s\n", Red, i, Bold, err, Reset)
									return true
								}

								if val, ok := safem.Get(ident.Name); ok {
									safem.Set(ident.Name, append(val.([]string), cstr))
								} else {
									safem.Set(ident.Name, []string{cstr})
								}
							}
						} else {
							safem.Set(ident.Name, kind)
						}
					}

					return true
				})
			}
		}(i)
	}

	err := filepath.WalkDir(*e, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") && !strings.HasSuffix(d.Name(), "_test.go") {
			filesChan <- path
		}

		return nil
	})

	if err != nil {
		log.Printf("%s%sError walking on directories %s%s\n", Red, Bold, err, Reset)
		os.Exit(1)
	}

	close(filesChan)
	wg.Wait()

	return &safem
}

func parseAndValidateArgs() (*string, *string, *bool, *bool, *argList) {
	exports := argList{}
	entry := flag.String("i", "", "Root directory of your Go worker")
	out := flag.String("o", "./bin", "Output directory")
	silent := flag.Bool("s", false, "Hide info logs")
	tiny := flag.Bool("tiny", false, "Use tinygo to compile the project")
	flag.Var(&exports, "ex", "Include a exports * from directory - the directory must contain a index.js(ts) file with the desired exports")

	flag.Parse()

	if *entry == "" {
		log.Printf("%s%sRoot directory is required%s\n", Red, Bold, Reset)
		os.Exit(1)
	}

	fp, _ := filepath.Abs(*entry)
	if strings.HasSuffix(fp, ".go") {
		fp = filepath.Join(fp, "..")
	}

	fo, _ := filepath.Abs(*out)
	_, err := os.Stat(fp)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		log.Printf("%sRoot directory %s %sdoes not exist%s\n", Red, fp, Bold, Reset)
		os.Exit(1)
	}

	return &fp, &fo, silent, tiny, &exports
}

func main() {
	start := time.Now()
	ctx := context.Background()
	log.SetFlags(0)

	fp, fo, silent, tiny, exports := parseAndValidateArgs()
	handlers := scandir(fp)

	if len(*exports) > 0 {
		for _, ex := range *exports {
			fex, _ := filepath.Abs(ex)
			rp, err := filepath.Rel(*fo, fex)

			if err != nil {
				log.Printf("%sError finding relative path: %s%s%s\n", Red, Bold, err, Reset)
				continue
			}

			if val, ok := handlers.Get("exports"); ok {
				handlers.Set("exports", append(val.([]string), rp))
			} else {
				handlers.Set("exports", []string{rp})
			}
		}
	}

	if !*silent {
		for key, v := range handlers.m {
			switch v.(type) {
			case bool:
				log.Printf("* Found: %s%s%s [ok]%s\n", Bold, Green, key, Reset)
			case []string:
				log.Printf("* Found: %s%s%s %s [ok]%s\n", Bold, Green, key, v, Reset)
			}
		}
	}

	if len(handlers.m) == 0 {
		log.Printf("* No `workers-go` usage found on: %s\n", *fp)
		os.Exit(0)
	}

	err := os.MkdirAll(*fo, os.ModePerm)
	if err != nil {
		log.Printf("%sError creating output directory: %s%s%s\n", Red, Bold, err, Reset)
		os.Exit(1)
	}

	maints := filepath.Join(*fo, "main.ts")
	wasmexecjs := filepath.Join(*fo, "wasm_exec.js")
	appwasm := filepath.Join(*fo, "app.wasm")

	file, err := os.Create(maints)
	if err != nil {
		log.Printf("%sError creating main.ts file: %s%s%s\n", Red, Bold, err, Reset)
		os.Exit(1)
	}

	tmpl := template.Must(template.New("main.ts").Parse(mainTSTmpl))
	err = tmpl.Execute(file, handlers.m)
	if err != nil {
		log.Printf("%sError populating template file: %s%s%s\n", Red, Bold, err, Reset)
		os.Exit(1)
	}

	var cmd *exec.Cmd
	switch *tiny {
	case true:
		log.Printf("%s%s⚠  Using tinygo might result in some unexpected bugs due compatibility issues ⚠%s\n", Bold, Yellow, Reset)
		//nolint:govet
		tinyroot, err := exec.CommandContext(ctx, "tinygo", "env", "TINYGOROOT").Output()
		if err != nil {
			log.Printf("%sError getting tinygo root path: %s%s%s\n", Red, Bold, err, Reset)
			os.Exit(1)
		}

		in, err := os.ReadFile(filepath.Join(strings.TrimSpace(string(tinyroot)), "targets/wasm_exec.js"))
		if err != nil {
			log.Printf("%sError reading wasm_exec.js file: %s%s%s\n", Red, Bold, err, Reset)
			os.Exit(1)
		}

		// Removes the global.require from the file, otherwise the worker wont startup
		fixedIn := strings.Replace(string(in), "global.require = require;", "// Polyfill removed due compatibility issues with Cloudflare Workers", 1)
		err = os.WriteFile(wasmexecjs, []byte(fixedIn), os.ModePerm)
		if err != nil {
			log.Printf("%sError writing wasm_exec.js file: %s%s%s\n", Red, Bold, err, Reset)
			os.Exit(1)
		}

		cmd = exec.CommandContext(ctx, "tinygo", "build", "-no-debug", "-o", appwasm, *fp)
	case false:
		//nolint:govet
		goroot, err := exec.CommandContext(ctx, "go", "env", "GOROOT").Output()
		if err != nil {
			log.Printf("%sError getting go root path: %s%s%s\n", Red, Bold, err, Reset)
			os.Exit(1)
		}

		file, err := os.Open(filepath.Join(strings.TrimSpace(string(goroot)), "/lib/wasm/wasm_exec.js"))
		if err != nil {
			log.Printf("%sError reading wasm_exec file: %s%s%s\n", Red, Bold, err, Reset)
			os.Exit(1)
		}

		dst, err := os.Create(wasmexecjs)
		if err != nil {
			log.Printf("%sError opening wasm_exec file: %s%s%s\n", Red, Bold, err, Reset)
			_ = file.Close()
			os.Exit(1)
		}

		_, err = io.Copy(dst, file)

		if err != nil {
			log.Printf("%sError writing wasm_exec file: %s%s%s\n", Red, Bold, err, Reset)
			os.Exit(1)
		}

		err = dst.Sync()
		if err != nil {
			log.Printf("%sError syncing wasm_exec file: %s%s%s\n", Red, Bold, err, Reset)
			_ = dst.Close()
			os.Exit(1)
		}

		_ = dst.Close()
		_ = file.Close()

		cmd = exec.CommandContext(ctx, "go", "build", "-trimpath", "-ldflags", "-s -w -buildid=", "-o", appwasm)
	}

	cmd.Dir = *fp
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS=js", "GOARCH=wasm")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%sError compiling app.wasm: %s%s%s%s\n", Red, Bold, err, string(out), Reset)
		os.Exit(1)
	}

	_, err = exec.LookPath("wasm-opt")

	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			log.Printf("wasm-opt not found - skipping extra compression step\n")
		} else {
			log.Printf("Error looking for cmd wasm-opt: %s\n", err)
		}
	} else {
		cmd := exec.CommandContext(ctx, "wasm-opt", "--all-features", "-Os", appwasm, "-o", appwasm)
		out, err = cmd.CombinedOutput()
		if err != nil {
			log.Printf("%sError compressing app.wasm: %s%s%s%s\n", Red, Bold, err, string(out), Reset)
			os.Exit(1)
		}
	}

	if !*silent {
		imn, _ := os.Stat(maints)
		iaw, _ := os.Stat(appwasm)
		iwe, _ := os.Stat(wasmexecjs)

		log.Printf("\nOutput:\n")
		log.Printf("  %s/\n", *fo)
		log.Printf("  ├─ main.ts (%s)\n", FmtBytes(imn.Size()))
		log.Printf("  ├─ app.wasm (%s)\n", FmtBytes(iaw.Size()))
		log.Printf("  └─ wasm_exec.js (%s)\n", FmtBytes(iwe.Size()))
		log.Printf("\nTook %s\n", FmtDuration(time.Since(start)))
	}
}
