<script setup lang="ts">
  defineProps<{ title: string; description: string }>();
</script>

<template>
  <div class="catalog-panel flex min-h-[240px] flex-col items-center justify-center rounded-2xl px-6 py-12 text-center">
    <div class="mb-4 grid h-14 w-14 place-items-center rounded-2xl bg-secondary text-secondary-content">
      <slot name="icon"><AppLogo class="h-7 w-7" /></slot>
    </div>
    <h3 class="text-lg font-bold">{{ title }}</h3>
    <p class="catalog-muted mt-2 max-w-md text-sm leading-6">{{ description }}</p>
    <div v-if="$slots.actions" class="mt-5 flex flex-wrap justify-center gap-2"><slot name="actions" /></div>
  </div>
</template>
