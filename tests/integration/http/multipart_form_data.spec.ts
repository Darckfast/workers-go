import { faker } from '@faker-js/faker';
import { SELF } from "cloudflare:test";
import { beforeAll, describe, expect, it } from "vitest";

describe("fetch", () => {
    describe('POST /multipart/form-data ', () => {
        let rs: Response
        let body
        let size = 0

        beforeAll(async () => {
            const formdata = new FormData()
            const url = faker.image.url()
            const img = await fetch(url).then(r => r.blob())

            size = img.size
            formdata.append('img', img)
            formdata.append("json", JSON.stringify({
                "$schema": "https://www.schemastore.org/tsconfig.json",
                "compilerOptions": {
                    "target": "esnext",
                    "module": "esnext",
                    "lib": [
                        "esnext",
                        "dom"
                    ],
                    "moduleResolution": "bundler",
                    "types": [
                        "@cloudflare/vitest-pool-workers",
                    ],
                },
                "extends": [
                    "../tsconfig.json"
                ],
                "include": [
                    "./**/*.ts",
                    "../vite.config.ts",
                ],
            }))

            rs = await SELF.fetch("http://example.com/multipart/form-data", { method: 'POST', body: formdata });
            body = await rs.json()
        })

        it('should have returned 200 status code', () => {
            expect(rs.status).toBe(200);
        })

        it('should have content-type "application/json"', () => {
            expect(rs.headers.get('content-type')).toBe("application/json")
        })

        it('should have returned with no error', () => {
            expect(body.has_error).toBe(false)
        })

        it('should have filesize', () => {
            expect(body['actual-size']).toBe(size)
            expect(body['size']).toBe(size)
        })

        it('should have filename', () => {
            expect(body.filename).toBe('img')
        })

        it('should have json as string', () => {
            expect(body.json).toBe('{"$schema":"https://www.schemastore.org/tsconfig.json","compilerOptions":{"lib":["esnext","dom"],"module":"esnext","moduleResolution":"bundler","target":"esnext","types":["@cloudflare/vitest-pool-workers"]},"extends":["../tsconfig.json"],"include":["./**/*.ts","../vite.config.ts"]}')
        })
    })
});
