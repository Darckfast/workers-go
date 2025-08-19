import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it, vi } from "vitest";

describe("produces and consumers queue message", () => {
	let res;
	beforeAll(async () => {
		res = await SELF.fetch("https://example.com/queue", {
			method: "POST",
			body: "value",
		});
	});

	it("should return status code 202", () => {
		expect(res.status).toBe(202);
	});

	it("should consume the message", async () => {
		const result = await vi.waitUntil(
			async () => {
				const response = await SELF.fetch("https://example.com/queue");
				const text = await response.json();
				if (response.ok) return text;
			},
			{ timeout: 7000 },
		);

		expect(result).toHaveProperty("result", "VALUE");
	});
});
