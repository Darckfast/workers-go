//go:build !js && !wasm

package main

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeCallExpr(receiver, method string, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: receiver},
			Sel: &ast.Ident{Name: method},
		},
		Args: args,
	}
}

func TestGetTypeName(t *testing.T) {
	cases := []struct {
		name     string
		input    ast.Expr
		expected string
	}{
		{"struct", &ast.Ident{Name: "Foo"}, "Foo"},
		{"pointer to struct", &ast.StarExpr{X: &ast.Ident{Name: "Foo"}}, "Foo"},
		{"selector", &ast.SelectorExpr{
			X:   &ast.Ident{Name: "pkg"},
			Sel: &ast.Ident{Name: "Bar"},
		}, "pkg.Bar"},
		{"channel", &ast.ChanType{Dir: ast.RECV, Value: ident("int")}, ""},
		{"array", &ast.ArrayType{Elt: &ast.Ident{Name: "string"}}, ""},    //TODO: should return
		{"map", &ast.MapType{Key: ident("key"), Value: ident("int")}, ""}, //TODO: should return
		{"func", &ast.FuncType{}, ""},                                     //TODO: return unknown
		{"interface", &ast.InterfaceType{}, ""},                           //TODO: return unknown
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, getTypeName(tc.input))
		})
	}
}

func TestLookupHandlers_NotACall(t *testing.T) {
	a2k := &A2K{}
	wm := &WM{}
	LookupHandlers(&ast.Ident{}, a2k, wm)
	assert.Empty(t, wm.m)
}

func TestLookupHandlers_RpcWithStringArg(t *testing.T) {
	a2k := &A2K{"rpc": struct{}{}}
	wm := &WM{}
	node := makeCallExpr("rpc", "RPCStub",
		&ast.BasicLit{Kind: token.STRING, Value: `"myMethod"`},
	)
	LookupHandlers(node, a2k, wm)
	v, ok := wm.Get("rpc")
	assert.True(t, ok)
	assert.Equal(t, []string{"myMethod"}, v)
}

func TestLookupHandlers(t *testing.T) {
	cases := []struct {
		name     string
		fName    string
		ok       bool
		expected any
	}{
		{"tail", "ConsumeNonBlock", true, []string{"tail"}},
		{"email", "ConsumeNonBlock", true, []string{"email"}},
		{"cron", "ScheduleTaskNonBlock", true, []string{"cron"}},
		{"queues", "ConsumeNonBlock", true, []string{"queues"}},
		{"fetch", "ServeNonBlock", true, []string{"fetch"}},
		{"worker", "UnknowMethod", false, []string(nil)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a2k := &A2K{handlerName(tc.name): struct{}{}}
			wm := &WM{}
			node := makeCallExpr(tc.name, tc.fName)
			LookupHandlers(node, a2k, wm)
			v, ok := wm.Get(tc.name)
			assert.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.expected, v)
		})
	}
}

func TestLookupHandlers_RpcMultipleFunctions(t *testing.T) {
	a2k := &A2K{"rpc": struct{}{}}
	wm := &WM{}
	LookupHandlers(makeCallExpr("rpc", "RPCStub", &ast.BasicLit{Kind: token.STRING, Value: `"a"`}), a2k, wm)
	LookupHandlers(makeCallExpr("rpc", "RPCStub", &ast.BasicLit{Kind: token.STRING, Value: `"b"`}), a2k, wm)
	LookupHandlers(makeCallExpr("rpc", "RPCStubStream", &ast.BasicLit{Kind: token.STRING, Value: `"c"`}), a2k, wm)
	v, _ := wm.Get("rpc")
	assert.Equal(t, []string{"a", "b", "c"}, v)
}

func TestLookupStructs_NonTypeSpec(t *testing.T) {
	fields := LookupStructs(&ast.FuncDecl{})
	assert.Empty(t, fields)
}

func TestLookupStructs_NonStruct(t *testing.T) {
	node := &ast.TypeSpec{
		Name: &ast.Ident{Name: "MyAlias"},
		Type: &ast.Ident{Name: "string"},
	}
	fields := LookupStructs(node)
	assert.Empty(t, fields)
}

func TestLookupStructs_StructWithoutDO(t *testing.T) {
	node := &ast.TypeSpec{
		Name: &ast.Ident{Name: "Plain"},
		Type: &ast.StructType{Fields: &ast.FieldList{
			List: []*ast.Field{
				{Names: []*ast.Ident{{Name: "X"}}, Type: &ast.Ident{Name: "int"}},
			},
		}},
	}
	fields := LookupStructs(node)
	assert.Empty(t, fields)
	_, inTypes := AllStructs.Get("Plain")
	assert.True(t, inTypes)
}

func TestLookupStructs_StructWithDO(t *testing.T) {
	node := &ast.TypeSpec{
		Name: &ast.Ident{Name: "MyDO"},
		Type: &ast.StructType{Fields: &ast.FieldList{
			List: []*ast.Field{
				{
					Type: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "durableobjects"},
						Sel: &ast.Ident{Name: "DurableObject"},
					},
				},
			},
		}},
	}
	fields := LookupStructs(node)
	assert.Equal(t, []string{"MyDO"}, fields)
}

func makeFuncDecl(recv, funcName string, params, results []*ast.Field) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{Name: funcName},
		Recv: &ast.FieldList{List: []*ast.Field{
			{Type: &ast.Ident{Name: recv}},
		}},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: params},
			Results: &ast.FieldList{List: results},
		},
	}
}

func TestInspectDOFuncs_NotFuncDecl(t *testing.T) {
	fields := []string{}
	dom := &DOM{}
	path := "file.go"
	LookupDurableObjects(&ast.Ident{}, fields, dom, &path)
	assert.Empty(t, dom.m)
}

func TestLookupDurableObjects_NoReceiver(t *testing.T) {
	fields := []string{}
	dom := &DOM{}
	path := "file.go"
	node := &ast.FuncDecl{
		Name: &ast.Ident{Name: "Standalone"},
		Recv: nil,
		Type: &ast.FuncType{Params: &ast.FieldList{}},
	}
	LookupDurableObjects(node, fields, dom, &path)
	assert.Empty(t, dom.m)
}

func TestLookupDurableObjects_ReceiverNotInFields(t *testing.T) {
	fields := []string{}
	dom := &DOM{}
	path := "file.go"
	LookupDurableObjects(makeFuncDecl("MyDO", "Greet", nil, nil), fields, dom, &path)
	assert.Empty(t, dom.m)
}

func TestLookupDurableObjects_SkipsUnderscorePrefix(t *testing.T) {
	fields := []string{"MyDO"}
	dom := &DOM{}
	path := "file.go"
	LookupDurableObjects(makeFuncDecl("MyDO", "_internal", nil, nil), fields, dom, &path)
	assert.Empty(t, dom.m)
}

func TestLookupDurableObjects_RecordsFunc(t *testing.T) {
	fields := []string{"MyDO"}
	dom := &DOM{}
	path := "file.go"

	params := []*ast.Field{
		{Type: &ast.Ident{Name: "string"}},
	}
	results := []*ast.Field{
		{Type: &ast.Ident{Name: "int"}},
	}
	LookupDurableObjects(makeFuncDecl("MyDO", "DoWork", params, results), fields, dom, &path)

	v, ok := dom.Get("MyDO")
	assert.True(t, ok)
	assert.Len(t, v, 1)
	assert.Equal(t, "DoWork", v[0].FuncName)
	assert.Equal(t, []string{"string"}, v[0].Args)
	assert.Equal(t, []string{"int"}, v[0].Rargs)
	assert.Equal(t, "file.go", v[0].Path)
}

func TestLookupDurableObjects_SkipsContextParam(t *testing.T) {
	fields := []string{"MyDO"}
	dom := &DOM{}
	path := "file.go"

	params := []*ast.Field{
		{Type: &ast.SelectorExpr{X: &ast.Ident{Name: "context"}, Sel: &ast.Ident{Name: "Context"}}},
		{Type: &ast.Ident{Name: "string"}},
	}
	LookupDurableObjects(makeFuncDecl("MyDO", "Handle", params, nil), fields, dom, &path)

	v, _ := dom.Get("MyDO")
	assert.Equal(t, []string{"string"}, v[0].Args)
}
