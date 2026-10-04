<script setup lang="ts">
  defineProps<{ eyebrow?: string; title: string; description?: string }>();
</script>

<template>
  <header class="mb-7 flex flex-col gap-5 md:mb-9 md:flex-row md:items-end md:justify-between">
    <div class="max-w-3xl">
      <p v-if="eyebrow" class="catalog-eyebrow mb-2">{{ eyebrow }}</p>
      <h1 class="text-3xl font-extrabold leading-tight text-base-content sm:text-4xl">{{ title }}</h1>
      <p v-if="description" class="catalog-muted mt-2 max-w-2xl text-sm leading-6 sm:text-base">{{ description }}</p>
    </div>
    <div v-if="$slots.actions" class="flex flex-wrap gap-2 md:justify-end">
      <slot name="actions" />
    </div>
  </header>
</template>
