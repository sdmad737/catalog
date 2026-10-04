<script setup lang="ts">
  definePageMeta({ middleware: ["auth"] });
  useHead({ title: "Add inventory item — Catalog" });

  const open = ref(true);
</script>

<template>
  <BaseContainer class="catalog-page max-w-none">
    <ItemCreateModal v-model="open" />
    <AppPageHeader
      eyebrow="Inventory"
      title="Add an inventory item"
      description="Capture the essentials now, then enrich the record with photos, documents, purchase details, and maintenance history."
    />
    <div class="catalog-panel mx-auto max-w-2xl rounded-2xl p-8 text-center">
      <AppLogo class="mx-auto h-12 w-12 text-primary" />
      <h2 class="mt-4 text-lg font-bold">Item creator</h2>
      <p class="catalog-muted mx-auto mt-2 max-w-md text-sm leading-6">
        The item form opens automatically. If you closed it by mistake, you can reopen it below.
      </p>
      <BaseButton class="btn-primary mt-5" @click="open = true">Open item form</BaseButton>
    </div>
  </BaseContainer>
</template>
