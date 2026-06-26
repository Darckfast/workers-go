//go:build !js && !wasm

package main

import (
	_ "embed"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type compile struct {
	In  string
	Out string
	Tag string
}

func InlineType(prefix string, ts, args []string) string {
	var str strings.Builder
	for i := range args {
		if i != 0 {
			str.WriteString(",")
		}
		str.WriteString(prefix)
		str.WriteString(strconv.Itoa(i))
		str.WriteString(":")
		str.WriteString(ts[i])
	}

	return str.String()
}

func Inline(prefix string, args []string) string {
	var str strings.Builder
	for i := range args {
		if i != 0 {
			str.WriteString(",")
		}
		str.WriteString(prefix)
		str.WriteString(strconv.Itoa(i))
	}

	return str.String()
}

func funcReturn(args []string) string {
	for i, v := range args {
		if v == "" {
			args[i] = "void"
		}
	}

	if len(args) == 1 {
		return args[0]
	}

	return `[` + strings.Join(args, ",") + `]`
}

var tmplFuncMap = template.FuncMap{
	"join":        strings.Join,
	"inline":      Inline,
	"inlinetypes": InlineType,
	"typereturn":  funcReturn,
}

func genDurableObjects(indir, outdir string, dos *DOM) ([]compile, []string, error) {
	if len(dos.m) == 0 {
		info("no durable objects found")
		return []compile{}, []string{}, nil
	}

	dospath := filepath.Join(outdir, "durable_objects")
	err := os.MkdirAll(dospath, os.ModePerm)
	if err != nil {
		erro("Error creating durable objects directory: {Bold}%s", err)
		return nil, nil, err
	}

	outputfiles := []string{}
	markedForComp := []compile{}

	for key, do := range dos.m {
		filename := strings.ToLower(key) + "_generated.go"
		info("[durable object] {Bold}{Green}%s", filename)

		dir := filepath.Dir(do[0].Path)
		dofile := filepath.Join(dir, filename)

		file, err := os.Create(dofile)
		if err != nil {
			erro("Error creating %s file: {Bold}%s", dofile, err)
			return nil, nil, err
		}

		tmpl := template.Must(template.New("durable_object.go").Funcs(tmplFuncMap).Parse(durableObjectGoTmpl))
		err = tmpl.Execute(file, map[string]any{key: do})
		if err != nil {
			erro("Error populating template file: {Bold}%s", err)
			return nil, nil, err
		}

		s, _ := file.Stat()
		file.Close()
		rl, _ := filepath.Rel(indir, dofile)
		outputfiles = append(outputfiles, "  ├─ "+rl+" ("+fmtBytes(s.Size())+")")
		markedForComp = append(markedForComp, compile{In: dir, Out: filepath.Join(dospath, key+".wasm"), Tag: "durableobject"})

		dots := filepath.Join(dospath, key+".ts")
		file, err = os.Create(dots)
		if err != nil {
			erro("Error creating %s file: {Bold}%s", err)
			return nil, nil, err
		}

		tmpl = template.Must(template.New(key + ".ts").Funcs(tmplFuncMap).Parse(durableObjectTSTmpl))
		err = tmpl.Execute(file, map[string]any{key: do})
		file.Close()
	}

	return markedForComp, outputfiles, nil
}
