import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it, test } from "vitest";

describe("durable object", () => {
  describe('write op', () => {
    let rs: Response
    let body

    beforeAll(async () => {

      rs = await SELF.fetch("http://service/do", {
        method: "GET",
      })

      body = await rs.json()
    })

    it('should have returned 200 status code', () => {
      expect(rs.status).toBe(200)
    })

    it('should have content-type "application/json"', () => {
      expect(rs.headers.get('content-type')).toBe('application/json')
    })

    it('should have no error', async () => {
      expect(body).toHaveProperty('has_error', false)
    })

    it('should have key count', async () => {
      expect(body).toHaveProperty('result', 'Hello, World!')
    })
  })
});
