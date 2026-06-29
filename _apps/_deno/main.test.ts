/// <reference lib="deno.ns" />
import { assertEquals } from "jsr:@std/assert";
import { afterAll, describe, it } from "jsr:@std/testing/bdd";
import app from "./main.ts";

describe("GET /hello", () => {
  afterAll(async () => {
    await app.shutdown();
  });

  it("should return hello from wasm", async () => {
    const res = await fetch("http://localhost:5173/");
    const data = await res.text();
    assertEquals(res.status, 200);
    assertEquals(data, "hello");
  });
});
