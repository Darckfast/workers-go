import { cloudflare } from "@cloudflare/vite-plugin";
import path from "path";
import { defineConfig } from "vite";
import { watch } from "vite-plugin-watch";

export default defineConfig({
  resolve: {
    alias: {
      $wrk: path.resolve(__dirname, "./worker"),
    },
  },
  plugins: [
    watch({
      pattern: "/**/*.go",
      command: "pnpm worker build",
    }),
    cloudflare({ configPath: "./_worker/wrangler.toml" }),
  ],
});
