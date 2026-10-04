<template>
  <div v-if="!inline" class="form-control w-full">
    <label class="label" :for="inputId">
      <span class="label-text">{{ label }} <span v-if="required" class="text-error" aria-hidden="true">*</span></span>
    </label>
    <input
      :id="inputId"
      ref="input"
      v-model="value"
      :placeholder="placeholder"
      :type="type"
      :required="required"
      :autocomplete="autocomplete"
      :autofocus="autofocus"
      class="input input-bordered w-full"
    />
    <p v-if="hint" class="catalog-muted mt-1 px-1 text-xs">{{ hint }}</p>
  </div>
  <div v-else class="sm:grid sm:grid-cols-4 sm:items-start sm:gap-4">
    <label class="label" :for="inputId">
      <span class="label-text">{{ label }}</span>
    </label>
    <input
      :id="inputId"
      v-model="value"
      :placeholder="placeholder"
      class="input input-bordered col-span-3 mt-2 w-full"
    />
  </div>
</template>

<script lang="ts" setup>
  const props = defineProps({
    label: {
      type: String,
      default: "",
    },
    modelValue: {
      type: [String, Number],
      default: null,
    },
    type: {
      type: String,
      default: "text",
    },
    triggerFocus: {
      type: Boolean,
      default: null,
    },
    inline: {
      type: Boolean,
      default: false,
    },
    placeholder: {
      type: String,
      default: "",
    },
    required: {
      type: Boolean,
      default: false,
    },
    hint: {
      type: String,
      default: "",
    },
    autocomplete: {
      type: String,
      default: "off",
    },
    autofocus: {
      type: Boolean,
      default: false,
    },
  });

  const input = ref<HTMLElement | null>(null);
  const inputId = useId();

  whenever(
    () => props.triggerFocus,
    () => {
      if (input.value) {
        input.value.focus();
      }
    }
  );

  const value = useVModel(props, "modelValue");
</script>
