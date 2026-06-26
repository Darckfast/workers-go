//go:build !js && !wasm

package main

import (
	"go/ast"
	"go/token"
	"slices"
	"strconv"
	"strings"
)

const doPackage = "durableobjects.DurableObject"

func LookupHandlers(n ast.Node, aliasToKind *A2K, workerMap *WM) {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	if !trgFuncs[sel.Sel.Name] {
		return
	}

	if _, ok := (*aliasToKind)[handlerName(ident.Name)]; ok {
		if ident.Name == "rpc" {
			if arg, ok := call.Args[0].(*ast.BasicLit); ok && arg.Kind == token.STRING {
				cstr, err := strconv.Unquote(arg.Value)
				if err != nil {
					erro("Error unquoting value: %s", err)
					return
				}

				v, _ := workerMap.Get(ident.Name)
				workerMap.Set(ident.Name, append(v, cstr))
			}
		} else {
			workerMap.Set(ident.Name, []string{ident.Name})
		}
	}
}

func LookupStructs(n ast.Node) []string {
	embStrct := []string{}
	typeSpec, ok := n.(*ast.TypeSpec)
	if !ok {
		return embStrct
	}

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return embStrct
	}

	AllStructs.Set(typeSpec.Name.Name, structType)
	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			if emb := getTypeName(field.Type); emb == doPackage {
				embStrct = append(embStrct, typeSpec.Name.Name)
			}
		}
	}

	return embStrct
}

func LookupDurableObjects(n ast.Node, fields []string, durableObject *DOM, path *string) {
	funcDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return
	}

	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return
	}

	recvType := funcDecl.Recv.List[0].Type

	t := getTypeName(recvType)

	var iargs []string
	var rargs []string
	var tsi []string
	var tsr []string
	if slices.Contains(fields, t) && !strings.HasPrefix(funcDecl.Name.Name, "_") {
		for _, p := range funcDecl.Type.Params.List {
			argtype := getTypeName(p.Type)

			if argtype == "context.Context" {
				continue
			}

			inargName := getTypeName(p.Type)
			iargs = append(iargs, inargName)
			tsi = append(tsi, mapGoTypeToTS(p.Type))
		}

		if funcDecl.Type.Results != nil {
			for _, p := range funcDecl.Type.Results.List {
				outargName := getTypeName(p.Type)
				rargs = append(rargs, outargName)
				tsr = append(tsr, mapGoTypeToTS(p.Type))
			}
		}

		v, _ := durableObject.Get(t)
		durableObject.Set(t, append(v, DurableObjectFunc{
			Path:     *path,
			FuncName: funcDecl.Name.Name,
			Args:     iargs,
			Rargs:    rargs,
			TSArgs:   tsi,
			TSRarg:   tsr,
		}))
	}
}

func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return getTypeName(t.X)
	case *ast.SelectorExpr:
		if pkg, ok := t.X.(*ast.Ident); ok {
			return pkg.Name + "." + t.Sel.Name
		}
		return ""
	default:
		return ""
	}
}
