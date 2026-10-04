<template>
  <div>
    <AppImportDialog v-model="modals.import" />
    <BaseContainer class="catalog-page max-w-none">
      <AppPageHeader
        eyebrow="Utilities"
        title="Tools & reports"
        description="Export your data, create printable asset labels, and run carefully scoped inventory maintenance."
      />
      <div class="grid gap-5 xl:grid-cols-2">
        <BaseCard>
          <template #title>
            <BaseSectionHeader>
              <MdiFileChart class="mr-2" />
              <span>Reports</span>
              <template #description>Create printable labels and inventory summaries.</template>
            </BaseSectionHeader>
          </template>
          <div class="border-t px-6 pb-3 border-gray-300 divide-gray-300 divide-y">
            <DetailAction @action="navigateTo('/reports/label-generator')">
              <template #title>Asset ID labels</template>
              Create a printable PDF for an asset ID range. Print labels ahead of time, then apply them as new equipment
              arrives.
              <template #button>
                Label Generator
                <MdiArrowRight class="ml-2" />
              </template>
            </DetailAction>
            <DetailAction @action="getBillOfMaterials()">
              <template #title>Bill of Materials</template>
              Download a spreadsheet-ready TSV summary with the essential item and pricing information.
              <template #button> Generate BOM </template>
            </DetailAction>
          </div>
        </BaseCard>
        <BaseCard>
          <template #title>
            <BaseSectionHeader>
              <MdiDatabase class="mr-2" />
              <span>Import & export</span>
              <template #description> Move Catalog inventory in or out with its standard CSV format. </template>
            </BaseSectionHeader>
          </template>
          <div class="border-t px-6 pb-3 border-gray-300 divide-gray-300 divide-y">
            <DetailAction @action="modals.import = true">
              <template #title>Import Inventory</template>
              Add items from a Catalog CSV without overwriting existing inventory.
            </DetailAction>
            <DetailAction @action="getExportTSV()">
              <template #title>Export Inventory</template>
              Download every inventory item in Catalog's standard CSV format.
            </DetailAction>
          </div>
        </BaseCard>
        <BaseCard>
          <template #title>
            <BaseSectionHeader>
              <MdiAlert class="mr-2" />
              <span>Inventory maintenance</span>
              <template #description>
                Repair or normalize inventory data in bulk. Review each confirmation carefully.
              </template>
            </BaseSectionHeader>
          </template>
          <div class="border-t px-6 pb-3 border-gray-300 divide-gray-300 divide-y">
            <DetailAction @action="ensureAssetIDs">
              <template #title>Ensure Asset IDs</template>
              Assign the next available asset ID to every item that does not have one, ordered by creation date.
            </DetailAction>
            <DetailAction @action="ensureImportRefs">
              <template #title>Ensure import references</template>
              Generate an eight-character import reference for every item that does not have one.
            </DetailAction>
            <DetailAction @action="resetItemDateTimes">
              <template #title> Zero Item Date Times</template>
              Normalizes stored date fields to the start of their selected day. Use this when imported dates display in
              an adjacent day because they contain an unintended time value.
            </DetailAction>
            <DetailAction @action="setPrimaryPhotos">
              <template #title> Set Primary Photos </template>
              Sets the first available photo as the primary image for items that do not currently have one.
            </DetailAction>
          </div>
        </BaseCard>
      </div>
    </BaseContainer>
  </div>
</template>

<script setup lang="ts">
  import MdiFileChart from "~icons/mdi/file-chart";
  import MdiArrowRight from "~icons/mdi/arrow-right";
  import MdiDatabase from "~icons/mdi/database";
  import MdiAlert from "~icons/mdi/alert";

  definePageMeta({
    middleware: ["auth"],
  });
  useHead({
    title: "Tools & reports — Catalog",
  });

  const modals = ref({
    item: false,
    location: false,
    label: false,
    import: false,
  });

  const api = useUserApi();
  const confirm = useConfirm();
  const notify = useNotifier();

  function getBillOfMaterials() {
    const url = api.reports.billOfMaterialsURL();
    window.open(url, "_blank");
  }

  function getExportTSV() {
    const url = api.items.exportURL();
    window.open(url, "_blank");
  }

  async function ensureAssetIDs() {
    const { isCanceled } = await confirm.open(
      "Are you sure you want to ensure all assets have an ID? This can take a while and cannot be undone."
    );

    if (isCanceled) {
      return;
    }

    const result = await api.actions.ensureAssetIDs();

    if (result.error) {
      notify.error("Failed to ensure asset IDs.");
      return;
    }

    notify.success(`${result.data.completed} assets have been updated.`);
  }

  async function ensureImportRefs() {
    const { isCanceled } = await confirm.open(
      "Are you sure you want to ensure all assets have an import_ref? This can take a while and cannot be undone."
    );

    if (isCanceled) {
      return;
    }

    const result = await api.actions.ensureImportRefs();

    if (result.error) {
      notify.error("Failed to ensure import refs.");
      return;
    }

    notify.success(`${result.data.completed} assets have been updated.`);
  }

  async function resetItemDateTimes() {
    const { isCanceled } = await confirm.open(
      "Are you sure you want to reset all date and time values? This can take a while and cannot be undone."
    );

    if (isCanceled) {
      return;
    }

    const result = await api.actions.resetItemDateTimes();

    if (result.error) {
      notify.error("Failed to reset date and time values.");
      return;
    }

    notify.success(`${result.data.completed} assets have been updated.`);
  }

  async function setPrimaryPhotos() {
    const { isCanceled } = await confirm.open(
      "Are you sure you want to set primary photos? This can take a while and cannot be undone."
    );

    if (isCanceled) {
      return;
    }

    const result = await api.actions.setPrimaryPhotos();

    if (result.error) {
      notify.error("Failed to set primary photos.");
      return;
    }

    notify.success(`${result.data.completed} assets have been updated.`);
  }
</script>

<style scoped></style>
