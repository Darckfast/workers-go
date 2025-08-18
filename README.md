# workers-go

This repo is a fork of https://github.com/syumai/workers ❤️

<!-- [![Go Reference](https://pkg.go.dev/badge/github.com/syumai/workers.svg)](https://pkg.go.dev/github.com/syumai/workers) -->
<!-- [![Discord Server](https://img.shields.io/discord/1095344956421447741?logo=discord&style=social)](https://discord.gg/tYhtatRqGs) -->

<!-- * `workers` is a package to run an HTTP server written in Go on [Cloudflare Workers](https://workers.cloudflare.com/). -->
<!-- * This package can easily serve *http.Handler* on Cloudflare Workers. -->
<!-- * Caution: This is an experimental project. -->

## Features

| Feature | Implemented | Notes |
|-|-|-|
|Handler `fetch`|✅| _At the moment all request use HTTP, RPC is not supported_. All functions uses either `http.Request` or `http.Response`|
|Handler `queue`|✅||
|Handler `email`|✅||
|Handler `scheduled`|✅||
|Handler `tail`|✅| **EXPERIMENTAL**: This has not been tested in production env yet|
|`Containers`| 🔵| Only the `containerFetch()` function has been implemented|
|`env`|✅|All Cloudflare Worker's are copied into `os.Environ()`, making them available at runtime with `os.Getenv()`|

## Installation

```bash
go get github.com/Darckfast/workers-go
```
## `main.ts`

`main.ts` is the entry point, declared in the `wrangler.toml`, and its where the wasm binary
will be loaded and used

Below is a (_non functional_) example, for a functional and complete one check `./worker/bin/main.ts`

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

<!-- For concrete examples, see `_examples` directory. -->

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

```console
pnpm create cloudflare@latest -- --template github.com/Darckfast/workers-go/worker
```

### Initialize the project

1. Navigate to your new project directory:

```console
cd my-app
```

2. Initialize Go modules:

```console
go mod init
go mod tidy
```

3. Start the development server:

```console
pnpm install
pnpm run dev
```

4. Verify the worker is running:

```console
curl http://localhost:5173/hello
```

<!-- You will see **"Hello!"** as the response. -->

<!-- If you want a more detailed description, please refer to the README.md file in the generated directory. -->

<!-- ## FAQ -->
<!---->
<!-- ### How do I deploy a worker implemented in this package? -->
<!---->
<!-- To deploy a Worker, the following steps are required. -->
<!---->
<!-- * Create a worker project using [wrangler](https://developers.cloudflare.com/workers/wrangler/). -->
<!-- * Build a Wasm binary. -->
<!-- * Upload a Wasm binary with a JavaScript code to load and instantiate Wasm (for entry point). -->
<!---->
<!-- The [worker-go template](https://github.com/syumai/workers/tree/main/_templates/cloudflare/worker-go) contains all the required files, so I recommend using this template. -->
<!---->
<!-- But Go (not TinyGo) with many dependencies may exceed the size limit of the Worker (3MB for free plan, 10MB for paid plan). In that case, you can use the [TinyGo template](https://github.com/syumai/workers/tree/main/_templates/cloudflare/worker-tinygo) instead. -->
<!---->
<!-- ### Where can I have discussions about contributions, or ask questions about how to use the library? -->
<!---->
<!-- You can do both through GitHub Issues. If you want to have a more casual conversation, please use the [Discord server](https://discord.gg/tYhtatRqGs). -->
