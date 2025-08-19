import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it, test } from "vitest";

describe("kv", () => {
	describe("delete op", () => {
		let rs: Response;

		beforeAll(async () => {
			rs = await SELF.fetch("http://service/kv", {
				method: "DELETE",
			});
		});

		it("should have returned 200 status code", () => {
			expect(rs.status).toBe(200);
		});

		it('should have content-type "application/json"', () => {
			expect(rs.headers.get("content-type")).toBe("application/json");
		});

		it("should have no error", async () => {
			expect(await rs.json()).toHaveProperty("has_error", false);
		});
	});

	describe("write op", () => {
		let rs: Response;

		beforeAll(async () => {
			rs = await SELF.fetch("http://service/kv", {
				method: "POST",
			});
		});

		it("should have returned 200 status code", () => {
			expect(rs.status).toBe(200);
		});

		it('should have content-type "application/json"', () => {
			expect(rs.headers.get("content-type")).toBe("application/json");
		});

		it('should return count "1"', async () => {
			const rbody = await rs.json();
			expect(rbody).toHaveProperty("count", "1");
		});
	});
});
