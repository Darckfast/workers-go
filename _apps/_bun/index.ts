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
