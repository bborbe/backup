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
  size: {
    type: String,
    required: false,
    default: "normal",
    validator: (value: string) => ["small", "normal"].includes(value),
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
      `btn-${props.size}`,
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
  padding: var(--spacing-sm) var(--spacing-md);
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
  font-size: var(--font-size-base);
  cursor: pointer;
  transition: all var(--transition-fast);
  line-height: var(--line-height-normal);
  position: relative;
  overflow: hidden;
  transform: translateZ(0);
}

.btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transition: left var(--transition-normal);
}

.btn:hover::before {
  left: 100%;
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.btn-primary {
  background-color: var(--accent-primary);
  color: var(--text-primary);
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--accent-primary-hover);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.btn-primary:active:not(:disabled) {
  background-color: var(--accent-primary-active);
  transform: translateY(0);
  box-shadow: var(--shadow-sm);
}

.btn-secondary {
  background-color: var(--bg-tertiary);
  color: var(--text-primary);
  border: 1px solid var(--border-primary);
}

.btn-secondary:hover:not(:disabled) {
  background-color: var(--bg-hover);
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
}

.btn-secondary:active:not(:disabled) {
  transform: translateY(0);
}

.btn-danger {
  background-color: var(--status-error);
  color: var(--text-primary);
}

.btn-danger:hover:not(:disabled) {
  background-color: var(--status-error);
  filter: brightness(1.1);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.btn-danger:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: var(--shadow-sm);
}

.btn-small {
  padding: var(--spacing-xs) var(--spacing-sm);
  font-size: var(--font-size-sm);
}

.btn-normal {
  padding: var(--spacing-sm) var(--spacing-md);
}
</style>