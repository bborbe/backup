// @ts-ignore
import { apiService } from "./ApiService.ts";
import type { BackupStatus, BackupSpec, ActionResult } from "./types";

export class BackupApiClient {
  /**
   * Get backup status for all hosts
   * Returns map of host -> last backup date
   */
  async getStatus(): Promise<BackupStatus> {
    return await apiService.get<BackupStatus>("/status");
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