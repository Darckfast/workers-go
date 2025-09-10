//go:build js && wasm

package d1

import (
	"math"
	"syscall/js"
	"testing"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
	"github.com/Darckfast/workers-go/platform/cloudflare/lifecycle"
	"github.com/stretchr/testify/assert"
)

func setupEnv() *js.Value {
	mock := jsclass.Object.New()
	mock.Set("prepare", js.Global().Get("Function").New("query", `
    this.query = query
    const that = this
    return { bind: function(values) {
        that.values = values
    }}
    `))

	lifecycle.Env = jsclass.Object.New()
	lifecycle.Env.Set("BINDING", mock)

	return &mock
}

func TestPrepareMethod(t *testing.T) {
	mock := setupEnv()

	db, _ := GetDB("BINDING")
	db.Prepare("SELECT * FROM mytable WHERE id = ?")

	query := mock.Get("query").String()

	assert.Equal(t, "SELECT * FROM mytable WHERE id = ?", query)
}

func TestMaxInt32Bind(t *testing.T) {
	mock := setupEnv()

	db, _ := GetDB("BINDING")
	db.Prepare("").Bind(math.MaxInt32)

	bindings := mock.Get("values").Int()

	assert.Equal(t, math.MaxInt32, bindings)
}

func TestMaxInt64Bind(t *testing.T) {
	mock := setupEnv()

	db, _ := GetDB("BINDING")
	db.Prepare("").Bind(math.MaxInt64)

	bindings := jsconv.MaybeInt64(mock.Get("values"))

	assert.Equal(t, int64(math.MaxInt64), bindings)
}
