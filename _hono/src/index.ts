import { serve } from '@hono/node-server';
import { Hono } from 'hono';
import { init } from './load-wasm.js';

init()

const app = new Hono()

app.get('*', async (c) => {
  await init()
  return cf.fetch(c.req.raw)
})

serve({
  fetch: app.fetch,
  port: 3000
}, (info) => {
  console.log(`Server is running on http://localhost:${info.port}`)
})
