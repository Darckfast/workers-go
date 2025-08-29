import {
  defineWorkersConfig,
  readD1Migrations,
} from "@cloudflare/vitest-pool-workers/config";
import path from "path";

export default defineWorkersConfig(async () => {
  // Read all migrations in the `migrations` directory
  const migrationsPath = path.join(__dirname, "./_worker/migrations");
  const migrations = await readD1Migrations(migrationsPath);

  return {
    resolve: {
      alias: {
        $wrk: path.resolve(__dirname, "./_worker"),
      },
    },
    test: {
      setupFiles: ["./tests/apply-migrations.ts"],
      testTimeout: 10000,
      poolOptions: {
        workers: {
          singleWorker: true,
          wrangler: { configPath: "./_worker/wrangler.toml" },
          miniflare: {
            kvNamespaces: ["TEST_NAMESPACE"],
            bindings: { TEST_MIGRATIONS: migrations },
            d1Databases: ["DB"],
            r2Buckets: ["TEST_BUCKET"],
            compatibilityFlags: ["service_binding_extra_handlers"],
            queueConsumers: {
              queue: { maxBatchTimeout: 0.05 /* 10ms */ },
            },
          },
        },
      },
    },
  };
});
