//go:build js && wasm

package d1

import (
	"math"
	"syscall/js"
	"testing"

	"codeberg.org/darckfast/workers-go/internal/jsclass"
	"codeberg.org/darckfast/workers-go/internal/jsconv"
	"github.com/stretchr/testify/assert"
)

func setupEnv() *js.Value {
	mock := jsclass.Object.New()
	mock.Set("prepare", js.Global().Get("Function").New("query", `
    this.query = query
    const that = this

    return {
      bind: function(values) {
        that.values = values
      },
      first: function() {
       return Promise.resolve( { "test": 2 })
      },
      run: function() {
       return Promise.resolve({
              success: true,
              results: [ { "test": 1 } ]
            })
      }
    }
    `))

	v := jsclass.Object.New()
	v.Set("BINDING", mock)
	jsclass.Env = jsclass.EnvBinding{}
	jsclass.Env.LoadEnvs(v)

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
	n := int64(math.MaxInt64)
	db.Prepare("").Bind(n)

	bindings := jsconv.MaybeInt64(mock.Get("values"))

	assert.Equal(t, int64(math.MaxInt64), bindings)
}

func TestReturnResultAsString(t *testing.T) {
	setupEnv()

	db, _ := GetDB("BINDING")
	r, err := db.Prepare("SELECT * FROM mytable WHERE id = ?").Run()

	assert.Nil(t, err)
	assert.Equal(t, `[{"test":1}]`, r.ResultsString)
}

func TestReturnFirstResultAsString(t *testing.T) {
	setupEnv()

	db, _ := GetDB("BINDING")
	r, err := db.Prepare("SELECT * FROM mytable WHERE id = ?").FirstAsString(nil)

	assert.Nil(t, err)
	assert.Equal(t, `{"test":2}`, r)
}

func TestReturnFirstResultAsStringWithColumnName(t *testing.T) {
	setupEnv()

	db, _ := GetDB("BINDING")
	col := "test"
	r, err := db.Prepare("SELECT * FROM mytable WHERE id = ?").FirstAsString(&col)

	assert.Nil(t, err)
	assert.Equal(t, `{"test":2}`, r)
}
