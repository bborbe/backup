<script setup lang="ts">
import { BooleanProp, StringProp } from "../lib/ComponentTypes";

const props = defineProps({
  label: StringProp,
  disabled: {
    ...BooleanProp,
    required: false,
    default: false,
  },
  variant: {
    type: String,
    required: false,
    default: "primary",
    validator: (value: string) => ["primary", "secondary", "danger"].includes(value),
  },
});

const emit = defineEmits<{
  click: [];
}>();

function handleClick() {
  if (!props.disabled) {
    emit("click");
  }
}
</script>

<template>
  <button
    :class="[
      'btn',
      `btn-${props.variant}`,
      { 'btn-disabled': props.disabled }
    ]"
    :disabled="props.disabled"
    @click="handleClick"
  >
    {{ props.label }}
  </button>
</template>

<style scoped>
.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.btn-primary {
  background-color: #2563eb;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #1d4ed8;
}

.btn-secondary {
  background-color: #6b7280;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #4b5563;
}

.btn-danger {
  background-color: #dc2626;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background-color: #b91c1c;
}
</style>