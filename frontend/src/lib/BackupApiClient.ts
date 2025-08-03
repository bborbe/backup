// @ts-ignore
import { apiService } from "./ApiService.ts";
import type { BackupStatus, BackupSpec, ActionResult } from "./types";

export class BackupApiClient {
  /**
   * Get backup status for all hosts
   * Returns map of host -> last backup date
   */
  async getStatus(): Promise<BackupStatus> {
    // Return mock data in development mode
    if ((import.meta as any).env.DEV) {
      return this.getMockBackupStatus();
    }
    return await apiService.get<BackupStatus>("/status");
  }

  /**
   * Generate mock backup status data for development
   */
  private getMockBackupStatus(): BackupStatus {
    const now = new Date();
    const today = now.toISOString().split('T')[0];
    const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString().split('T')[0];
    const weekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];

    return {
      "burn.hm.benjamin-borbe.de": weekAgo,
      "co2hz.hm.benjamin-borbe.de": today,
      "co2wz.hm.benjamin-borbe.de": today,
      "fire.hm.benjamin-borbe.de": weekAgo,
      "hetzner-1.benjamin-borbe.de": today,
      "nuke-k3s-agent-0.hm.benjamin-borbe.de": today,
      "nuke-k3s-agent-1.hm.benjamin-borbe.de": today,
      "nuke-k3s-dev-0.hm.benjamin-borbe.de": weekAgo,
      "nuke-k3s-kafka-0.hm.benjamin-borbe.de": weekAgo,
      "nuke-k3s-kafka-1.hm.benjamin-borbe.de": weekAgo,
      "nuke-k3s-kafka-2.hm.benjamin-borbe.de": weekAgo,
      "nuke-k3s-longhorn-0.hm.benjamin-borbe.de": weekAgo,
      "nuke-k3s-longhorn-1.hm.benjamin-borbe.de": weekAgo,
      "nuke-k3s-master-0.hm.benjamin-borbe.de": weekAgo,
      "nuke-k3s-master-1.hm.benjamin-borbe.de": weekAgo,
      "nuke-k3s-master-2.hm.benjamin-borbe.de": weekAgo,
      "postgres-server.benjamin-borbe.de": yesterday,
      "mysql-db.benjamin-borbe.de": today,
      "redis-cache.benjamin-borbe.de": yesterday,
      "monitoring.benjamin-borbe.de": today
    };
  }

  /**
   * List all backup targets
   */
  async listTargets(): Promise<BackupSpec[]> {
    return await apiService.get<BackupSpec[]>("/list");
  }

  /**
   * Trigger backup for all targets
   */
  async triggerBackupAll(): Promise<ActionResult> {
    const response = await apiService.post<{}, string>("/backup/all", {});
    return {
      success: true,
      message: typeof response === "string" ? response : "Backup started successfully",
    };
  }

  /**
   * Trigger backup for specific target
   */
  async triggerBackup(name: string): Promise<ActionResult> {
    const response = await apiService.post<{}, string>(`/backup/${name}`, {});
    return {
      success: true,
      message: typeof response === "string" ? response : `Backup for ${name} started successfully`,
    };
  }

  /**
   * Trigger cleanup for all targets
   */
  async triggerCleanupAll(): Promise<ActionResult> {
    const response = await apiService.post<{}, string>("/cleanup/all", {});
    return {
      success: true,
      message: typeof response === "string" ? response : "Cleanup started successfully",
    };
  }

  /**
   * Trigger cleanup for specific target
   */
  async triggerCleanup(name: string): Promise<ActionResult> {
    const response = await apiService.post<{}, string>(`/cleanup/${name}`, {});
    return {
      success: true,
      message: typeof response === "string" ? response : `Cleanup for ${name} started successfully`,
    };
  }

  /**
   * Check service health
   */
  async getHealth(): Promise<string> {
    return await apiService.get<string>("/healthz");
  }

  /**
   * Check service readiness
   */
  async getReadiness(): Promise<string> {
    return await apiService.get<string>("/readiness");
  }
}

// Create and export client instance
export const backupApiClient = new BackupApiClient();
export default backupApiClient;