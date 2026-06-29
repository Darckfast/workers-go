import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { Readable } from "node:stream";
import { fileURLToPath } from "node:url";
import express from "ultimate-express";
import "./bin/wasm_exec.js";

const CURRENT_DIR = dirname(fileURLToPath(import.meta.url));
const binary = readFileSync(resolve(CURRENT_DIR, "./bin/app.wasm"));
const go = new Go();
let initiliazed = false;

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

export async function init() {
  if (go.exited) {
    initiliazed = false;
  }

  if (!initiliazed) {
    const { instance } = await WebAssembly.instantiate(binary, go.importObject);

    go.run(instance).finally(() => {
      initiliazed = false;
    });
    initiliazed = true;
  }
}

// warm up the init
init();

const app = express();

// the body must be raw, the parsing will happen within
// our Go code
app.use(express.raw({ type: "*/*" }));

let encoder = new TextEncoder()
let decoder = new TextDecoder()

app.all("/", async (_req, res) => {
  // if the first request comes before the initialization
  // is done, or if the Go process exit, this will guarantee
  // it's up and running
  await init();
  let stream = Readable.toWeb(_req)
  let { writable, readable } = new TransformStream();
  let keys = Object.keys(_req.headers);
  let selHeaders = "";
  for (let i = 0; i < keys.length; i++) {
    const key = keys[i];
    const value = _req.headers[key];
    selHeaders += `${key}: ${value}\n`;
  }
  let rawHeaders = await cf.fetch(
    stream,
    encoder.encode(_req.method),
    encoder.encode(_req.url),
    encoder.encode(selHeaders),
    writable,
    _req.signal,
  );
  let parts = decoder.decode(rawHeaders).split("\n");
  let [, status, statusText] = parts[0].split(" ");
  let headers = new Headers();
  for (let i = 1; i < parts.length; i++) {
    let [key, val] = parts[i].split(":");
    if (key) {
      headers.append(key, val);
    }
  }

  res.writeHead(
    status,
    statusText,
    Object.fromEntries(headers),
  );

  Readable.fromWeb(readable).pipe(res);
});

app.listen(5173, () => {
  console.log("Server is running on port 5173");
});

export default app
