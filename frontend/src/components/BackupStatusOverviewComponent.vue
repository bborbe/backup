<script setup lang="ts">
import { ref, onMounted, defineExpose, computed } from "vue";
// @ts-ignore
import backupApiClient from "../lib/BackupApiClient.ts";
import type { BackupStatus, LoadingState } from "../lib/types";
import ActionButtonComponent from "./ActionButtonComponent.vue";

interface FilterState {
  total: boolean;
  healthy: boolean;
  warning: boolean;
  critical: boolean;
}

interface Props {
  filters?: FilterState;
}

const props = withDefaults(defineProps<Props>(), {
  filters: () => ({
    total: true,
    healthy: true,
    warning: true,
    critical: true
  })
});

const status = ref<BackupStatus>({});
const loadingState = ref<LoadingState>({
  isLoading: false,
  error: null,
});

const actionState = ref({
  isBackingUp: false,
  isCleaningUp: false,
  message: null as string | null,
});

const individualActionState = ref<Record<string, { isBackingUp: boolean; isCleaningUp: boolean }>>({});

onMounted(async () => {
  await loadStatus();
  // Auto-refresh every 30 seconds
  setInterval(loadStatus, 30000);
});

async function loadStatus() {
  loadingState.value.isLoading = true;
  loadingState.value.error = null;
  
  try {
    const result = await backupApiClient.getStatus();
    console.log("Backup status API response:", result);
    console.log("Result type:", typeof result, "Is array:", Array.isArray(result));
    
    // Ensure we have a proper object, not an array
    if (result && typeof result === 'object' && !Array.isArray(result)) {
      status.value = result;
      // Initialize individual action states for all hosts
      for (const host of Object.keys(result)) {
        if (!individualActionState.value[host]) {
          individualActionState.value[host] = { isBackingUp: false, isCleaningUp: false };
        }
      }
    } else {
      console.warn("API returned unexpected data format:", result);
      status.value = {};
    }
  } catch (err) {
    console.error("Failed to load backup status:", err);
    loadingState.value.error = err instanceof Error ? err.message : "Failed to load status";
  } finally {
    loadingState.value.isLoading = false;
  }
}

function getStatusClass(date: string): string {
  if (!date || date === "") return "status-error";
  
  const backupDate = new Date(date);
  if (isNaN(backupDate.getTime())) return "status-error";
  
  const now = new Date();
  const daysDiff = Math.floor((now.getTime() - backupDate.getTime()) / (1000 * 60 * 60 * 24));
  
  if (daysDiff === 0) return "status-success";
  if (daysDiff === 1) return "status-warning";
  return "status-error";
}

function getStatusText(date: string): string {
  if (!date || date === "") return "Failed";
  
  const backupDate = new Date(date);
  if (isNaN(backupDate.getTime())) return "Failed";
  
  const now = new Date();
  const daysDiff = Math.floor((now.getTime() - backupDate.getTime()) / (1000 * 60 * 60 * 24));
  
  if (daysDiff === 0) return "Healthy";
  if (daysDiff === 1) return "Warning";
  return "Critical";
}

function getStatusCategory(date: string): 'healthy' | 'warning' | 'critical' {
  if (!date || date === "") return "critical";
  
  const backupDate = new Date(date);
  if (isNaN(backupDate.getTime())) return "critical";
  
  const now = new Date();
  const daysDiff = Math.floor((now.getTime() - backupDate.getTime()) / (1000 * 60 * 60 * 24));
  
  if (daysDiff === 0) return "healthy";
  if (daysDiff === 1) return "warning";
  return "critical";
}

const filteredStatus = computed(() => {
  const entries = Object.entries(status.value);
  const filtered = entries.filter(([, lastBackup]) => {
    const category = getStatusCategory(lastBackup);
    
    // Check if this category should be shown based on filters
    switch (category) {
      case 'healthy':
        return props.filters.healthy;
      case 'warning':
        return props.filters.warning;
      case 'critical':
        return props.filters.critical;
      default:
        return true; // fallback
    }
  });
  
  return Object.fromEntries(filtered);
});

function formatDate(dateString: string): string {
  if (!dateString || dateString === "") return "No backup yet";
  
  const date = new Date(dateString);
  if (isNaN(date.getTime())) return "Invalid date";
  
  const now = new Date();
  const daysDiff = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60 * 24));
  
  if (daysDiff === 0) return "Today";
  if (daysDiff === 1) return "Yesterday";
  if (daysDiff < 0) return "Future date (error)";
  return `${daysDiff} days ago`;
}

async function triggerBackupAll() {
  actionState.value.isBackingUp = true;
  actionState.value.message = null;
  
  try {
    const result = await backupApiClient.triggerBackupAll();
    actionState.value.message = result.message;
    setTimeout(() => {
      actionState.value.message = null;
    }, 5000);
  } catch (err) {
    actionState.value.message = err instanceof Error ? err.message : "Failed to trigger backup";
    setTimeout(() => {
      actionState.value.message = null;
    }, 5000);
  } finally {
    actionState.value.isBackingUp = false;
  }
}

async function triggerCleanupAll() {
  actionState.value.isCleaningUp = true;
  actionState.value.message = null;
  
  try {
    const result = await backupApiClient.triggerCleanupAll();
    actionState.value.message = result.message;
    setTimeout(() => {
      actionState.value.message = null;
    }, 5000);
  } catch (err) {
    actionState.value.message = err instanceof Error ? err.message : "Failed to trigger cleanup";
    setTimeout(() => {
      actionState.value.message = null;
    }, 5000);
  } finally {
    actionState.value.isCleaningUp = false;
  }
}

async function triggerBackup(host: string) {
  if (!individualActionState.value[host]) {
    individualActionState.value[host] = { isBackingUp: false, isCleaningUp: false };
  }
  
  individualActionState.value[host].isBackingUp = true;
  actionState.value.message = null;
  
  try {
    const result = await backupApiClient.triggerBackup(host);
    actionState.value.message = result.message;
    setTimeout(() => {
      actionState.value.message = null;
    }, 5000);
  } catch (err) {
    actionState.value.message = err instanceof Error ? err.message : `Failed to trigger backup for ${host}`;
    setTimeout(() => {
      actionState.value.message = null;
    }, 5000);
  } finally {
    individualActionState.value[host].isBackingUp = false;
  }
}

async function triggerCleanup(host: string) {
  if (!individualActionState.value[host]) {
    individualActionState.value[host] = { isBackingUp: false, isCleaningUp: false };
  }
  
  individualActionState.value[host].isCleaningUp = true;
  actionState.value.message = null;
  
  try {
    const result = await backupApiClient.triggerCleanup(host);
    actionState.value.message = result.message;
    setTimeout(() => {
      actionState.value.message = null;
    }, 5000);
  } catch (err) {
    actionState.value.message = err instanceof Error ? err.message : `Failed to trigger cleanup for ${host}`;
    setTimeout(() => {
      actionState.value.message = null;
    }, 5000);
  } finally {
    individualActionState.value[host].isCleaningUp = false;
  }
}

// Expose functions and state for parent component
defineExpose({
  triggerBackupAll,
  triggerCleanupAll,
  actionState
});
</script>

<template>
  <div class="status-overview">
    
    <div v-if="actionState.message" class="action-message">
      {{ actionState.message }}
    </div>
    
    <div v-if="loadingState.isLoading" class="loading">
      Loading backup status...
    </div>
    
    <div v-else-if="loadingState.error" class="error">
      Error: {{ loadingState.error }}
      <button @click="loadStatus" class="retry-btn">Retry</button>
    </div>
    
    <div v-else-if="Object.keys(status).length === 0" class="empty">
      No backup targets found or no backups have been created yet
    </div>
    
    <div v-else-if="Object.keys(filteredStatus).length === 0" class="empty">
      No backup targets match the current filter
    </div>
    
    <div v-else class="status-grid">
      <div
        v-for="(lastBackup, host) in filteredStatus"
        :key="host"
        :class="['status-card', getStatusClass(lastBackup)]"
      >
        <div class="card-header">
          <div class="host-name">{{ host }}</div>
          <div class="status-indicator" :class="getStatusClass(lastBackup)"></div>
        </div>
        
        <div class="card-body">
          <div class="status-details">
            <div class="detail-row">
              <span class="detail-label">Last Backup</span>
              <span class="detail-value" :class="getStatusClass(lastBackup)">{{ formatDate(lastBackup) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Status</span>
              <span class="status-pill" :class="getStatusClass(lastBackup)">
                {{ getStatusText(lastBackup) }}
              </span>
            </div>
          </div>
          
          <div class="card-actions">
            <ActionButtonComponent
              label="Backup"
              variant="primary"
              size="small"
              :disabled="
                (individualActionState[host]?.isBackingUp || individualActionState[host]?.isCleaningUp) ||
                actionState.isBackingUp || actionState.isCleaningUp
              "
              @click="triggerBackup(String(host))"
            />
            <ActionButtonComponent
              label="Cleanup"
              variant="danger"
              size="small"
              :disabled="
                (individualActionState[host]?.isBackingUp || individualActionState[host]?.isCleaningUp) ||
                actionState.isBackingUp || actionState.isCleaningUp
              "
              @click="triggerCleanup(String(host))"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.status-overview {
  padding: var(--spacing-md) var(--spacing-xl);
  background-color: var(--bg-primary);
}


.action-message {
  padding: var(--spacing-sm) var(--spacing-md);
  margin-bottom: var(--spacing-lg);
  background-color: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-normal);
}

.loading, .error, .empty {
  padding: var(--spacing-xl);
  text-align: center;
  background-color: var(--bg-tertiary);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-secondary);
  color: var(--text-secondary);
  font-size: var(--font-size-base);
  line-height: var(--line-height-normal);
}

.error {
  color: var(--status-error);
}

.retry-btn {
  margin-left: var(--spacing-md);
  padding: var(--spacing-xs) var(--spacing-sm);
  background-color: var(--accent-primary);
  color: var(--text-primary);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: var(--transition-fast);
}

.retry-btn:hover {
  background-color: var(--accent-primary-hover);
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 0.75rem;
}

@media (min-width: 1024px) {
  .status-grid {
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  }
}

@media (min-width: 1280px) {
  .status-grid {
    grid-template-columns: repeat(auto-fill, minmax(170px, 1fr));
  }
}

@media (min-width: 1536px) {
  .status-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  }
}

.status-card {
  border-radius: var(--radius-md);
  border: 1px solid var(--border-primary);
  background-color: var(--bg-tertiary);
  overflow: hidden;
  transition: all var(--transition-normal);
  cursor: pointer;
  position: relative;
}

.status-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.02), transparent);
  opacity: 0;
  transition: opacity var(--transition-normal);
  pointer-events: none;
}

.status-card:hover {
  border-color: var(--border-accent);
  box-shadow: var(--shadow-hover);
  transform: translateY(-4px);
}

.status-card:hover::before {
  opacity: 1;
}

.card-header {
  padding: 0.75rem;
  background-color: var(--bg-quaternary);
  border-bottom: 1px solid var(--border-secondary);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-body {
  padding: 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.status-details {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 1.5rem;
}

.detail-label {
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-value {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.detail-value.status-success {
  color: var(--financial-positive);
}

.detail-value.status-warning {
  color: var(--status-warning);
}

.detail-value.status-error {
  color: var(--financial-negative);
}

.status-pill {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.status-pill.status-success {
  background-color: var(--status-success-bg);
  color: var(--status-success);
  border: 1px solid var(--status-success-border);
}

.status-pill.status-warning {
  background-color: var(--status-warning-bg);
  color: var(--status-warning);
  border: 1px solid var(--status-warning-border);
}

.status-pill.status-error {
  background-color: var(--status-error-bg);
  color: var(--status-error);
  border: 1px solid var(--status-error-border);
}

.status-indicator {
  width: 0.625rem;
  height: 0.625rem;
  border-radius: 50%;
  flex-shrink: 0;
  position: relative;
  transition: all var(--transition-normal);
}

.status-indicator::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  transform: translate(-50%, -50%);
  opacity: 0;
  transition: all var(--transition-normal);
}

.status-card:hover .status-indicator::after {
  width: 150%;
  height: 150%;
  opacity: 0.3;
}

.status-success .status-indicator::after {
  background-color: var(--status-success);
}

.status-warning .status-indicator::after {
  background-color: var(--status-warning);
}

.status-error .status-indicator::after {
  background-color: var(--status-error);
}

.card-actions {
  display: flex;
  gap: 0.375rem;
  padding-top: 0.5rem;
  border-top: 1px solid var(--border-secondary);
  margin-top: 0.5rem;
}

.status-success .status-indicator {
  background-color: var(--status-success);
}

.status-warning .status-indicator {
  background-color: var(--status-warning);
}

.status-error .status-indicator {
  background-color: var(--status-error);
}

.host-name {
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  word-wrap: break-word;
  overflow-wrap: break-word;
  hyphens: auto;
  line-height: var(--line-height-tight);
  flex: 1;
}

</style>