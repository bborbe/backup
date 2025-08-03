<script setup lang="ts">
import { RouterView } from "vue-router";
</script>

<template>
  <RouterView />
</template>

<style>
:root {
  /* Core Backgrounds */
  --bg-primary: #1a1a1a;
  --bg-secondary: #252525;
  --bg-tertiary: #2d2d2d;
  --bg-quaternary: #363636;
  
  /* Interactive Backgrounds */
  --bg-hover: #3a3a3a;
  --bg-active: #404040;
  --bg-selected: #2d4a6b;
  
  /* Text Colors */
  --text-primary: #ffffff;
  --text-secondary: #b4b4b4;
  --text-muted: #7a7a7a;
  --text-inverse: #1a1a1a;
  
  /* Border Colors */
  --border-primary: #3a3a3a;
  --border-secondary: #2d2d2d;
  --border-accent: #4a4a4a;
  
  /* Status Colors */
  --status-success: #00d4aa;
  --status-success-bg: #0d2818;
  --status-success-border: #1a4731;
  
  --status-warning: #ff9500;
  --status-warning-bg: #2d1a00;
  --status-warning-border: #4a2c00;
  
  --status-error: #ff4757;
  --status-error-bg: #2d0a0d;
  --status-error-border: #4a141a;
  
  /* Financial Colors (positive/negative values) */
  --financial-positive: #00d4aa;
  --financial-negative: #ff4757;
  --financial-neutral: #6c757d;
  
  /* Interactive Colors */
  --accent-primary: #4a9eff;
  --accent-primary-hover: #3d8bef;
  --accent-primary-active: #2d7adf;
  
  --accent-secondary: #6c757d;
  --accent-secondary-hover: #5a6268;
  
  /* Spacing System */
  --spacing-xs: 0.25rem;
  --spacing-sm: 0.5rem;
  --spacing-md: 1rem;
  --spacing-lg: 1.5rem;
  --spacing-xl: 2rem;
  --spacing-2xl: 3rem;
  --spacing-3xl: 4rem;
  
  /* Typography */
  --font-size-xs: 0.75rem;
  --font-size-sm: 0.875rem;
  --font-size-base: 1rem;
  --font-size-lg: 1.125rem;
  --font-size-xl: 1.25rem;
  --font-size-2xl: 1.5rem;
  --font-size-3xl: 1.875rem;
  --font-size-4xl: 2.25rem;
  
  --font-weight-normal: 400;
  --font-weight-medium: 500;
  --font-weight-semibold: 600;
  --font-weight-bold: 700;
  
  --line-height-tight: 1.2;
  --line-height-normal: 1.5;
  --line-height-relaxed: 1.6;
  
  /* Border Radius */
  --radius-sm: 0.25rem;
  --radius-md: 0.375rem;
  --radius-lg: 0.5rem;
  --radius-xl: 1rem;
  
  /* Transitions */
  --transition-fast: 150ms ease-in-out;
  --transition-normal: 250ms ease-in-out;
  --transition-slow: 350ms ease-in-out;
  
  /* Shadows */
  --shadow-sm: 0 2px 4px rgba(0, 0, 0, 0.1);
  --shadow-md: 0 4px 8px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 8px 16px rgba(0, 0, 0, 0.1);
  --shadow-hover: 0 4px 12px rgba(0, 0, 0, 0.15);
  
  /* Component Spacing */
  --card-padding: var(--spacing-md);
  --card-gap: var(--spacing-sm);
  --section-gap: var(--spacing-lg);
  --element-gap: var(--spacing-sm);
  
  /* Grid System Variables */
  --grid-base: 8px;
  --grid-1: calc(var(--grid-base) * 1);  /* 8px */
  --grid-2: calc(var(--grid-base) * 2);  /* 16px */
  --grid-3: calc(var(--grid-base) * 3);  /* 24px */
  --grid-4: calc(var(--grid-base) * 4);  /* 32px */
  --grid-5: calc(var(--grid-base) * 5);  /* 40px */
  --grid-6: calc(var(--grid-base) * 6);  /* 48px */
  --grid-8: calc(var(--grid-base) * 8);  /* 64px */
  --grid-10: calc(var(--grid-base) * 10); /* 80px */
  
  /* Table Variables */
  --table-row-height: 2.5rem;
  --table-header-height: 3rem;
  --table-cell-padding: var(--spacing-sm) var(--spacing-md);
  --table-border: 1px solid var(--border-secondary);
  
  /* Metric Card Variables */
  --metric-card-padding: var(--spacing-lg);
  --metric-number-size: var(--font-size-3xl);
  --metric-label-size: var(--font-size-sm);
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  min-height: 100vh;
}

#app {
  min-height: 100vh;
  background-color: var(--bg-primary);
}
</style>
