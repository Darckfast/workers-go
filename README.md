
![banner](.github/images/banner.png)

Powered by <img src="https://vite.dev/logo.svg" style="height: 1rem"/> Vite and <img src="https://workers.cloudflare.com/logo.svg" style="height: 1rem"/> Cloudflare Workers

# workers-go

This repository is a fork of https://github.com/syumai/workers ❤️

`workers-go` is a pure Go library, made to help interface Go's WASM with [Cloudflare Workers](https://workers.cloudflare.com/).
It implements a series of handlers, helpers and bindings, making easier to integrate Go with Workers

## Quick Start

This project has only been tested on **Go 1.23+** with **NodeJS 22+**

### Create a new Worker project

Run the following command:

```bash
pnpm create cloudflare@latest -- --template github.com/Darckfast/workers-go/worker
```

### Initialize the project

1. Navigate to your new project directory:

```bash
cd my-app
```

2. Initialize Go modules and NodeJS packages:

```bash
pnpm run init # this will run go mod tidy && pnpm install
```

3. Start the development server:

```bash
pnpm run dev
```

4. Verify the worker is running:

```bash
curl http://localhost:5173/hello
```

## Installation

```bash
go get github.com/Darckfast/workers-go
```

## Features

Below is a list of implemented, and not implemented Cloudflare features

| Feature                      | Implemented | Notes                                                                                                                                            |   |
|------------------------------|-------------|--------------------------------------------------------------------------------------------------------------------------------------------------|---|
| `fetch`                      | ✅           | All functions uses either `http.Request` or `http.Response`                          |   |
| `queue`                      | ✅           |                                                                                                                                                  |   |
| `email`                      | ✅           |                                                                                                                                                  |   |
| `scheduled`                  | ✅           |                                                                                                                                                  |   |
| `tail`                       | ✅           | **EXPERIMENTAL**: This has not been tested in production env yet                                                                                 |   |
| Env                          | ✅           | All Cloudflare Worker's env are copied into `os.Environ()`, making them available at runtime with `os.Getenv()`. Only string typed values are copied |   |
| Containers                   | 🔵          | Only the `containerFetch()` function has been implemented                                                                                        |   |
| R2                           | 🔵          | _Options for R2 methods still not implementd_                                                                                                    |   |
| D1                           | 🔵          |                                                                                                                                                  |   |
| KV                           | 🔵          | _Options for KV methods still not implemented_                                                                                                   |   |
| Cache API                    | ✅           |                                                                                                                                                  |   |
| Durable Objects              | 🔵          | _Only stub calls have been implemented_                                                                                                          |   |
| RPC                          | ❌           | _Not implemented_                                                                                                                                |   |
| Service binding              | ✅           | `fetch.Client{}.WithBinding(serviceName)`. only works for `fetch` or HTTP requests                                                                                                                          |   |
| HTTP                         | ✅           | native fetch interface using `fetch.Client{}.Do(req)`                                                                                            |   |
| HTTP Timeout                 | ✅           | Implemented using the same interface as `http.Client{ Timeout: 20 * time.Second }`                                                               |   |
| HTTP RequestInitCfProperties | ✅           | Implemented all but the `image` property, they must be set on the `http.Client{ CF: &RequestInitCF{} }`                                          |   |
| FetchEvent                   | ✅          |                                                                                                                                                  |   |
| TCP Sockets                  | ✅          |                                                                                                                                                  |   |
| Queue producer               | ✅          |                                                                                                                                                  |   |

## Implementing `fetch` handler

Implement your `http.Handler` and give it to `fetch.ServeNonBlock()`.

```go
//go:build js && wasm

package main

func main() {
	var handler http.HandlerFunc = func (w http.ResponseWriter, req *http.Request) {
    //...
  }
	fetch.ServeNonBlock(handler)

  <-make(chan struct{})
}
```

or just call `http.Handle` and `http.HandleFunc`, then invoke `workers.Serve()` with nil.

```go
//go:build js && wasm

package main

func main() {
	http.HandleFunc("/hello", func (w http.ResponseWriter, req *http.Request) {
    //...
  })

	fetch.ServeNonBlock(handler)// if nil is given, http.DefaultServeMux is used.

  <-make(chan struct{})
}
```

## Making HTTP Request
For compatability reasons, you must use the `fetch.Client{}` to make http request, as it interfaces Go's http with Cloudflare Worker `fetch()` API

```go
r, _ := http.NewRequest("GET", "https://google.com", nil)
c := fetch.Client{
  Timeout: 5 * time.Second,
}

// Timeouts return error
rs, err := c.Do(r)

defer rs.Body.Close()
b, _ := io.ReadAll(rs.Body)

fmt.Println(string(b))
```

## `main.ts`

`main.ts` is the entry point, declared in the `wrangler.toml`, and its where the wasm binary
will be loaded and used

Below is a (_non functional_) example, for a functional and complete example check `./worker/bin/main.ts`

```ts
import app from "./bin/app.wasm"; // Compiled wasm binary
import "./bin/wasm_exec.js"; // cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .

/**
  * This function is what initialize your Go's compiled WASM binary
  * only after this function has finished, that the handlers will be
  * defined in the globalThis scope
  *
  * At the moment, due limitations with the getRandomData(), this block
  * cannot be executed at top level, it must be contained within the handlers
  * scope
  *
  * It's REQUIRED and needs to be called before using the globalThis.cf.<handler>()
  */
function init() {
  const go = new Go()

  /*
  * This will execute the binary, and all Go's `init()` will run and instantiate
  * the callbacks. They all will be within the globalThis.cf object
  */
  go.run(new WebAssembly.Instance(app, go.importObject))
}

async function fetch(req: Request, env: Env, ctx: ExecutionContext) {
  init()
  return await globalThis.cf.fetch(req, env, ctx);
}

export default {
  fetch,
} satisfies ExportedHandler<Env>;
```

## Caveats

### ▶️ C Binding
IF you use any library or package that depends or use any C binding, or C compiled code, compiling to WASM is not possible

Some examples

| Package | Compatible |
|-|-|
|https://github.com/anthonynsimon/bild|✅|
|https://github.com/nfnt/resize|✅|
|https://github.com/bamiaux/rez|✅|
|https://github.com/kolesa-team/go-webp|❌|
|https://github.com/Kagami/go-avif|❌|
|https://github.com/h2non/bimg|❌|
|https://github.com/davidbyttow/govips|❌|
|https://github.com/gographics/imagick|❌|

### HTTP Requests
When making http request, the `fetch.NewClient()` must be used, as it implements the Cloudflare Worker native `fetch()` call

### Queues
Cloudflare Queue locally is incredibly slow to produce events (up to 7 seconds)

### TinyGo
Go's compiled binary can exceed the Free 3MB Cloudflare Worker's limit, in which case one suggestion is to use TinyGo to compile, but for performance reasons `workers-go` uses the `encoding/json` from the std Go's library, which makes this package incompatible with the current build of TinyGo

Another possible fix is related to this issue https://github.com/golang/go/issues/63904

### Errors
Although we can wrap JavaScript errors in Go, at the moment there is no source maps available in wasm, meaning we can get errors messages, but not a useful stack trace

### Build constraint
For [gopls](https://github.com/golang/tools/tree/master/gopls) to show `syscall/js` method's signature and auto complete, either `export GOOS=js && export GOARCH=wasm` or add the comment `//go:build js && wasm` at the top of your Go files

