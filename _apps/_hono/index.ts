import { serve } from "@hono/node-server";
import { readFileSync } from "fs";
import { Hono } from "hono";
import { dirname, resolve } from "path";
import { fileURLToPath } from "url";
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

let go = new Go();
let initiliazed = false;

export async function init() {
  if (go.exited) {
    initiliazed = false;
    go = new Go();
  }

  if (!initiliazed) {
    const CURRENT_DIR = dirname(fileURLToPath(import.meta.url));
    const app = readFileSync(resolve(CURRENT_DIR, "./bin/app.wasm"));
    let { instance } = await WebAssembly.instantiate(app, go.importObject);

    go.run(instance).finally(() => {
      initiliazed = false;
    });

    initiliazed = true;
  }
}

//warm up init
init();

const app = new Hono();
let encoder = new TextEncoder()
let decoder = new TextDecoder()

app.all("*", async (c) => {
  await init();
  let { writable, readable } = new TransformStream();
  let keys = Object.keys(c.req.raw.headers);
  let selHeaders = "";
  for (let i = 0; i < keys.length; i++) {
    const key = keys[i];
    const value = c.req.raw.headers.get(key);
    selHeaders += `${key}: ${value}\n`;
  }
  let rawHeaders = await cf.fetch(
    c.body,
    encoder.encode(c.req.method),
    encoder.encode(c.req.url),
    encoder.encode(selHeaders),
    writable,
    c.req.raw.signal,
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

export default serve({
  fetch: app.fetch,
  port: 5173,
}, (info) => {
  console.log(`Server is running on http://localhost:${info.port}`);
});

