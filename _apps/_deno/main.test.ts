import { assertEquals } from "jsr:@std/assert";
import { afterAll, beforeAll, describe, it } from "jsr:@std/testing/bdd";

describe("GET /hello", () => {
  let serverProcess: Deno.ChildProcess;

  beforeAll(() => {
    serverProcess = new Deno.Command("deno", {
      args: ["run", "dev"],
      stdin: "piped",
      stdout: "piped",
    }).spawn();

    return new Promise((resolve) => setTimeout(resolve, 100));
  });

  afterAll(async () => {
    serverProcess.stdin.close();
    serverProcess.stdout.cancel();
    serverProcess.kill();
    await serverProcess.status;
  });

  it("should return hello from wasm", async () => {
    const res = await fetch("http://localhost:5173/hello");
    const data = await res.text();
    assertEquals(res.status, 200);
    assertEquals(data, "hello");
  });
});
