<template>
  <div class="border-t border-base-content/5 px-4 py-2 sm:px-6">
    <dl class="divide-y divide-base-content/5">
      <div v-for="(detail, i) in details" :key="i" class="group py-4 sm:grid sm:grid-cols-3 sm:gap-5">
        <dt class="text-xs font-bold uppercase tracking-wide text-base-content/45">
          {{ detail.name }}
        </dt>
        <dd class="mt-1 text-start text-sm text-base-content sm:col-span-2 sm:mt-0">
          <slot :name="detail.slot || detail.name" v-bind="{ detail }">
            <DateTime
              v-if="detail.type == 'date'"
              :date="detail.text"
              :datetime-type="detail.date ? 'date' : 'datetime'"
            />
            <Currency v-else-if="detail.type == 'currency'" :amount="detail.text" />
            <template v-else-if="detail.type === 'link'">
              <div class="tooltip tooltip-primary tooltip-right" :data-tip="detail.href">
                <a class="btn btn-primary btn-xs" :href="detail.href" target="_blank">
                  <MdiOpenInNew class="mr-2 swap-on" />
                  {{ detail.text }}
                </a>
              </div>
            </template>
            <template v-else-if="detail.type === 'markdown'">
              <ClientOnly>
                <Markdown :source="detail.text" />
              </ClientOnly>
            </template>
            <template v-else>
              <span class="flex items-center">
                {{ detail.text }}
                <span
                  v-if="detail.copyable"
                  class="ml-3 opacity-50 transition-opacity hover:opacity-100 focus-within:opacity-100"
                >
                  <CopyText
                    v-if="detail.text.toString()"
                    :text="detail.text.toString()"
                    :icon-size="16"
                    class="btn btn-xs btn-ghost btn-circle"
                  />
                </span>
              </span>
            </template>
          </slot>
        </dd>
      </div>
    </dl>
  </div>
</template>

<script setup lang="ts">
  import type { AnyDetail, Detail } from "./types";
  import MdiOpenInNew from "~icons/mdi/open-in-new";

  defineProps({
    details: {
      type: Object as () => (Detail | AnyDetail)[],
      required: true,
    },
  });
</script>
