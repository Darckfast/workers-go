import assert from "node:assert/strict";
import { after, describe, it } from "node:test";
import app from './index.ts';

describe("GET /hello", () => {
  after(() => {
    app.close()
  });

  it("should return hello from wasm", async () => {
    const res = await fetch("http://localhost:5173/hello");
    const data = await res.text();
    assert.equal(res.status, 200);
    assert.equal(data, "hello");
  });
});
