import { faker } from "@faker-js/faker";
import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it } from "vitest";

describe("fetch", () => {
	describe("POST /application/x-www-form-urlencoded ", () => {
		let rs: Response;
		let body;
		let params: URLSearchParams;

		beforeAll(async () => {
			params = new URLSearchParams();

			params.set("id", faker.string.uuid());
			params.set("alpha", faker.string.alphanumeric());
			params.set("url", faker.image.url());
			params.set("fullname", faker.person.fullName());
			params.set("number", faker.number.bigInt().toString());

			rs = await SELF.fetch("https://api/application/x-www-form-urlencoded ", {
				method: "POST",
				body: params,
			});
			body = await rs.json();
		});

		it("should have returned 200 status code", () => {
			expect(rs.status).toBe(200);
		});

		it.each(["id", "alpha", "url", "fullname", "number"])(
			"should return %s",
			(key) => {
				expect(body[key]).toBe(params.get(key));
			},
		);
	});
});
