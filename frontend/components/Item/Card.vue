<script setup lang="ts">
  import type { ItemOut, ItemSummary } from "~~/lib/api/types/data-contracts";
  import MdiArrowTopRight from "~icons/mdi/arrow-top-right";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";
  import MdiShieldCheckOutline from "~icons/mdi/shield-check-outline";

  const props = defineProps<{ item: ItemOut | ItemSummary }>();
  const api = useUserApi();
  const imageUrl = computed(() =>
    props.item.imageId
      ? api.authURL(`/items/${props.item.id}/attachments/${props.item.imageId}`)
      : "/inventory-placeholder.svg"
  );
  const topLabels = computed(() => props.item.labels.slice(0, 2));
</script>

<template>
  <NuxtLink
    class="group catalog-panel flex min-h-full flex-col overflow-hidden rounded-2xl transition duration-200 hover:-translate-y-1 hover:border-primary/25 hover:shadow-xl"
    :to="`/item/${item.id}`"
  >
    <div class="relative h-44 overflow-hidden bg-base-200">
      <img
        class="h-full w-full object-cover transition duration-300 group-hover:scale-[1.025]"
        :src="imageUrl"
        :alt="`${item.name} inventory image`"
      />
      <div class="absolute inset-x-0 top-0 flex items-start justify-between p-3">
        <span v-if="item.archived" class="badge border-0 bg-neutral/85 text-neutral-content">Archived</span>
        <span v-else-if="item.status === 'missing'" class="badge badge-error border-0">Missing</span>
        <span v-else-if="item.status === 'checked_out'" class="badge badge-info border-0">Checked out</span>
        <span v-else class="badge badge-success border-0">Available</span>
        <span
          class="grid h-9 w-9 place-items-center rounded-xl bg-base-100/90 text-base-content opacity-0 shadow-lg backdrop-blur transition group-hover:opacity-100"
          ><MdiArrowTopRight class="h-5 w-5"
        /></span>
      </div>
    </div>
    <div class="flex flex-1 flex-col p-4">
      <div class="flex items-start justify-between gap-3">
        <h2 class="two-line text-base font-bold leading-6">{{ item.name }}</h2>
        <span class="badge badge-secondary shrink-0">×{{ item.quantity }}</span>
      </div>
      <div class="catalog-muted mt-2 flex min-h-[22px] items-center gap-1.5 text-xs">
        <MdiMapMarkerOutline class="h-4 w-4 shrink-0" /><span class="truncate">{{
          item.location?.name || "No location assigned"
        }}</span>
        <MdiShieldCheckOutline v-if="item.insured" class="ml-auto h-4 w-4 shrink-0 text-success" aria-label="Insured" />
      </div>
      <p class="mt-2 text-xs font-bold capitalize text-base-content/55">Condition: {{ item.condition }}</p>
      <Markdown
        v-if="item.description"
        class="three-line catalog-muted mt-3 text-sm leading-5"
        :source="item.description"
      />
      <p v-else class="catalog-muted mt-3 text-sm italic">No description yet</p>
      <div class="mt-auto flex flex-wrap gap-1.5 pt-4">
        <LabelChip v-for="label in topLabels" :key="label.id" :label="label" size="sm" />
        <span v-if="item.labels.length > 2" class="badge badge-ghost badge-sm">+{{ item.labels.length - 2 }}</span>
      </div>
    </div>
  </NuxtLink>
</template>

<style scoped>
  .three-line,
  .two-line {
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-box-orient: vertical;
  }
  .three-line {
    -webkit-line-clamp: 3;
    line-clamp: 3;
  }
  .two-line {
    -webkit-line-clamp: 2;
    line-clamp: 2;
  }
</style>
