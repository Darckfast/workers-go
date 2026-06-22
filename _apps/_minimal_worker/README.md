# Minimal Worker

First install [mise](https://mise.jdx.dev/getting-started.html) then run:

```bash
pnpx wrangler dev
```

*Access or cURL http://localhost:8787*

## Deploy

```bash
pnpx wrangler deploy
```

## Using TinyGo

Change the `workers-go` build command to use the argument `-tiny`
```toml
[build]
command = "go install codeberg.org/darckfast/workers-go/cmd/workers-go && workers-go -i . -tiny"
```

PS: TinyGo is incompatible with Cloudflare's build images, you will need to use `wrangler deploy` to deploy it
