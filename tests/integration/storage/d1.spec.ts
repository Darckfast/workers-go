import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it, test } from "vitest";

describe("d1", () => {
  describe('read timestamp', () => {
    let rs: Response
    let body

    beforeAll(async () => {
      rs = await SELF.fetch("http://service/d1", {
        method: "GET",
      })

      body = await rs.json()
    })

    it('should have returned 200 status code', () => {
      expect(rs.status).toBe(200)
    })

    it('should have key count', async () => {
      expect(body).toHaveProperty('result')
      expect(body.result.includes('2025')).toBe(true)
    })
  })

});
