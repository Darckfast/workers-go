# Changelog

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
