<script setup lang="ts">
  import type { ItemOut } from "~~/lib/api/types/data-contracts";
  import type { ItemTag } from "~~/lib/api/classes/tags";
  import { route as apiRoute } from "~~/lib/api/base";
  import MdiApple from "~icons/mdi/apple";
  import MdiAndroid from "~icons/mdi/android";
  import MdiCellphoneNfc from "~icons/mdi/cellphone-nfc";
  import MdiCheckCircle from "~icons/mdi/check-circle";
  import MdiChevronLeft from "~icons/mdi/chevron-left";
  import MdiContentCopy from "~icons/mdi/content-copy";
  import MdiDownload from "~icons/mdi/download";
  import MdiHelpCircleOutline from "~icons/mdi/help-circle-outline";
  import MdiPrinter from "~icons/mdi/printer";
  import MdiRefresh from "~icons/mdi/refresh";
  import MdiShieldKeyOutline from "~icons/mdi/shield-key-outline";

  const props = defineProps<{ item: ItemOut }>();
  const api = useUserApi();
  const toast = useNotifier();
  const busy = ref(false);
  const tags = ref<ItemTag[]>([]);
  const origin = ref("");
  const canWriteNfc = ref(false);
  const setupOpen = ref(false);
  const setupStage = ref<"device" | "program" | "verify">("device");
  const device = ref<"iphone" | "android" | "manual">();
  let pollTimer: ReturnType<typeof setInterval> | undefined;

  const activeTag = computed(() => tags.value.find(tag => !tag.revokedAt));
  const isReady = computed(() => Boolean(activeTag.value?.verifiedAt));
  const publicUrl = computed(() => (activeTag.value ? `${origin.value}/t/${activeTag.value.token}` : ""));
  const qrUrl = computed(() =>
    activeTag.value
      ? apiRoute(`/tags/${activeTag.value.token}/qrcode`, {
          origin: origin.value,
        })
      : ""
  );
  const shortToken = computed(() => activeTag.value?.token.slice(-8).toUpperCase() || "");

  async function load(showVerifiedToast = false) {
    const wasVerified = isReady.value;
    const response = await api.tags.list(props.item.id);
    if (response.error) return toast.error("Could not load the physical tag status.");
    tags.value = response.data.items || [];
    if (showVerifiedToast && !wasVerified && isReady.value) {
      toast.success("Physical tag verified and ready to use.");
    }
  }

  async function createLink(replace = false) {
    if (replace && !window.confirm("Replace this tag? Its current NFC and QR link will stop working immediately."))
      return;
    busy.value = true;
    const response = replace ? await api.tags.replace(props.item.id) : await api.tags.assign(props.item.id);
    busy.value = false;
    if (response.error) return toast.error("Could not create the secure tag link.");
    await load();
    device.value = undefined;
    setupStage.value = "device";
    setupOpen.value = true;
    toast.success(replace ? "Replacement link created. Reprogram the physical tag." : "Secure link created.");
  }

  async function revoke() {
    if (!activeTag.value || !window.confirm("Revoke this tag? Its NFC and QR link will stop working immediately."))
      return;
    busy.value = true;
    const response = await api.tags.revoke(props.item.id, activeTag.value.id);
    busy.value = false;
    if (response.error) return toast.error("Could not revoke the physical tag.");
    setupOpen.value = false;
    await load();
    toast.success("Physical tag revoked.");
  }

  async function copyUrl() {
    await navigator.clipboard.writeText(publicUrl.value);
    toast.success("NFC link copied.");
  }

  async function writeNfc() {
    try {
      const Writer = (
        window as unknown as {
          NDEFReader?: new () => { write: (value: unknown) => Promise<void> };
        }
      ).NDEFReader;
      if (!Writer) return;
      await new Writer().write({
        records: [{ recordType: "url", data: publicUrl.value }],
      });
      setupStage.value = "verify";
      toast.success("Link written. Remove the tag, then tap it again to verify.");
    } catch {
      toast.error("The NFC tag was not written. Keep it near the phone and try again.");
    }
  }

  function chooseDevice(value: "iphone" | "android" | "manual") {
    device.value = value;
    setupStage.value = "program";
  }

  function printLabel() {
    window.print();
  }

  onMounted(() => {
    const appOrigin = new URL(window.location.origin);
    if (appOrigin.protocol === "http:" && appOrigin.hostname !== "localhost" && appOrigin.hostname !== "127.0.0.1") {
      appOrigin.protocol = "https:";
    }
    origin.value = appOrigin.origin;
    canWriteNfc.value = "NDEFReader" in window;
    load();
    pollTimer = setInterval(() => {
      if (activeTag.value && !isReady.value) load(true);
    }, 2500);
  });

  onUnmounted(() => {
    if (pollTimer) clearInterval(pollTimer);
  });
</script>

<template>
  <BaseCard class="overflow-hidden">
    <template #title>Physical tag</template>
    <div class="border-t border-base-content/10 p-5 sm:p-6">
      <div v-if="!activeTag" class="flex flex-col gap-5 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex gap-4">
          <div class="grid h-12 w-12 shrink-0 place-items-center rounded-2xl bg-primary/10 text-primary">
            <MdiCellphoneNfc class="h-6 w-6" />
          </div>
          <div>
            <h3 class="font-extrabold">Connect this item to the physical world</h3>
            <p class="mt-1 max-w-xl text-sm text-base-content/60">
              Set up a revocable NFC link with a printable QR fallback. No private item details are stored on the tag.
            </p>
          </div>
        </div>
        <BaseButton class="btn-primary shrink-0" :loading="busy" @click="createLink(false)">Set up tag</BaseButton>
      </div>

      <div v-else>
        <div class="grid gap-6 lg:grid-cols-[1fr_auto]">
          <div>
            <div class="flex flex-wrap items-center gap-2">
              <span v-if="isReady" class="badge badge-success gap-1.5"><MdiCheckCircle />Ready</span>
              <span v-else class="badge badge-warning gap-1.5"
                ><span class="loading loading-ring loading-xs"></span>Setup incomplete</span
              >
              <span class="text-xs font-semibold text-base-content/45">TAG •••• {{ shortToken }}</span>
            </div>
            <h3 class="mt-3 text-lg font-extrabold">
              {{ isReady ? "Tap or scan to open this item" : "Program and verify the physical tag" }}
            </h3>
            <dl class="mt-4 grid gap-3 text-sm sm:grid-cols-3">
              <div>
                <dt class="text-xs font-bold uppercase tracking-wider text-base-content/40">Secure link</dt>
                <dd class="mt-1">Created <DateTime :date="activeTag.createdAt" /></dd>
              </div>
              <div>
                <dt class="text-xs font-bold uppercase tracking-wider text-base-content/40">Physical verification</dt>
                <dd class="mt-1">
                  <DateTime v-if="activeTag.verifiedAt" :date="activeTag.verifiedAt" /><span v-else>Waiting</span>
                </dd>
              </div>
              <div>
                <dt class="text-xs font-bold uppercase tracking-wider text-base-content/40">Last opened</dt>
                <dd class="mt-1">
                  <DateTime v-if="activeTag.lastScannedAt" :date="activeTag.lastScannedAt" /><span v-else>Never</span>
                </dd>
              </div>
            </dl>
            <div class="mt-5 flex flex-wrap gap-2">
              <BaseButton v-if="!isReady" class="btn-primary btn-sm" @click="setupOpen = true"
                ><template #icon><MdiCellphoneNfc /></template>Continue setup</BaseButton
              >
              <BaseButton v-else class="btn-secondary btn-sm" @click="copyUrl"
                ><template #icon><MdiContentCopy /></template>Copy link</BaseButton
              >
              <button
                v-if="isReady"
                class="btn btn-ghost btn-sm"
                @click="
                  setupOpen = true;
                  setupStage = 'device';
                "
              >
                <MdiCellphoneNfc />Reprogram instructions
              </button>
              <a class="btn btn-ghost btn-sm" :href="qrUrl" :download="`catalog-${item.assetId || item.id}-qr.jpg`"
                ><MdiDownload />Download QR</a
              >
              <button class="btn btn-ghost btn-sm" @click="printLabel"><MdiPrinter />Print label</button>
            </div>
            <div class="mt-5 flex flex-wrap gap-2 border-t border-base-content/10 pt-4">
              <button class="btn btn-ghost btn-xs" :disabled="busy" @click="createLink(true)">
                <MdiRefresh />Replace tag
              </button>
              <button class="btn btn-ghost btn-xs text-error" :disabled="busy" @click="revoke">
                <MdiShieldKeyOutline />Revoke
              </button>
            </div>
          </div>

          <div
            class="catalog-tag-print flex w-full max-w-[270px] items-center gap-3 rounded-2xl border border-base-content/15 bg-white p-3 text-slate-900 shadow-sm"
          >
            <img :src="qrUrl" class="h-24 w-24" alt="Item tag QR code" />
            <div class="min-w-0">
              <p class="text-[10px] font-black uppercase tracking-[0.18em] text-violet-700">Catalog</p>
              <p class="mt-1 truncate text-sm font-black">{{ item.name }}</p>
              <p v-if="item.assetId && item.assetId !== '000-000'" class="text-xs font-semibold text-slate-500">
                Asset {{ item.assetId }}
              </p>
              <p class="mt-2 text-xs font-extrabold">Tap or scan</p>
            </div>
          </div>
        </div>

        <section v-if="setupOpen && !isReady" class="mt-6 rounded-2xl border border-primary/20 bg-primary/5 p-4 sm:p-6">
          <div class="mb-5 flex items-start justify-between gap-3">
            <div>
              <p class="catalog-eyebrow">Physical setup</p>
              <h4 class="mt-1 text-xl font-black">
                {{
                  setupStage === "device"
                    ? "Choose your setup device"
                    : setupStage === "program"
                      ? "Program the NFC tag"
                      : "Verify the tag"
                }}
              </h4>
            </div>
            <button class="btn btn-ghost btn-sm" @click="setupOpen = false">Close</button>
          </div>

          <div class="mb-6 grid grid-cols-5 gap-2" aria-label="Tag setup progress">
            <div
              v-for="(label, index) in ['Link', 'Device', 'Program', 'Verify', 'Ready']"
              :key="label"
              class="text-center"
            >
              <div
                class="mx-auto h-1.5 rounded-full"
                :class="
                  index === 0 ||
                  (setupStage === 'device' ? index <= 1 : setupStage === 'program' ? index <= 2 : index <= 3)
                    ? 'bg-primary'
                    : 'bg-base-300'
                "
              ></div>
              <span class="mt-1 hidden text-[10px] font-bold uppercase tracking-wide text-base-content/45 sm:block">{{
                label
              }}</span>
            </div>
          </div>

          <div v-if="setupStage === 'device'" class="grid gap-3 sm:grid-cols-3">
            <button
              class="rounded-2xl border border-base-content/15 bg-base-100 p-5 text-left transition hover:border-primary hover:shadow-md"
              @click="chooseDevice('iphone')"
            >
              <MdiApple class="h-7 w-7 text-primary" /><strong class="mt-3 block">iPhone</strong
              ><span class="mt-1 block text-xs text-base-content/55">Use a compatible NFC-writing app.</span>
            </button>
            <button
              class="rounded-2xl border border-base-content/15 bg-base-100 p-5 text-left transition hover:border-primary hover:shadow-md"
              @click="chooseDevice('android')"
            >
              <MdiAndroid class="h-7 w-7 text-primary" /><strong class="mt-3 block">Android</strong
              ><span class="mt-1 block text-xs text-base-content/55"
                >Use Web NFC when available, or any NFC writer.</span
              >
            </button>
            <button
              class="rounded-2xl border border-base-content/15 bg-base-100 p-5 text-left transition hover:border-primary hover:shadow-md"
              @click="chooseDevice('manual')"
            >
              <MdiHelpCircleOutline class="h-7 w-7 text-primary" /><strong class="mt-3 block">Other / Manual</strong
              ><span class="mt-1 block text-xs text-base-content/55">Copy the URL as an NDEF URL record.</span>
            </button>
          </div>

          <div v-else-if="setupStage === 'program'" class="grid gap-5 lg:grid-cols-[1fr_auto]">
            <div>
              <ol v-if="device === 'iphone'" class="list-decimal space-y-2 pl-5 text-sm leading-6">
                <li>Open a compatible NFC-writing app such as NFC Tools.</li>
                <li>Select <strong>Write</strong>.</li>
                <li>Add a <strong>URL / URI</strong> record.</li>
                <li>Paste the Catalog link below.</li>
                <li>Tap <strong>Write</strong> and hold the top of the iPhone near the NFC tag.</li>
              </ol>
              <div v-else-if="device === 'android'" class="space-y-3 text-sm leading-6">
                <p v-if="canWriteNfc">
                  This browser supports NFC writing. Press
                  <strong>Write NFC now</strong>, approve access, and hold the blank tag against your phone.
                </p>
                <p v-else>
                  This browser cannot write NFC directly. Copy the link and write it as a URL/URI record using a
                  compatible NFC-writing app.
                </p>
              </div>
              <p v-else class="text-sm leading-6">
                Write the link below to the tag as a standard NDEF URL/URI record. Do not permanently lock the tag while
                testing.
              </p>
              <div class="mt-4 rounded-xl border border-base-content/10 bg-base-100 p-3 font-mono text-xs break-all">
                {{ publicUrl }}
              </div>
              <div class="mt-3 flex flex-wrap gap-2">
                <BaseButton class="btn-primary" @click="copyUrl"
                  ><template #icon><MdiContentCopy /></template>Copy NFC link</BaseButton
                ><BaseButton v-if="device === 'android' && canWriteNfc" class="btn-secondary" @click="writeNfc"
                  ><template #icon><MdiCellphoneNfc /></template>Write NFC now</BaseButton
                ><button class="btn btn-ghost" @click="setupStage = 'verify'">I programmed the tag</button>
              </div>
            </div>
            <img :src="qrUrl" class="h-36 w-36 rounded-xl bg-white p-2" alt="QR fallback for this physical tag" />
          </div>

          <div v-else class="text-center">
            <div class="mx-auto grid h-16 w-16 place-items-center rounded-full bg-warning/15 text-warning">
              <MdiCellphoneNfc class="h-8 w-8" />
            </div>
            <h5 class="mt-4 text-lg font-black">Tap the programmed tag now</h5>
            <p class="mx-auto mt-2 max-w-lg text-sm leading-6 text-base-content/60">
              Remove the tag from your phone, then tap it again. Opening its Catalog page verifies the physical setup.
              You can also scan the printed QR fallback.
            </p>
            <div class="mt-4 inline-flex items-center gap-2 rounded-full bg-base-100 px-4 py-2 text-sm font-bold">
              <span class="loading loading-ring loading-sm text-primary"></span>Waiting for a real tag open…
            </div>
            <div class="mt-5">
              <button class="btn btn-ghost btn-sm" @click="setupStage = 'program'">
                <MdiChevronLeft />Back to programming
              </button>
            </div>
          </div>
        </section>

        <section
          v-else-if="setupOpen && isReady"
          class="mt-6 rounded-2xl border border-success/25 bg-success/10 p-6 text-center"
        >
          <MdiCheckCircle class="mx-auto h-12 w-12 text-success" />
          <p class="catalog-eyebrow mt-3">Tag verified</p>
          <h4 class="mt-1 text-xl font-black">{{ item.name }}</h4>
          <p v-if="item.assetId && item.assetId !== '000-000'" class="mt-1 font-mono text-sm">
            Asset {{ item.assetId }}
          </p>
          <p class="mt-3 font-bold text-success">Ready to use</p>
          <button class="btn btn-success btn-sm mt-4" @click="setupOpen = false">Done</button>
        </section>
      </div>
    </div>
  </BaseCard>
</template>

<style>
  @media print {
    body * {
      visibility: hidden !important;
    }
    .catalog-tag-print,
    .catalog-tag-print * {
      visibility: visible !important;
    }
    .catalog-tag-print {
      position: fixed;
      inset: 0 auto auto 0;
      width: 2.75in !important;
      box-shadow: none !important;
    }
  }
</style>
