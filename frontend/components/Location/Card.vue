<script setup lang="ts">
  import type { LocationOut, LocationOutCount, LocationSummary } from "~~/lib/api/types/data-contracts";
  import MdiArrowRight from "~icons/mdi/arrow-right";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";

  const props = defineProps<{ location: LocationOutCount | LocationOut | LocationSummary; dense?: boolean }>();
  const count = computed(() => (props.location as LocationOutCount).itemCount || 0);
</script>

<template>
  <NuxtLink
    :to="`/location/${location.id}`"
    class="group catalog-panel flex items-center gap-3 rounded-2xl p-4 transition hover:-translate-y-0.5 hover:border-primary/25"
  >
    <div class="grid h-10 w-10 shrink-0 place-items-center rounded-xl bg-info/10 text-info">
      <MdiMapMarkerOutline class="h-5 w-5" />
    </div>
    <div class="min-w-0">
      <h3 class="truncate text-sm font-bold">{{ location.name }}</h3>
      <p class="catalog-muted mt-0.5 text-xs">{{ count }} {{ count === 1 ? "item" : "items" }}</p>
    </div>
    <MdiArrowRight
      class="ml-auto h-4 w-4 text-base-content/25 transition group-hover:translate-x-1 group-hover:text-primary"
    />
  </NuxtLink>
</template>
