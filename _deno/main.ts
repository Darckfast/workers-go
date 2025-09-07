import { init } from "./load-wasm.ts";

init();

Deno.serve({ port: 5173 }, async (_req) => {
  await init();

  return cf.fetch(_req);
});
