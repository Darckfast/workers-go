import worker from "$wrk/main";
import {
  createExecutionContext,
  env,
  waitOnExecutionContext,
} from "cloudflare:test";
import { beforeAll, describe, expect, it } from "vitest";

describe("tail()", () => {
  let resultSaveInKV;
  beforeAll(async () => {
    const ctx = createExecutionContext();
    const headers = new Headers();
    headers.append("content-type", "application/json");
    headers.append("x-custom-header", "true");

    const events = [
      {
        wallTime: 0,
        cpuTime: 0,
        truncated: false,
        executionModel: "stateless",
        outcome: "ok",
        scriptName: null,
        diagnosticsChannelEvents: [],
        exceptions: [],
        logs: [
          {
            message: ["2025/08/24 08:47:15 this is a log 751000000"],
            level: "log",
            timestamp: 1756036035754,
          },
          {
            message: ["2025/08/24 08:47:15 my error this is a error"],
            level: "log",
            timestamp: 1756036035771,
          },
          {
            message: ["panic: this is a error"],
            level: "log",
            timestamp: 1756036035772,
          },
          {
            message: [""],
            level: "log",
            timestamp: 1756036035772,
          },
          {
            message: ["goroutine 9 [running]:"],
            level: "log",
            timestamp: 1756036035773,
          },
          {
            message: [
              "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/errors.init.func1({0xb3888, 0x47ab40}, 0x498140)",
            ],
            level: "log",
            timestamp: 1756036035773,
          },
          {
            message: [
              "\t/home/v/workers-go/worker/pkg/fetchhandler/errors/get-with-error.go:19 +0x13",
            ],
            level: "log",
            timestamp: 1756036035773,
          },
          {
            message: [
              "net/http.HandlerFunc.ServeHTTP(0x78d58, {0xb3888, 0x47ab40}, 0x498140)",
            ],
            level: "log",
            timestamp: 1756036035773,
          },
          {
            message: ["\t/usr/local/go/src/net/http/server.go:2294 +0x4"],
            level: "log",
            timestamp: 1756036035773,
          },
          {
            message: [
              "net/http.(*ServeMux).ServeHTTP(0x2ffe60, {0xb3888, 0x47ab40}, 0x498140)",
            ],
            level: "log",
            timestamp: 1756036035774,
          },
          {
            message: ["\t/usr/local/go/src/net/http/server.go:2822 +0x2e"],
            level: "log",
            timestamp: 1756036035774,
          },
          {
            message: [
              "github.com/Darckfast/workers-go/cloudflare/fetch.handler.func1()",
            ],
            level: "log",
            timestamp: 1756036035774,
          },
          {
            message: [
              "\t/home/v/workers-go/cloudflare/fetch/handler.go:82 +0x4",
            ],
            level: "log",
            timestamp: 1756036035774,
          },
          {
            message: [
              "created by github.com/Darckfast/workers-go/cloudflare/fetch.handler in goroutine 8",
            ],
            level: "log",
            timestamp: 1756036035775,
          },
          {
            message: [
              "\t/home/v/workers-go/cloudflare/fetch/handler.go:72 +0x28",
            ],
            level: "log",
            timestamp: 1756036035775,
          },
          {
            message: ["exit code:", 2],
            level: "warn",
            timestamp: 1756036035775,
          },
        ],
        eventTimestamp: 1756036035640,
        event: {
          request: {
            url: "http://localhost:5173/error",
            method: "GET",
            headers: {
              accept: "*/*",
              "accept-encoding": "br, gzip",
              "cf-connecting-ip": "127.0.0.1",
              host: "localhost:5173",
              "user-agent": "yaak",
            },
            cf: {
              clientTcpRtt: 17,
              requestHeaderNames: {},
              httpProtocol: "HTTP/1.1",
              tlsCipher: "AEAD-AES256-GCM-SHA384",
            },
          },
          response: {
            status: 500,
          },
        },
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
        wallTime: 0,
        cpuTime: 0,
        truncated: false,
        executionModel: "stateless",
        outcome: "ok",
        scriptName: "", // null becomes empty string
        diagnosticsChannelEvents: [],
        exceptions: [],
        logs: [
          {
            message: ["2025/08/24 08:47:15 this is a log 751000000"],
            level: "log",
            timestamp: 1756036035754,
          },
          {
            message: ["2025/08/24 08:47:15 my error this is a error"],
            level: "log",
            timestamp: 1756036035771,
          },
          {
            message: ["panic: this is a error"],
            level: "log",
            timestamp: 1756036035772,
          },
          {
            message: [""],
            level: "log",
            timestamp: 1756036035772,
          },
          {
            message: ["goroutine 9 [running]:"],
            level: "log",
            timestamp: 1756036035773,
          },
          {
            message: [
              "github.com/Darckfast/workers-go/worker/pkg/fetchhandler/errors.init.func1({0xb3888, 0x47ab40}, 0x498140)",
            ],
            level: "log",
            timestamp: 1756036035773,
          },
          {
            message: [
              "\t/home/v/workers-go/worker/pkg/fetchhandler/errors/get-with-error.go:19 +0x13",
            ],
            level: "log",
            timestamp: 1756036035773,
          },
          {
            message: [
              "net/http.HandlerFunc.ServeHTTP(0x78d58, {0xb3888, 0x47ab40}, 0x498140)",
            ],
            level: "log",
            timestamp: 1756036035773,
          },
          {
            message: ["\t/usr/local/go/src/net/http/server.go:2294 +0x4"],
            level: "log",
            timestamp: 1756036035773,
          },
          {
            message: [
              "net/http.(*ServeMux).ServeHTTP(0x2ffe60, {0xb3888, 0x47ab40}, 0x498140)",
            ],
            level: "log",
            timestamp: 1756036035774,
          },
          {
            message: ["\t/usr/local/go/src/net/http/server.go:2822 +0x2e"],
            level: "log",
            timestamp: 1756036035774,
          },
          {
            message: [
              "github.com/Darckfast/workers-go/cloudflare/fetch.handler.func1()",
            ],
            level: "log",
            timestamp: 1756036035774,
          },
          {
            message: [
              "\t/home/v/workers-go/cloudflare/fetch/handler.go:82 +0x4",
            ],
            level: "log",
            timestamp: 1756036035774,
          },
          {
            message: [
              "created by github.com/Darckfast/workers-go/cloudflare/fetch.handler in goroutine 8",
            ],
            level: "log",
            timestamp: 1756036035775,
          },
          {
            message: [
              "\t/home/v/workers-go/cloudflare/fetch/handler.go:72 +0x28",
            ],
            level: "log",
            timestamp: 1756036035775,
          },
          {
            message: [
              "exit code:",
              "", // this was a number, but we are parsing only string values
            ],
            level: "warn",
            timestamp: 1756036035775,
          },
        ],
        eventTimestamp: 1756036035640,
        event: {
          request: {
            url: "http://localhost:5173/error",
            method: "GET",
            headers: {
              accept: "*/*",
              "accept-encoding": "br, gzip",
              "cf-connecting-ip": "127.0.0.1",
              host: "localhost:5173",
              "user-agent": "yaak",
            },
            cf: {
              clientTcpRtt: 17,
              requestHeaderNames: {},
              httpProtocol: "HTTP/1.1",
              tlsCipher: "AEAD-AES256-GCM-SHA384",
            },
          },
          response: {
            status: 500,
          },
        },
      },
    ]);
  });
});
