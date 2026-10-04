<script setup lang="ts">
  import { useLabelStore } from "~~/stores/labels";
  import { useLocationStore } from "~~/stores/locations";
  import MdiViewDashboardOutline from "~icons/mdi/view-dashboard-outline";
  import MdiPackageVariantClosed from "~icons/mdi/package-variant-closed";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";
  import MdiAccountOutline from "~icons/mdi/account-outline";
  import MdiTools from "~icons/mdi/tools";
  import MdiPlus from "~icons/mdi/plus";
  import MdiMagnify from "~icons/mdi/magnify";
  import MdiMenu from "~icons/mdi/menu";
  import MdiClose from "~icons/mdi/close";
  import MdiChevronDown from "~icons/mdi/chevron-down";
  import MdiLogoutVariant from "~icons/mdi/logout-variant";
  import MdiWeatherNight from "~icons/mdi/weather-night";
  import MdiWhiteBalanceSunny from "~icons/mdi/white-balance-sunny";
  import MdiTagOutline from "~icons/mdi/tag-outline";
  import MdiAccountGroupOutline from "~icons/mdi/account-group-outline";
  import MdiAccountMultipleOutline from "~icons/mdi/account-multiple-outline";

  useFormatCurrency();

  const route = useRoute();
  const router = useRouter();
  const authCtx = useAuthContext();
  const api = useUserApi();
  const { theme, setTheme } = useTheme();
  const searchInput = ref<HTMLInputElement>();
  const searchQuery = ref("");
  const mobileOpen = ref(false);
  const username = computed(() => authCtx.user?.name || "Catalog user");
  const initials = computed(() =>
    username.value
      .split(/\s+/)
      .map(part => part[0])
      .join("")
      .slice(0, 2)
      .toUpperCase()
  );
  const isDark = computed(() => theme.value === "catalog-dark");

  const modals = reactive({ item: false, location: false, label: false });

  const nav = computed(() => [
    {
      name: "Dashboard",
      to: "/home",
      icon: MdiViewDashboardOutline,
      match: (path: string) => path === "/home",
    },
    {
      name: "Inventory",
      to: "/items",
      icon: MdiPackageVariantClosed,
      match: (path: string) => path.startsWith("/item") || path === "/items",
    },
    ...(authCtx.user?.role === "owner" || authCtx.user?.role === "admin"
      ? [
          {
            name: "Staff",
            to: "/staff",
            icon: MdiAccountGroupOutline,
            match: (path: string) => path === "/staff",
          },
        ]
      : []),
    {
      name: "People",
      to: "/people",
      icon: MdiAccountMultipleOutline,
      match: (path: string) => path === "/people",
    },
    {
      name: "Locations",
      to: "/locations",
      icon: MdiMapMarkerOutline,
      match: (path: string) => path.startsWith("/location") || path === "/locations",
    },
    {
      name: "Profile",
      to: "/profile",
      icon: MdiAccountOutline,
      match: (path: string) => path === "/profile",
    },
    {
      name: "Tools",
      to: "/tools",
      icon: MdiTools,
      match: (path: string) => path.startsWith("/tools") || path.startsWith("/reports"),
    },
  ]);
  const mobileNav = computed(() =>
    nav.value.filter(item => ["Dashboard", "Inventory", "People", "Locations", "Profile"].includes(item.name))
  );

  function submitSearch() {
    const q = searchQuery.value.trim();
    router.push({ path: "/items", query: q ? { q } : {} });
    mobileOpen.value = false;
  }

  function handleShortcut(event: KeyboardEvent) {
    if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k") {
      event.preventDefault();
      searchInput.value?.focus();
    }
  }

  onMounted(() => window.addEventListener("keydown", handleShortcut));
  onBeforeUnmount(() => window.removeEventListener("keydown", handleShortcut));

  async function logout() {
    await authCtx.logout(api);
    navigateTo("/");
  }

  const labelStore = useLabelStore();
  const locationStore = useLocationStore();
  onServerEvent(ServerEvent.LabelMutation, () => labelStore.refresh());
  onServerEvent(ServerEvent.LocationMutation, () => {
    locationStore.refreshChildren();
    locationStore.refreshParents();
    locationStore.refreshTree();
  });
  onServerEvent(ServerEvent.ItemMutation, () => {
    locationStore.refreshChildren();
    locationStore.refreshParents();
  });
</script>

<template>
  <div class="min-h-screen bg-transparent lg:flex">
    <ModalConfirm />
    <ItemCreateModal v-model="modals.item" />
    <LabelCreateModal v-model="modals.label" />
    <LocationCreateModal v-model="modals.location" />
    <AppToast />

    <button
      v-if="mobileOpen"
      class="fixed inset-0 z-40 bg-slate-950/45 backdrop-blur-sm lg:hidden"
      aria-label="Close navigation"
      @click="mobileOpen = false"
    />

    <aside
      class="fixed inset-y-0 left-0 z-50 flex w-[280px] -translate-x-full flex-col border-r border-white/10 bg-neutral px-4 py-5 text-neutral-content transition-transform duration-200 lg:sticky lg:translate-x-0"
      :class="{ 'translate-x-0': mobileOpen }"
    >
      <div class="flex items-center gap-3 px-2">
        <div
          class="grid h-11 w-11 place-items-center rounded-xl bg-primary text-primary-content shadow-lg shadow-primary/20"
        >
          <AppLogo class="h-7 w-7" />
        </div>
        <div>
          <p class="font-display text-xl font-extrabold leading-none">Catalog</p>
          <p class="mt-1 text-xs text-neutral-content/55">Inventory workspace</p>
        </div>
        <button
          class="btn btn-ghost btn-sm btn-square ml-auto text-neutral-content lg:hidden"
          aria-label="Close navigation"
          @click="mobileOpen = false"
        >
          <MdiClose class="h-5 w-5" />
        </button>
      </div>

      <button class="btn btn-primary mt-8 w-full justify-start gap-3" @click="modals.item = true">
        <MdiPlus class="h-5 w-5" /> Add inventory item
      </button>

      <p class="mb-2 mt-8 px-3 text-[11px] font-bold uppercase tracking-[0.16em] text-neutral-content/40">Workspace</p>
      <nav class="space-y-1" aria-label="Primary navigation">
        <NuxtLink
          v-for="item in nav"
          :key="item.to"
          :to="item.to"
          class="group flex min-h-[44px] items-center gap-3 rounded-xl px-3 text-sm font-semibold text-neutral-content/65 transition hover:bg-white/[0.07] hover:text-neutral-content"
          :class="{
            'bg-white/10 text-white shadow-inner': item.match(route.path),
          }"
          @click="mobileOpen = false"
        >
          <component :is="item.icon" class="h-5 w-5" :class="{ 'text-primary': item.match(route.path) }" />
          {{ item.name }}
          <span v-if="item.match(route.path)" class="ml-auto h-1.5 w-1.5 rounded-full bg-primary" />
        </NuxtLink>
      </nav>

      <div class="mt-auto">
        <div class="mb-3 rounded-xl border border-white/10 bg-white/5 p-3">
          <p class="text-xs text-neutral-content/45">Signed in as</p>
          <p class="mt-1 truncate text-sm font-semibold">{{ username }}</p>
        </div>
        <button
          class="flex min-h-[44px] w-full items-center gap-3 rounded-xl px-3 text-sm font-semibold text-neutral-content/60 transition hover:bg-white/[0.07] hover:text-white"
          @click="logout"
        >
          <MdiLogoutVariant class="h-5 w-5" /> Sign out
        </button>
        <p class="mt-5 px-3 text-[11px] text-neutral-content/35">Designed by Shahin Madantschi</p>
      </div>
    </aside>

    <div class="min-w-0 flex-1">
      <header class="sticky top-0 z-30 border-b border-base-content/5 bg-base-100/80 backdrop-blur-xl">
        <div class="flex h-[72px] items-center gap-3 px-4 md:px-7">
          <button class="btn btn-ghost btn-square lg:hidden" aria-label="Open navigation" @click="mobileOpen = true">
            <MdiMenu class="h-6 w-6" />
          </button>
          <NuxtLink to="/home" class="mr-auto flex items-center gap-2 font-display text-lg font-extrabold lg:hidden">
            <AppLogo class="h-7 w-7 text-primary" /> Catalog
          </NuxtLink>

          <form class="relative hidden w-full max-w-xl md:block" role="search" @submit.prevent="submitSearch">
            <MdiMagnify
              class="pointer-events-none absolute left-3.5 top-1/2 h-5 w-5 -translate-y-1/2 text-base-content/40"
            />
            <input
              ref="searchInput"
              v-model="searchQuery"
              class="input h-11 w-full border-base-content/10 bg-base-200/65 pl-11 pr-20 text-sm"
              type="search"
              aria-label="Search inventory"
              placeholder="Search inventory, serials, and locations…"
            />
            <kbd
              class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 rounded-md border border-base-content/10 bg-base-100 px-2 py-1 text-[10px] font-bold text-base-content/40"
              >Ctrl K</kbd
            >
          </form>

          <div class="ml-auto flex items-center gap-2">
            <button
              class="btn btn-ghost btn-square"
              :aria-label="isDark ? 'Use light theme' : 'Use dark theme'"
              @click="setTheme(isDark ? 'catalog' : 'catalog-dark')"
            >
              <MdiWhiteBalanceSunny v-if="isDark" class="h-5 w-5" />
              <MdiWeatherNight v-else class="h-5 w-5" />
            </button>
            <div class="dropdown dropdown-end">
              <label tabindex="0" class="btn btn-primary gap-2">
                <MdiPlus class="h-5 w-5" /><span class="hidden sm:inline">Create</span
                ><MdiChevronDown class="h-4 w-4" />
              </label>
              <ul
                tabindex="0"
                class="dropdown-content menu mt-2 w-56 rounded-xl border border-base-content/10 bg-base-100 p-2 shadow-2xl"
              >
                <li>
                  <button @click="modals.item = true"><MdiPackageVariantClosed />Inventory item</button>
                </li>
                <li>
                  <button @click="modals.location = true"><MdiMapMarkerOutline />Location</button>
                </li>
                <li>
                  <button @click="modals.label = true"><MdiTagOutline />Label</button>
                </li>
              </ul>
            </div>
            <NuxtLink
              to="/profile"
              class="grid h-10 w-10 place-items-center rounded-full bg-secondary text-xs font-bold text-secondary-content"
              aria-label="Open profile"
              >{{ initials }}</NuxtLink
            >
          </div>
        </div>
      </header>

      <main><slot /></main>
    </div>

    <nav
      class="fixed inset-x-3 bottom-3 z-30 grid grid-cols-5 rounded-2xl border border-base-content/10 bg-base-100/95 p-1.5 shadow-2xl backdrop-blur-xl lg:hidden"
      aria-label="Mobile navigation"
    >
      <NuxtLink
        v-for="item in mobileNav"
        :key="item.to"
        :to="item.to"
        class="flex min-h-[50px] flex-col items-center justify-center gap-1 rounded-xl text-[10px] font-semibold text-base-content/45"
        :class="{
          'bg-secondary text-secondary-content': item.match(route.path),
        }"
      >
        <component :is="item.icon" class="h-5 w-5" />{{ item.name }}
      </NuxtLink>
    </nav>
  </div>
</template>
