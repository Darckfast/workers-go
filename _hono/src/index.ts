import { serve } from "@hono/node-server";
import { Hono } from "hono";
import { init } from "./load-wasm.js";

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
