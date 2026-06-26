//go:build !js && !wasm

package main

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func arrayOf(elt ast.Expr) *ast.ArrayType {
	return &ast.ArrayType{Elt: elt}
}

func ptrTo(x ast.Expr) *ast.StarExpr {
	return &ast.StarExpr{X: x}
}

func selector(pkg, name string) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   ident(pkg),
		Sel: ident(name),
	}
}

func structType(fields ...struct {
	name string
	typ  ast.Expr
}) *ast.StructType {
	var list []*ast.Field
	for _, f := range fields {
		list = append(list, &ast.Field{
			Names: []*ast.Ident{ident(f.name)},
			Type:  f.typ,
		})
	}
	return &ast.StructType{
		Fields: &ast.FieldList{List: list},
	}
}

func TestMapGoTypeToTS_Primitives(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"string", "string"},
		{"bool", "boolean"},
		{"int", "number"},
		{"int8", "number"},
		{"int16", "number"},
		{"int32", "number"},
		{"int64", "number"},
		{"uint", "number"},
		{"uint8", "number"},
		{"uint16", "number"},
		{"uint32", "number"},
		{"uint64", "number"},
		{"float32", "number"},
		{"float64", "number"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.expected, mapGoTypeToTS(ident(tc.input)))
		})
	}
}

func TestMapGoTypeToTS_Array(t *testing.T) {
	assert.Equal(t, "string[]", mapGoTypeToTS(arrayOf(ident("string"))))
	assert.Equal(t, "number[]", mapGoTypeToTS(arrayOf(ident("int"))))
	assert.Equal(t, "boolean[]", mapGoTypeToTS(arrayOf(ident("bool"))))
}

func TestMapGoTypeToTS_NestedArray(t *testing.T) {
	assert.Equal(t, "string[][]", mapGoTypeToTS(arrayOf(arrayOf(ident("string")))))
}

func TestMapGoTypeToTS_Pointer(t *testing.T) {
	assert.Equal(t, "string", mapGoTypeToTS(ptrTo(ident("string"))))
	assert.Equal(t, "number", mapGoTypeToTS(ptrTo(ident("int"))))
}

func TestMapGoTypeToTS_TimeTime(t *testing.T) {
	assert.Equal(t, "string", mapGoTypeToTS(selector("time", "Time")))
}

func TestMapGoTypeToTS_SelectorNonTime(t *testing.T) {
	assert.Equal(t, "MyType", mapGoTypeToTS(selector("pkg", "MyType")))
}

func TestMapGoTypeToTS_InlineStruct(t *testing.T) {
	st := structType(
		struct {
			name string
			typ  ast.Expr
		}{"ID", ident("int")},
		struct {
			name string
			typ  ast.Expr
		}{"Name", ident("string")},
	)
	result := mapGoTypeToTS(st)
	assert.Equal(t, "{ID:number;Name:string;}", result)
}

func TestMapGoTypeToTS_Unknown_ReturnsAny(t *testing.T) {
	unknown := &ast.BasicLit{Kind: token.INT, Value: "42"}
	assert.Equal(t, "any", mapGoTypeToTS(unknown))
}

func TestMapGoStructToTS_DirectStructType(t *testing.T) {
	st := structType(
		struct {
			name string
			typ  ast.Expr
		}{"Foo", ident("string")},
		struct {
			name string
			typ  ast.Expr
		}{"Bar", ident("bool")},
	)
	assert.Equal(t, "{Foo:string;Bar:boolean;}", mapGoStructToTS(st))
}

func TestMapGoStructToTS_EmptyStruct(t *testing.T) {
	st := &ast.StructType{Fields: &ast.FieldList{}}
	assert.Equal(t, "{}", mapGoStructToTS(st))
}

func TestMapGoStructToTS_NestedStruct(t *testing.T) {
	inner := structType(struct {
		name string
		typ  ast.Expr
	}{"X", ident("float64")})

	outer := structType(
		struct {
			name string
			typ  ast.Expr
		}{"Inner", inner},
		struct {
			name string
			typ  ast.Expr
		}{"Label", ident("string")},
	)
	assert.Equal(t, "{Inner:{X:number;};Label:string;}", mapGoStructToTS(outer))
}

func TestMapGoStructToTS_LookupByIdent(t *testing.T) {
	st := structType(struct {
		name string
		typ  ast.Expr
	}{"Age", ident("int")})
	AllStructs.Set("Person", st)

	result := mapGoStructToTS(ident("Person"))
	assert.Equal(t, "{Age:number;}", result)
}

func TestMapGoStructToTS_UnknownIdentPanics(t *testing.T) {
	assert.Panics(t, func() {
		mapGoStructToTS(ident("DoesNotExist"))
	})
}

func TestMapGoStructToTS_ArrayField(t *testing.T) {
	st := structType(struct {
		name string
		typ  ast.Expr
	}{"Tags", arrayOf(ident("string"))})
	assert.Equal(t, "{Tags:string[];}", mapGoStructToTS(st))
}

func TestMapGoStructToTS_PointerField(t *testing.T) {
	st := structType(struct {
		name string
		typ  ast.Expr
	}{"Score", ptrTo(ident("float64"))})
	assert.Equal(t, "{Score:number;}", mapGoStructToTS(st))
}
