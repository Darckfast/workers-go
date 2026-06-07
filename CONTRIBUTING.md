# Developing

First install [mise](https://mise.jdx.dev/getting-started.html) then run:

```bash
# you will have to `mise trust` on your first time running this code
mise //_apps/_worker:dev
```

*Access or cURL http://localhost:5173/env*

`_apps/_worker` is the only app fully integrated with `workers-go` lib

At the moment, there is no live-reload for changes outside the workers own dir, meaning every change made inside the lib itself will require a recompile before being usable in dev

# Issues

Feel free to open issues on either Github or Codeberg
