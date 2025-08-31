<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import BackupStatusOverviewComponent from "../components/BackupStatusOverviewComponent.vue";
import ActionButtonComponent from "../components/ActionButtonComponent.vue";
// @ts-ignore
import backupApiClient from "../lib/BackupApiClient.ts";
import type { BackupStatus } from "../lib/types";

interface FilterState {
  total: boolean;
  healthy: boolean;
  warning: boolean;
  critical: boolean;
}

const currentTime = ref("");
const backupOverviewRef = ref<InstanceType<typeof BackupStatusOverviewComponent>>();
const status = ref<BackupStatus>({});
const baseFilters = ref({
  healthy: true,
  warning: true,
  critical: true
});

const activeFilters = computed<FilterState>({
  get() {
    const allSelected = baseFilters.value.healthy && baseFilters.value.warning && baseFilters.value.critical;
    return {
      total: allSelected,
      healthy: baseFilters.value.healthy,
      warning: baseFilters.value.warning,
      critical: baseFilters.value.critical
    };
  },
  set(newValue) {
    baseFilters.value = {
      healthy: newValue.healthy,
      warning: newValue.warning,
      critical: newValue.critical
    };
  }
});

onMounted(async () => {
  updateTime();
  setInterval(updateTime, 1000);
  await loadStatus();
  // Auto-refresh status every 30 seconds
  setInterval(loadStatus, 30000);
});

function updateTime() {
  currentTime.value = new Date().toISOString().split(".")[0].split("T")[1];
}

async function loadStatus() {
  try {
    const result = await backupApiClient.getStatus();
    if (result && typeof result === 'object' && !Array.isArray(result)) {
      status.value = result;
    }
  } catch (err) {
    console.error("Failed to load backup status:", err);
  }
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

const metrics = computed(() => {
  const hosts = Object.entries(status.value);
  const total = hosts.length;
  
  const counts = hosts.reduce((acc, [, lastBackup]) => {
    const category = getStatusCategory(lastBackup);
    acc[category]++;
    return acc;
  }, { healthy: 0, warning: 0, critical: 0 });
  
  return {
    total,
    healthy: counts.healthy,
    warning: counts.warning,
    critical: counts.critical
  };
});

function toggleFilter(category: keyof FilterState) {
  if (category === 'total') {
    // If total is clicked, toggle all three categories to the same state
    const currentTotalState = activeFilters.value.total;
    const newState = !currentTotalState;
    
    baseFilters.value.healthy = newState;
    baseFilters.value.warning = newState;
    baseFilters.value.critical = newState;
  } else {
    // Toggle individual category
    baseFilters.value[category] = !baseFilters.value[category];
  }
}
</script>

<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <h1>Backup Service Dashboard</h1>
      
      <div class="header-actions">
        <div class="action-buttons">
          <ActionButtonComponent
            label="Backup All"
            variant="primary"
            :disabled="backupOverviewRef?.actionState?.isBackingUp || backupOverviewRef?.actionState?.isCleaningUp"
            @click="backupOverviewRef?.triggerBackupAll()"
          />
          <ActionButtonComponent
            label="Cleanup All"
            variant="danger"
            :disabled="backupOverviewRef?.actionState?.isBackingUp || backupOverviewRef?.actionState?.isCleaningUp"
            @click="backupOverviewRef?.triggerCleanupAll()"
          />
        </div>
        <div class="time">{{ currentTime }}</div>
      </div>
    </header>
    
    <main class="dashboard-content">
      <div class="metrics-section">
        <div class="metrics-grid">
          <div 
            class="metric-card"
            :class="{ active: activeFilters.total, inactive: !activeFilters.total }"
            @click="toggleFilter('total')"
          >
            <div class="metric-number">{{ metrics.total }}</div>
            <div class="metric-label">Total Hosts</div>
          </div>
          <div 
            class="metric-card"
            :class="{ active: activeFilters.healthy, inactive: !activeFilters.healthy }"
            @click="toggleFilter('healthy')"
          >
            <div class="metric-number success">{{ metrics.healthy }}</div>
            <div class="metric-label">Healthy</div>
          </div>
          <div 
            class="metric-card"
            :class="{ active: activeFilters.warning, inactive: !activeFilters.warning }"
            @click="toggleFilter('warning')"
          >
            <div class="metric-number warning">{{ metrics.warning }}</div>
            <div class="metric-label">Warning</div>
          </div>
          <div 
            class="metric-card"
            :class="{ active: activeFilters.critical, inactive: !activeFilters.critical }"
            @click="toggleFilter('critical')"
          >
            <div class="metric-number error">{{ metrics.critical }}</div>
            <div class="metric-label">Critical</div>
          </div>
        </div>
      </div>
      
      <BackupStatusOverviewComponent ref="backupOverviewRef" :filters="activeFilters" />
    </main>
  </div>
</template>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: var(--bg-primary);
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md) var(--spacing-xl);
  background-color: var(--bg-secondary);
  border-bottom: 1px solid var(--border-primary);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
}

.action-buttons {
  display: flex;
  gap: var(--element-gap);
}

.dashboard-header h1 {
  margin: 0;
  color: var(--text-primary);
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  line-height: var(--line-height-tight);
}

.time {
  font-family: monospace;
  font-size: var(--font-size-lg);
  color: var(--text-secondary);
  background-color: var(--bg-tertiary);
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-secondary);
  font-weight: var(--font-weight-medium);
}

.dashboard-content {
  padding: 0;
  max-width: none;
  margin: 0;
}

.metrics-section {
  padding: var(--spacing-lg) var(--spacing-xl);
  background-color: var(--bg-primary);
  border-bottom: 1px solid var(--border-secondary);
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-lg);
  max-width: 1200px;
  margin: 0 auto;
}

.metric-card {
  background-color: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  padding: var(--metric-card-padding);
  text-align: center;
  transition: var(--transition-fast);
  cursor: pointer;
  user-select: none;
}

.metric-card:hover {
  border-color: var(--border-accent);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.metric-card.active {
  border-color: var(--accent-primary);
  background-color: var(--bg-secondary);
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.metric-card.inactive {
  border-color: var(--border-secondary);
  background-color: var(--bg-quaternary);
  opacity: 0.5;
  transform: none;
}

.metric-card.inactive:hover {
  border-color: var(--border-accent);
  opacity: 0.7;
  transform: translateY(-1px);
}

.metric-card.active:hover {
  border-color: var(--accent-primary-hover);
}

.metric-number {
  font-size: var(--metric-number-size);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  line-height: var(--line-height-tight);
  margin-bottom: var(--spacing-xs);
}

.metric-number.success {
  color: var(--financial-positive);
}

.metric-number.warning {
  color: var(--status-warning);
}

.metric-number.error {
  color: var(--financial-negative);
}

.metric-label {
  font-size: var(--metric-label-size);
  font-weight: var(--font-weight-medium);
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

@media (max-width: 768px) {
  .dashboard-header {
    flex-direction: column;
    gap: var(--spacing-md);
    text-align: center;
  }
  
  .header-actions {
    flex-direction: column;
    gap: var(--spacing-md);
  }
  
  .dashboard-content {
    padding: 0;
  }
}
</style>