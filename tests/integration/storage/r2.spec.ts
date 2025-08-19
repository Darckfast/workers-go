import { faker } from "@faker-js/faker";
import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it, test } from "vitest";

describe("r2", () => {
	describe("write op", () => {
		let rs: Response;
		let params: URLSearchParams;
		let body;

		beforeAll(async () => {
			params = new URLSearchParams();

			params.set("b64", faker.image.dataUri().split("base64,")[1]);

			rs = await SELF.fetch("http://service/r2", {
				method: "POST",
				body: params,
			});

			body = await rs.json();
		});

		it("should have returned 200 status code", () => {
			expect(rs.status).toBe(200);
		});

		it('should have content-type "application/json"', () => {
			expect(rs.headers.get("content-type")).toBe("application/json");
		});

		it("should have no error", async () => {
			expect(body).toHaveProperty("has_error", false);
		});
		it("should have key count", async () => {
			expect(body).toHaveProperty("result");
			expect(body.result.Key).toBe("count");
		});
	});

	describe("read op", () => {
		let rs: Response;
		let body;
		let params;

		beforeAll(async () => {
			params = new URLSearchParams();
			params.set("b64", btoa(faker.lorem.paragraph()));

			await SELF.fetch("http://service/r2", {
				method: "POST",
				body: params,
			});

			rs = await SELF.fetch("http://service/r2", {
				method: "GET",
			});

			body = await rs.json();
		});

		it("should have returned 200 status code", () => {
			expect(rs.status).toBe(200);
		});

		it('should have content-type "application/json"', () => {
			expect(rs.headers.get("content-type")).toBe("application/json");
		});

		it("should have body b64", async () => {
			expect(body).toHaveProperty("body", params.get("b64"));
		});
	});
});
