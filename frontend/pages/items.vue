<script setup lang="ts">
  import type { ItemSummary, LabelSummary, LocationOutCount } from "~~/lib/api/types/data-contracts";
  import { useLabelStore } from "~~/stores/labels";
  import { useLocationStore } from "~~/stores/locations";
  import MdiLoading from "~icons/mdi/loading";
  import MdiMagnify from "~icons/mdi/magnify";
  import MdiDelete from "~icons/mdi/delete";
  import MdiChevronRight from "~icons/mdi/chevron-right";
  import MdiChevronLeft from "~icons/mdi/chevron-left";
  import MdiFilterVariant from "~icons/mdi/filter-variant";
  import MdiPackageVariantClosed from "~icons/mdi/package-variant-closed";
  import MdiPlus from "~icons/mdi/plus";

  definePageMeta({
    middleware: ["auth"],
  });

  useHead({
    title: "Inventory — Catalog",
  });

  const searchLocked = ref(false);
  const queryParamsInitialized = ref(false);
  const initialSearch = ref(true);

  const api = useUserApi();
  const loading = useMinLoader(500);
  const items = ref<ItemSummary[]>([]);
  const total = ref(0);

  const page1 = useRouteQuery("page", 1);

  const page = computed({
    get: () => page1.value,
    set: value => {
      page1.value = value;
    },
  });

  const pageSize = useRouteQuery("pageSize", 21);
  const query = useRouteQuery("q", "");
  const advanced = useRouteQuery("advanced", false);
  const includeArchived = useRouteQuery("archived", false);
  const fieldSelector = useRouteQuery("fieldSelector", false);

  const totalPages = computed(() => Math.ceil(total.value / pageSize.value));
  const hasNext = computed(() => page.value * pageSize.value < total.value);
  const hasPrev = computed(() => page.value > 1);

  function prev() {
    page.value = Math.max(1, page.value - 1);
  }

  function next() {
    page.value = Math.min(Math.ceil(total.value / pageSize.value), page.value + 1);
  }

  const route = useRoute();
  const router = useRouter();

  onMounted(async () => {
    loading.value = true;
    // Wait until locations and labels are loaded
    let maxRetry = 10;
    while (!labels.value || !locations.value) {
      await new Promise(resolve => setTimeout(resolve, 100));
      if (maxRetry-- < 0) {
        break;
      }
    }
    searchLocked.value = true;
    const qLoc = route.query.loc as string[];
    if (qLoc) {
      selectedLocations.value = locations.value.filter(l => qLoc.includes(l.id));
    }

    const qLab = route.query.lab as string[];
    if (qLab) {
      selectedLabels.value = labels.value.filter(l => qLab.includes(l.id));
    }

    queryParamsInitialized.value = true;
    searchLocked.value = false;

    const qFields = route.query.fields as string[];
    if (qFields) {
      fieldTuples.value = qFields.map(f => f.split("=") as [string, string]);

      for (const t of fieldTuples.value) {
        if (t[0] && t[1]) {
          await fetchValues(t[0]);
        }
      }
    }

    // trigger search if no changes
    if (!qLab && !qLoc) {
      search();
    }

    loading.value = false;
    window.scroll({
      top: 0,
      left: 0,
      behavior: "smooth",
    });
  });

  const locationsStore = useLocationStore();

  const locationFlatTree = await useFlatLocations();

  const locations = computed(() => locationsStore.allLocations);

  const labelStore = useLabelStore();
  const labels = computed(() => labelStore.labels);

  const selectedLocations = ref<LocationOutCount[]>([]);
  const selectedLabels = ref<LabelSummary[]>([]);

  const locIDs = computed(() => selectedLocations.value.map(l => l.id));
  const labIDs = computed(() => selectedLabels.value.map(l => l.id));

  function parseAssetIDString(d: string) {
    d = d.replace(/"/g, "").replace(/-/g, "");

    const aidInt = parseInt(d);
    if (isNaN(aidInt)) {
      return [-1, false];
    }

    return [aidInt, true];
  }

  const byAssetId = computed(() => query.value?.startsWith("#") || false);
  const parsedAssetId = computed(() => {
    if (!byAssetId.value) {
      return "";
    } else {
      const [aid, valid] = parseAssetIDString(query.value.replace("#", ""));
      if (!valid) {
        return "Invalid Asset ID";
      } else {
        return aid;
      }
    }
  });

  const fieldTuples = ref<[string, string][]>([]);
  const fieldValuesCache = ref<Record<string, string[]>>({});

  const { data: allFields } = useAsyncData(async () => {
    const { data, error } = await api.items.fields.getAll();

    if (error) {
      return [];
    }

    return data;
  });

  watch(fieldSelector, (newV, oldV) => {
    if (newV === false && oldV === true) {
      fieldTuples.value = [];
    }
  });

  async function fetchValues(field: string): Promise<string[]> {
    if (fieldValuesCache.value[field]) {
      return fieldValuesCache.value[field];
    }

    const { data, error } = await api.items.fields.getAllValues(field);

    if (error) {
      return [];
    }

    fieldValuesCache.value[field] = data;

    return data;
  }

  watch(advanced, (v, lv) => {
    if (v === false && lv === true) {
      selectedLocations.value = [];
      selectedLabels.value = [];
      fieldTuples.value = [];

      console.log("advanced", advanced.value);

      router.push({
        query: {
          advanced: route.query.advanced,
          q: query.value,
          page: page.value,
          pageSize: pageSize.value,
          includeArchived: includeArchived.value ? "true" : "false",
        },
      });
    }
  });

  async function search() {
    if (searchLocked.value) {
      return;
    }

    loading.value = true;

    const fields = [];

    for (const t of fieldTuples.value) {
      if (t[0] && t[1]) {
        fields.push(`${t[0]}=${t[1]}`);
      }
    }

    const toast = useNotifier();

    const { data, error } = await api.items.getAll({
      q: query.value || "",
      locations: locIDs.value,
      labels: labIDs.value,
      includeArchived: includeArchived.value,
      page: page.value,
      pageSize: pageSize.value,
      fields,
    });

    function resetItems() {
      page.value = Math.max(1, page.value - 1);
      loading.value = false;
      total.value = 0;
      items.value = [];
    }

    if (error) {
      resetItems();
      toast.error("Failed to search items");
      return;
    }

    if (!data.items || data.items.length === 0) {
      resetItems();
      return;
    }

    total.value = data.total;
    items.value = data.items;

    loading.value = false;
    initialSearch.value = false;
  }

  watchDebounced([page, pageSize, query, selectedLabels, selectedLocations], search, { debounce: 250, maxWait: 1000 });

  async function submit() {
    // Set URL Params
    const fields = [];
    for (const t of fieldTuples.value) {
      if (t[0] && t[1]) {
        fields.push(`${t[0]}=${t[1]}`);
      }
    }

    // Push non-reactive query fields
    await router.push({
      query: {
        // Reactive
        advanced: "true",
        archived: includeArchived.value ? "true" : "false",
        fieldSelector: fieldSelector.value ? "true" : "false",
        pageSize: pageSize.value,
        page: page.value,
        q: query.value,

        // Non-reactive
        loc: locIDs.value,
        lab: labIDs.value,
        fields,
      },
    });

    // Reset Pagination
    page.value = 1;

    // Perform Search
    await search();
  }

  async function reset() {
    // Set URL Params
    const fields = [];
    for (const t of fieldTuples.value) {
      if (t[0] && t[1]) {
        fields.push(`${t[0]}=${t[1]}`);
      }
    }

    await router.push({
      query: {
        archived: "false",
        fieldSelector: "false",
        pageSize: 10,
        page: 1,
        q: "",
        loc: [],
        lab: [],
        fields,
      },
    });

    await search();
  }
</script>

<template>
  <BaseContainer class="catalog-page max-w-none">
    <AppPageHeader
      eyebrow="Inventory"
      title="Everything you track"
      description="Search, filter, and open any item without digging through nested menus."
    >
      <template #actions>
        <BaseButton to="/item/new" class="btn-primary"
          ><template #icon><MdiPlus class="h-5 w-5" /></template>Add item</BaseButton
        >
      </template>
    </AppPageHeader>

    <section
      v-if="locations && labels"
      class="catalog-panel rounded-2xl p-4 sm:p-5"
      aria-label="Inventory search and filters"
    >
      <form class="flex flex-col gap-3 md:flex-row" @submit.prevent="submit">
        <label class="relative flex-1">
          <span class="sr-only">Search inventory</span>
          <MdiMagnify
            class="pointer-events-none absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-base-content/35"
          />
          <input
            v-model="query"
            type="search"
            class="input h-12 w-full bg-base-200/60 pl-12"
            placeholder="Search by name, asset ID, description, or serial…"
          />
        </label>
        <BaseButton type="submit" class="btn-primary md:min-w-[120px]">
          <template #icon
            ><MdiLoading v-if="loading" class="h-5 w-5 animate-spin" /><MdiMagnify v-else class="h-5 w-5"
          /></template>
          Search
        </BaseButton>
      </form>

      <p v-if="byAssetId" class="mt-3 rounded-lg bg-primary/5 px-3 py-2 text-xs font-semibold text-primary">
        Looking up asset ID {{ parsedAssetId }}
      </p>

      <div class="mt-4 flex flex-wrap items-center gap-2 border-t border-base-content/5 pt-4">
        <span
          class="mr-1 inline-flex items-center gap-1.5 text-xs font-bold uppercase tracking-wide text-base-content/45"
          ><MdiFilterVariant class="h-4 w-4" />Filters</span
        >
        <SearchFilter v-model="selectedLocations" label="Locations" :options="locationFlatTree">
          <template #display="{ item }"
            ><div>
              <div>{{ item.name }}</div>
              <div v-if="item.name != item.treeString" class="mt-1 text-xs opacity-60">{{ item.treeString }}</div>
            </div></template
          >
        </SearchFilter>
        <SearchFilter v-model="selectedLabels" label="Labels" :options="labels" />
        <label
          class="flex min-h-[36px] cursor-pointer items-center gap-2 rounded-lg border border-base-content/10 bg-base-100 px-3 text-xs font-semibold"
        >
          <input v-model="includeArchived" type="checkbox" class="checkbox checkbox-primary checkbox-xs" />Archived
        </label>
        <label
          class="flex min-h-[36px] cursor-pointer items-center gap-2 rounded-lg border border-base-content/10 bg-base-100 px-3 text-xs font-semibold"
        >
          <input v-model="fieldSelector" type="checkbox" class="checkbox checkbox-primary checkbox-xs" />Custom fields
        </label>
        <button
          type="button"
          class="ml-auto min-h-[36px] px-2 text-xs font-bold text-primary hover:underline"
          @click="reset"
        >
          Clear all
        </button>
      </div>

      <div v-if="fieldSelector" class="mt-4 rounded-xl bg-base-200/60 p-4">
        <div class="mb-3">
          <h3 class="text-sm font-bold">Custom field filters</h3>
          <p class="catalog-muted mt-1 text-xs">Match an item using values stored in your custom fields.</p>
        </div>
        <div v-for="(f, idx) in fieldTuples" :key="idx" class="mb-3 grid gap-2 sm:grid-cols-[1fr_1fr_auto]">
          <select
            v-model="fieldTuples[idx][0]"
            class="select select-bordered w-full"
            aria-label="Custom field"
            @change="fetchValues(f[0])"
          >
            <option value="" disabled>Select a field</option>
            <option v-for="(fv, _, i) in allFields" :key="i" :value="fv">{{ fv }}</option>
          </select>
          <select v-model="fieldTuples[idx][1]" class="select select-bordered w-full" aria-label="Custom field value">
            <option value="" disabled>Select a value</option>
            <option v-for="v in fieldValuesCache[f[0]]" :key="v" :value="v">{{ v }}</option>
          </select>
          <button
            type="button"
            class="btn btn-ghost btn-square"
            aria-label="Remove custom field filter"
            @click="fieldTuples.splice(idx, 1)"
          >
            <MdiDelete class="h-5 w-5" />
          </button>
        </div>
        <BaseButton type="button" class="btn-sm btn-ghost" @click="() => fieldTuples.push(['', ''])"
          ><template #icon><MdiPlus class="h-4 w-4" /></template>Add field filter</BaseButton
        >
      </div>
    </section>

    <section class="mt-8" aria-labelledby="inventory-results">
      <div class="mb-4 flex items-end justify-between gap-3">
        <div>
          <p class="catalog-eyebrow">Results</p>
          <h2 id="inventory-results" ref="itemsTitle" class="mt-1 text-xl font-bold">
            {{ total.toLocaleString() }} {{ total === 1 ? "item" : "items" }}
          </h2>
        </div>
        <p v-if="totalPages > 0" class="catalog-muted text-xs font-semibold">Page {{ page }} of {{ totalPages }}</p>
      </div>

      <div v-if="loading" class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
        <div v-for="i in 8" :key="i" class="catalog-panel overflow-hidden rounded-2xl">
          <div class="catalog-skeleton h-44" />
          <div class="space-y-3 p-4">
            <div class="catalog-skeleton h-5 w-2/3 rounded" />
            <div class="catalog-skeleton h-4 rounded" />
            <div class="catalog-skeleton h-4 w-4/5 rounded" />
          </div>
        </div>
      </div>
      <BaseEmptyState
        v-else-if="!items.length"
        title="No inventory matches this search"
        description="Try a broader search, clear a filter, or add a new item to your Catalog."
      >
        <template #icon><MdiPackageVariantClosed class="h-7 w-7" /></template>
        <template #actions
          ><BaseButton class="btn-ghost" @click="reset">Clear filters</BaseButton
          ><BaseButton to="/item/new" class="btn-primary">Add item</BaseButton></template
        >
      </BaseEmptyState>
      <div v-else ref="cardgrid" class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
        <ItemCard v-for="item in items" :key="item.id" :item="item" />
      </div>

      <nav
        v-if="items.length > 0 && (hasNext || hasPrev)"
        class="mt-8 flex flex-col items-center justify-between gap-3 rounded-2xl border border-base-content/10 bg-base-100 p-3 sm:flex-row"
        aria-label="Inventory pages"
      >
        <p class="catalog-muted px-2 text-xs font-semibold">Showing page {{ page }} of {{ totalPages }}</p>
        <div class="flex gap-2">
          <button :disabled="!hasPrev" class="btn btn-sm btn-ghost" @click="prev">
            <MdiChevronLeft class="h-5 w-5" />Previous</button
          ><button :disabled="!hasNext" class="btn btn-sm btn-ghost" @click="next">
            Next<MdiChevronRight class="h-5 w-5" />
          </button>
        </div>
      </nav>
    </section>
  </BaseContainer>
</template>

<style lang="css">
  .list-move,
  .list-enter-active,
  .list-leave-active {
    transition: all 0.25s ease;
  }

  .list-enter-from,
  .list-leave-to {
    opacity: 0;
    transform: translateY(30px);
  }

  .list-leave-active {
    position: absolute;
  }
</style>
