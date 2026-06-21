//go:build !js && !wasm

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fetchImpl = `
//go:build js && wasm

package main

import (
	"net/http"

	"codeberg.org/darckfast/workers-go/platform/cloudflare/fetch"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	})

	fetch.ServeNonBlock(mux)
}`

const modfile = `
module main

go 1.21

require codeberg.org/darckfast/workers-go v0.4.2

require (
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.9.2 // indirect
)`

const modsum = `
codeberg.org/darckfast/workers-go v0.4.2 h1:V7nMYy8cQ7F4iFgfwXc+Qy6Lwu6PwbRxehN8QwWJZM0=
codeberg.org/darckfast/workers-go v0.4.2/go.mod h1:rPkFKyaO1ZYPDE+tM2Tfc/+VC4V1DEEuG25azff401k=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/josharian/intern v1.0.0 h1:vlS4z54oSdjm0bgjRigI+G1HpF+tI+9rE5LLzOg8HmY=
github.com/josharian/intern v1.0.0/go.mod h1:5DoeVV0s6jJacbCEi61lwdGj/aVlrQvzHFFd8Hwg//Y=
github.com/mailru/easyjson v0.9.2 h1:dX8U45hQsZpxd80nLvDGihsQ/OxlvTkVUXH2r/8cb2M=
github.com/mailru/easyjson v0.9.2/go.mod h1:1+xMtQp2MRNVL/V1bOzuP3aP8VNwRW55fQUto+XFtTU=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=`

func TestCLICompile(t *testing.T) {
	tmpdir := t.TempDir()
	mf, _ := os.Create(filepath.Join(tmpdir, "main.go"))
	modf, _ := os.Create(filepath.Join(tmpdir, "go.mod"))
	sumf, _ := os.Create(filepath.Join(tmpdir, "go.sum"))

	mf.Write([]byte(fetchImpl))
	modf.Write([]byte(modfile))
	sumf.Write([]byte(modsum))

	os.Args = []string{"./workers-go", "-o", filepath.Join(tmpdir, "dist"), "-i", tmpdir}

	main()

	outwasmfile, _ := os.Stat(filepath.Join(tmpdir, "dist", "app.wasm"))
	outwasmjs, _ := os.Stat(filepath.Join(tmpdir, "dist", "wasm_exec.js"))
	outmaints, _ := os.Stat(filepath.Join(tmpdir, "dist", "main.ts"))

	assert.NotNil(t, outwasmfile)
	assert.Greater(t, outwasmfile.Size(), int64(1024))
	assert.NotNil(t, outwasmjs)
	assert.Greater(t, outwasmjs.Size(), int64(1024))
	assert.NotNil(t, outmaints)
	assert.Greater(t, outmaints.Size(), int64(1024))
}
