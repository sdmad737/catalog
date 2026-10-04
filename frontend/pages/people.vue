<script setup lang="ts">
  import type { PersonOut } from "~~/lib/api/types/data-contracts";
  import MdiAccountPlusOutline from "~icons/mdi/account-plus-outline";
  import MdiAccountGroupOutline from "~icons/mdi/account-group-outline";

  definePageMeta({ middleware: ["auth"] });
  useHead({ title: "People & borrowers — Catalog" });

  const api = useUserApi();
  const auth = useAuthContext();
  const toast = useNotifier();
  const search = ref("");
  const modalOpen = ref(false);
  const editing = ref<PersonOut>();
  const busy = ref(false);
  const form = reactive({
    name: "",
    email: "",
    identifier: "",
    notes: "",
    active: true,
  });
  const canManage = computed(() => ["owner", "admin", "staff"].includes(auth.user?.role || ""));

  const { data: people, refresh } = useAsyncData("borrowers", async () => {
    const response = await api.people.list();
    if (response.error) {
      toast.error("Could not load borrowers.");
      return [];
    }
    return response.data;
  });

  const filteredPeople = computed(() => {
    const term = search.value.trim().toLowerCase();
    if (!term) return people.value || [];
    return (people.value || []).filter(person =>
      [person.name, person.email, person.identifier].some(value => value?.toLowerCase().includes(term))
    );
  });

  function openCreate() {
    editing.value = undefined;
    Object.assign(form, {
      name: "",
      email: "",
      identifier: "",
      notes: "",
      active: true,
    });
    modalOpen.value = true;
  }

  function openEdit(person: PersonOut) {
    editing.value = person;
    Object.assign(form, {
      name: person.name,
      email: person.email || "",
      identifier: person.identifier || "",
      notes: person.notes || "",
      active: person.active,
    });
    modalOpen.value = true;
  }

  async function save() {
    busy.value = true;
    const response = editing.value ? await api.people.update(editing.value.id, form) : await api.people.create(form);
    busy.value = false;
    if (response.error) return toast.error("Could not save this person.");
    modalOpen.value = false;
    await refresh();
    toast.success(editing.value ? "Borrower updated." : "Borrower added.");
  }
</script>

<template>
  <BaseContainer class="catalog-page max-w-none">
    <AppPageHeader
      eyebrow="Inventory assignments"
      title="People & borrowers"
      description="Manage students, members, or other people who can receive inventory. They do not get Catalog sign-in access."
    />

    <div
      class="mb-5 grid gap-4 rounded-2xl border border-base-content/10 bg-base-100 p-5 md:grid-cols-[1fr_auto] md:items-center"
    >
      <div class="flex gap-3">
        <div class="grid h-11 w-11 shrink-0 place-items-center rounded-xl bg-secondary text-secondary-content">
          <MdiAccountGroupOutline class="h-6 w-6" />
        </div>
        <div>
          <p class="font-black">Borrowers are separate from staff accounts</p>
          <p class="mt-1 text-sm text-base-content/55">
            A student can be assigned an instrument without a login or administrative access.
          </p>
        </div>
      </div>
      <BaseButton v-if="canManage" class="btn-primary" @click="openCreate"
        ><template #icon><MdiAccountPlusOutline /></template>Add person</BaseButton
      >
    </div>

    <BaseCard>
      <template #title>Borrower directory</template>
      <div class="border-t border-base-content/10 p-4">
        <input
          v-model="search"
          type="search"
          class="input input-bordered w-full max-w-lg"
          placeholder="Search name, email, or student/member ID"
        />
      </div>
      <div v-if="filteredPeople.length" class="overflow-x-auto border-t border-base-content/10">
        <table class="table">
          <thead>
            <tr>
              <th>Person</th>
              <th>Identifier</th>
              <th>Status</th>
              <th>Added</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="person in filteredPeople" :key="person.id">
              <td>
                <p class="font-black">{{ person.name }}</p>
                <p v-if="person.email" class="text-xs text-base-content/50">
                  {{ person.email }}
                </p>
              </td>
              <td>{{ person.identifier || "—" }}</td>
              <td>
                <span class="badge" :class="person.active ? 'badge-success' : 'badge-neutral'">{{
                  person.active ? "Active" : "Inactive"
                }}</span>
              </td>
              <td>
                <DateTime :date="person.createdAt" datetime-type="date" />
              </td>
              <td class="text-right">
                <button v-if="canManage" class="btn btn-ghost btn-sm" @click="openEdit(person)">Edit</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="border-t border-base-content/10 p-10 text-center">
        <p class="font-black">No borrowers yet</p>
        <p class="mx-auto mt-2 max-w-lg text-sm text-base-content/55">
          Add students, members, or clients here so staff can check inventory out to them from an item or NFC tap.
        </p>
        <BaseButton v-if="canManage" class="btn-primary btn-sm mt-4" @click="openCreate"
          >Add your first person</BaseButton
        >
      </div>
    </BaseCard>

    <BaseModal v-model="modalOpen">
      <template #title>{{ editing ? "Edit borrower" : "Add borrower" }}</template>
      <form class="space-y-3" @submit.prevent="save">
        <FormTextField v-model="form.name" label="Name" required />
        <FormTextField v-model="form.identifier" label="Student/member ID (optional)" />
        <FormTextField v-model="form.email" type="email" label="Email (optional)" />
        <FormTextArea v-model="form.notes" label="Private notes (optional)" />
        <label v-if="editing" class="flex cursor-pointer items-center gap-3 rounded-xl bg-base-200 p-3"
          ><input v-model="form.active" type="checkbox" class="toggle toggle-primary" /><span
            ><span class="block font-bold">Active borrower</span
            ><span class="text-xs text-base-content/50"
              >Inactive people remain in past checkout history but cannot receive new items.</span
            ></span
          ></label
        >
        <div class="flex justify-end">
          <BaseButton type="submit" class="btn-primary" :loading="busy">{{
            editing ? "Save changes" : "Add person"
          }}</BaseButton>
        </div>
      </form>
    </BaseModal>
  </BaseContainer>
</template>
