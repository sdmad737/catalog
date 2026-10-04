<template>
  <div class="grid grid-cols-1 gap-4 py-5 md:grid-cols-[minmax(0,1fr)_auto] md:items-center">
    <div>
      <h4 class="mb-1 text-base font-bold">
        <slot name="title"></slot>
      </h4>
      <p class="catalog-muted text-sm leading-6">
        <slot></slot>
      </p>
    </div>
    <BaseButton class="btn-primary md:min-w-[150px]" @click="$emit('action')">
      <slot name="button">
        <slot name="title"></slot>
      </slot>
    </BaseButton>
  </div>
</template>
