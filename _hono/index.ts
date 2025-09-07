import { serve } from "@hono/node-server";
import { readFileSync } from "fs";
import { Hono } from "hono";
import { dirname, resolve } from "path";
import { fileURLToPath } from "url";
import "./bin/wasm_exec.js";

// @ts-ignore
globalThis.tryCatch = (fn) => {
  try {
    return { data: fn() };
  } catch (err) {
    if (!(err instanceof Error)) {
      if (err instanceof Object) {
        err = JSON.stringify(err);
      }

      // @ts-ignore
      err = new Error(err || "no error message");
    }

    return { error: err };
  }
};

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

app.all("*", async (c) => {
  await init();
  return cf.fetch(c.req.raw);
});

serve({
  fetch: app.fetch,
  port: 5173,
}, (info) => {
  console.log(`Server is running on http://localhost:${info.port}`);
});
