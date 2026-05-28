import { connect } from "cloudflare:sockets";
import { WorkerEntrypoint } from "cloudflare:workers";
import app from "./bin/app.wasm";
import "./bin/wasm_exec.js";

export { GoContainer } from "./durable_objects/go_container";
export { TestDurableObject } from "./durable_objects/test";

/*
 * This import is only used with the sockets lib in Go
 */
globalThis.cf = {
  connect,
};

/**
 * A REQUIRED nice to have function, since errors thrown within the JS runtime
 * inside Go's will cause the process to exit
 *
 * It's just a try...catch with error normalization
 */
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
/**
 * This function is what initialize your Go's compiled WASM binary
 * only after this function has finished, that the handlers will be
 * defined in the globalThis scope
 *
 * At the moment, due limitations with the getRandomData(), it
 * cannot be executed at top level, it must be contained within the handlers
 * scope
 *
 * It's REQUIRED and needs to be called before using the globalThis.cf.<handler>()
 */
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

/**
 * All handlers are available in this template, but feel free to keep
 * only the ones that will be used
 */
export default class extends WorkerEntrypoint {
  constructor(ctx, env) {
    super(ctx, env);

    globalThis.workerapp = this;
    init();

    // Required to make RPC stubs available
    const prototype = Object.getPrototypeOf(this);
    for (const [k, v] of Object.entries(cf.rpc)) {
      prototype[k] = v;
    }
  }

  async email(message: ForwardableEmailMessage) {
    return await cf.email(message, this.env, this.ctx);
  }

  async scheduled(controler: ScheduledController) {
    return await cf.scheduled(controler, this.env, this.ctx);
  }

  async queue(batch: MessageBatch) {
    return await cf.queue(batch, this.env, this.ctx);
  }

  async tail(events: TraceItem[]) {
    return await cf.tail(events, this.env, this.ctx);
  }

  async fetch(request: Request): Response | Promise<Response> {
    if (request.url.endsWith("rpc")) {
      const data = await this.echo(new Uint8Array(await request.arrayBuffer()));
      return new Response(data);
    } else {
      return await cf.fetch(request, this.env, this.ctx);
    }
  }
}
