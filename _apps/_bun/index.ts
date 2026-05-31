import "./bin/wasm_exec.js";

/**
 * A REQUIRED function, since errors thrown within the JS runtime
 * inside Go's will cause the process to exit
 *
 * It's just a try...catch with error normalization
 *
 * This cannot be initialized within Go code, due Cloudflare workers
 * limits
 */
globalThis.tryCatch = function(o, fn, args) {
  try {
    if (fn) {
      return { data: o[fn](...args) };
    }

    return { data: o(...args) };
  } catch (err) {
    if (!(err instanceof Error)) {
      if (err instanceof Object) {
        err = JSON.stringify(err);
      }
      err = new Error(err || "no error message");
    }
    return { error: err };
  }
}

const go = new Go();
let initiliazed = false;
const app = await Bun.file("./bin/app.wasm").arrayBuffer();

export async function init() {
  if (go.exited) {
    initiliazed = false;
  }

  if (!initiliazed) {
    const { instance } = await WebAssembly.instantiate(app, go.importObject);

    go.run(instance).finally(() => {
      initiliazed = false;
    });
    initiliazed = true;
  }
}

init();

const server = Bun.serve({
  port: 5173,
  fetch: async (req) => {
    await init();
    return cf.fetch(req);
  },
});

console.log(`Listening on http://localhost:${server.port} ...`);
