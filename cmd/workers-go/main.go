//go:build !js && !wasm

package main

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"strings"
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

	indir, outdir, silent, tiny, exports, err := args()
	if err != nil {
		os.Exit(1)
	}

	handlers, dos, err := scandir(indir)
	if err != nil {
		os.Exit(1)
	}

	if len(*exports) > 0 {
		for _, ex := range *exports {
			fex, _ := filepath.Abs(ex)
			rp, err := filepath.Rel(*outdir, fex)

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

	if !*silent {
		for _, v := range handlers.m {
			info("[worker] {Bold}{Green}%s", strings.Join(v, ", "))
		}
	}

	err = os.MkdirAll(*outdir, os.ModePerm)
	if err != nil {
		erro("Error creating output directory: {Bold}%s", err)
		os.Exit(1)
	}

	compileList, outputfiles, err := genDurableObjects(*indir, *outdir, dos)
	if err != nil {
		os.Exit(1)
	}
	if len(handlers.m) == 0 && len(dos.m) == 0 {
		warn("* No `workers-go` handlers usage found on: %s", *indir)
		os.Exit(0)
	}

	for k := range dos.m {
		if val, ok := handlers.Get("exports"); ok {
			handlers.Set("exports", append(val, "./durable_objects/"+k+".ts"))
		} else {
			handlers.Set("exports", []string{"./durable_objects/" + k + ".ts"})
		}
	}

	maints := filepath.Join(*outdir, "main.ts")
	compileList = append(compileList, compile{In: *indir, Out: filepath.Join(*outdir, "app.wasm")})

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
	rl, _ := filepath.Rel(*indir, file.Name())
	outputfiles = append(outputfiles, "  ├─ "+rl+" ("+fmtBytes(s.Size())+")")

	o, err := compileGo(ctx, *indir, *outdir, compileList, *tiny)
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
