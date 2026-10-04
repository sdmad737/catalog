<script setup lang="ts">
  import { useLabelStore } from "~~/stores/labels";
  import { useLocationStore } from "~~/stores/locations";
  import MdiPackageVariantClosed from "~icons/mdi/package-variant-closed";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";
  import MdiTagMultipleOutline from "~icons/mdi/tag-multiple-outline";
  import MdiCashMultiple from "~icons/mdi/cash-multiple";
  import MdiArrowRight from "~icons/mdi/arrow-right";
  import MdiPlus from "~icons/mdi/plus";
  import MdiMagnify from "~icons/mdi/magnify";
  import MdiQrcodeScan from "~icons/mdi/qrcode-scan";
  import MdiFileDocumentOutline from "~icons/mdi/file-document-outline";
  import MdiLightbulbOnOutline from "~icons/mdi/lightbulb-on-outline";
  import MdiCheckCircleOutline from "~icons/mdi/check-circle-outline";
  import MdiAlertCircleOutline from "~icons/mdi/alert-circle-outline";
  import MdiAccountMultipleOutline from "~icons/mdi/account-multiple-outline";
  import MdiClockAlertOutline from "~icons/mdi/clock-alert-outline";

  definePageMeta({ middleware: ["auth"] });
  useHead({ title: "Dashboard — Catalog" });

  const api = useUserApi();
  const auth = useAuthContext();
  const locationStore = useLocationStore();
  const labelStore = useLabelStore();
  const breakpoints = useBreakpoints();
  const onboardingHidden = ref(false);

  const { data: statistics, pending: statsPending } = useAsyncData("dashboard-statistics", async () => {
    const response = await api.stats.group();
    return response.error ? null : response.data;
  });

  const {
    data: recentItems,
    pending: itemsPending,
    refresh: refreshItems,
  } = useAsyncData("dashboard-recent-items", async () => {
    const response = await api.items.getAll({
      page: 1,
      pageSize: 5,
      orderBy: "createdAt",
    });
    return response.error ? [] : response.data.items;
  });

  onServerEvent(ServerEvent.ItemMutation, () => refreshItems());

  const firstName = computed(() => auth.user?.name?.split(" ")[0] || "there");
  const greeting = computed(() => {
    const hour = new Date().getHours();
    if (hour < 12) return "Good morning";
    if (hour < 18) return "Good afternoon";
    return "Good evening";
  });

  const stats = computed(() => [
    {
      label: "Inventory items",
      value: statistics.value?.totalItems || 0,
      icon: MdiPackageVariantClosed,
      tone: "bg-primary/10 text-primary",
    },
    {
      label: "Available",
      value: statistics.value?.totalAvailable || 0,
      icon: MdiCheckCircleOutline,
      tone: "bg-success/10 text-success",
    },
    {
      label: "Checked out",
      value: statistics.value?.totalCheckedOut || 0,
      icon: MdiAccountMultipleOutline,
      tone: "bg-info/10 text-info",
    },
    {
      label: "Overdue",
      value: statistics.value?.totalOverdue || 0,
      icon: MdiClockAlertOutline,
      tone: "bg-warning/10 text-warning",
    },
    {
      label: "Missing",
      value: statistics.value?.totalMissing || 0,
      icon: MdiAlertCircleOutline,
      tone: "bg-error/10 text-error",
    },
    {
      label: "Total value",
      value: statistics.value?.totalItemPrice || 0,
      currency: true,
      icon: MdiCashMultiple,
      tone: "bg-secondary text-secondary-content",
    },
  ]);

  const quickActions = [
    {
      label: "Add an item",
      description: "Create a new inventory record",
      to: "/item/new",
      icon: MdiPlus,
      primary: true,
    },
    {
      label: "Search inventory",
      description: "Find gear, serials, or asset IDs",
      to: "/items",
      icon: MdiMagnify,
    },
    {
      label: "Add a borrower",
      description: "Create a student or member record",
      to: "/people",
      icon: MdiAccountMultipleOutline,
    },
    {
      label: "Scan an asset",
      description: "Open an item from its QR label",
      to: "/items",
      icon: MdiQrcodeScan,
    },
    {
      label: "Export or report",
      description: "Create labels and inventory files",
      to: "/tools",
      icon: MdiFileDocumentOutline,
    },
  ];

  const setupSteps = computed(() => [
    {
      label: "Review your locations",
      done: (statistics.value?.totalLocations || 0) > 0,
      to: "/locations",
    },
    {
      label: "Add your first inventory item",
      done: (statistics.value?.totalItems || 0) > 0,
      to: "/item/new",
    },
    {
      label: "Add a person or borrower",
      done: (statistics.value?.totalPeople || 0) > 0,
      to: "/people",
    },
    {
      label: "Program and verify a physical tag",
      done: (statistics.value?.totalVerifiedTags || 0) > 0,
      to: "/items",
    },
  ]);
  const setupComplete = computed(() => setupSteps.value.every(step => step.done));
</script>

<template>
  <BaseContainer class="catalog-page max-w-none">
    <AppPageHeader
      eyebrow="Overview"
      :title="`${greeting}, ${firstName}`"
      description="Here’s a clear view of your inventory and the fastest ways to get work done."
    >
      <template #actions>
        <BaseButton to="/item/new" class="btn-primary"
          ><template #icon><MdiPlus class="h-5 w-5" /></template>Add item</BaseButton
        >
      </template>
    </AppPageHeader>

    <section
      v-if="statistics && !setupComplete && !onboardingHidden"
      class="mb-7 rounded-2xl border border-primary/20 bg-primary/5 p-5"
    >
      <div class="flex items-start justify-between gap-4">
        <div>
          <p class="catalog-eyebrow">First-time setup</p>
          <h2 class="mt-1 text-xl font-black">Welcome to Catalog</h2>
          <p class="mt-1 text-sm text-base-content/55">
            Complete these practical steps in any order. Nothing is blocked.
          </p>
        </div>
        <button class="btn btn-ghost btn-sm" @click="onboardingHidden = true">Hide for now</button>
      </div>
      <div class="mt-5 grid gap-2 sm:grid-cols-2 xl:grid-cols-4">
        <NuxtLink
          v-for="step in setupSteps"
          :key="step.label"
          :to="step.to"
          class="flex min-h-[64px] items-center gap-3 rounded-xl border p-3 transition hover:bg-base-100"
          :class="step.done ? 'border-success/20 bg-success/5' : 'border-base-content/10 bg-base-100/70'"
          ><MdiCheckCircleOutline
            class="h-6 w-6 shrink-0"
            :class="step.done ? 'text-success' : 'text-base-content/20'"
          /><span class="text-sm font-bold" :class="step.done ? 'text-base-content/55 line-through' : ''">{{
            step.label
          }}</span></NuxtLink
        >
      </div>
    </section>

    <section aria-labelledby="overview-heading">
      <h2 id="overview-heading" class="sr-only">Inventory overview</h2>
      <div class="grid grid-cols-2 gap-3 lg:grid-cols-3 xl:grid-cols-6">
        <div v-for="stat in stats" :key="stat.label" class="catalog-panel rounded-2xl p-4 sm:p-5">
          <div class="flex items-start justify-between gap-3">
            <div>
              <p class="catalog-muted text-xs font-bold uppercase tracking-[0.08em]">
                {{ stat.label }}
              </p>
              <div v-if="statsPending" class="catalog-skeleton mt-3 h-9 w-24 rounded-lg" />
              <p v-else class="mt-2 text-2xl font-extrabold sm:text-3xl">
                <Currency v-if="stat.currency" :amount="stat.value" />
                <span v-else>{{ stat.value.toLocaleString() }}</span>
              </p>
            </div>
            <div class="grid h-10 w-10 shrink-0 place-items-center rounded-xl" :class="stat.tone">
              <component :is="stat.icon" class="h-5 w-5" />
            </div>
          </div>
        </div>
      </div>
    </section>

    <section
      v-if="statistics && (statistics.totalMissing || statistics.totalOverdue)"
      class="mt-6 rounded-2xl border border-warning/25 bg-warning/10 p-5"
    >
      <div class="flex items-start gap-3">
        <MdiAlertCircleOutline class="h-6 w-6 shrink-0 text-warning" />
        <div>
          <p class="font-black">Attention needed</p>
          <p class="mt-1 text-sm text-base-content/60">
            <span v-if="statistics.totalMissing"
              >{{ statistics.totalMissing }} missing item<span v-if="statistics.totalMissing !== 1">s</span></span
            ><span v-if="statistics.totalMissing && statistics.totalOverdue"> • </span
            ><span v-if="statistics.totalOverdue"
              >{{ statistics.totalOverdue }} overdue checkout<span v-if="statistics.totalOverdue !== 1">s</span></span
            >
          </p>
          <NuxtLink to="/items" class="mt-2 inline-block text-sm font-bold text-primary hover:underline"
            >Review inventory</NuxtLink
          >
        </div>
      </div>
    </section>

    <section class="mt-8 grid gap-6 xl:grid-cols-[minmax(0,1.65fr)_minmax(300px,0.8fr)]">
      <div class="min-w-0">
        <div class="mb-4 flex items-end justify-between gap-4">
          <div>
            <p class="catalog-eyebrow">Recently added</p>
            <h2 class="mt-1 text-xl font-bold">Latest inventory</h2>
          </div>
          <NuxtLink
            to="/items"
            class="inline-flex min-h-[44px] items-center gap-1 text-sm font-bold text-primary hover:underline"
            >View all <MdiArrowRight class="h-4 w-4"
          /></NuxtLink>
        </div>

        <div v-if="itemsPending" class="catalog-panel space-y-3 rounded-2xl p-5">
          <div v-for="i in 4" :key="i" class="catalog-skeleton h-14 rounded-xl" />
        </div>
        <BaseEmptyState
          v-else-if="!recentItems?.length"
          title="Your inventory is ready to grow"
          description="Add your first instrument, accessory, or piece of stage equipment to get started."
        >
          <template #icon><MdiPackageVariantClosed class="h-7 w-7" /></template>
          <template #actions><BaseButton to="/item/new" class="btn-primary">Add your first item</BaseButton></template>
        </BaseEmptyState>
        <ItemViewTable v-else-if="breakpoints.lg" :items="recentItems" />
        <div v-else class="grid gap-4 sm:grid-cols-2">
          <ItemCard v-for="item in recentItems" :key="item.id" :item="item" />
        </div>
      </div>

      <aside>
        <div class="mb-4">
          <p class="catalog-eyebrow">Shortcuts</p>
          <h2 class="mt-1 text-xl font-bold">Quick actions</h2>
        </div>
        <div class="catalog-panel overflow-hidden rounded-2xl">
          <NuxtLink
            v-for="action in quickActions"
            :key="action.label"
            :to="action.to"
            class="group flex min-h-[76px] items-center gap-4 border-b border-base-content/5 px-4 transition last:border-0 hover:bg-base-200/70"
          >
            <div
              class="grid h-11 w-11 shrink-0 place-items-center rounded-xl"
              :class="action.primary ? 'bg-primary text-primary-content' : 'bg-secondary text-secondary-content'"
            >
              <component :is="action.icon" class="h-5 w-5" />
            </div>
            <div class="min-w-0">
              <p class="text-sm font-bold">{{ action.label }}</p>
              <p class="catalog-muted mt-0.5 truncate text-xs">
                {{ action.description }}
              </p>
            </div>
            <MdiArrowRight
              class="ml-auto h-4 w-4 text-base-content/25 transition group-hover:translate-x-1 group-hover:text-primary"
            />
          </NuxtLink>
        </div>

        <div class="mt-5 rounded-2xl border border-primary/15 bg-primary/5 p-5">
          <MdiLightbulbOnOutline class="h-6 w-6 text-primary" />
          <h3 class="mt-3 text-sm font-bold">Make Catalog faster</h3>
          <p class="catalog-muted mt-1 text-xs leading-5">
            Press <kbd class="kbd kbd-xs">Ctrl K</kbd> from anywhere to jump to inventory search.
          </p>
        </div>
      </aside>
    </section>

    <section class="mt-10 grid gap-6 lg:grid-cols-2">
      <div>
        <div class="mb-4 flex items-end justify-between">
          <div>
            <p class="catalog-eyebrow">Organization</p>
            <h2 class="mt-1 text-xl font-bold">Top-level locations</h2>
          </div>
          <NuxtLink to="/locations" class="text-sm font-bold text-primary hover:underline">Manage</NuxtLink>
        </div>
        <BaseEmptyState
          v-if="!locationStore.parentLocations.length"
          title="No locations yet"
          description="Create places such as a band room, storage case, shelf, or stage area."
          ><template #icon><MdiMapMarkerOutline class="h-7 w-7" /></template
        ></BaseEmptyState>
        <div v-else class="grid gap-3 sm:grid-cols-2">
          <LocationCard
            v-for="location in locationStore.parentLocations.slice(0, 6)"
            :key="location.id"
            :location="location"
          />
        </div>
      </div>
      <div>
        <div class="mb-4 flex items-end justify-between">
          <div>
            <p class="catalog-eyebrow">Classification</p>
            <h2 class="mt-1 text-xl font-bold">Labels</h2>
          </div>
        </div>
        <BaseEmptyState
          v-if="!labelStore.labels.length"
          title="No labels yet"
          description="Labels make it easy to group instruments, accessories, condition, or ownership."
          ><template #icon><MdiTagMultipleOutline class="h-7 w-7" /></template
        ></BaseEmptyState>
        <div v-else class="catalog-panel flex min-h-[180px] flex-wrap content-start gap-2 rounded-2xl p-5">
          <LabelChip v-for="label in labelStore.labels" :key="label.id" size="lg" :label="label" />
        </div>
      </div>
    </section>
  </BaseContainer>
</template>
