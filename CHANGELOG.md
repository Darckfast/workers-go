# Changelog

## [0.5.2](https://codeberg.org/darckfast/workers-go/compare/v0.5.1...v0.5.2) (2026-07-02)

### 🚀 Features

- Refactored `workers-go` client/transport to be compatible with Go `net/http` module (https://codeberg.org/darckfast/workers-go/pulls/39)
- Add -nocompression opt in the workers-go` cli (https://codeberg.org/darckfast/workers-go/pulls/40)

### 🐛 Bug Fixes

- Fixed workers not restarting in wasm after a crash (https://codeberg.org/darckfast/workers-go/pulls/38)

## [0.5.1](https://codeberg.org/darckfast/workers-go/compare/v0.5.0...v0.5.1) (2026-06-28)

### 🚀 Features
- `bind.Ctx` and `bind.Env` are accessible everywhere (https://codeberg.org/darckfast/workers-go/pulls/37)
- Fetch now uses a wriatable passed from Js to minimize the serialization and deserialization overhead when transferring data between JS -> Go (https://codeberg.org/darckfast/workers-go/pulls/37)
- Removed `RPCStream()` _this stub will be reworked_

### 🐛 Bug Fixes
- `procces.env` copy to `os.Environ()` now copies only env variables (https://codeberg.org/darckfast/workers-go/pulls/37)

## [0.5.0](https://codeberg.org/darckfast/workers-go/compare/v0.4.3...v0.5.0) (2026-06-26)

### 🚀 Features

- Durable objects and types generation (https://codeberg.org/darckfast/workers-go/pulls/34)
  - Durable Objects now can be interfaced through Go code. See `_apps/_durable_object` for an working exampls (Not all bindings were implemented yet)
- All JavaScript classes and methods are now acquired through lazy unwrapping using `jshelper.LazyJSVal` (https://codeberg.org/darckfast/workers-go/pulls/35)
- `Env` and `Ctx` now are loaded only once with `sync.Once` at its first usage (https://codeberg.org/darckfast/workers-go/pulls/35)

### 🐛 Bug Fixes
- When producers sends `Uint8Array` it gets transformed into `ArrayBuffer` on the receiving consumer, this checks and transforms it back to `Uint8Array`, this allows `CopyBytesToGo` to be used - fixed https://github.com/Darckfast/workers-go/issues/45 (https://codeberg.org/darckfast/workers-go/pulls/33)
- `SQLDurationMs` and `Duration` were both using int64, while this works locally, when deployed on Cloudflare, D1 returns a float, which was resulting in a conversion error being returned - fixes https://github.com/Darckfast/workers-go/issues/46 (https://codeberg.org/darckfast/workers-go/pulls/33)
- Remove legacy lifecycle functions and variables (https://codeberg.org/darckfast/workers-go/pulls/35)

### ⚙️ Miscellaneous Tasks

- Updated go install cmd to use @latest tag (https://codeberg.org/darckfast/workers-go/pulls/30)
- Updated README.md to use github.com url for templates (https://codeberg.org/darckfast/workers-go/pulls/31)

## [0.4.3](https://codeberg.org/darckfast/workers-go/compare/v0.4.2...v0.4.3) (2026-06-20)

### 🚀 Features

- Add `workers-go` cli to create `wrangler` entrypoint typescript file with the expected typing for better integration, compile your Go project using
either Go or TinyGo (`-tiny`), and compress using `wasm-opt` if available
- Updated `_apps` to use the `workers-go` cli for building
- Add `minimumReleaseAge` (and equivalents to each package manager), blocked scripts runs on by default in this repository, and modified all `mise`
install shortcut to use frozen lockfile to help prevent supply chain attacks

### 🐛 Bug Fixes

- Replaced global call with `jsclass`
- Removed envs `GOOS` and `GOARCH` from CI due issues with `workers-go` cli

### 💼 Other

- Add new linting rules, and fixed all warning and errors

### 📚 Documentation

- Updated README.md with a small snippet to get started using this library
- Updated `_apps` with a cleaner templates using `workers-go` cli

### 🧪 Testing

- Updated tests commands, and queues to run with concurrency
- Fixed test timestamp greater than it should be
- Add basic compiling testing for `workers-go` cli

### ⚙️ Miscellaneous Tasks

- Add CI jobs/testing filtering according to its respective direct dependencies, so it runs only
the necessary jobs according to what have changed in the PR
- Add `node` on root mise for `ts_ls` lsp to work properly
- Changed testing queue `batch_timeout` to 0
- Renamed `DEVELOPING.md` to `CONTRIBUTING md`
- Removed `bun install` when running `mise test-lib` which are only Go testing files
- Increased tests timeout to 60s for CI due constraint with runners in Codeberg
- Updated .gitignore file to include `bin/` diretory
- Updated repository banner image
- Change CI instance type due timeout when building
- Removed `@faker` dependency that was used only for testing
- Removed unnecessary configuration from root `tsonfig.json`
- Add `GOOS=js` and `GOARCH=wasm` on linting command
- Updated lock files
- Add `cmd/` folder to the linting rules in CI


## [0.4.2](https://codeberg.org/darckfast/workers-go/compare/v0.4.1...v0.4.2) (2026-06-06)

### 🚀 Features
- Bump Golang version to `1.24` to be compatible with `tinygo@0.41.1`
- Add `.metrics()` function, and return correct values when calling the queue producer `.send()`
- Updated app `worker`, `bun`, `deno`, `ultimate-express` and `minimal-queue` tooling to use `mise`
- Removed hard-coded `wasm_exec.js` and `bin` directory
- Updated lib and apps `mise` commands to setup and copy the necessary files at build time
- Add DEVELOPING.md with basic instructions to run and develop this project
- Add `--tiny` arg when using `_worker` or `_minimal_queues_consumer`, this will compile and use `tinygo`
```bash
mise dev --tiny
```

### 🐛 Bug Fixes

- Queue function `BytesBody()` now accepts `ArrayBuffer`
- Fixed wrong conversion from 64 to 32 when compiled using `tinygo`
- Fixed wrong value comparison when running tests with `tinygo`
- Fixed `mise` command path on CI
- Create missing directory during build

### ⚙️ Miscellaneous Tasks
- Remove `biome` and `vite`
- Re-enabled uws and hono testing suit on CI
- Removed unused configs and files
- Updated `_apps` `.gitignore` and `.dockerignore`
- Add env MISE_EXPERIMENTAL=1 on CI
- Updated CI workflow's actions and run commands

## [0.4.1](https://codeberg.org/darckfast/workers-go/compare/v0.4.0...v0.4.1) (2026-05-31)

### 🐛 Bug Fixes

- Fixed function `tryCatch` not being created correctly

## [0.4.0](https://codeberg.org/darckfast/workers-go/compare/v0.3.0...v0.4.0) (2026-05-31)
### 🚀 Features

- Add truthy check on args before trying to copying JS bytes into Go - this allows null and undefined values to be used
- Create and pass context.Background on handlers calls
- Improved `ctx` and `env` copy on RPC stub creation
- Add `getPrototypeOf` js func call

### 🐛 Bug Fixes

- Fixed spread on args

### 🧪 Testing

- Add tests cases for rpc stub creation

### ⚙️ Miscellaneous Tasks

- Updated handlers to accept `ctx` from workers-go inner libs
- Add fixed node ver on mise.toml
- Removed unused code block
- Cleanup commented code
- Fixed linting error

## [0.3.0](https://codeberg.org/darckfast/workers-go/compare/v0.2.7...v0.3.0) (2026-05-29)

### ⚙️ Miscellaneous Tasks

- Migrated package to codeberg
- Bump workers-go ver no queue template
- Removed unused mod dependencies

## [0.2.7](https://github.com/Darckfast/workers-go/compare/v0.2.2...v0.2.7) (2026-05-29)

### 🚀 Features

- Changed the return type from RPCStub to [][]byte, this allows to return multiple values

### ⚡ Performance

- Replaced `.Call('push',...)` with `.SetIndex(i,...)` to reduced the overhead when inserting new elements in a list

### ⚙️ Miscellaneous Tasks

- Removed unused workflows and files

## [0.2.6](https://github.com/Darckfast/workers-go/compare/v0.2.2...v0.2.6) (2026-05-28)

### 🐛 Bug Fixes

- Fixed set-cookie headers not being copied correctly
- Fixed call with no column name
- Fixed broken test
- Fixed typo
- Fixed missing args on durable objects stub rpc call
- Fixed pointer reference with calling JS binding
- Fixed worker test CI
- Fixed tests and CI
- Fixed broken confs on yml
- Add depends on install before running tests

### 🧪 Testing

- Fixed test result
- Add test for rpc stubs
- Fixed broken tests
- Added RPC stub test

### ⚙️ Miscellaneous Tasks

- Added missing package.json
- Fixed arg name
- Updated workflows
- Bump easyjson dep
- Removed labaler ci
- Updated deps and configs
- Reduce queue timeout for testing
- Bump workers-go version and imports
- Updated test:go workflow to use mise
- Updated workflows and removes unused .js
- Updated README.md

## [0.2.0](https://github.com/Darckfast/workers-go/compare/v0.1.0...v0.2.0) (2025-09-07)

### Features

* add new testing endpoint for r2 with encryption ([#28](https://github.com/Darckfast/workers-go/issues/28)) ([aca215e](https://github.com/Darckfast/workers-go/commit/aca215e9a4be51c3446fe7112844ebd310f3ca69))
* added dockerfile for bun app ([#31](https://github.com/Darckfast/workers-go/issues/31)) ([1d20239](https://github.com/Darckfast/workers-go/commit/1d202393356ded7d5bd79050fa0fa4fd7ab4edf0))
* added signals support + perf work + more examples ([#30](https://github.com/Darckfast/workers-go/issues/30)) ([dceaaeb](https://github.com/Darckfast/workers-go/commit/dceaaeb87e4bb03dc6abb638e2940eebf97e601f))
* simplified dockerfiles ([#35](https://github.com/Darckfast/workers-go/issues/35)) ([52ab54b](https://github.com/Darckfast/workers-go/commit/52ab54bd2c26910c7e14a508281f6283b626822e))


### Bug Fixes

* fixed poor performance when using express ([#34](https://github.com/Darckfast/workers-go/issues/34)) ([9b83c2c](https://github.com/Darckfast/workers-go/commit/9b83c2c962564e29ff3ef1bf9edc8b15d60cf8e9))
* fixed stream copy when using express like lib (http/https) ([#33](https://github.com/Darckfast/workers-go/issues/33)) ([7150c8b](https://github.com/Darckfast/workers-go/commit/7150c8b6358b258e2f549c8fdaad5fd3c3485103))
* minor fixes to dockerfiles ([#32](https://github.com/Darckfast/workers-go/issues/32)) ([5b4bc6e](https://github.com/Darckfast/workers-go/commit/5b4bc6e898bcea8c0c2c34932279da882c5f5b54))

## [0.1.0](https://github.com/Darckfast/workers-go/compare/v0.0.4...v0.1.0) (2025-08-30)


### Features

* add d1 bindings using js ([#14](https://github.com/Darckfast/workers-go/issues/14)) ([4155842](https://github.com/Darckfast/workers-go/commit/41558426dfe50bd6f2e521ae2785730477fe82fa))
* add global types ([#12](https://github.com/Darckfast/workers-go/issues/12)) ([156d218](https://github.com/Darckfast/workers-go/commit/156d2185f74990563298908936cbcd5206238444))
* add image props on cf ([#27](https://github.com/Darckfast/workers-go/issues/27)) ([9587ada](https://github.com/Darckfast/workers-go/commit/9587adaedc5b277266d3726faf3ec77a20d9152f))
* replaced `encoding/json` with `easyjson` ([#22](https://github.com/Darckfast/workers-go/issues/22)) ([3b4bac4](https://github.com/Darckfast/workers-go/commit/3b4bac42b149086356c3f4b3d84355c30569e982))
* replaced KV `io.ReadAll` with streaming ([#11](https://github.com/Darckfast/workers-go/issues/11)) ([6d25eba](https://github.com/Darckfast/workers-go/commit/6d25eba6a1502d9364edd6b65fbf538fbd74b4b0))
* replaced kv methods ([#8](https://github.com/Darckfast/workers-go/issues/8)) ([4e2c8b8](https://github.com/Darckfast/workers-go/commit/4e2c8b89591b508d369bac3c82ea24543a63727b))
* updated main.ts to allow re-initialization on process exit ([#18](https://github.com/Darckfast/workers-go/issues/18)) ([2e643d7](https://github.com/Darckfast/workers-go/commit/2e643d7ddc6c0b5d6fa13ac3b6caa4776ea63443))


### Bug Fixes

* fixed devalue vuln ([#21](https://github.com/Darckfast/workers-go/issues/21)) ([3a25a92](https://github.com/Darckfast/workers-go/commit/3a25a9298ff40a6c02870c2e6a48ac1436a55850))
* fixed labeler permission to write labels ([#16](https://github.com/Darckfast/workers-go/issues/16)) ([11ef591](https://github.com/Darckfast/workers-go/commit/11ef59112bddc972cf3224b1d6b575763d445041))
* fixed package scripts ([#26](https://github.com/Darckfast/workers-go/issues/26)) ([9bed794](https://github.com/Darckfast/workers-go/commit/9bed7947da8170dba5de9d86ec61c96d10fab943))
* fixed tail handler not parsing events correctly ([#13](https://github.com/Darckfast/workers-go/issues/13)) ([85ab5da](https://github.com/Darckfast/workers-go/commit/85ab5dab6d26e866ac88017f68d51845bf24f402))
* fixed workflow permissions ([#10](https://github.com/Darckfast/workers-go/issues/10)) ([18f220f](https://github.com/Darckfast/workers-go/commit/18f220ff2dd004d0d99d000b569e3820ef181a46))


### Performance Improvements

* minor perf changes to `main.ts` file ([#15](https://github.com/Darckfast/workers-go/issues/15)) ([400e580](https://github.com/Darckfast/workers-go/commit/400e580aba64cf9c8d7958b2f9a15ded21142529))
