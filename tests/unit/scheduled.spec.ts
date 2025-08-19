import worker from "$wrk/main";
import {
  createExecutionContext,
  createScheduledController,
  env,
  waitOnExecutionContext,
} from "cloudflare:test";
import { beforeAll, describe, expect, it } from "vitest";

describe("scheduled()", () => {
	let cronResult;
	let time;

	beforeAll(async () => {
		const ctx = createExecutionContext();
		const controller = createScheduledController();
		time = controller.scheduledTime;
		await worker.scheduled(controller, env, ctx);

		const request = new Request("http://example.com/kv?key=cron:result");
		const response: Response = await worker.fetch(request, env, ctx);
		await waitOnExecutionContext(ctx);
		cronResult = await response.json();
	});

	it("should run scheduled task with no errors", () => {
		expect(cronResult).toHaveProperty("data");
		expect(cronResult.data).toEqual(`${time}`);
	});
});
