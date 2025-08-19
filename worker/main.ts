import { connect } from "cloudflare:sockets";
import { catchThis } from "catch-this";
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
 * A REQUIRED nice to have lib, since errors thrown within the JS runtime
 * inside Go's will cause the process to exit
 *
 * This functions is just a try...catch with error normalization, that work both
 * with sync (callback) and async (promises) functions
 */
globalThis.catchThis = catchThis.auto;

/**
 * This function is what initialize your Go's compiled WASM binary
 * only after this function has finished, that the handlers will be
 * defined in the globalThis scope
 *
 * At the moment, due limitations with the getRandomData(), this block
 * cannot be executed at top level, it must be contained within the handlers
 * scope
 *
 * It's REQUIRED and needs to be called before using the globalThis.cf.<handler>()
 */
function init() {
	const go = new Go();
	go.run(new WebAssembly.Instance(app, go.importObject));
}

async function fetch(req: Request, env: Env, ctx: ExecutionContext) {
	init();
	return await globalThis.cf.fetch(req, env, ctx);
}

async function email(
	message: ForwardableEmailMessage,
	env: Env,
	ctx: ExecutionContext,
) {
	init();
	return await globalThis.cf.email(message, env, ctx);
}

async function scheduled(
	controler: ScheduledController,
	env: Env,
	ctx: ExecutionContext,
) {
	init();
	return await globalThis.cf.scheduled(controler, env, ctx);
}

async function queue(batch: MessageBatch, env: Env, ctx: ExecutionContext) {
	init();
	return await globalThis.cf.queue(batch, env, ctx);
}

async function tail(events: TraceItem[], env: Env, ctx: ExecutionContext) {
	init();
	return await globalThis.cf.tail(events, env, ctx);
}

/**
 * All handlers are available in this template, but feel free to keep
 * only the ones that will be used
 */
export default {
	fetch,
	email,
	scheduled,
	queue,
	tail,
} satisfies ExportedHandler<Env>;
