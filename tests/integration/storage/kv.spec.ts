import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it, test } from "vitest";

describe("kv", () => {
	describe("delete kv entry", () => {
		let rs: Response;

		beforeAll(async () => {
			await SELF.fetch("http://service/kv?key=count", {
				method: "POST",
				body: JSON.stringify({ test: true }),
			}).then((r) => r.text());

			rs = await SELF.fetch("http://service/kv?key=count", {
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
			expect(await rs.json()).toHaveProperty("error", null);
		});
	});

	describe("create kv entry", () => {
		let rs: Response;

		beforeAll(async () => {
			rs = await SELF.fetch("http://service/kv?key=entry", {
				method: "POST",
				body: JSON.stringify({
					$schema: "https://www.schemastore.org/tsconfig.json",
					compilerOptions: {
						target: "esnext",
						module: "esnext",
						lib: ["esnext", "dom"],
						moduleResolution: "bundler",
						types: [
							"node",
							"@cloudflare/workers-types",
							"@cloudflare/vitest-pool-workers",
						],
						paths: {
							"$wrk/*": ["./worker/*"],
						},
					},
					include: ["./**/*.ts", "./worker/types/*.ts", "./**/*.spec.ts"],
				}),
			});
		});

		it("should have returned 200 status code", () => {
			expect(rs.status).toBe(200);
		});

		it('should have content-type "application/json"', () => {
			expect(rs.headers.get("content-type")).toBe("application/json");
		});

		it("should have no error", async () => {
			const rbody = await rs.json();
			expect(rbody).toHaveProperty("error", null);
		});
	});

	describe("list entries from kv", () => {
		let rs;

		beforeAll(async () => {
			await SELF.fetch("http://service/kv?key=entry1", {
				method: "POST",
				body: "test",
			}).then((r) => r.text());
			await SELF.fetch("http://service/kv?key=entry2", {
				method: "POST",
				body: "test",
			}).then((r) => r.text());
			await SELF.fetch("http://service/kv?key=entry3", {
				method: "POST",
				body: "test",
			}).then((r) => r.text());
			rs = await SELF.fetch("http://service/kv/list").then((r) => r.json());
		});

		it("should have returned list with 3 entries", () => {
			expect(rs).toHaveProperty("data");
			expect(rs.data).toHaveProperty("Keys");
			expect(rs.data.Keys).toHaveLength(3);
		});
	});

	describe("get entry from kv", () => {
		let rs;

		beforeAll(async () => {
			await SELF.fetch("http://service/kv?key=get:entry", {
				method: "POST",
				body: "test",
			}).then((r) => r.text());

			rs = await SELF.fetch("http://service/kv?key=get:entry").then((r) =>
				r.text(),
			);
		});

		it("should have returned list with 3 entries", () => {
			expect(rs).toEqual("test");
		});
	});

	describe("get non-existent entry from kv", () => {
		let rs;

		beforeAll(async () => {
			rs = await SELF.fetch("http://service/kv?key=get:entry:1").then((r) =>
				r.json(),
			);
		});

		it("should have returned list with 3 entries", () => {
			expect(rs).toStrictEqual({ error: "key has no value" });
		});
	});
});
