<p align="center">
  <a href="https://darckfast.com/docs/workers-go">
    <img alt="workers-go" src=".github/images/workers-go.png">
    <h1 align="center">workers-go</h1>
  </a>
</p>

<p align="center">
  workers-go is fork of <a href="https://github.com/syumai/workers">syumai's workers</a> ❤️ — a lightweight package for building and running Go on <a href="https://workers.cloudflare.com/">Cloudflare Workers</a> using WebAssembly. Also compatible with NodeJS, Bun, and Deno
</p>

<p align="center">
  <a href="https://darckfast.com/docs/workers-go"><strong>📜 docs</strong></a>
  <a href="https://codeberg.org/darckfast/workers-go"><strong>⭐ main repository</strong></a>
  <a href="https://github.com/Darckfast/workers-go"><strong>🪞 mirror repository</strong></a>
</p>

<p align="center">
  <i>GitHub is a mirror, all development is centered on codeberg.org. Issues are welcomed on both</i>
</p>

---

## Getting Started
This library uses [mise](https://mise.jdx.dev/) to prepare the dev environment, otherwise you must install Go, Node +24, pnpm, and TinyGo (if applicable)

```sh
# minimal worker with only GET /hello
pnpm create cloudflare@latest --template=github.com/darckfast/workers-go/_apps/_minimal_worker

# tidy modules
go mod tidy

# dev
pnpx wrangler dev

# deploy
pnpx wrangler deploy
```

---

[Check the complete documentation for more details](https://darckfast.com/docs/workers-go)
