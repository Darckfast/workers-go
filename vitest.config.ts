import { defineWorkersConfig } from "@cloudflare/vitest-pool-workers/config";
import path from 'path';

export default defineWorkersConfig({
  resolve: {
    alias: {
      '$wrk': path.resolve(__dirname, './worker'),
    }
  },
  test: {
    testTimeout: 10000,
    poolOptions: {
      workers: {
        singleWorker: true,
        wrangler: { configPath: "./worker/wrangler.toml" },
        miniflare: {
          kvNamespaces: ["TEST_NAMESPACE"],
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
});
