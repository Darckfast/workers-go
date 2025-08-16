import worker from '$wrk/main';
import {
  createExecutionContext,
  env,
  waitOnExecutionContext,
} from "cloudflare:test";
import { describe, expect, it } from "vitest";

describe("tail handler", () => {
  it("should proccess the event", async () => {
    const ctx = createExecutionContext();
    const headers = new Headers()
    headers.append('content-type', 'application/json')
    headers.append('x-custom-header', 'true')

    const events = [
      {
        type: "tail",
        traces: [{
          scriptName: "worker-producer",
          truncated: false,
          wallTime: 3713091757396727,
          cpuTime: 5175330951149985,
          dispatchNamespace: "dispatch-namespace",
          entrypoint: 'main.ts',
          scriptTags: ["tags"],
          executionModel: "exec-mode",
          diagnosticsChannelEvents: [{
            message: "some",
            timestamp: 1755357868138,
            channel: 'chan'
          }],
          scriptVersion: {
            id: "e3a18c84-4622-4a11-b703-455b38d50001",
            message: "script-message",
            tag: "01K2SQPW3AKDG2R0QZC2YNCSSB"
          },
          event: {
            request: {
              cf: null,
              headers: headers,
              method: "GET",
              url: "http://service-url"
            },
            response: {
              status: 200,
            }
          },
          eventTimestamp: 1755357868138,
          logs: [
            {
              level: 'info',
              timestamp: 1755357868138,
              message: "start"
            },
            {
              level: 'info',
              timestamp: 1755357868138,
              message: "end"
            },
          ],
          exceptions: [],
          outcome: "ok"
        }],
        waitUntil(promise) {
          return promise
        },
      }
    ]

    await worker.tail(events, env, ctx);
    await waitOnExecutionContext(ctx);
    // expect(response.status).toBe(404);
    const request = new Request("http://example.com/tail");
    // Create an empty context to pass to `worker.fetch()`
    const response: Response = await worker.fetch(request, env, ctx);
    // Wait for all `Promise`s passed to `ctx.waitUntil()` to settle before running test assertions
    await waitOnExecutionContext(ctx);
    const result = await response.json()

    expect(JSON.parse(result.result)).toStrictEqual([
      {
        "Type": "tail",
        "Events": [],
        "Traces": [
          {
            "ScriptName": "worker-producer",
            "Entrypoint": "main.ts",
            "Event": {
              "Type": "fetch",
              "RpcMethod": "",
              "MailFrom": "",
              "RcptTo": "",
              "RawSize": 0,
              "Queue": "",
              "BatchSize": 0,
              "ScheduledTime": 0,
              "Cron": "",
              "ConsumedEvents": null,
              "Response": {
                "Status": 200
              },
              "Request": {
                "Cf": null,
                "Headers": {
                  "Content-Type": [
                    "application/json"
                  ],
                  "X-Custom-Header": [
                    "true"
                  ]
                },
                "Method": "GET",
                "Url": "http://service-url"
              },
              "GetWebSocketEvent": null
            },
            "EventTimeStamp": 1755357868138,
            "Logs": [
              {
                "Timestamp": 1755357868138,
                "Level": "info",
                "Message": "start"
              },
              {
                "Timestamp": 1755357868138,
                "Level": "info",
                "Message": "end"
              }
            ],
            "Exceptions": [],
            "DiagnosticsChannelEvents": [
              {
                "Timestamp": 1755357868138,
                "Channel": "chan",
                "Message": "some"
              }
            ],
            "Outcome": "ok",
            "Truncated": false,
            "CpuTime": 5175330951149985,
            "WallTime": 3713091757396727,
            "ExecutionModel": "exec-mode",
            "ScriptTags": [
              "tags"
            ],
            "DispatchNamespace": "dispatch-namespace",
            "ScriptVersion": {
              "Id": "e3a18c84-4622-4a11-b703-455b38d50001",
              "Tag": "01K2SQPW3AKDG2R0QZC2YNCSSB",
              "Message": "script-message"
            }
          }
        ]
      }
    ])
  });
});
