<template>
  <NuxtLayout>
    <Html lang="en" :data-theme="theme || 'catalog'" />
    <Link rel="icon" type="image/svg" href="/favicon.svg"></Link>
    <Meta name="theme-color" content="#111827" />
    <Meta name="description" content="Catalog keeps instruments, equipment, locations, and maintenance organized." />
    <Link rel="manifest" href="/manifest.webmanifest" />
    <NuxtPage />
  </NuxtLayout>
</template>

<script lang="ts" setup>
  const { theme } = useTheme();
</script>
