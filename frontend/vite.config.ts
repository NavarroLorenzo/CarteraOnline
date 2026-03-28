import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, ".", "");
  const port = Number(env.VITE_PORT || 5173);
  const proxyTarget = (env.VITE_API_PROXY_TARGET || "").trim();

  return {
    plugins: [react()],
    server: {
      host: "0.0.0.0",
      port,
      proxy: proxyTarget
        ? {
            "/api": {
              target: proxyTarget,
              changeOrigin: true,
            },
          }
        : undefined,
    },
  };
});
