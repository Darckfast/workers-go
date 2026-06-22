import { exports } from "cloudflare:workers";
import { beforeAll, describe, expect, it } from "vitest";

const payload = JSON.stringify({
  bool: true,
  number: 1,
  string: "my super string",
  list: [1, 2, 3, 4, 5, 6, 7],
});

describe("RPC", () => {
  describe("RPC and HTTP calls", () => {
    let rHttp, rRpc;

    beforeAll(async () => {
      let rs = await exports.default.fetch("http://dummy/echo", {
        body: payload,
        method: "POST",
        headers: {
          "content-type": "application/json",
        },
      });

      rHttp = await rs.text();
      rs = await exports.default.echo(new TextEncoder().encode(payload));
      rRpc = new TextDecoder().decode(rs[0]);
    });

    it("should yeild the same output for both calls", () => {
      expect(rRpc).to.be.eq(rHttp);
      expect(rRpc).to.be.eq(payload);
    });
  });
});
