import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it } from "vitest";

describe("fetch", () => {
	describe("GET /application/json should return JSON", () => {
		let rs: Response;
		beforeAll(async () => {
			rs = await SELF.fetch("https://example.com/application/json");
		});

		it("should have returned 200 status code", () => {
			expect(rs.status).toBe(200);
		});

		it('should have content-type "application/json"', () => {
			expect(rs.headers.get("content-type")).toBe("application/json");
		});

		it('should have body {"vitest":true}', async () => {
			expect(await rs.text()).toBe(`{"vitest":true}\n`);
		});
	});

	describe("POST /application/json should return JSON with payload info", () => {
		const testid = crypto.randomUUID();
		let rs: Response;
		let rsbody;

		beforeAll(async () => {
			const body = {
				bool: true,
				number: 1,
				string: "my super string",
				list: [1, 2, 3, 4, 5, 6],
			};

			rs = await SELF.fetch(
				"http://service/application/json?id=0.5385553010283278&uuid=6430d5f6-a1a8-48fb-9fe3-747c1d5d9ecb",
				{
					method: "POST",
					headers: {
						"x-test-id": testid,
						"content-type": "application/json",
					},
					body: JSON.stringify(body),
				},
			);
			rsbody = await rs.json();
		});

		it("should have returned 200 status code", () => {
			expect(rs.status).toBe(200);
		});

		it.each([
			["size", 84],
			[
				"query",
				"id=0.5385553010283278&uuid=6430d5f6-a1a8-48fb-9fe3-747c1d5d9ecb",
			],
			["header", testid],
			[
				"raw",
				'{"bool":true,"list":[1,2,3,4,5,6],"number":1,"string":"my super string"}',
			],
		])("should have property %s with value %s", async (prop, expected) => {
			expect(rsbody).toHaveProperty(prop, expected);
		});
	});
});
