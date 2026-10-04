import { defineNuxtConfig } from "nuxt/config";

// https://v3.nuxtjs.org/api/configuration/nuxt.config
export default defineNuxtConfig({
  ssr: false,
  modules: [
    "@nuxtjs/tailwindcss",
    "@pinia/nuxt",
    "@vueuse/nuxt",
    "@vite-pwa/nuxt",
    "./nuxt.proxyoverride.ts",
    "unplugin-icons/nuxt",
  ],
  nitro: {
    devProxy: {
      "/api": {
        target: "http://localhost:7745/api",
        ws: true,
        changeOrigin: true,
      },
    },
  },
  css: ["@/assets/css/main.css"],
  pwa: {
    workbox: {
      navigateFallbackDenylist: [/^\/api/],
    },
    injectRegister: "script",
    injectManifest: {
      swSrc: "sw.js",
    },
    devOptions: {
      // Enable to troubleshoot during development
      enabled: false,
    },
    manifest: {
      name: "Catalog",
      short_name: "Catalog",
      description: "Band equipment and inventory management",
      theme_color: "#111827",
      background_color: "#f8fafc",
      start_url: "/home",
      icons: [
        {
          src: "favicon.svg",
          sizes: "any",
          type: "image/svg+xml",
          purpose: "any",
        },
      ],
    },
  },
});
