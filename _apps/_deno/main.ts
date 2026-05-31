import "./bin/wasm_exec.js";

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

Deno.serve({ port: 5173 }, async (_req) => {
  await init();

  return cf.fetch(_req);
});
