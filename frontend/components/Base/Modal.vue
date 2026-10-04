<template>
  <div class="z-[999]">
    <input :id="modalId" v-model="modal" type="checkbox" class="modal-toggle" />
    <div
      class="modal modal-bottom overflow-visible sm:modal-middle"
      role="dialog"
      aria-modal="true"
      :aria-labelledby="`${modalId}-title`"
    >
      <button class="absolute inset-0 cursor-default" aria-label="Close dialog" @click="close" />
      <div class="modal-box relative overflow-visible">
        <button class="btn btn-ghost btn-sm btn-circle absolute right-3 top-3" aria-label="Close dialog" @click="close">
          ✕
        </button>

        <h3 :id="`${modalId}-title`" class="pr-10 text-lg font-bold">
          <slot name="title"></slot>
        </h3>
        <slot> </slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  const emit = defineEmits(["cancel", "update:modelValue"]);
  const props = defineProps({
    modelValue: {
      type: Boolean,
      required: true,
    },
    /**
     * in readonly mode the modal only `emits` a "cancel" event to indicate
     * that the modal was closed via the "x" button. The parent component is
     * responsible for closing the modal.
     */
    readonly: {
      type: Boolean,
      default: false,
    },
  });

  function escClose(e: KeyboardEvent) {
    if (e.key === "Escape") {
      close();
    }
  }

  function close() {
    if (props.readonly) {
      emit("cancel");
      return;
    }
    modal.value = false;
  }

  const modalId = useId();
  const modal = useVModel(props, "modelValue", emit);

  watchEffect(() => {
    if (modal.value) {
      document.addEventListener("keydown", escClose);
    } else {
      document.removeEventListener("keydown", escClose);
    }
  });
</script>
