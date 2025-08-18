
![banner](.github/images/banner.png)

Powered by !<img src="https://vite.dev/logo.svg" style="height: 1rem"/> Vite and <img src="https://workers.cloudflare.com/logo.svg" style="height: 1rem"/> Cloudflare Workers

# workers-go

This repository is a fork of https://github.com/syumai/workers ❤️

<!-- [![Go Reference](https://pkg.go.dev/badge/github.com/syumai/workers.svg)](https://pkg.go.dev/github.com/syumai/workers) -->
<!-- [![Discord Server](https://img.shields.io/discord/1095344956421447741?logo=discord&style=social)](https://discord.gg/tYhtatRqGs) -->

`workers-go` is a pure Go library, made to help interface Go's WASM with [Cloudflare Workers](https://workers.cloudflare.com/).
It implements a series of handlers, helpers and bindings, making easier to integrate Go with Workers

## Features

| Feature | Implemented | Notes |
|-|-|-|
|`fetch`|✅| _At the moment all request use HTTP, RPC is not supported_. All functions uses either `http.Request` or `http.Response`|
|`queue`|✅||
|`email`|✅||
|`scheduled`|✅||
|`tail`|✅| **EXPERIMENTAL**: This has not been tested in production env yet|
|`env`|✅|All Cloudflare Worker's are copied into `os.Environ()`, making them available at runtime with `os.Getenv()`. Only string typed values are copied|
|Containers| 🔵| Only the `containerFetch()` function has been implemented|
|R2| 🔵|_Options for R2 methods still not implementd_|
|D1|🔵||
|KV|🔵|_Options for KV methods still not implemented_|
|Cache API|✅||
|Durable Objects|🔵|_Only stub calls have been implemented_|
|RPC|❌|_Not implemented_|
|Service binding|❌|_Not implemented_|
|HTTP|🔵|_RequestInitCfProperties still not implemented_|
|FetchEvent|🔵||
|TCP Sockets|🔵||
|Queue producer|🔵||

## Installation

```bash
go get github.com/Darckfast/workers-go
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
## Caveats

### C Binding
If any Go lib that depends on a C lib, e.g. vips, it wont work

### HTTP Requests
When making http request, the `fetch.NewClient()` must be used, as it implements the Cloudflare Worker native `fetch()` call

### Queues
Cloudflare Queue locally is incredibly slow to produce events (up to 7 seconds)

### TinyGo
Go's compiled binary can exceed the Free Cloudflare Worker's limit, in which case one suggestion is to use TinyGo to compile, but for performance reasons, this package makes use of the `encoding/json` from the std Go's library, which makes this package incompatible with the current build of TinyGo

### Errors

Although we can wrap JavaScript errors in Go, at the moment there is no source maps available in wasm, meaning we can get errors messages, but not a useful stack trace

### Build constraint

For Gopls show `syscall/js` methods signature and autocomplete, either `export GOOS=js && export GOARCH=wasm` or add the comment `//go:build js && wasm` at the top of your Go file

## Quick Start
<!---->
<!-- * You can easily create and deploy a project from `Deploy to Cloudflare` button. -->
<!---->
<!-- <!-- [![Deploy to Cloudflare](https://deploy.workers.cloudflare.com/button)](https://deploy.workers.cloudflare.com/?url=https%3A%2F%2Fgithub.com%2Fsyumai%2Fworker-go-deploy) --> -->
<!---->
<!-- * If you want to create a project manually, please follow the guide below. -->

### Requirements

* NodeJS v22+
* Go 1.24+

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

2. Initialize Go modules:

```bash
go mod init
go mod tidy
```

3. Start the development server:

The development server is powered by Vite and Cloudflare Worker's plugin

```bash
pnpm install
pnpm run dev
```

4. Verify the worker is running:

```bash
curl http://localhost:5173/hello
```
