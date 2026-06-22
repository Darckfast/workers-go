import { afterAll, expect, test } from "bun:test";
import app from './index';

afterAll(async () => {
  await app.stop()
});

test("should return hello from wasm", async () => {
  const response = await fetch("http://localhost:5173/hello");
  expect(response.status).toBe(200);
  expect(await response.text()).toBe("hello");
});
