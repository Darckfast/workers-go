import app from "./bin/app.wasm";
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

let initiliazed = false;

let go = new Go();
let instance = new WebAssembly.Instance(app, go.importObject);

function init() {
  if (!initiliazed) {
    go.run(instance).finally(() => {
      initiliazed = false;
      instance = new WebAssembly.Instance(app, go.importObject);
    });
    initiliazed = true;
  }

  if (go.exited) {
    go = new Go();
    go.run(instance).finally(() => {
      instance = new WebAssembly.Instance(app, go.importObject);
    });
  }
}

async function queue(batch: MessageBatch, env: Env, ctx: ExecutionContext) {
  init();
  return await cf.queue(batch, env, ctx);
}

export default {
  queue,
} satisfies ExportedHandler<Env>;
