<script setup lang="ts">
  import type { TreeItem } from "~~/lib/api/types/data-contracts";

  type Props = {
    locs: TreeItem[];
    treeId: string;
  };

  defineProps<Props>();
</script>

<template>
  <div class="root space-y-1 rounded-xl bg-base-200/45 p-2 sm:p-3">
    <LocationTreeNode v-for="item in locs" :key="item.id" :item="item" :tree-id="treeId" />
  </div>
</template>
