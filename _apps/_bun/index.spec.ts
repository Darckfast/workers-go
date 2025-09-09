import { afterAll, beforeAll, expect, test } from "bun:test";

let serverProcess: ReturnType<typeof Bun.spawn>;

beforeAll(() => {
  serverProcess = Bun.spawn(["bun", "run", "index.ts"], {
    stdio: ["inherit", "inherit", "inherit"],
  });

  return new Promise((resolve) => setTimeout(resolve, 500));
});

afterAll(() => {
  serverProcess.kill();
});

test("should return hello from wasm", async () => {
  const response = await fetch("http://localhost:5173/hello");
  expect(response.status).toBe(200);
  expect(await response.text()).toBe("hello");
});
