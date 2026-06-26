//go:build !js && !wasm

package main

import (
	"go/ast"
	"strings"
)

func mapGoStructToTS(expr ast.Expr) string {
	var structType *ast.StructType

	switch t := expr.(type) {
	case *ast.StructType:
		structType = t
	case *ast.Ident:
		var ok bool
		if structType, ok = AllStructs.Get(t.Name); !ok {
			panic("struct not found: " + t.Name)
		}
	}

	innerTypes := strings.Builder{}
	innerTypes.WriteString("{")

	for _, f := range structType.Fields.List {
		for _, n := range f.Names {
			innerTypes.WriteString(n.Name)
			innerTypes.WriteString(":")
			innerTypes.WriteString(mapGoTypeToTS(f.Type))
			innerTypes.WriteString(";")
		}
	}

	innerTypes.WriteString("}")
	return innerTypes.String()
}

func mapGoTypeToTS(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "string":
			return "string"
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			return "number"
		case "bool":
			return "boolean"
		default:
			return mapGoStructToTS(expr)
		}
	case *ast.ArrayType:
		return mapGoTypeToTS(t.Elt) + "[]"
	case *ast.StarExpr:
		return mapGoTypeToTS(t.X)
	case *ast.SelectorExpr:
		if t.X.(*ast.Ident).Name == "time" && t.Sel.Name == "Time" {
			return "string"
		}
		return t.Sel.Name
	case *ast.StructType:
		return mapGoStructToTS(expr)
	default:
		warn("undefined mapping Go -> JS: %+v", t)
		return "any"
	}
}
