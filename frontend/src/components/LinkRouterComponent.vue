<script setup lang="ts">
const props = defineProps({
  to: {
    type: Object,
    required: true,
  },
  title: {
    type: String,
    required: true,
  },
});
</script>

<template>
  <router-link
    class="link"
    :to="props.to"
    :title="props.title"
  >
    <slot />
  </router-link>
</template>

<style scoped>
.link {
  color: #2563eb;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  transition: background-color 0.2s;
}

.link:hover {
  background-color: #f3f4f6;
  color: #1d4ed8;
}

.link.router-link-active {
  background-color: #2563eb;
  color: white;
}
</style>