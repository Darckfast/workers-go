import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import "./bin/wasm_exec.js";

globalThis.tryCatch = (fn) => {
  try {
    return { data: fn() };
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
    const CURRENT_DIR = dirname(fileURLToPath(import.meta.url));
    const app = readFileSync(resolve(CURRENT_DIR, "./bin/app.wasm"));
    const { instance } = await WebAssembly.instantiate(app, go.importObject);

    go.run(instance).finally(() => {
      initiliazed = false;
    });
    initiliazed = true;
  }
}
