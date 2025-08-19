import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it, test } from "vitest";

describe("d1", () => {
	describe("create entry", () => {
		let body;
		let id;

		beforeAll(async () => {
			const r = await SELF.fetch("http://service/d1", {
				method: "POST",
				body: "my test data",
			}).then((r) => r.json());

			id = r.data.id;
			body = await SELF.fetch(`http://service/d1?id=${id}`, {
				method: "GET",
			}).then((r) => r.json());
		});

		it("should return entry with id", async () => {
			expect(body.data).toHaveProperty("id", id);
		});

		it("should return entry with data", async () => {
			expect(body.data).toHaveProperty("data", "my test data");
		});

		it("should return entry with created_at", async () => {
			expect(body.data).toHaveProperty("created_at");
		});

		it("should return entry with updated_at", async () => {
			expect(body.data).toHaveProperty("updated_at");
		});
	});

	describe("delete entry", () => {
		let body;
		let id;

		beforeAll(async () => {
			const r = await SELF.fetch("http://service/d1", {
				method: "POST",
				body: "my test data",
			}).then((r) => r.json());
			id = r.data.id;

			await SELF.fetch(`http://service/d1?id=${id}`, {
				method: "DELETE",
			}).then((r) => r.json());

			body = await SELF.fetch(`http://service/d1?id=${id}`, {
				method: "GET",
			}).then((r) => r.json());
		});

		it("should return error", async () => {
			expect(body).toHaveProperty("error");
		});
	});

	describe("update entry", () => {
		let body;
		let id;

		beforeAll(async () => {
			const r = await SELF.fetch("http://service/d1", {
				method: "POST",
				body: "my test data",
			}).then((r) => r.json());
			id = r.data.id;
			await SELF.fetch("http://service/d1?id=" + id, {
				method: "PUT",
				body: "my new test data",
			}).then((r) => r.json());

			body = await SELF.fetch(`http://service/d1?id=${id}`, {
				method: "GET",
			}).then((r) => r.json());
		});

		it("should return entry with id", async () => {
			expect(body.data).toHaveProperty("id", id);
		});

		it("should return entry with data", async () => {
			expect(body.data).toHaveProperty("data", "my new test data");
		});

		it("should return entry with created_at", async () => {
			expect(body.data).toHaveProperty("created_at");
		});

		it("should return entry with updated_at", async () => {
			expect(body.data).toHaveProperty("updated_at");
		});
	});
});
