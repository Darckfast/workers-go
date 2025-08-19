import worker from "$wrk/main";
import {
	createExecutionContext,
	env,
	waitOnExecutionContext,
} from "cloudflare:test";
import { describe, expect, it } from "vitest";

// For now, you'll need to do something like this to get a correctly-typed
// `Request` to pass to `worker.fetch()`.
const IncomingRequest = Request;

describe("fetch handler", () => {
	it("should return JSON body", async () => {
		const request = new IncomingRequest("http://example.com/");
		// Create an empty context to pass to `worker.fetch()`
		const ctx = createExecutionContext();
		const response: Response = await worker.fetch(request, env, ctx);
		// Wait for all `Promise`s passed to `ctx.waitUntil()` to settle before running test assertions
		await waitOnExecutionContext(ctx);
		expect(response.status).toBe(404);
	});
});
