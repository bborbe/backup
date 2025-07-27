<script setup lang="ts">
import { ref, onMounted } from "vue";
// @ts-ignore
import backupApiClient from "../lib/BackupApiClient.ts";
import type { BackupStatus, LoadingState } from "../lib/types";
import ActionButtonComponent from "./ActionButtonComponent.vue";

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
</script>

<template>
  <div class="status-overview">
    <div class="header">
      <h2>Backup Status Overview</h2>
      
      <div class="action-buttons">
        <ActionButtonComponent
          label="Backup All"
          variant="primary"
          :disabled="actionState.isBackingUp || actionState.isCleaningUp"
          @click="triggerBackupAll"
        />
        <ActionButtonComponent
          label="Cleanup All"
          variant="danger"
          :disabled="actionState.isBackingUp || actionState.isCleaningUp"
          @click="triggerCleanupAll"
        />
      </div>
    </div>
    
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
    
    <div v-else class="status-grid">
      <div
        v-for="(lastBackup, host) in status"
        :key="host"
        :class="['status-card', getStatusClass(lastBackup)]"
      >
        <div class="card-content">
          <div class="host-name">{{ host }}</div>
          <div class="last-backup">{{ formatDate(lastBackup) }}</div>
          <div class="backup-date">{{ lastBackup || 'No date available' }}</div>
        </div>
        
        <div class="card-actions">
          <ActionButtonComponent
            label="Backup"
            variant="primary"
            :disabled="
              (individualActionState[host]?.isBackingUp || individualActionState[host]?.isCleaningUp) ||
              actionState.isBackingUp || actionState.isCleaningUp
            "
            @click="triggerBackup(host)"
          />
          <ActionButtonComponent
            label="Cleanup"
            variant="danger"
            :disabled="
              (individualActionState[host]?.isBackingUp || individualActionState[host]?.isCleaningUp) ||
              actionState.isBackingUp || actionState.isCleaningUp
            "
            @click="triggerCleanup(host)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.status-overview {
  padding: 1rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.status-overview h2 {
  margin: 0;
  color: #1f2937;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
}

.action-message {
  padding: 0.75rem 1rem;
  margin-bottom: 1rem;
  background-color: #dbeafe;
  border: 1px solid #60a5fa;
  border-radius: 0.5rem;
  color: #1e40af;
}

.loading, .error, .empty {
  padding: 2rem;
  text-align: center;
  background-color: #f9fafb;
  border-radius: 0.5rem;
  border: 1px solid #e5e7eb;
}

.error {
  color: #dc2626;
}

.retry-btn {
  margin-left: 1rem;
  padding: 0.25rem 0.5rem;
  background-color: #2563eb;
  color: white;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1rem;
}

.status-card {
  padding: 1rem;
  border-radius: 0.5rem;
  border: 2px solid;
  background-color: white;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.card-content {
  flex: 1;
}

.card-actions {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
  padding-top: 0.5rem;
  border-top: 1px solid #e5e7eb;
}

.status-success {
  border-color: #10b981;
  background-color: #f0fdf4;
}

.status-warning {
  border-color: #f59e0b;
  background-color: #fffbeb;
}

.status-error {
  border-color: #ef4444;
  background-color: #fef2f2;
}

.host-name {
  font-weight: 600;
  font-size: 1.125rem;
  margin-bottom: 0.5rem;
  color: #1f2937;
}

.last-backup {
  font-weight: 500;
  margin-bottom: 0.25rem;
}

.backup-date {
  font-size: 0.875rem;
  color: #6b7280;
}
</style>