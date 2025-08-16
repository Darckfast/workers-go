import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it, test } from "vitest";

describe("cache", () => {
  describe('GET /cache', () => {
    let rs: Response
    let body

    beforeAll(async () => {
      rs = await SELF.fetch("http://service/cache", {
        method: "GET",
      })

      body = await rs.text()
      rs = await SELF.fetch("http://service/cache", {
        method: "GET",
      })

      body = await rs.text()
    })

    it('should have returned 200 status code', () => {
      expect(rs.status).toBe(200)
    })

    it('should have no error', async () => {
      expect(body).not.toBeUndefined()
    })

    it('should have key count', async () => {
      expect(rs.headers.get('x-cache')).toBe('hit')
    })
  })
});
