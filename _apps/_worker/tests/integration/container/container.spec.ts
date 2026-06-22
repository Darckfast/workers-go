import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it, test } from "vitest";

// cloudflare vitest plugin currently does not support containers
describe.skip("container", () => {
	describe("fetch container", () => {
		let rs: Response;
		let body;

		beforeAll(async () => {
			rs = await SELF.fetch("http://service/container", {
				method: "GET",
			});

			body = await rs.json();
		});

		it("should have returned 200 status code", () => {
			expect(rs.status).toBe(200);
		});

		it("should have key count", async () => {
			expect(body).toHaveProperty("result", "200 OK");
		});
	});
});
