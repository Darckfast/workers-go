import worker from "$wrk/main";
import {
  createExecutionContext,
  env,
  waitOnExecutionContext,
} from "cloudflare:test";
import { beforeAll, describe, expect, it } from "vitest";

describe("email()", () => {
	let email;
	const called: any = {};

	beforeAll(async () => {
		email = {
			headers: new Headers({
				"return-path": "<mlemos@acm.org>",
				to: "Manuel Lemos <mlemos@linux.local>",
				subject:
					"Testing Manuel Lemos' MIME E-mail composing and sending PHP class: HTML message",
				from: "mlemos <mlemos@acm.org>",
				"reply-to": "mlemos <mlemos@acm.org>",
				sender: "mlemos@acm.org",
				"x-mailer":
					"http://www.phpclasses.org/mimemessage $Revision: 1.63 $ (mail)",
				"mime-version": "1.0",
				"content-type":
					'multipart/mixed; boundary="652b8c4dcb00cdcdda1e16af36781caf"',
				"message-id": "<20050430192829.0489.mlemos@acm.org>",
				date: "Sat, 30 Apr 2005 19:28:29 -0300",
			}),
			from: "mlemos@acm.org",
			to: "mlemos@linux.local",
			forward(rcptTo, headers) {
				called.forward = { rcptTo, headers };
				return Promise.resolve();
			},
			reply(message) {
				called.reply = { message };
				return Promise.resolve();
			},
			setReject(reason) {
				called.setReject = { reason };
			},
			raw: undefined,
			rawSize: 0,
		};

		const r = await fetch(
			"https://gist.githubusercontent.com/billsinc/967795/raw/8c7c36615f33380f923c575d4e27a5ae03f10ef7/Test%2520Email%2520from%2520PHPClasses",
		);

		email.raw = r.body;
		email.rawSize = Number(r.headers.get("content-length"));

		const ctx = createExecutionContext();
		await worker.email(email, env, ctx);
		await waitOnExecutionContext(ctx);
	});

	it("should have called all 3 functions", () => {
		expect(Object.keys(called)).toHaveLength(3);
	});

	it("should have returned headers x-test-id", () => {
		expect(called.forward.headers.get("x-test-id")).toBe(
			"12345-asdfg-56789-ghjkl",
		);
	});

	it("should have forwarded to rcptTo ", () => {
		expect(called.forward.rcptTo).toBe("<YOUR_EMAIL>");
	});

	it("should have returned the reply from", () => {
		expect(called.reply.message.from).toBe("me");
	});

	it("should have returned the reply to", () => {
		expect(called.reply.message.to).toBe("you");
	});

	it("should have returned reply raw", async () => {
		expect(called.reply.message.raw).toBeInstanceOf(ReadableStream);

		const rs = new Response(called.reply.message.raw);
		expect(await rs.text()).toBe("this is a test, and this email has 8878");
	});

	it("should have returned reject reason", async () => {
		expect(called.setReject.reason).toBe("this reject is just for testing");
	});
});
