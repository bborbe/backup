<script setup lang="ts">
import { ref } from "vue";
// @ts-ignore
import backupApiClient from "../lib/BackupApiClient.ts";
import ActionButtonComponent from "./ActionButtonComponent.vue";
import type { ActionResult } from "../lib/types";

const isBackupLoading = ref(false);
const isCleanupLoading = ref(false);
const backupResult = ref<ActionResult | null>(null);
const cleanupResult = ref<ActionResult | null>(null);

async function triggerBackupAll() {
  isBackupLoading.value = true;
  backupResult.value = null;
  
  try {
    const result = await backupApiClient.triggerBackupAll();
    backupResult.value = result;
    setTimeout(() => {
      backupResult.value = null;
    }, 10000);
  } catch (err) {
    backupResult.value = {
      success: false,
      message: err instanceof Error ? err.message : "Backup failed",
    };
  } finally {
    isBackupLoading.value = false;
  }
}

async function triggerCleanupAll() {
  isCleanupLoading.value = true;
  cleanupResult.value = null;
  
  try {
    const result = await backupApiClient.triggerCleanupAll();
    cleanupResult.value = result;
    setTimeout(() => {
      cleanupResult.value = null;
    }, 10000);
  } catch (err) {
    cleanupResult.value = {
      success: false,
      message: err instanceof Error ? err.message : "Cleanup failed",
    };
  } finally {
    isCleanupLoading.value = false;
  }
}
</script>

<template>
  <div class="action-panel">
    <h2>Bulk Actions</h2>
    
    <div class="actions">
      <div class="action-group">
        <ActionButtonComponent
          label="Backup All Targets"
          :disabled="isBackupLoading"
          @click="triggerBackupAll"
        />
        <div
          v-if="backupResult"
          :class="[
            'action-result',
            backupResult.success ? 'success' : 'error'
          ]"
        >
          {{ backupResult.message }}
        </div>
      </div>
      
      <div class="action-group">
        <ActionButtonComponent
          label="Cleanup All Targets"
          variant="secondary"
          :disabled="isCleanupLoading"
          @click="triggerCleanupAll"
        />
        <div
          v-if="cleanupResult"
          :class="[
            'action-result',
            cleanupResult.success ? 'success' : 'error'
          ]"
        >
          {{ cleanupResult.message }}
        </div>
      </div>
    </div>
    
    <div class="action-info">
      <p><strong>Backup All:</strong> Triggers backup for all configured targets</p>
      <p><strong>Cleanup All:</strong> Removes old backups according to retention policy</p>
    </div>
  </div>
</template>

<style scoped>
.action-panel {
  padding: 1rem;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  margin-bottom: 2rem;
}

.action-panel h2 {
  margin: 0 0 1rem 0;
  color: #1f2937;
}

.actions {
  display: flex;
  gap: 2rem;
  margin-bottom: 1rem;
}

.action-group {
  flex: 1;
}

.action-result {
  margin-top: 0.5rem;
  padding: 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.875rem;
}

.action-result.success {
  background-color: #d1fae5;
  color: #065f46;
  border: 1px solid #a7f3d0;
}

.action-result.error {
  background-color: #fee2e2;
  color: #991b1b;
  border: 1px solid #fca5a5;
}

.action-info {
  font-size: 0.875rem;
  color: #6b7280;
}

.action-info p {
  margin: 0.25rem 0;
}

@media (max-width: 768px) {
  .actions {
    flex-direction: column;
    gap: 1rem;
  }
}
</style>