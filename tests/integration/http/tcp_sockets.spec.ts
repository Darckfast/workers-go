import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it } from "vitest";

describe("sockets", () => {
  describe('create conn', () => {
    let rs: Response

    beforeAll(async () => {
      rs = await SELF.fetch("http://service/socket", {
        method: "GET",
      })
    })

    it('should have returned 200 status code', () => {
      expect(rs.status).toBe(200)
    })

    it('should have no error', async () => {
      expect(await rs.text()).toBe('hello.')
    })
  })
});
