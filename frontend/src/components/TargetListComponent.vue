<script setup lang="ts">
import { ref, onMounted } from "vue";
// @ts-ignore
import backupApiClient from "../lib/BackupApiClient.ts";
import ActionButtonComponent from "./ActionButtonComponent.vue";
import type { BackupSpec, LoadingState, ActionResult } from "../lib/types";

const targets = ref<BackupSpec[]>([]);
const loadingState = ref<LoadingState>({
  isLoading: false,
  error: null,
});
const actionStates = ref<{ [host: string]: { isLoading: boolean; result: ActionResult | null } }>({});

onMounted(async () => {
  await loadTargets();
});

async function loadTargets() {
  loadingState.value.isLoading = true;
  loadingState.value.error = null;
  
  try {
    targets.value = await backupApiClient.listTargets();
  } catch (err) {
    loadingState.value.error = err instanceof Error ? err.message : "Failed to load targets";
  } finally {
    loadingState.value.isLoading = false;
  }
}

async function triggerBackup(host: string) {
  if (!actionStates.value[host]) {
    actionStates.value[host] = { isLoading: false, result: null };
  }
  
  actionStates.value[host].isLoading = true;
  actionStates.value[host].result = null;
  
  try {
    const result = await backupApiClient.triggerBackup(host);
    actionStates.value[host].result = result;
    setTimeout(() => {
      if (actionStates.value[host]) {
        actionStates.value[host].result = null;
      }
    }, 5000);
  } catch (err) {
    actionStates.value[host].result = {
      success: false,
      message: err instanceof Error ? err.message : "Backup failed",
    };
  } finally {
    actionStates.value[host].isLoading = false;
  }
}

async function triggerCleanup(host: string) {
  if (!actionStates.value[host]) {
    actionStates.value[host] = { isLoading: false, result: null };
  }
  
  actionStates.value[host].isLoading = true;
  actionStates.value[host].result = null;
  
  try {
    const result = await backupApiClient.triggerCleanup(host);
    actionStates.value[host].result = result;
    setTimeout(() => {
      if (actionStates.value[host]) {
        actionStates.value[host].result = null;
      }
    }, 5000);
  } catch (err) {
    actionStates.value[host].result = {
      success: false,
      message: err instanceof Error ? err.message : "Cleanup failed",
    };
  } finally {
    actionStates.value[host].isLoading = false;
  }
}
</script>

<template>
  <div class="target-list">
    <h2>Backup Targets</h2>
    
    <div v-if="loadingState.isLoading" class="loading">
      Loading targets...
    </div>
    
    <div v-else-if="loadingState.error" class="error">
      Error: {{ loadingState.error }}
      <button @click="loadTargets" class="retry-btn">Retry</button>
    </div>
    
    <div v-else-if="targets.length === 0" class="empty">
      No backup targets configured
    </div>
    
    <div v-else class="targets">
      <div
        v-for="target in targets"
        :key="target.host"
        class="target-card"
      >
        <div class="target-header">
          <h3>{{ target.host }}</h3>
          <div class="target-connection">
            {{ target.user }}@{{ target.host }}:{{ target.port }}
          </div>
        </div>
        
        <div class="target-dirs">
          <h4>Directories:</h4>
          <ul>
            <li v-for="dir in target.dirs" :key="dir">{{ dir }}</li>
          </ul>
        </div>
        
        <div v-if="target.excludes.length > 0" class="target-excludes">
          <h4>Excludes:</h4>
          <ul>
            <li v-for="exclude in target.excludes" :key="exclude">{{ exclude }}</li>
          </ul>
        </div>
        
        <div class="target-actions">
          <ActionButtonComponent
            label="Backup"
            :disabled="actionStates[target.host]?.isLoading || false"
            @click="triggerBackup(target.host)"
          />
          <ActionButtonComponent
            label="Cleanup"
            variant="secondary"
            :disabled="actionStates[target.host]?.isLoading || false"
            @click="triggerCleanup(target.host)"
          />
        </div>
        
        <div
          v-if="actionStates[target.host]?.result"
          :class="[
            'action-result',
            actionStates[target.host]?.result?.success ? 'success' : 'error'
          ]"
        >
          {{ actionStates[target.host]?.result?.message }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.target-list {
  padding: 1rem;
}

.target-list h2 {
  margin: 0 0 1rem 0;
  color: #1f2937;
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

.targets {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}

.target-card {
  padding: 1.5rem;
  background-color: white;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.target-header h3 {
  margin: 0 0 0.5rem 0;
  color: #1f2937;
  font-size: 1.25rem;
}

.target-connection {
  color: #6b7280;
  font-family: monospace;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.target-dirs, .target-excludes {
  margin-bottom: 1rem;
}

.target-dirs h4, .target-excludes h4 {
  margin: 0 0 0.5rem 0;
  color: #374151;
  font-size: 0.875rem;
  font-weight: 600;
}

.target-dirs ul, .target-excludes ul {
  margin: 0;
  padding-left: 1rem;
  color: #6b7280;
  font-family: monospace;
  font-size: 0.875rem;
}

.target-actions {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.action-result {
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
</style>