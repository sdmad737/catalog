<template>
  <BaseModal v-model="dialog">
    <template #title> Import CSV File </template>
    <p>
      Import a CSV file containing your items, labels, and locations. See documentation for more information on the
      required format.
    </p>
    <div class="alert alert-warning shadow-lg mt-4">
      <div>
        <MdiAlertOutline class="mb-auto h-6 w-6 flex-shrink-0" />
        <span class="text-sm">
          Behavior for imports with existing import_refs has changed. If an import_ref is present in the CSV file, the
          item will be updated with the values in the CSV file.
        </span>
      </div>
    </div>

    <form @submit.prevent="submitCsvFile">
      <div class="flex flex-col gap-2 py-6">
        <input ref="importRef" type="file" class="hidden" accept=".csv,.tsv" @change="setFile" />

        <BaseButton type="button" @click="uploadCsv">
          <MdiUpload class="h-5 w-5 mr-2" />
          Upload
        </BaseButton>
        <p class="text-center pt-4 -mb-5">
          {{ importCsv?.name }}
        </p>
      </div>

      <div class="modal-action">
        <BaseButton type="submit" :disabled="!importCsv"> Submit </BaseButton>
      </div>
    </form>
  </BaseModal>
</template>

<script setup lang="ts">
  import MdiUpload from "~icons/mdi/upload";
  import MdiAlertOutline from "~icons/mdi/alert-outline";
  type Props = {
    modelValue: boolean;
  };

  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
  });

  const emit = defineEmits(["update:modelValue"]);

  const dialog = useVModel(props, "modelValue", emit);

  const api = useUserApi();
  const toast = useNotifier();

  const importCsv = ref<File | null>(null);
  const importLoading = ref(false);
  const importRef = ref<HTMLInputElement>();
  whenever(
    () => !dialog.value,
    () => {
      importCsv.value = null;
    }
  );

  function setFile(e: Event) {
    const result = e.target as HTMLInputElement;
    if (!result.files || result.files.length === 0) {
      return;
    }

    importCsv.value = result.files[0];
  }

  function uploadCsv() {
    importRef.value?.click();
  }

  async function submitCsvFile() {
    if (!importCsv.value) {
      toast.error("Please select a file to import.");
      return;
    }

    importLoading.value = true;

    const { error } = await api.items.import(importCsv.value);

    if (error) {
      toast.error("Import failed. Please try again later.");
    }

    // Reset
    dialog.value = false;
    importLoading.value = false;
    importCsv.value = null;

    if (importRef.value) {
      importRef.value.value = "";
    }

    toast.success("Import successful!");
  }
</script>
