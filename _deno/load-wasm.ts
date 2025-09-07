import "./bin/wasm_exec.js";

// @ts-expect-error
globalThis.tryCatch = (fn) => {
  try {
    return { data: fn() };
  } catch (err) {
    if (!(err instanceof Error)) {
      if (err instanceof Object) {
        err = JSON.stringify(err);
      }

      // @ts-expect-error
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
