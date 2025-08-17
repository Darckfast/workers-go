import { connect } from 'cloudflare:sockets';
import app from "./bin/app.wasm";
import "./bin/wasm_exec.js";

export { GoContainer } from './durable_objects/go_container';
export { TestDurableObject } from './durable_objects/test';

globalThis.cf = {
  connect
}

globalThis.tryCatch = (fn) => {
  try {
    return {
      result: fn(),
    };
  } catch (error) {
    if (!(error instanceof Error)) {
      if (error instanceof Object) {
        error = JSON.stringify(error)
      }

      error = new Error(error || 'no error message')
    }
    return {
      error,
    };
  }

}

function init() {
  const go = new Go()
  go.run(new WebAssembly.Instance(app, go.importObject))
}

async function fetch(req: Request, env: Env, ctx: ExecutionContext) {
  init()
  return await globalThis.cf.fetch(req, env, ctx);
}

async function email(message: ForwardableEmailMessage, env: Env, ctx: ExecutionContext) {
  init()
  return await globalThis.cf.email(message, env, ctx)
}

async function scheduled(controler: ScheduledController, env: Env, ctx: ExecutionContext) {
  init()
  return await globalThis.cf.scheduled(controler, env, ctx)
}

async function queue(batch: MessageBatch, env: Env, ctx: ExecutionContext) {
  init()
  return await globalThis.cf.queue(batch, env, ctx)
}

async function tail(events: TailEvent[], env: Env, ctx: ExecutionContext) {
  init()
  return await globalThis.cf.tail(events, env, ctx)
}

export default {
  fetch,
  email,
  scheduled,
  queue,
  tail,
};
