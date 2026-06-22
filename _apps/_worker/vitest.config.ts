import {
  cloudflareTest,
  readD1Migrations,
} from "@cloudflare/vitest-pool-workers";
import path from "node:path";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [
    cloudflareTest(async () => {
      const migrationsPath = path.join(__dirname, "./migrations");
      const migrations = await readD1Migrations(migrationsPath);

      return {
        singleWorker: true,
        wrangler: { configPath: "./wrangler.toml" },
        miniflare: {
          kvNamespaces: ["TEST_NAMESPACE"],
          bindings: { TEST_MIGRATIONS: migrations },
          d1Databases: ["DB"],
          r2Buckets: ["TEST_BUCKET"],
          compatibilityFlags: ["service_binding_extra_handlers"],
          queueConsumers: {
            queue: { maxBatchTimeout: 0 },
          },
        },
      };
    }),
  ],
  resolve: {
    alias: {
      $wrk: path.resolve(__dirname, "."),
    },
  },
  test: {
    sequence: {
      concurrent: true
    },
    dir: "./tests",
    setupFiles: ["./tests/apply-migrations.ts"],
    testTimeout: 10000,
  },
});
