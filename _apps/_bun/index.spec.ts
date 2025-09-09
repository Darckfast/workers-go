import { afterAll, beforeAll, expect, test } from "bun:test";

let serverProcess: ReturnType<typeof Bun.spawn>;

beforeAll(() => {
  serverProcess = Bun.spawn(["bun", "run", "index.ts"], {
    stdio: ["inherit", "inherit", "inherit"],
  });
  // Give the server a moment to start up
  return new Promise((resolve) => setTimeout(resolve, 1000));
});

afterAll(() => {
  serverProcess.kill();
});

test("server responds to root path", async () => {
  const response = await fetch("http://localhost:5173/hello");
  expect(response.status).toBe(200);
  expect(await response.text()).toBe("hello");
});
