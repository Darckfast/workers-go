# Minimal Queues Consumer Worker

While developing locally, this worker needs to run in the same process as the producer

```bash
pnpx wrangler dev -c ../path-to-producer/wrangler.toml -c wrangler.toml --persist-to .wrangler/state
```

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
