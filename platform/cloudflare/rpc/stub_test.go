//go:build js && wasm

package rpc

import (
	"context"
	"io"
	"net/http"
	"os"
	"syscall/js"
	"testing"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"github.com/stretchr/testify/assert"
)

var app js.Value

func TestMain(m *testing.M) {
	app = jsclass.Object.New()

	js.Global().Set("workerapp", app)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreateNewRPCStub(t *testing.T) {
	RPCStub("test", func(c context.Context, args [][]byte) [][]byte {
		assert.Len(t, args, 5)
		assert.Nil(t, args[0])
		assert.Nil(t, args[1])
		assert.Equal(t, []byte("arg1"), args[2])
		assert.Equal(t, []byte("arg2"), args[3])
		assert.Equal(t, []byte("arg3"), args[4])

		return [][]byte{[]byte("result from rpc")}
	})

	b1 := []byte("arg1")
	b2 := []byte("arg2")
	b3 := []byte("arg3")

	arg1 := jsclass.Uint8Array.New(len(b1))
	arg2 := jsclass.Uint8Array.New(len(b2))
	arg3 := jsclass.Uint8Array.New(len(b3))

	js.CopyBytesToJS(arg1, b1)
	js.CopyBytesToJS(arg2, b2)
	js.CopyBytesToJS(arg3, b3)

	r, _ := jsclass.Await(app.Call("test", js.Null(), js.Undefined(), arg1, arg2, arg3))

	var dst = make([]byte, r.Index(0).Length())
	js.CopyBytesToGo(dst, r.Index(0))

	assert.Equal(t, "result from rpc", string(dst))
}

func TestCreateNewRPCStubStream(t *testing.T) {
	RPCStubStream("test-stream", func(c context.Context, w http.ResponseWriter, body io.ReadCloser, args [][]byte) {
		b, _ := io.ReadAll(body)
		assert.Equal(t, `{"test":2}`, string(b))

		_, err := w.Write([]byte("writer response from rpc"))
		assert.Nil(t, err)
	})

	req := jsclass.Request.New("http://dummy", map[string]any{
		"method": "post",
		"body":   `{"test":2}`,
	})

	b1 := []byte("arg1")

	arg1 := jsclass.Uint8Array.New(len(b1))

	js.CopyBytesToJS(arg1, b1)

	r, _ := jsclass.Await(app.Call("test-stream", req.Get("body"), js.Null(), js.Undefined(), arg1))
	result, _ := jsclass.Await(r.Call("getReader").Call("read"))

	br := result.Get("value")
	var dst = make([]byte, br.Length())
	js.CopyBytesToGo(dst, br)

	assert.Equal(t, "writer response from rpc", string(dst))
}
