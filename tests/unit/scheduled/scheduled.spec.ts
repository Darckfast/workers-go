import worker from '$wrk/main';
import {
  createExecutionContext,
  createScheduledController,
  env,
  waitOnExecutionContext,
} from "cloudflare:test";
import { describe, it } from "vitest";

describe("scheduled", () => {
  it("should run with no errors", async () => {
    const ctx = createExecutionContext();
    const controller = createScheduledController()
    await worker.scheduled(controller, env, ctx);
    await waitOnExecutionContext(ctx);
  });
});
