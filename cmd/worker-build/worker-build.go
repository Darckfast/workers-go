//go:build !js && !wasm

package main

import (
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
	"strconv"
	"strings"
	"text/template"
	"time"
)

type handlerName string

//go:embed main.ts.tmpl
var mainTsTmpl string

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

func scandir(e string) *map[string]any {
	fset := token.NewFileSet()
	results := make(map[string]any)
	filepath.WalkDir(e, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") && !strings.HasSuffix(d.Name(), "_test.go") {
			file, err := parser.ParseFile(fset, path, nil, 0)

			if err != nil {
				log.Printf("%sError parsing file: %s%s%s\n", Red, Bold, err, Reset)
				return err
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
						arg, ok := call.Args[0].(*ast.BasicLit)
						if ok {
							switch arg.Kind {
							case token.STRING:
								cstr, err := strconv.Unquote(arg.Value)
								if err != nil {
									log.Printf("%sError unquoting value: %s%s%s\n", Red, Bold, err, Reset)
									break
								}

								if results[ident.Name] == nil {
									results[ident.Name] = []string{
										cstr,
									}
								} else {
									results[ident.Name] = append(results[ident.Name].([]string), cstr)
								}
							}
						}
					} else {
						results[ident.Name] = kind
					}
				}

				return true
			})
		}

		return nil
	})

	return &results
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

func Duration(d time.Duration) string {
	type unit struct {
		dur  time.Duration
		name string
	}

	units := []unit{
		{time.Minute, "m"},
		{time.Second, "s"},
		{time.Millisecond, "ms"},
		{time.Microsecond, "µs"},
		{time.Nanosecond, "ns"},
	}

	abs := d
	if abs < 0 {
		abs = -abs
	}

	for _, u := range units {
		if abs >= u.dur {
			return strconv.FormatFloat(float64(d)/float64(u.dur), 'f', 2, 64) + " " + u.name
		}
	}

	return "0 ns"
}
func Bytes(b int64) string {
	const unit = 1024

	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

	v := float64(b)
	i := 0

	for v >= unit && i < len(units)-1 {
		v /= unit
		i++
	}

	return strconv.FormatFloat(v, 'f', 2, 64) + " " + units[i]
}
func main() {
	start := time.Now()
	log.SetFlags(0)

	exports := argList{}
	entry := flag.String("i", ".", "Root directory of your Go worker")
	out := flag.String("o", "./bin", "Output directory")
	silent := flag.Bool("s", false, "Hide info logs")
	tiny := flag.Bool("tiny", false, "Use tinygo to compile the project")
	flag.Var(&exports, "ex", "Include a exports * from directory - the directory must contain a index.js(ts) file with the desired exports")

	flag.Parse()

	fp, _ := filepath.Abs(*entry)
	if strings.HasSuffix(fp, ".go") {
		fp = filepath.Join(fp, "..")
	}

	fo, _ := filepath.Abs(*out)

	_, err := os.Stat(fp)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		log.Printf("%sRoot directory %s %sdoes not exist%s\n", Red, fp, Bold, Reset)
		return
	}

	handlers := scandir(fp)

	if len(exports) > 0 {
		for _, ex := range exports {
			fex, _ := filepath.Abs(ex)
			rp, err := filepath.Rel(fo, fex)

			if err != nil {
				log.Printf("%sError finding relative path: %s%s%s\n", Red, Bold, err, Reset)
				continue
			}

			if (*handlers)["exports"] == nil {
				(*handlers)["exports"] = []string{
					rp,
				}
			} else {
				(*handlers)["exports"] = append((*handlers)["exports"].([]string), rp)
			}
		}
	}
	if !*silent {
		for key, v := range *handlers {
			switch v.(type) {
			case bool:
				log.Printf("* Found: %s%s%s [ok]%s\n", Bold, Green, key, Reset)
			case []string:
				log.Printf("* Found: %s%s%s %s [ok]%s\n", Bold, Green, key, v, Reset)
			}
		}
	}
	if len(*handlers) == 0 {
		log.Printf("* No `workers-go` usage found on: %s\n", fp)
		return
	}

	err = os.MkdirAll(fo, os.ModePerm)
	if err != nil {
		log.Printf("%sError creating output directory: %s%s%s\n", Red, Bold, err, Reset)
		return
	}

	file, err := os.Create(filepath.Join(fo, "main.ts"))
	if err != nil {
		log.Printf("%sError creating main.ts file: %s%s%s\n", Red, Bold, err, Reset)
		return
	}

	tmpl := template.Must(template.New("main.ts").Parse(mainTsTmpl))
	err = tmpl.Execute(file, *handlers)
	if err != nil {
		log.Printf("%sError populating template file: %s%s%s\n", Red, Bold, err, Reset)
		return
	}

	wasmOut := filepath.Join(fo, "app.wasm")
	if *tiny {
		log.Printf("%s%s⚠ [WARN] using tinygo might result in some unexpected bugs due compatibility issues ⚠%s\n", Bold, Yellow, Reset)
		tinyroot, err := exec.Command("tinygo", "env", "TINYGOROOT").Output()
		if err != nil {
			log.Printf("%sError getting tinygo root path: %s%s%s\n", Red, Bold, err, Reset)
			return
		}

		in, err := os.ReadFile(filepath.Join(strings.TrimSpace(string(tinyroot)), "targets/wasm_exec.js"))
		if err != nil {
			log.Printf("%sError reading wasm_exec.js file: %s%s%s\n", Red, Bold, err, Reset)
			return
		}

		// Removes the global.require from the file, otherwise the worker wont startup
		fixedIn := strings.Replace(string(in), "global\\.require = require;", "", 1)
		err = os.WriteFile(filepath.Join(fp, "wasm_exec.js"), []byte(fixedIn), os.ModePerm)
		if err != nil {
			log.Printf("%sError writing wasm_exec.js file: %s%s%s\n", Red, Bold, err, Reset)
			return
		}

		cmd := exec.Command("tinygo", "build", "-no-debug", "-o", wasmOut, fp)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOOS=js", "GOARCH=wasm")

		err = cmd.Run()
		if err != nil {
			log.Printf("%sError compiling app.wasm: %s%s%s\n", Red, Bold, err, Reset)
			return
		}
	} else {
		goroot, err := exec.Command("go", "env", "GOROOT").Output()
		if err != nil {
			log.Printf("%sError getting go root path: %s%s%s\n", Red, Bold, err, Reset)
			return
		}

		file, err := os.Open(filepath.Join(strings.TrimSpace(string(goroot)), "/lib/wasm/wasm_exec.js"))
		if err != nil {
			log.Printf("%sError reading wasm_exec file: %s%s%s\n", Red, Bold, err, Reset)
			return
		}

		defer file.Close()
		dst, err := os.Create(filepath.Join(fo, "wasm_exec.js"))
		if err != nil {
			log.Printf("%sError opening wasm_exec file: %s%s%s\n", Red, Bold, err, Reset)
			return
		}

		defer dst.Close()
		_, err = io.Copy(dst, file)

		if err != nil {
			log.Printf("%sError writing wasm_exec file: %s%s%s\n", Red, Bold, err, Reset)
			return
		}

		err = dst.Sync()
		if err != nil {
			log.Printf("%sError syncing wasm_exec file: %s%s%s\n", Red, Bold, err, Reset)
			return
		}

		cmd := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w -buildid=", "-o", wasmOut, fp)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOOS=js", "GOARCH=wasm")

		err = cmd.Run()
		if err != nil {
			log.Printf("%sError compiling app.wasm: %s%s%s\n", Red, Bold, err, Reset)
			return
		}
	}
	_, err = exec.LookPath("wasm-opt")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			log.Printf("wasm-opt not found - skipping extra compression step\n")
		} else {
			log.Printf("Error looking for cmd wasm-opt: %s\n", err)
		}
	} else {
		cmd := exec.Command("wasm-opt", "--all-features", "-Os", filepath.Join(fo, "app.wasm"), "-o", wasmOut)
		err = cmd.Run()
		if err != nil {
			log.Printf("%sError compressing app.wasm: %s%s%s\n", Red, Bold, err, Reset)
			return
		}
	}

	if !*silent {
		imn, _ := os.Stat(filepath.Join(fo, "main.ts"))
		iaw, _ := os.Stat(filepath.Join(fo, "app.wasm"))
		iwe, _ := os.Stat(filepath.Join(fo, "wasm_exec.js"))

		log.Printf("\nOutput:\n")
		log.Printf("  %s/\n", fo)
		log.Printf("  ├─ main.ts (%s)\n", Bytes(imn.Size()))
		log.Printf("  ├─ app.wasm (%s)\n", Bytes(iaw.Size()))
		log.Printf("  └─ wasm_exec.js (%s)\n", Bytes(iwe.Size()))
		log.Printf("\nTook %s\n", Duration(time.Since(start)))
	}
}
