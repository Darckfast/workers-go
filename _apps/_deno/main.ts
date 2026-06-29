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
globalThis.tryCatch = function (o, fn, args) {
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
};

const go = new Go();
let initiliazed = false;

export async function init() {
  if (go.exited) {
    initiliazed = false;
  }

  if (!initiliazed) {
    const app = await Deno.readFile("./bin/app.wasm");
    const { instance } = await WebAssembly.instantiate(app, go.importObject);

    go.run(instance).finally(() => {
      initiliazed = false;
    });
    initiliazed = true;
  }
}

init();
let decoder = new TextDecoder();
let encoder = new TextEncoder();

export default Deno.serve({ port: 5173 }, async (_req) => {
  await init();
  let { writable, readable } = new TransformStream();
  let keys = Object.keys(_req.headers);
  let selHeaders = "";
  for (let i = 0; i < keys.length; i++) {
    const key = keys[i];
    const value = _req.headers.get(key);
    selHeaders += `${key}: ${value}\n`;
  }
  let rawHeaders = await cf.fetch(
    _req.body,
    encoder.encode(_req.method),
    encoder.encode(_req.url),
    encoder.encode(selHeaders),
    writable,
    _req.signal,
  );
  let parts = decoder.decode(rawHeaders).split("\n");
  let [, status] = parts[0].split(" ");
  let headers = new Headers();
  for (let i = 1; i < parts.length; i++) {
    let [key, val] = parts[i].split(":");
    if (key) {
      headers.append(key, val);
    }
  }
  return new Response(readable, { status, headers });
});
