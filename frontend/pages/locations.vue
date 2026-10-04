<script setup lang="ts">
  import { useTreeState } from "~~/components/Location/Tree/tree-state";
  import MdiCollapseAllOutline from "~icons/mdi/collapse-all-outline";
  import MdiMapMarkerPlusOutline from "~icons/mdi/map-marker-plus-outline";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";

  definePageMeta({
    middleware: ["auth"],
  });

  useHead({
    title: "Locations — Catalog",
  });

  const api = useUserApi();
  const createLocation = ref(false);

  const { data: tree } = useAsyncData(async () => {
    const { data, error } = await api.locations.getTree({
      withItems: true,
    });

    if (error) {
      return [];
    }

    return data;
  });

  const locationTreeId = "locationTree";

  const treeState = useTreeState(locationTreeId);

  const route = useRouter();

  onMounted(() => {
    // set tree state from query params
    const query = route.currentRoute.value.query;

    if (query && query[locationTreeId]) {
      console.debug("setting tree state from query params");
      const data = JSON.parse(query[locationTreeId] as string);

      for (const key in data) {
        treeState.value[key] = data[key];
      }
    }
  });

  watch(
    treeState,
    () => {
      // Push the current state to the URL
      route.replace({ query: { [locationTreeId]: JSON.stringify(treeState.value) } });
    },
    { deep: true }
  );

  function closeAll() {
    for (const key in treeState.value) {
      treeState.value[key] = false;
    }
  }
</script>

<template>
  <BaseContainer class="catalog-page max-w-none">
    <LocationCreateModal v-model="createLocation" />
    <AppPageHeader
      eyebrow="Organization"
      title="Locations"
      description="Mirror the way your equipment is actually stored—from buildings and rooms down to cases and shelves."
    >
      <template #actions
        ><BaseButton class="btn-primary" @click="createLocation = true"
          ><template #icon><MdiMapMarkerPlusOutline class="h-5 w-5" /></template>Add location</BaseButton
        ></template
      >
    </AppPageHeader>
    <BaseEmptyState
      v-if="tree && !tree.length"
      title="Create your first location"
      description="Start with a broad place such as Band Room, Storage, or Stage, then nest more specific areas inside it."
    >
      <template #icon><MdiMapMarkerOutline class="h-7 w-7" /></template>
      <template #actions
        ><BaseButton class="btn-primary" @click="createLocation = true">Add location</BaseButton></template
      >
    </BaseEmptyState>
    <BaseCard v-else>
      <div class="flex items-center justify-between border-b border-base-content/5 px-5 py-4">
        <div>
          <h2 class="text-sm font-bold">Location hierarchy</h2>
          <p class="catalog-muted mt-0.5 text-xs">Expand a location to see its contents and nested spaces.</p>
        </div>
        <button
          class="btn btn-sm btn-ghost tooltip tooltip-left"
          data-tip="Collapse all"
          aria-label="Collapse all locations"
          @click="closeAll"
        >
          <MdiCollapseAllOutline class="h-5 w-5" />
        </button>
      </div>
      <div class="p-4 sm:p-6"><LocationTreeRoot v-if="tree" :locs="tree" :tree-id="locationTreeId" /></div>
    </BaseCard>
  </BaseContainer>
</template>
