//go:build !js && !wasm

package main

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

//go:embed templates/main.ts.tmpl
var mainTSTmpl string

//go:embed templates/durable_object.go.tmpl
var durableObjectGoTmpl string

//go:embed templates/durable_object.ts.tmpl
var durableObjectTSTmpl string

func main() {
	start := time.Now()
	ctx := context.Background()

	a, err := args()
	if err != nil {
		os.Exit(1)
	}

	handlers, dos, err := scandir(&a.EntryDir)
	if err != nil {
		os.Exit(1)
	}

	if len(a.Exports) > 0 {
		for _, ex := range a.Exports {
			fex, _ := filepath.Abs(ex)
			rp, err := filepath.Rel(a.OutDir, fex)

			if err != nil {
				erro("Error finding relative path: %s", err)
				continue
			}

			if val, ok := handlers.Get("exports"); ok {
				handlers.Set("exports", append(val, rp))
			} else {
				handlers.Set("exports", []string{rp})
			}
		}
	}

	err = os.MkdirAll(a.OutDir, os.ModePerm)
	if err != nil {
		erro("Error creating output directory: {Bold}%s", err)
		os.Exit(1)
	}

	compileList, outputfiles, err := genDurableObjects(a.EntryDir, a.OutDir, dos)
	if err != nil {
		os.Exit(1)
	}
	if len(handlers.m) == 0 && len(dos.m) == 0 {
		warn("* No `workers-go` handlers usage found on: %s", a.EntryDir)
		os.Exit(0)
	}

	for k := range dos.m {
		if val, ok := handlers.Get("exports"); ok {
			handlers.Set("exports", append(val, "./durable_objects/"+k+".ts"))
		} else {
			handlers.Set("exports", []string{"./durable_objects/" + k + ".ts"})
		}
	}

	maints := filepath.Join(a.OutDir, "main.ts")
	compileList = append(compileList, compile{In: a.EntryDir, Out: filepath.Join(a.OutDir, "app.wasm")})

	file, err := os.Create(maints)
	if err != nil {
		erro("Error creating main.ts file: {Bold}%s", err)
		os.Exit(1)
	}

	tmpl := template.Must(template.New("main.ts").Parse(mainTSTmpl))
	err = tmpl.Execute(file, handlers.m)
	if err != nil {
		erro("Error populating template file: {Bold}%s", err)
		os.Exit(1)
	}

	s, _ := file.Stat()
	rl, _ := filepath.Rel(a.EntryDir, file.Name())
	outputfiles = append(outputfiles, "  ├─ "+rl+" ("+fmtBytes(s.Size())+")")

	o, err := compileGo(ctx, a, compileList)
	if err != nil {
		os.Exit(1)
	}
	outputfiles = append(outputfiles, o...)

	info("\nOutput:")
	for _, s := range outputfiles {
		info(s)
	}
	info("\nTook {Bold}{Green}%s", fmtDuration(time.Since(start)))
}
