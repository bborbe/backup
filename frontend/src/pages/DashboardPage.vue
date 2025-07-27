<script setup lang="ts">
import { onMounted, ref } from "vue";
import BackupStatusOverviewComponent from "../components/BackupStatusOverviewComponent.vue";
import ActionPanelComponent from "../components/ActionPanelComponent.vue";
import TargetListComponent from "../components/TargetListComponent.vue";

const currentTime = ref("");

onMounted(() => {
  updateTime();
  setInterval(updateTime, 1000);
});

function updateTime() {
  currentTime.value = new Date().toISOString().split(".")[0].split("T")[1];
}
</script>

<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <h1>Backup Service Dashboard</h1>
      <div class="time">{{ currentTime }}</div>
    </header>
    
    <main class="dashboard-content">
      <ActionPanelComponent />
      <BackupStatusOverviewComponent />
      <TargetListComponent />
    </main>
  </div>
</template>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f8fafc;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: white;
  border-bottom: 1px solid #e5e7eb;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.dashboard-header h1 {
  margin: 0;
  color: #1f2937;
  font-size: 1.875rem;
  font-weight: 700;
}

.time {
  font-family: monospace;
  font-size: 1.125rem;
  color: #6b7280;
  background-color: #f3f4f6;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
}

.dashboard-content {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

@media (max-width: 768px) {
  .dashboard-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .dashboard-content {
    padding: 1rem;
  }
}
</style>