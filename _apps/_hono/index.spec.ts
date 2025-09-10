import assert from "node:assert/strict";
import { spawn, type ChildProcessWithoutNullStreams } from "node:child_process";
import { after, before, describe, it } from "node:test";

describe("GET /hello", () => {
  let serverProcess: ChildProcessWithoutNullStreams;

  before(() => {
    serverProcess = spawn("tsx", ["index.ts"]);
    return new Promise((resolve) => setTimeout(resolve, 500));
  });

  after(() => {
    serverProcess.kill();
  });

  it("should return hello from wasm", async () => {
    const res = await fetch("http://localhost:5173/hello");
    const data = await res.text();
    assert.equal(res.status, 200);
    assert.equal(data, "hello");
  });
});
