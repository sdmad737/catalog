<template>
  <div class="force-above fixed right-3 top-3 w-[min(360px,calc(100vw-1.5rem))]" aria-live="polite" aria-atomic="true">
    <TransitionGroup name="notify" tag="div">
      <div
        v-for="(notify, index) in notifications.slice(0, 4)"
        :key="notify.id"
        class="my-2 w-full cursor-pointer rounded-xl border border-white/15 p-4 text-sm text-white shadow-2xl"
        :class="{
          'bg-primary': notify.type === 'info',
          'bg-red-600': notify.type === 'error',
          'bg-green-600': notify.type === 'success',
        }"
        @click="dropNotification(index)"
      >
        <div class="flex items-start gap-2 font-semibold">
          <template v-if="notify.type == 'success'">
            <MdiCheckboxMarkedCircle class="h-5 w-5" />
          </template>
          <template v-if="notify.type == 'info'">
            <MdiInformationSlabCircle class="h-5 w-5" />
          </template>

          <template v-if="notify.type == 'error'">
            <MdiAlert class="h-5 w-5" />
          </template>
          {{ notify.message }}
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
  import MdiCheckboxMarkedCircle from "~icons/mdi/checkbox-marked-circle";
  import MdiInformationSlabCircle from "~icons/mdi/information-slab-circle";
  import MdiAlert from "~icons/mdi/alert";

  import { useNotifications } from "@/composables/use-notifier";

  const { notifications, dropNotification } = useNotifications();
</script>

<style scoped>
  .force-above {
    z-index: 9999;
  }

  .notify-move,
  .notify-enter-active,
  .notify-leave-active {
    transition: all 0.5s ease;
  }
  .notify-enter-from,
  .notify-leave-to {
    opacity: 0;
    transform: translateY(-30px);
  }
  .notify-leave-active {
    position: absolute;
    transform: translateY(30px);
  }
</style>
