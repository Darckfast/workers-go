import { init } from "./load-wasm.ts";

init();
const server = Bun.serve({
  port: 5173,
  fetch: async (req) => {
    await init();
    return cf.fetch(req);
  },
});

console.log(`Listening on http://localhost:${server.port} ...`);

