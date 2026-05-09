# Changelog

## [0.3.0](https://github.com/Darckfast/workers-go/compare/v0.2.0...v0.3.0) (2026-05-09)


### Features

* add bench and some perf changes ([de989e1](https://github.com/Darckfast/workers-go/commit/de989e1d2b0515630c4bd8c8a3d6902571fa8eb3))
* add containerFetch ([d917b93](https://github.com/Darckfast/workers-go/commit/d917b934791a006136bac14033b11a72c3378882))
* add d1 bindings using js ([#14](https://github.com/Darckfast/workers-go/issues/14)) ([4155842](https://github.com/Darckfast/workers-go/commit/41558426dfe50bd6f2e521ae2785730477fe82fa))
* add developing segment on readme ([d9ad4cf](https://github.com/Darckfast/workers-go/commit/d9ad4cf0cd8b3bab1aa2bd8f8c7f2c50db014b95))
* add email handler ([a9bd93f](https://github.com/Darckfast/workers-go/commit/a9bd93f5a39bfd7d0eba5f768e29d393205839a0))
* add global types ([#12](https://github.com/Darckfast/workers-go/issues/12)) ([156d218](https://github.com/Darckfast/workers-go/commit/156d2185f74990563298908936cbcd5206238444))
* add http client and transport interface ([38a8feb](https://github.com/Darckfast/workers-go/commit/38a8febde7685f9dc9d54e834a966b6a4b0d3747))
* add image props on cf ([#27](https://github.com/Darckfast/workers-go/issues/27)) ([9587ada](https://github.com/Darckfast/workers-go/commit/9587adaedc5b277266d3726faf3ec77a20d9152f))
* add linting and static check workflows ([67920d1](https://github.com/Darckfast/workers-go/commit/67920d10758ff76861bf329050a2e6dad0fcf27f))
* add new testing endpoint for r2 with encryption ([#28](https://github.com/Darckfast/workers-go/issues/28)) ([aca215e](https://github.com/Darckfast/workers-go/commit/aca215e9a4be51c3446fe7112844ebd310f3ca69))
* add results as string ([3d7adce](https://github.com/Darckfast/workers-go/commit/3d7adce273d2a73a713766f7d37054b0539c1f80))
* add some more perf funcs ([09a19dc](https://github.com/Darckfast/workers-go/commit/09a19dc9fbaaf2dd3083d48a3f3586c47294c43e))
* add tail handler ([daec390](https://github.com/Darckfast/workers-go/commit/daec3901cd675bedf3a5a3b9b12800744770d7dc))
* add woodpecker ci ([1199afa](https://github.com/Darckfast/workers-go/commit/1199afa26473347921dfd034a14b44409f519abe))
* add worker testing with better local dev experience ([f400a0f](https://github.com/Darckfast/workers-go/commit/f400a0f350d3160eff55bc90b7a6ae3870dab809))
* added dockerfile for bun app ([#31](https://github.com/Darckfast/workers-go/issues/31)) ([1d20239](https://github.com/Darckfast/workers-go/commit/1d202393356ded7d5bd79050fa0fa4fd7ab4edf0))
* added signals support + perf work + more examples ([#30](https://github.com/Darckfast/workers-go/issues/30)) ([dceaaeb](https://github.com/Darckfast/workers-go/commit/dceaaeb87e4bb03dc6abb638e2940eebf97e601f))
* auto load env into go envs ([4f3f3f5](https://github.com/Darckfast/workers-go/commit/4f3f3f558bec52dbf8a4093d9906c487b14b8606))
* changed cron to pass the event instead of ctx ([d177f53](https://github.com/Darckfast/workers-go/commit/d177f53099e705f8b960d0ba43d67e283bfffcbd))
* embeded tryCatch fn with normalized errors ([871d6fb](https://github.com/Darckfast/workers-go/commit/871d6fb972aa3cef38e618d637feb086285b75ae))
* export `ctx` and `env` ([f837285](https://github.com/Darckfast/workers-go/commit/f837285534555be57d1d861de8862e800657f3d8))
* export Env and Ctx ([c6cbdc0](https://github.com/Darckfast/workers-go/commit/c6cbdc0bfcd397fe8d0ffac708c3524d8eaf80c7))
* fetch now returns plain http.Response and consumes http.Request ([fd055c0](https://github.com/Darckfast/workers-go/commit/fd055c0c2bc9499055019166eab90c137129239e))
* implemented options with metadata for kv ([7e1066f](https://github.com/Darckfast/workers-go/commit/7e1066f08f4cb6c1989e1d1429a81a3a76586bd6))
* improved some streams and convs ([25aa0b4](https://github.com/Darckfast/workers-go/commit/25aa0b414ed4ee19ed3075309a804edf04208c53))
* more kv methods and fixes ([9560698](https://github.com/Darckfast/workers-go/commit/9560698dace4c6b50dd60bf9d6562be98e6a7567))
* more reworks, perf changes and testing ([d9816e0](https://github.com/Darckfast/workers-go/commit/d9816e08203bd3f07b3fc18b03558ee02e285cb6))
* moved readme to its own page ([2586150](https://github.com/Darckfast/workers-go/commit/2586150828579496a0ee6e01522e6f7edd4c2d5d))
* replace fmt.Errof with errors.New due its size ([99fc7f5](https://github.com/Darckfast/workers-go/commit/99fc7f541e8e256b5cf3d5cfb03e32f166ae3bbf))
* replaced `encoding/json` with `easyjson` ([#22](https://github.com/Darckfast/workers-go/issues/22)) ([3b4bac4](https://github.com/Darckfast/workers-go/commit/3b4bac42b149086356c3f4b3d84355c30569e982))
* replaced KV `io.ReadAll` with streaming ([#11](https://github.com/Darckfast/workers-go/issues/11)) ([6d25eba](https://github.com/Darckfast/workers-go/commit/6d25eba6a1502d9364edd6b65fbf538fbd74b4b0))
* replaced kv methods ([#8](https://github.com/Darckfast/workers-go/issues/8)) ([4e2c8b8](https://github.com/Darckfast/workers-go/commit/4e2c8b89591b508d369bac3c82ea24543a63727b))
* revamp some APIs and add missing handlers ([f3e501e](https://github.com/Darckfast/workers-go/commit/f3e501edb141fbde5a02a4a54822fc19a140dc41))
* reworked http client interface ([4c1948e](https://github.com/Darckfast/workers-go/commit/4c1948efe52c2977592a5360b27ce23016f19a41))
* simplified dockerfiles ([#35](https://github.com/Darckfast/workers-go/issues/35)) ([52ab54b](https://github.com/Darckfast/workers-go/commit/52ab54bd2c26910c7e14a508281f6283b626822e))
* updated `README.md` and renamed module ([e3324c1](https://github.com/Darckfast/workers-go/commit/e3324c124f17f8f86d91ddcb3c6e72ab61da7f50))
* updated main.ts to allow re-initialization on process exit ([#18](https://github.com/Darckfast/workers-go/issues/18)) ([2e643d7](https://github.com/Darckfast/workers-go/commit/2e643d7ddc6c0b5d6fa13ac3b6caa4776ea63443))


### Bug Fixes

* add install cmd on workflow ([e3b6825](https://github.com/Darckfast/workers-go/commit/e3b682534b9664347458ff970dfd8ef4af4f1ab3))
* fixed code exmaples ([94415d6](https://github.com/Darckfast/workers-go/commit/94415d67ac45f6ce86e9a13621e65b0798807faa))
* fixed Date conv ([a10e2ba](https://github.com/Darckfast/workers-go/commit/a10e2ba0243708eae54b911f3462d391887ec337))
* fixed deadlock on cache while passing a Response from fetch ([540239a](https://github.com/Darckfast/workers-go/commit/540239a57c1e0925b09dc538d2d9d774516886a8))
* fixed devalue vuln ([#21](https://github.com/Darckfast/workers-go/issues/21)) ([3a25a92](https://github.com/Darckfast/workers-go/commit/3a25a9298ff40a6c02870c2e6a48ac1436a55850))
* fixed envs not being loaded ([8b7d009](https://github.com/Darckfast/workers-go/commit/8b7d0096ca78c6c5be5adf57cdbcef2008a5a7dc))
* fixed fmt not printing ([2cfeb36](https://github.com/Darckfast/workers-go/commit/2cfeb36bc6cc0e4f20741788b04e9020938e9e88))
* fixed int and float conversion ([803161f](https://github.com/Darckfast/workers-go/commit/803161f65bfec1727e7fbff404aa84e645a3153b))
* fixed int conversion and changes not returning ([aab3e32](https://github.com/Darckfast/workers-go/commit/aab3e329a7af76fca0b56988d11425a04d3e7deb))
* fixed internals being exported ([999a455](https://github.com/Darckfast/workers-go/commit/999a4550716249a6854a3556804b53f7e8ec17bb))
* fixed labeler permission to write labels ([#16](https://github.com/Darckfast/workers-go/issues/16)) ([11ef591](https://github.com/Darckfast/workers-go/commit/11ef59112bddc972cf3224b1d6b575763d445041))
* fixed multipart test ([d337444](https://github.com/Darckfast/workers-go/commit/d33744491dde4c9341701b4ff1028c7beda2b1a9))
* fixed package scripts ([#26](https://github.com/Darckfast/workers-go/issues/26)) ([9bed794](https://github.com/Darckfast/workers-go/commit/9bed7947da8170dba5de9d86ec61c96d10fab943))
* fixed poor performance when using express ([#34](https://github.com/Darckfast/workers-go/issues/34)) ([9b83c2c](https://github.com/Darckfast/workers-go/commit/9b83c2c962564e29ff3ef1bf9edc8b15d60cf8e9))
* fixed producer test ([7145f0b](https://github.com/Darckfast/workers-go/commit/7145f0b1d43334b6114e0addb3a23ea6cab2e470))
* fixed scheduledTime losing its precision ([4ed7b85](https://github.com/Darckfast/workers-go/commit/4ed7b853afdedafbfe9efd7a4b0ca451da4af8cd))
* fixed stream copy when using express like lib (http/https) ([#33](https://github.com/Darckfast/workers-go/issues/33)) ([7150c8b](https://github.com/Darckfast/workers-go/commit/7150c8b6358b258e2f549c8fdaad5fd3c3485103))
* fixed tail handler not parsing events correctly ([#13](https://github.com/Darckfast/workers-go/issues/13)) ([85ab5da](https://github.com/Darckfast/workers-go/commit/85ab5dab6d26e866ac88017f68d51845bf24f402))
* fixed testing function not loading ([4198a37](https://github.com/Darckfast/workers-go/commit/4198a37e34405ad99941c49714e0322ec495345c))
* fixed trycatch not being set on init ([5ca0947](https://github.com/Darckfast/workers-go/commit/5ca0947ca1b4049d0ca7da5fd85fb9c03429d631))
* fixed trycatch not loading on tests ([d0d3ddf](https://github.com/Darckfast/workers-go/commit/d0d3ddf7c1d5ee8d87088936aae059f0cd42a0a5))
* fixed types not being loaded ([fea647f](https://github.com/Darckfast/workers-go/commit/fea647f419a1633261eb8ae944e0a90a1cba0796))
* fixed typing on tail ([d5eafbb](https://github.com/Darckfast/workers-go/commit/d5eafbb6f5dc505477080c3a663cdb679cbad1d1))
* fixed workflow export path ([a5da831](https://github.com/Darckfast/workers-go/commit/a5da831dffa1cf15632f2c650e4f5ecb540d35da))
* fixed workflow permissions ([#10](https://github.com/Darckfast/workers-go/issues/10)) ([18f220f](https://github.com/Darckfast/workers-go/commit/18f220ff2dd004d0d99d000b569e3820ef181a46))
* fixed workflow pnpm version ([3e08fc8](https://github.com/Darckfast/workers-go/commit/3e08fc832d7da83cbf5822a2671534c29518ccca))
* fixed wrong arg to func ([35f4bd8](https://github.com/Darckfast/workers-go/commit/35f4bd82c8356c12b0646ee82e514ac91a10481d))
* ignore path on wasm ([ad6dcf2](https://github.com/Darckfast/workers-go/commit/ad6dcf25809c3ba9247d072f507e80a35b273942))
* minor fixes to dockerfiles ([#32](https://github.com/Darckfast/workers-go/issues/32)) ([5b4bc6e](https://github.com/Darckfast/workers-go/commit/5b4bc6e898bcea8c0c2c34932279da882c5f5b54))
* set auto by default ([3213d02](https://github.com/Darckfast/workers-go/commit/3213d026278185b1fa93c232788aed1ddd01e33b))
* tmp fix to a test ([92e9ff9](https://github.com/Darckfast/workers-go/commit/92e9ff9106288f5cbc0c83cdda12ec375f29d987))


### Performance Improvements

* minor perf changes to `main.ts` file ([#15](https://github.com/Darckfast/workers-go/issues/15)) ([400e580](https://github.com/Darckfast/workers-go/commit/400e580aba64cf9c8d7958b2f9a15ded21142529))

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
