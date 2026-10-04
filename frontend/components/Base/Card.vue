<template>
  <div class="card overflow-hidden rounded-2xl bg-base-100">
    <div v-if="$slots.title" class="border-b border-base-content/5 px-4 py-4 sm:px-6">
      <component :is="collapsable ? 'button' : 'div'" v-on="collapsable ? { click: toggle } : {}">
        <h3 class="text-base font-bold leading-6 flex items-center">
          <slot name="title"></slot>
          <template v-if="collapsable">
            <span class="ml-2 swap swap-rotate" :class="`${collapsed ? 'swap-active' : ''}`">
              <MdiChevronRight class="h-6 w-6 swap-on" />
              <MdiChevronDown class="h-6 w-6 swap-off" />
            </span>
          </template>
        </h3>
      </component>
      <div>
        <p v-if="$slots.subtitle" class="catalog-muted mt-1 max-w-2xl text-sm">
          <slot name="subtitle"></slot>
        </p>
        <template v-if="$slots['title-actions']">
          <slot name="title-actions"></slot>
        </template>
      </div>
    </div>
    <div
      :class="{
        'max-h-[9000px]': collapsable && !collapsed,
        'max-h-0 overflow-hidden': collapsed,
      }"
      class="transition-[max-height] duration-200"
    >
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
  import MdiChevronDown from "~icons/mdi/chevron-down";
  import MdiChevronRight from "~icons/mdi/chevron-right";

  defineProps<{
    collapsable?: boolean;
  }>();

  function toggle() {
    collapsed.value = !collapsed.value;
  }

  const collapsed = ref(false);
</script>
