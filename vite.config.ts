import { cloudflare } from "@cloudflare/vite-plugin";
import path from "path";
import { defineConfig } from "vite";
import { watch } from "vite-plugin-watch";

export default defineConfig({
  resolve: {
    alias: {
      $wrk: path.resolve(__dirname, "./_apps/_worker"),
    },
  },
  plugins: [
    watch({
      pattern: "/**/*.go",
      command: "bun worker build",
    }),
    cloudflare({ configPath: "./_apps/_worker/wrangler.toml" }),
  ],
});
