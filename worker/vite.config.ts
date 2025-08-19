import { cloudflare } from "@cloudflare/vite-plugin";
import { defineConfig } from "vite";
import { watch } from "vite-plugin-watch";

export default defineConfig({
	plugins: [
		watch({
			pattern: "/**/*.go",
			command: "pnpm run build",
		}),
		cloudflare(),
	],
});
