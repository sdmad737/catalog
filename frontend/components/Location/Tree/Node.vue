<script setup lang="ts">
  import { useTreeState } from "./tree-state";
  import type { TreeItem } from "~~/lib/api/types/data-contracts";
  import MdiChevronDown from "~icons/mdi/chevron-down";
  import MdiChevronRight from "~icons/mdi/chevron-right";
  import MdiMapMarker from "~icons/mdi/map-marker";
  import MdiPackageVariant from "~icons/mdi/package-variant";

  type Props = {
    treeId: string;
    item: TreeItem;
  };
  const props = withDefaults(defineProps<Props>(), {});

  const link = computed(() => {
    return props.item.type === "location" ? `/location/${props.item.id}` : `/item/${props.item.id}`;
  });

  const hasChildren = computed(() => {
    return props.item.children.length > 0;
  });

  const state = useTreeState(props.treeId);

  const openRef = computed({
    get() {
      return state.value[nodeHash.value] ?? false;
    },
    set(value: boolean) {
      state.value[nodeHash.value] = value;
    },
  });

  const nodeHash = computed(() => {
    // converts a UUID to a short hash
    return props.item.id.replace(/-/g, "").substring(0, 8);
  });
</script>

<template>
  <div>
    <div
      class="node group flex min-h-[44px] items-center gap-2 rounded-lg px-2 py-1 text-sm transition"
      :class="{
        'cursor-pointer hover:bg-base-100': hasChildren,
      }"
      @click="openRef = !openRef"
    >
      <div class="flex h-7 w-7 items-center justify-center rounded-md">
        <div v-if="!hasChildren" class="h-5 w-5" />
        <label
          v-else
          class="swap swap-rotate"
          :class="{
            'swap-active': openRef,
          }"
        >
          <MdiChevronRight name="mdi-chevron-right" class="h-5 w-5 swap-off" />
          <MdiChevronDown name="mdi-chevron-down" class="h-5 w-5 swap-on" />
        </label>
      </div>
      <span class="grid h-8 w-8 place-items-center rounded-lg bg-primary/10 text-primary">
        <MdiMapMarker v-if="item.type === 'location'" class="h-4 w-4" />
        <MdiPackageVariant v-else class="h-4 w-4" />
      </span>
      <NuxtLink class="min-w-0 flex-1 truncate font-semibold hover:text-primary" :to="link" @click.stop>
        {{ item.name }}
      </NuxtLink>
      <span v-if="item.children.length" class="badge badge-ghost badge-sm">{{ item.children.length }}</span>
    </div>
    <div v-if="openRef" class="ml-5 border-l border-base-content/10 pl-3">
      <LocationTreeNode v-for="child in item.children" :key="child.id" :item="child" :tree-id="treeId" />
    </div>
  </div>
</template>
