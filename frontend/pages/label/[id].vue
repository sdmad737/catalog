<script setup lang="ts">
  import MdiPencil from "~icons/mdi/pencil";
  import MdiDelete from "~icons/mdi/delete";
  import MdiTagOutline from "~icons/mdi/tag-outline";

  definePageMeta({
    middleware: ["auth"],
  });

  const route = useRoute();
  const api = useUserApi();
  const toast = useNotifier();

  const labelId = computed<string>(() => route.params.id as string);

  const { data: label } = useAsyncData(labelId.value, async () => {
    const { data, error } = await api.labels.get(labelId.value);
    if (error) {
      toast.error("Failed to load label");
      navigateTo("/home");
      return;
    }
    return data;
  });

  const confirm = useConfirm();

  async function confirmDelete() {
    const { isCanceled } = await confirm.open(
      "Are you sure you want to delete this label? This action cannot be undone."
    );

    if (isCanceled) {
      return;
    }

    const { error } = await api.labels.delete(labelId.value);

    if (error) {
      toast.error("Failed to delete label");
      return;
    }
    toast.success("Label deleted");
    navigateTo("/home");
  }

  const updateModal = ref(false);
  const updating = ref(false);
  const updateData = reactive({
    name: "",
    description: "",
    color: "",
  });

  function openUpdate() {
    updateData.name = label.value?.name || "";
    updateData.description = label.value?.description || "";
    updateModal.value = true;
  }

  async function update() {
    updating.value = true;
    const { error, data } = await api.labels.update(labelId.value, updateData);

    if (error) {
      updating.value = false;
      toast.error("Failed to update label");
      return;
    }

    toast.success("Label updated");
    label.value = data;
    updateModal.value = false;
    updating.value = false;
  }

  const items = computedAsync(async () => {
    if (!label.value) {
      return [];
    }

    const resp = await api.items.getAll({
      labels: [label.value.id],
    });

    if (resp.error) {
      toast.error("Failed to load items");
      return [];
    }

    return resp.data.items;
  });
</script>

<template>
  <BaseContainer class="catalog-page max-w-none">
    <BaseModal v-model="updateModal">
      <template #title> Update Label </template>
      <form v-if="label" @submit.prevent="update">
        <FormTextField v-model="updateData.name" :autofocus="true" label="Label Name" />
        <FormTextArea v-model="updateData.description" label="Label Description" />
        <div class="modal-action">
          <BaseButton type="submit" :loading="updating"> Update </BaseButton>
        </div>
      </form>
    </BaseModal>

    <BaseContainer v-if="label" cmp="div" class="max-w-none px-0">
      <div class="catalog-panel rounded-2xl p-5 sm:p-6">
        <header>
          <div class="flex flex-wrap items-start gap-4">
            <div class="avatar placeholder mb-auto">
              <div class="grid h-12 w-12 place-items-center rounded-2xl bg-primary/10 text-primary">
                <MdiTagOutline class="h-6 w-6" />
              </div>
            </div>
            <div>
              <p class="catalog-eyebrow">Label</p>
              <h1 class="mt-1 text-2xl font-extrabold">
                {{ label ? label.name : "" }}
              </h1>
              <div class="catalog-muted mt-2 flex flex-wrap gap-1 text-xs">
                <div>
                  Created
                  <DateTime :date="label?.createdAt" />
                </div>
              </div>
            </div>
            <div class="ml-auto mt-2 flex flex-wrap items-center justify-between gap-3">
              <div class="flex gap-2">
                <PageQRCode class="dropdown-left" />
                <BaseButton size="sm" @click="openUpdate">
                  <MdiPencil class="mr-1" />
                  Edit
                </BaseButton>
              </div>
              <BaseButton class="btn btn-ghost btn-sm text-error" @click="confirmDelete()">
                <MdiDelete class="mr-2" />
                Delete
              </BaseButton>
            </div>
          </div>
        </header>
        <Markdown
          v-if="label && label.description"
          class="catalog-muted mt-5 border-t border-base-content/5 pt-4 text-sm"
          :source="label.description"
        >
        </Markdown>
      </div>
      <section v-if="label && items">
        <ItemViewSelectable :items="items" />
      </section>
    </BaseContainer>
  </BaseContainer>
</template>
