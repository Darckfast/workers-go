//go:build !js && !wasm

package main

import (
	_ "embed"
	"errors"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
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

func main() {
	start := time.Now()
	log.SetFlags(0)

	exports := argList{}
	entry := flag.String("i", ".", "Root directory of your Go worker")
	out := flag.String("o", "./bin", "Output directory")
	silent := flag.Bool("s", false, "Hide info logs")
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

	if !*silent {
		log.Printf("Output:\n")
		log.Printf("  %s/\n", fo)
		log.Printf("  └─ main.ts\n")
		log.Printf("Took %dμs\n", time.Since(start).Microseconds())
	}
}
