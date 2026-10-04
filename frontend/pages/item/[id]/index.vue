<script setup lang="ts">
  import type { AnyDetail, Detail, Details } from "~~/components/global/DetailsSection/types";
  import { filterZeroValues } from "~~/components/global/DetailsSection/types";
  import type { ItemAttachment } from "~~/lib/api/types/data-contracts";
  import MdiClose from "~icons/mdi/close";
  import MdiPackageVariant from "~icons/mdi/package-variant";
  import MdiPlus from "~icons/mdi/plus";
  import MdiMinus from "~icons/mdi/minus";
  import MdiDownload from "~icons/mdi/download";
  import MdiPencilOutline from "~icons/mdi/pencil-outline";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";

  definePageMeta({
    middleware: ["auth"],
  });

  const route = useRoute();
  const api = useUserApi();
  const toast = useNotifier();

  const itemId = computed<string>(() => route.params.id as string);
  const preferences = useViewPreferences();

  const hasNested = computed<boolean>(() => {
    return route.fullPath.split("/").at(-1) !== itemId.value;
  });

  const { data: item, refresh } = useAsyncData(itemId.value, async () => {
    const { data, error } = await api.items.get(itemId.value);
    if (error) {
      toast.error("Failed to load item");
      navigateTo("/home");
      return;
    }
    return data;
  });
  useHead({
    title: computed(() => (item.value ? `${item.value.name} — Catalog` : "Item — Catalog")),
  });
  onMounted(() => {
    refresh();
  });

  const lastRoute = ref(route.fullPath);
  watchEffect(() => {
    if (lastRoute.value.endsWith("edit")) {
      refresh();
    }

    lastRoute.value = route.fullPath;
  });

  async function adjustQuantity(amount: number) {
    if (!item.value) {
      return;
    }

    const newQuantity = item.value.quantity + amount;
    if (newQuantity < 0) {
      toast.error("Quantity cannot be negative");
      return;
    }

    const resp = await api.items.patch(item.value.id, {
      id: item.value.id,
      quantity: newQuantity,
    });

    if (resp.error) {
      toast.error("Failed to adjust quantity");
      return;
    }

    item.value.quantity = newQuantity;
  }

  type FilteredAttachments = {
    attachments: ItemAttachment[];
    warranty: ItemAttachment[];
    manuals: ItemAttachment[];
    receipts: ItemAttachment[];
  };

  type Photo = {
    src: string;
  };

  const photos = computed<Photo[]>(() => {
    return (
      item.value?.attachments.reduce((acc, cur) => {
        if (cur.type === "photo") {
          acc.push({
            // @ts-expect-error - it's impossible for this to be null at this point
            src: api.authURL(`/items/${item.value.id}/attachments/${cur.id}`),
          });
        }
        return acc;
      }, [] as Photo[]) || []
    );
  });

  const attachments = computed<FilteredAttachments>(() => {
    if (!item.value) {
      return {
        attachments: [],
        manuals: [],
        warranty: [],
        receipts: [],
      };
    }

    return item.value.attachments.reduce(
      (acc, attachment) => {
        if (attachment.type === "photo") {
          return acc;
        }
        if (attachment.type === "warranty") {
          acc.warranty.push(attachment);
        } else if (attachment.type === "manual") {
          acc.manuals.push(attachment);
        } else if (attachment.type === "receipt") {
          acc.receipts.push(attachment);
        } else {
          acc.attachments.push(attachment);
        }
        return acc;
      },
      {
        attachments: [] as ItemAttachment[],
        warranty: [] as ItemAttachment[],
        manuals: [] as ItemAttachment[],
        receipts: [] as ItemAttachment[],
      }
    );
  });

  const assetID = computed<Details>(() => {
    if (!item.value) {
      return [];
    }

    if (item.value?.assetId === "000-000") {
      return [];
    }

    return [
      {
        name: "Asset ID",
        text: item.value?.assetId,
      },
    ];
  });

  const itemDetails = computed<Details>(() => {
    if (!item.value) {
      return [];
    }

    const ret: Details = [
      {
        name: "Quantity",
        text: item.value?.quantity,
        slot: "quantity",
      },
      {
        name: "Serial Number",
        text: item.value?.serialNumber,
        copyable: true,
      },
      {
        name: "Model Number",
        text: item.value?.modelNumber,
        copyable: true,
      },
      {
        name: "Manufacturer",
        text: item.value?.manufacturer,
        copyable: true,
      },
      {
        name: "Insured",
        text: item.value?.insured ? "Yes" : "No",
      },
      {
        name: "Notes",
        type: "markdown",
        text: item.value?.notes,
      },
      ...assetID.value,
      ...item.value.fields.map(field => {
        /**
         * Support Special URL Syntax
         */
        const url = maybeUrl(field.textValue);
        if (url.isUrl) {
          return {
            type: "link",
            name: field.name,
            text: url.text,
            href: url.url,
          } as AnyDetail;
        }

        return {
          name: field.name,
          text: field.textValue,
        };
      }),
    ];

    if (!preferences.value.showEmpty) {
      return filterZeroValues(ret);
    }

    return ret;
  });

  const showAttachments = computed(() => {
    if (preferences.value?.showEmpty) {
      return true;
    }

    return (
      attachments.value.attachments.length > 0 ||
      attachments.value.warranty.length > 0 ||
      attachments.value.manuals.length > 0 ||
      attachments.value.receipts.length > 0
    );
  });

  const attachmentDetails = computed(() => {
    const details: Detail[] = [];

    const push = (name: string) => {
      details.push({
        name,
        text: "",
        slot: name.toLowerCase(),
      });
    };

    if (attachments.value.attachments.length > 0) {
      push("Attachments");
    }

    if (attachments.value.warranty.length > 0) {
      push("Warranty");
    }

    if (attachments.value.manuals.length > 0) {
      push("Manuals");
    }

    if (attachments.value.receipts.length > 0) {
      push("Receipts");
    }

    return details;
  });

  const showWarranty = computed(() => {
    if (preferences.value.showEmpty) {
      return true;
    }
    return validDate(item.value?.warrantyExpires);
  });

  const warrantyDetails = computed(() => {
    const details: Details = [
      {
        name: "Lifetime Warranty",
        text: item.value?.lifetimeWarranty ? "Yes" : "No",
      },
    ];

    if (item.value?.lifetimeWarranty) {
      details.push({
        name: "Warranty Expires",
        text: "N/A",
      });
    } else {
      details.push({
        name: "Warranty Expires",
        text: item.value?.warrantyExpires || "",
        type: "date",
        date: true,
      });
    }

    details.push({
      name: "Warranty Details",
      type: "markdown",
      text: item.value?.warrantyDetails || "",
    });

    if (!preferences.value.showEmpty) {
      return filterZeroValues(details);
    }

    return details;
  });

  const showPurchase = computed(() => {
    if (preferences.value.showEmpty) {
      return true;
    }
    return item.value?.purchaseFrom || item.value?.purchasePrice !== "0";
  });

  const purchaseDetails = computed<Details>(() => {
    const v: Details = [
      {
        name: "Purchased From",
        text: item.value?.purchaseFrom || "",
      },
      {
        name: "Purchase Price",
        text: item.value?.purchasePrice || "",
        type: "currency",
      },
      {
        name: "Purchase Date",
        text: item.value?.purchaseTime || "",
        type: "date",
        date: true,
      },
    ];

    if (!preferences.value.showEmpty) {
      return filterZeroValues(v);
    }

    return v;
  });

  const showSold = computed(() => {
    if (preferences.value.showEmpty) {
      return true;
    }
    return item.value?.soldTo || item.value?.soldPrice !== "0";
  });

  const soldDetails = computed<Details>(() => {
    const v: Details = [
      {
        name: "Sold To",
        text: item.value?.soldTo || "",
      },
      {
        name: "Sold Price",
        text: item.value?.soldPrice || "",
        type: "currency",
      },
      {
        name: "Sold At",
        text: item.value?.soldTime || "",
        type: "date",
        date: true,
      },
    ];

    if (!preferences.value.showEmpty) {
      return filterZeroValues(v);
    }

    return v;
  });

  const refDialog = ref<HTMLDialogElement>();
  const dialoged = reactive({
    src: "",
  });

  function openDialog(img: Photo) {
    // @ts-ignore - I don't know why this is happening
    refDialog.value?.showModal();
    dialoged.src = img.src;
  }

  function closeDialog() {
    // @ts-ignore - I don't know why this is happening
    refDialog.value?.close();
  }

  const refDialogBody = ref<HTMLDivElement>();
  onClickOutside(refDialogBody, () => {
    closeDialog();
  });

  const currentPath = computed(() => {
    return route.path;
  });

  const tabs = computed(() => {
    return [
      {
        id: "details",
        name: "Details",
        to: `/item/${itemId.value}`,
      },
      {
        id: "log",
        name: "Maintenance",
        to: `/item/${itemId.value}/maintenance`,
      },
      {
        id: "edit",
        name: "Edit",
        to: `/item/${itemId.value}/edit`,
      },
    ];
  });

  const fullpath = computedAsync(async () => {
    if (!item.value) {
      return [];
    }

    const resp = await api.items.fullpath(item.value.id);
    if (resp.error) {
      toast.error("Failed to load item");
      return [];
    }

    return resp.data;
  });

  const items = computedAsync(async () => {
    if (!item.value) {
      return [];
    }

    const resp = await api.items.getAll({
      parentIds: [item.value.id],
    });

    if (resp.error) {
      toast.error("Failed to load items");
      return [];
    }

    return resp.data.items;
  });
</script>

<template>
  <BaseContainer v-if="item" class="catalog-page max-w-none">
    <Title>{{ item.name }}</Title>
    <dialog ref="refDialog" class="z-[999] fixed bg-transparent">
      <div ref="refDialogBody" class="relative">
        <div class="absolute right-0 -mt-3 -mr-3 sm:-mt-4 sm:-mr-4 space-x-1">
          <a class="btn btn-sm sm:btn-md btn-primary btn-circle" :href="dialoged.src" download>
            <MdiDownload class="h-5 w-5" />
          </a>
          <button class="btn btn-sm sm:btn-md btn-primary btn-circle" @click="closeDialog()">
            <MdiClose class="h-5 w-5" />
          </button>
        </div>

        <img class="max-w-[80vw] max-h-[80vh]" :src="dialoged.src" />
      </div>
    </dialog>

    <section class="mb-7">
      <div v-if="fullpath && fullpath.length > 0" class="breadcrumbs mb-3 py-0 text-xs text-base-content/50">
        <ul>
          <li><NuxtLink to="/items">Inventory</NuxtLink></li>
          <li v-for="part in fullpath" :key="part.id">
            <NuxtLink :to="`/${part.type}/${part.id}`">{{ part.name }}</NuxtLink>
          </li>
        </ul>
      </div>
      <div class="catalog-panel rounded-2xl p-5 sm:p-6">
        <header class="flex flex-col gap-5 sm:flex-row sm:items-start">
          <div class="grid h-14 w-14 shrink-0 place-items-center rounded-2xl bg-primary/10 text-primary">
            <MdiPackageVariant class="h-7 w-7" />
          </div>
          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div>
                <p class="catalog-eyebrow">Inventory item</p>
                <h1 class="mt-1 text-2xl font-extrabold sm:text-3xl">
                  {{ item.name }}
                </h1>
              </div>
              <BaseButton :to="`/item/${item.id}/edit`" class="btn-primary btn-sm"
                ><template #icon><MdiPencilOutline class="h-4 w-4" /></template>Edit item</BaseButton
              >
            </div>
            <div class="catalog-muted mt-3 flex flex-wrap items-center gap-x-4 gap-y-2 text-xs">
              <span v-if="item.location" class="inline-flex items-center gap-1.5"
                ><MdiMapMarkerOutline class="h-4 w-4" />{{ item.location.name }}</span
              >
              <span>Updated <DateTime :date="item.updatedAt" /></span>
              <span v-if="item.archived" class="badge badge-neutral badge-sm">Archived</span>
              <span class="badge badge-sm" :class="item.status === 'missing' ? 'badge-error' : 'badge-success'">{{
                item.status === "missing" ? "Missing" : "Available"
              }}</span>
              <span
                v-if="item.condition && item.condition !== 'unknown'"
                class="badge badge-outline badge-sm capitalize"
                >{{ item.condition }}</span
              >
            </div>
            <Markdown
              v-if="item.description"
              class="catalog-muted mt-4 max-w-3xl text-sm leading-6"
              :source="item.description"
            />
            <p v-else class="catalog-muted mt-4 text-sm italic">No description has been added.</p>
          </div>
        </header>
      </div>

      <div class="mt-4 catalog-panel rounded-2xl p-4 sm:p-5">
        <ItemOperations :item="item" @updated="item = $event" />
      </div>

      <div class="mt-4 overflow-x-auto">
        <div class="inline-flex rounded-xl border border-base-content/10 bg-base-100 p-1">
          <NuxtLink
            v-for="t in tabs"
            :key="t.id"
            :to="t.to"
            class="min-h-[40px] whitespace-nowrap rounded-lg px-4 py-2 text-sm font-bold transition"
            :class="
              t.to === currentPath
                ? 'bg-secondary text-secondary-content'
                : 'text-base-content/55 hover:bg-base-200 hover:text-base-content'
            "
          >
            {{ t.name }}
          </NuxtLink>
        </div>
      </div>
    </section>

    <section>
      <div class="space-y-6">
        <BaseCard v-if="!hasNested" collapsable>
          <template #title> Details </template>
          <template #title-actions>
            <div class="flex flex-wrap justify-between items-center mt-2 gap-4">
              <label class="label cursor-pointer">
                <input v-model="preferences.showEmpty" type="checkbox" class="toggle toggle-primary" />
                <span class="label-text ml-4"> Show Empty </span>
              </label>
              <PageQRCode />
            </div>
          </template>
          <DetailsSection :details="itemDetails">
            <template #quantity="{ detail }">
              {{ detail.text }}
              <span
                class="opacity-0 group-hover:opacity-100 ml-4 my-0 duration-75 transition-opacity inline-flex gap-2"
              >
                <button class="btn btn-circle btn-xs" @click="adjustQuantity(-1)">
                  <MdiMinus class="h-3 w-3" />
                </button>
                <button class="btn btn-circle btn-xs" @click="adjustQuantity(1)">
                  <MdiPlus class="h-3 w-3" />
                </button>
              </span>
            </template>
          </DetailsSection>
        </BaseCard>

        <NuxtPage :item="item" :page-key="itemId" />
        <template v-if="!hasNested">
          <BaseCard v-if="photos && photos.length > 0">
            <template #title> Photos </template>
            <div
              class="container border-t border-gray-300 p-4 flex flex-wrap gap-2 mx-auto max-h-[500px] overflow-y-scroll scroll-bg"
            >
              <button v-for="(img, i) in photos" :key="i" @click="openDialog(img)">
                <img class="rounded max-h-[200px]" :src="img.src" />
              </button>
            </div>
          </BaseCard>

          <ItemPhysicalTag :item="item" />

          <BaseCard v-if="showAttachments" collapsable>
            <template #title> Attachments </template>
            <DetailsSection v-if="attachmentDetails.length > 0" :details="attachmentDetails">
              <template #manuals>
                <ItemAttachmentsList
                  v-if="attachments.manuals.length > 0"
                  :attachments="attachments.manuals"
                  :item-id="item.id"
                />
              </template>
              <template #attachments>
                <ItemAttachmentsList
                  v-if="attachments.attachments.length > 0"
                  :attachments="attachments.attachments"
                  :item-id="item.id"
                />
              </template>
              <template #warranty>
                <ItemAttachmentsList
                  v-if="attachments.warranty.length > 0"
                  :attachments="attachments.warranty"
                  :item-id="item.id"
                />
              </template>
              <template #receipts>
                <ItemAttachmentsList
                  v-if="attachments.receipts.length > 0"
                  :attachments="attachments.receipts"
                  :item-id="item.id"
                />
              </template>
            </DetailsSection>
            <div v-else>
              <p class="text-base-content/70 px-6 pb-4">No attachments found</p>
            </div>
          </BaseCard>

          <BaseCard v-if="showPurchase" collapsable>
            <template #title> Purchase Details </template>
            <DetailsSection :details="purchaseDetails" />
          </BaseCard>

          <BaseCard v-if="showWarranty" collapsable>
            <template #title> Warranty Details </template>
            <DetailsSection :details="warrantyDetails" />
          </BaseCard>

          <BaseCard v-if="showSold" collapsable>
            <template #title> Sold Details </template>
            <DetailsSection :details="soldDetails" />
          </BaseCard>
        </template>
      </div>
    </section>

    <section v-if="items && items.length > 0" class="my-6">
      <ItemViewSelectable :items="items" />
    </section>
  </BaseContainer>
</template>

<style lang="css" scoped>
  /* Style dialog background */
  dialog::backdrop {
    background: rgba(0, 0, 0, 0.5);
  }

  .scroll-bg::-webkit-scrollbar {
    width: 0.5rem;
  }

  .scroll-bg::-webkit-scrollbar-thumb {
    border-radius: 0.25rem;
    @apply bg-base-300;
  }
</style>
