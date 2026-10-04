<script setup lang="ts">
  import type { PublicItemTag, ResolvedItemTag } from "~~/lib/api/classes/tags";
  import MdiAlertCircleOutline from "~icons/mdi/alert-circle-outline";
  import MdiArrowRight from "~icons/mdi/arrow-right";
  import MdiCellphoneNfc from "~icons/mdi/cellphone-nfc";
  import MdiCheckDecagram from "~icons/mdi/check-decagram";
  import MdiLockOutline from "~icons/mdi/lock-outline";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";
  import MdiPackageVariantClosed from "~icons/mdi/package-variant-closed";

  definePageMeta({ layout: "empty" });

  const route = useRoute();
  const token = computed(() => String(route.params.token || ""));
  const publicApi = usePublicApi();
  const auth = useAuthContext();
  const userApi = useUserApi();
  const loading = ref(true);
  const state = ref<"active" | "revoked" | "invalid">("active");
  const publicItem = ref<PublicItemTag>();
  const resolved = ref<ResolvedItemTag>();

  useHead({
    title: computed(() => (publicItem.value?.name ? `${publicItem.value.name} — Catalog` : "Physical tag — Catalog")),
    meta: [{ name: "robots", content: "noindex, nofollow" }],
  });

  async function load() {
    loading.value = true;
    const response = await publicApi.tag(token.value);
    if (response.error) {
      state.value = "invalid";
      loading.value = false;
      return;
    }
    publicItem.value = response.data;
    state.value = response.data.status;
    if (state.value === "active" && auth.isAuthorized()) {
      const staff = await userApi.tags.resolve(token.value);
      if (!staff.error) resolved.value = staff.data;
    }
    loading.value = false;
  }

  onMounted(load);
</script>

<template>
  <main class="catalog-tag-page catalog-grid-pattern min-h-screen px-4 py-6 sm:grid sm:place-items-center sm:py-12">
    <div class="mx-auto w-full max-w-lg">
      <NuxtLink
        to="/"
        class="mb-6 inline-flex items-center gap-2 rounded-xl bg-base-100/70 px-3 py-2 text-sm font-black shadow-sm backdrop-blur"
      >
        <span class="grid h-7 w-7 place-items-center rounded-lg bg-primary text-primary-content">C</span>Catalog
      </NuxtLink>

      <section v-if="loading" class="catalog-panel overflow-hidden rounded-[1.75rem] p-6 sm:p-8">
        <div class="catalog-skeleton h-14 w-14 rounded-2xl"></div>
        <div class="catalog-skeleton mt-6 h-4 w-28 rounded"></div>
        <div class="catalog-skeleton mt-3 h-9 w-3/4 rounded-lg"></div>
        <div class="catalog-skeleton mt-8 h-24 rounded-2xl"></div>
      </section>

      <section v-else-if="state === 'invalid'" class="catalog-panel rounded-[1.75rem] p-6 sm:p-8">
        <div class="grid h-14 w-14 place-items-center rounded-2xl bg-error/10 text-error">
          <MdiAlertCircleOutline class="h-7 w-7" />
        </div>
        <p class="catalog-eyebrow mt-6">Tag not found</p>
        <h1 class="mt-2 text-3xl font-black tracking-tight">This tag isn’t connected</h1>
        <p class="catalog-muted mt-3 leading-7">
          The link may be incomplete, or the tag has never been assigned. Ask the equipment manager to check it.
        </p>
      </section>

      <section v-else-if="state === 'revoked'" class="catalog-panel rounded-[1.75rem] p-6 sm:p-8">
        <div class="grid h-14 w-14 place-items-center rounded-2xl bg-warning/15 text-warning">
          <MdiLockOutline class="h-7 w-7" />
        </div>
        <p class="catalog-eyebrow mt-6">Inactive tag</p>
        <h1 class="mt-2 text-3xl font-black tracking-tight">This tag was replaced</h1>
        <p class="catalog-muted mt-3 leading-7">
          This link is no longer active. Use the item’s newer NFC tag or QR label, or ask the equipment manager for
          help.
        </p>
      </section>

      <section v-else class="catalog-panel overflow-hidden rounded-[1.75rem]">
        <div class="bg-gradient-to-br from-primary/15 via-secondary/5 to-transparent p-6 sm:p-8">
          <div class="flex items-start justify-between gap-4">
            <div
              class="grid h-14 w-14 place-items-center rounded-2xl bg-primary text-primary-content shadow-lg shadow-primary/20"
            >
              <MdiPackageVariantClosed class="h-7 w-7" />
            </div>
            <span class="badge badge-success gap-1.5"><MdiCheckDecagram />Verified item</span>
          </div>
          <p class="catalog-eyebrow mt-6">Catalog equipment</p>
          <h1 class="mt-2 text-3xl font-black tracking-tight sm:text-4xl">
            {{ resolved?.item.name || publicItem?.name }}
          </h1>
          <p
            v-if="publicItem?.assetId && publicItem.assetId !== '000-000'"
            class="mt-2 font-mono text-sm font-bold text-base-content/45"
          >
            Asset {{ publicItem.assetId }}
          </p>
          <p v-if="publicItem?.organization" class="mt-4 text-sm font-semibold text-base-content/55">
            Belongs to {{ publicItem.organization }}
          </p>
        </div>

        <div v-if="resolved" class="space-y-5 border-t border-base-content/10 p-6 sm:p-8">
          <div class="flex flex-wrap gap-2">
            <span
              class="badge"
              :class="
                resolved.item.status === 'missing'
                  ? 'badge-error'
                  : resolved.item.archived
                    ? 'badge-neutral'
                    : 'badge-success'
              "
              >{{
                resolved.item.status === "missing" ? "Missing" : resolved.item.archived ? "Archived" : "Available"
              }}</span
            >
            <span v-if="resolved.item.insured" class="badge badge-outline">Insured</span>
            <span class="badge badge-outline">Qty {{ resolved.item.quantity }}</span>
          </div>
          <div v-if="resolved.item.location" class="flex items-center gap-3 rounded-2xl bg-base-200/70 p-4">
            <MdiMapMarkerOutline class="h-6 w-6 text-primary" />
            <div>
              <p class="text-xs font-bold uppercase tracking-wider text-base-content/40">Current location</p>
              <p class="mt-0.5 font-extrabold">
                {{ resolved.item.location.name }}
              </p>
            </div>
          </div>
          <p v-if="resolved.item.description" class="catalog-muted text-sm leading-6">
            {{ resolved.item.description }}
          </p>
          <ItemOperations
            :item="resolved.item"
            compact
            :can-operate="resolved.canOperate"
            @updated="resolved.item = $event"
          />
          <NuxtLink :to="`/item/${resolved.item.id}`" class="btn btn-primary btn-block"
            >Open full record <MdiArrowRight
          /></NuxtLink>
        </div>

        <div v-else class="border-t border-base-content/10 p-6 sm:p-8">
          <div
            v-if="publicItem?.itemStatus === 'missing'"
            class="mb-4 rounded-2xl border border-error/25 bg-error/10 p-4"
          >
            <p class="font-black text-error">This item has been reported missing.</p>
            <p class="mt-1 text-sm leading-6 text-base-content/60">
              If you found it, please return it to the organization shown above. Catalog does not expose private staff
              contact details.
            </p>
          </div>
          <div class="flex gap-3 rounded-2xl bg-base-200/70 p-4">
            <MdiLockOutline class="mt-0.5 h-5 w-5 shrink-0 text-primary" />
            <p class="text-sm leading-6 text-base-content/65">
              Item location and internal details are only visible to authorized staff.
            </p>
          </div>
          <NuxtLink :to="{ path: '/', query: { redirect: route.fullPath } }" class="btn btn-primary btn-block mt-5"
            >Staff sign in <MdiArrowRight
          /></NuxtLink>
        </div>
      </section>

      <p class="mt-5 flex items-center justify-center gap-2 text-center text-xs font-semibold text-base-content/40">
        <MdiCellphoneNfc />This tag stores only a secure Catalog link.
      </p>
    </div>
  </main>
</template>
