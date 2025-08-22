import {
	createExecutionContext,
	env,
	waitOnExecutionContext,
} from "cloudflare:test";
import { beforeAll, describe, expect, it } from "vitest";
import worker from "$wrk/main";

describe("tail()", () => {
	let resultSaveInKV;
	beforeAll(async () => {
		const ctx = createExecutionContext();
		const headers = new Headers();
		headers.append("content-type", "application/json");
		headers.append("x-custom-header", "true");

		const events = [
			{
				scriptName: "worker-producer",
				truncated: false,
				wallTime: 3713091757396727,
				cpuTime: 5175330951149985,
				dispatchNamespace: "dispatch-namespace",
				entrypoint: "main.ts",
				scriptTags: ["tags"],
				executionModel: "exec-mode",
				diagnosticsChannelEvents: [
					{
						message: "some",
						timestamp: 1755357868138,
						channel: "chan",
					},
				],
				scriptVersion: {
					id: "e3a18c84-4622-4a11-b703-455b38d50001",
					message: "script-message",
					tag: "01K2SQPW3AKDG2R0QZC2YNCSSB",
				},
				event: {
					request: {
						cf: null,
						headers: headers,
						method: "GET",
						url: "http://service-url",
					},
					response: {
						status: 200,
					},
				},
				eventTimestamp: 1755357868138,
				logs: [
					{
						level: "info",
						timestamp: 1755357868138,
						message: "start",
					},
					{
						level: "info",
						timestamp: 1755357868138,
						message: "end",
					},
				],
				exceptions: [],
				outcome: "ok",
			},
		];

		await worker.tail(events, env, ctx);
		await waitOnExecutionContext(ctx);
		const request = new Request("http://example.com/tail");
		const response: Response = await worker.fetch(request, env, ctx);
		await waitOnExecutionContext(ctx);
		resultSaveInKV = await response.json();
	});

	it("should serialize and proccess the event", async () => {
		expect(JSON.parse(resultSaveInKV.result["tail:result"])).toStrictEqual([
			{
				scriptName: "worker-producer",
				entrypoint: "main.ts",
				event: {
					response: { status: 200 },
					request: {
						headers: {
							"Content-Type": ["application/json"],
							"X-Custom-Header": ["true"],
						},
						method: "GET",
						url: "http://service-url",
					},
				},
				eventTimestamp: 1755357868138,
				logs: [
					{ timestamp: 1755357868138, level: "info", message: "start" },
					{ timestamp: 1755357868138, level: "info", message: "end" },
				],
				diagnosticsChannelEvents: [
					{ timestamp: 1755357868138, channel: "chan", message: "some" },
				],
				outcome: "ok",
				cpuTime: 5175330951149985,
				wallTime: 3713091757396727,
				executionModel: "exec-mode",
				scriptTags: ["tags"],
				dispatchNamespace: "dispatch-namespace",
				scriptVersion: {
					id: "e3a18c84-4622-4a11-b703-455b38d50001",
					tag: "01K2SQPW3AKDG2R0QZC2YNCSSB",
					message: "script-message",
				},
			},
		]);
	});
});
