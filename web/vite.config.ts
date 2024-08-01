import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    proxy: {
      "^/api/.*": {
        target: "http://localhost:1323/",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, ""),
      },
      "^/api/ws": {
        target: "ws://localhost:1323/",
        ws: true,
        rewrite: (p) => p.replace(/^\/api/, ""),
        rewriteWsOrigin: true,
      },
    },
  },
});
