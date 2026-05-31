import "./bin/wasm_exec.js";

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
