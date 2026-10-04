<script setup lang="ts">
  import type { ItemOut, LocationSummary, PersonOut } from "~~/lib/api/types/data-contracts";
  import MdiAlertOutline from "~icons/mdi/alert-outline";
  import MdiCheckCircleOutline from "~icons/mdi/check-circle-outline";
  import MdiHistory from "~icons/mdi/history";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";
  import MdiPencilOutline from "~icons/mdi/pencil-outline";
  import MdiWrenchOutline from "~icons/mdi/wrench-outline";
  import MdiAccountArrowRightOutline from "~icons/mdi/account-arrow-right-outline";
  import MdiKeyboardReturn from "~icons/mdi/keyboard-return";

  const props = defineProps<{ item: ItemOut; compact?: boolean; canOperate?: boolean }>();
  const emit = defineEmits<{ updated: [item: ItemOut] }>();
  const api = useUserApi();
  const toast = useNotifier();
  const auth = useAuthContext();
  const dialog = ref<"missing" | "found" | "checkout" | "return">();
  const busy = ref(false);
  const note = ref("");
  const lastKnownLocation = ref("");
  const foundLocation = ref<LocationSummary | null>(props.item.location || null);
  const condition = ref(props.item.condition || "unknown");
  const people = ref<PersonOut[]>([]);
  const selectedPersonId = ref("");
  const dueDate = ref("");
  const historyOpen = ref(false);
  const { data: history, refresh: refreshHistory } = useAsyncData(
    `item-history-${props.item.id}`,
    async () => {
      const response = await api.items.history(props.item.id);
      return response.error ? [] : response.data;
    },
    { immediate: false }
  );

  const canOperate = computed(() => props.canOperate ?? ["owner", "admin", "staff"].includes(auth.user?.role || ""));

  function openMissing() {
    lastKnownLocation.value = props.item.location?.name || "";
    note.value = "";
    dialog.value = "missing";
  }

  function openFound() {
    foundLocation.value = props.item.location || null;
    condition.value = props.item.condition || "unknown";
    note.value = "";
    dialog.value = "found";
  }

  async function openCheckout() {
    const response = await api.people.list();
    if (response.error) return toast.error("Could not load borrowers.");
    people.value = response.data.filter(person => person.active);
    selectedPersonId.value = "";
    dueDate.value = "";
    condition.value = props.item.condition || "unknown";
    note.value = "";
    dialog.value = "checkout";
  }

  function openReturn() {
    foundLocation.value = props.item.location || null;
    condition.value = props.item.condition || "unknown";
    note.value = "";
    dialog.value = "return";
  }

  async function markMissing() {
    busy.value = true;
    const response = await api.items.markMissing(props.item.id, {
      lastKnownLocation: lastKnownLocation.value,
      note: note.value,
    });
    busy.value = false;
    if (response.error) return toast.error("Could not mark this item missing.");
    dialog.value = undefined;
    emit("updated", response.data);
    toast.success("Item marked missing and recorded in its history.");
  }

  async function markFound() {
    busy.value = true;
    const response = await api.items.markFound(props.item.id, {
      locationId: foundLocation.value?.id,
      condition: condition.value,
      note: note.value,
    });
    busy.value = false;
    if (response.error) return toast.error("Could not mark this item found.");
    dialog.value = undefined;
    emit("updated", response.data);
    toast.success("Item marked found and returned to inventory.");
  }

  async function checkoutItem() {
    if (!selectedPersonId.value) return toast.error("Choose who is receiving this item.");
    busy.value = true;
    const response = await api.people.checkout(props.item.id, {
      personId: selectedPersonId.value,
      dueAt: dueDate.value ? `${dueDate.value}T23:59:59Z` : undefined,
      condition: condition.value,
      note: note.value,
    });
    busy.value = false;
    if (response.error) return toast.error("This item could not be checked out.");
    dialog.value = undefined;
    emit("updated", response.data);
    toast.success(`Checked out to ${response.data.activeCheckout?.person.name}.`);
  }

  async function returnItem() {
    busy.value = true;
    const response = await api.people.returnItem(props.item.id, {
      locationId: foundLocation.value?.id,
      condition: condition.value,
      note: note.value,
    });
    busy.value = false;
    if (response.error) return toast.error("This item could not be returned.");
    dialog.value = undefined;
    emit("updated", response.data);
    toast.success("Item returned and made available.");
  }

  async function showHistory() {
    await refreshHistory();
    historyOpen.value = true;
  }
</script>

<template>
  <div>
    <div v-if="item.status === 'missing'" class="mb-4 rounded-2xl border border-error/25 bg-error/10 p-4">
      <div class="flex gap-3">
        <MdiAlertOutline class="h-6 w-6 shrink-0 text-error" />
        <div>
          <p class="font-black">This item is marked missing</p>
          <p class="mt-1 text-sm text-base-content/60">
            <span v-if="item.lastKnownLocation">Last known at {{ item.lastKnownLocation }}.</span>
            Reported
            <DateTime v-if="item.missingSince" :date="item.missingSince" />.
          </p>
          <p v-if="item.missingNote" class="mt-2 text-sm">
            {{ item.missingNote }}
          </p>
        </div>
      </div>
    </div>

    <div
      v-else-if="item.status === 'checked_out' && item.activeCheckout"
      class="mb-4 rounded-2xl border border-info/25 bg-info/10 p-4"
    >
      <p class="text-xs font-black uppercase tracking-wide text-info">Checked out</p>
      <p class="mt-1 font-black">Assigned to {{ item.activeCheckout.person.name }}</p>
      <p class="mt-1 text-sm text-base-content/60">
        <span v-if="item.activeCheckout.dueAt"
          >Due <DateTime :date="item.activeCheckout.dueAt" datetime-type="date" /></span
        ><span v-else>No due date</span>
      </p>
    </div>

    <div v-if="canOperate" class="flex flex-wrap gap-2">
      <BaseButton v-if="item.status === 'missing'" class="btn-success" size="sm" @click="openFound"
        ><template #icon><MdiCheckCircleOutline /></template>Mark found</BaseButton
      >
      <BaseButton v-else-if="item.status === 'available'" class="btn-primary" size="sm" @click="openCheckout"
        ><template #icon><MdiAccountArrowRightOutline /></template>Check out</BaseButton
      >
      <BaseButton v-else-if="item.status === 'checked_out'" class="btn-primary" size="sm" @click="openReturn"
        ><template #icon><MdiKeyboardReturn /></template>Return item</BaseButton
      >
      <BaseButton v-if="item.status !== 'missing'" class="btn-error btn-outline" size="sm" @click="openMissing"
        ><template #icon><MdiAlertOutline /></template>Mark missing</BaseButton
      >
      <BaseButton :to="`/item/${item.id}/edit`" size="sm"
        ><template #icon><MdiMapMarkerOutline /></template>Move location</BaseButton
      >
      <BaseButton :to="`/item/${item.id}/maintenance?new=damage`" size="sm"
        ><template #icon><MdiAlertOutline /></template>Report damage</BaseButton
      >
      <BaseButton :to="`/item/${item.id}/maintenance?new=maintenance`" size="sm"
        ><template #icon><MdiWrenchOutline /></template>Add maintenance</BaseButton
      >
      <BaseButton :to="`/item/${item.id}/edit`" size="sm"
        ><template #icon><MdiPencilOutline /></template>Edit item</BaseButton
      >
      <BaseButton size="sm" @click="showHistory"
        ><template #icon><MdiHistory /></template>View history</BaseButton
      >
    </div>

    <BaseModal :model-value="dialog === 'missing'" @update:model-value="dialog = undefined">
      <template #title>Mark {{ item.name }} missing</template>
      <form class="space-y-3" @submit.prevent="markMissing">
        <FormTextField v-model="lastKnownLocation" label="Last known location" /><FormTextArea
          v-model="note"
          label="Note (optional)"
          placeholder="What was checked, and when was it last seen?"
        />
        <div class="flex justify-end">
          <BaseButton type="submit" class="btn-error" :loading="busy">Confirm missing</BaseButton>
        </div>
      </form>
    </BaseModal>

    <BaseModal :model-value="dialog === 'found'" @update:model-value="dialog = undefined">
      <template #title>Mark {{ item.name }} found</template>
      <form class="space-y-3" @submit.prevent="markFound">
        <LocationSelector v-model="foundLocation" /><label class="form-control"
          ><span class="label-text mb-2 text-sm font-bold">Condition when found</span
          ><select v-model="condition" class="select select-bordered">
            <option value="unknown">Not checked</option>
            <option value="excellent">Excellent</option>
            <option value="good">Good</option>
            <option value="fair">Fair</option>
            <option value="poor">Poor</option>
            <option value="damaged">Damaged</option>
          </select></label
        ><FormTextArea v-model="note" label="Note (optional)" placeholder="Where and how was it found?" />
        <div class="flex justify-end">
          <BaseButton type="submit" class="btn-success" :loading="busy">Confirm found</BaseButton>
        </div>
      </form>
    </BaseModal>

    <BaseModal :model-value="dialog === 'checkout'" @update:model-value="dialog = undefined">
      <template #title>Check out {{ item.name }}</template>
      <form class="space-y-4" @submit.prevent="checkoutItem">
        <div class="rounded-xl bg-base-200 p-3 text-sm">
          <span class="font-bold">Asset {{ item.assetId }}</span
          ><span class="mx-2 text-base-content/30">•</span>Current condition:
          <span class="capitalize">{{ item.condition }}</span>
        </div>
        <label class="form-control"
          ><span class="label-text mb-2 text-sm font-bold">Borrower</span
          ><select v-model="selectedPersonId" class="select select-bordered" required>
            <option value="" disabled>Choose a person</option>
            <option v-for="person in people" :key="person.id" :value="person.id">
              {{ person.name }}<template v-if="person.identifier"> — {{ person.identifier }}</template>
            </option>
          </select></label
        >
        <p v-if="!people.length" class="rounded-xl border border-warning/25 bg-warning/10 p-3 text-sm">
          No active borrowers yet.
          <NuxtLink to="/people" class="font-bold underline">Add a person</NuxtLink>
          first.
        </p>
        <FormTextField v-model="dueDate" type="date" label="Due date (optional)" />
        <label class="form-control"
          ><span class="label-text mb-2 text-sm font-bold">Condition at checkout</span
          ><select v-model="condition" class="select select-bordered">
            <option value="unknown">Not checked</option>
            <option value="excellent">Excellent</option>
            <option value="good">Good</option>
            <option value="fair">Fair</option>
            <option value="poor">Poor</option>
            <option value="damaged">Damaged</option>
          </select></label
        >
        <FormTextArea
          v-model="note"
          label="Note (optional)"
          placeholder="Case, accessories, or other checkout details"
        />
        <div class="flex items-center justify-between">
          <NuxtLink to="/people" class="text-sm font-bold text-primary">Manage borrowers</NuxtLink
          ><BaseButton type="submit" class="btn-primary" :loading="busy" :disabled="!people.length"
            >Confirm checkout</BaseButton
          >
        </div>
      </form>
    </BaseModal>

    <BaseModal :model-value="dialog === 'return'" @update:model-value="dialog = undefined">
      <template #title>Return {{ item.name }}</template>
      <form class="space-y-4" @submit.prevent="returnItem">
        <div v-if="item.activeCheckout" class="rounded-xl bg-base-200 p-3 text-sm">
          Returning from
          <span class="font-bold">{{ item.activeCheckout.person.name }}</span>
        </div>
        <LocationSelector v-model="foundLocation" />
        <label class="form-control"
          ><span class="label-text mb-2 text-sm font-bold">Condition on return</span
          ><select v-model="condition" class="select select-bordered">
            <option value="unknown">Not checked</option>
            <option value="excellent">Excellent</option>
            <option value="good">Good</option>
            <option value="fair">Fair</option>
            <option value="poor">Poor</option>
            <option value="damaged">Damaged</option>
          </select></label
        >
        <FormTextArea
          v-model="note"
          label="Return note (optional)"
          placeholder="Condition, missing accessories, or follow-up needed"
        />
        <div class="flex justify-end">
          <BaseButton type="submit" class="btn-primary" :loading="busy">Confirm return</BaseButton>
        </div>
      </form>
    </BaseModal>

    <BaseModal v-model="historyOpen"
      ><template #title>Item history</template>
      <div v-if="history?.length" class="space-y-4">
        <article v-for="event in history" :key="event.id" class="border-l-2 border-primary/30 pl-4">
          <p class="font-bold">{{ event.summary }}</p>
          <p class="mt-1 text-xs text-base-content/45">
            <DateTime :date="event.createdAt" />
            <span v-if="event.actorName">• {{ event.actorName }}</span>
          </p>
          <p v-if="event.note" class="mt-2 text-sm text-base-content/65">
            {{ event.note }}
          </p>
        </article>
      </div>
      <p v-else class="text-sm text-base-content/55">No operational history has been recorded yet.</p></BaseModal
    >
  </div>
</template>
