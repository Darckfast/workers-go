import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it } from "vitest";

describe("fetch", () => {
  describe("POST /application/x-www-form-urlencoded ", () => {
    let rs: Response;
    let body;
    let params: URLSearchParams;

    beforeAll(async () => {
      params = new URLSearchParams();

      params.set("id", crypto.randomUUID());
      params.set("alpha", "asdlkjaslkdiuzxc");
      params.set("url", "https://darckfast.com");
      params.set("fullname", "Jonh doe");
      params.set("number", "12039812938210938210938");

      rs = await SELF.fetch("https://api/application/x-www-form-urlencoded ", {
        method: "POST",
        body: params,
      });
      body = await rs.json();
    });

    it("should have returned 200 status code", () => {
      expect(rs.status).toBe(200);
    });

    it.each(["id", "alpha", "url", "fullname", "number"])(
      "should return %s",
      (key) => {
        expect(body[key]).toBe(params.get(key));
      },
    );
  });
});
