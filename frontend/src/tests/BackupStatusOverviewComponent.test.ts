import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import { AxiosError } from 'axios';
import BackupStatusOverviewComponent from '../components/BackupStatusOverviewComponent.vue';

vi.mock("../lib/BackupApiClient.ts", () => {
  const client = {
    getStatus: vi.fn(),
    triggerBackup: vi.fn(),
    triggerCleanup: vi.fn(),
    triggerBackupAll: vi.fn(),
    triggerCleanupAll: vi.fn(),
  };
  return { default: client, backupApiClient: client, BackupApiClient: class {} };
});

describe('BackupStatusOverviewComponent', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockBackupAlreadyRunningError = (host: string) =>
    new AxiosError(
      "Request failed with status code 409",
      "ERR_BAD_REQUEST",
      undefined,
      undefined,
      {
        status: 409,
        statusText: "Conflict",
        headers: {},
        config: {},
        data: {
          error: {
            code: "BACKUP_ALREADY_RUNNING",
            message: `backup for ${host} is already running`,
          },
        },
      } as never,
    );

  const mockCleanupAlreadyRunningError = (host: string) =>
    new AxiosError(
      "Request failed with status code 409",
      "ERR_BAD_REQUEST",
      undefined,
      undefined,
      {
        status: 409,
        statusText: "Conflict",
        headers: {},
        config: {},
        data: {
          error: {
            code: "CLEANUP_ALREADY_RUNNING",
            message: `cleanup for ${host} is already running`,
          },
        },
      } as never,
    );

  const mockInternalError = (message: string) =>
    new AxiosError(
      "Request failed with status code 500",
      "ERR_INTERNAL_SERVER_ERROR",
      undefined,
      undefined,
      {
        status: 500,
        statusText: "Internal Server Error",
        headers: {},
        config: {},
        data: {
          error: {
            code: "INTERNAL_ERROR",
            message,
          },
        },
      } as never,
    );

  const mockTextPlainError = () =>
    new AxiosError(
      "Request failed with status code 500",
      "ERR_INTERNAL_SERVER_ERROR",
      undefined,
      undefined,
      {
        status: 500,
        statusText: "Internal Server Error",
        headers: {},
        config: {},
        data: "request failed: something broke\n",
      } as never,
    );

  it('shows warning severity for BACKUP_ALREADY_RUNNING conflict on per-host trigger', async () => {
    const { backupApiClient } = await import('../lib/BackupApiClient.ts');
    vi.mocked(backupApiClient.getStatus).mockResolvedValue({ "host1.example.com": "2024-01-01" });
    vi.mocked(backupApiClient.triggerBackup).mockRejectedValue(mockBackupAlreadyRunningError("host1.example.com"));

    const wrapper = mount(BackupStatusOverviewComponent);
    await flushPromises();

    await wrapper.findAll(".card-actions button")[0].trigger("click");
    await flushPromises();

    const messageEl = wrapper.get(".action-message");
    const text = messageEl.text();
    const classes = messageEl.classes();

    expect(text).toContain("host1.example.com");
    expect(text).toContain("already running");
    expect(text).not.toContain("Request failed with status code");
    expect(text).not.toContain("500");
    expect(classes).toContain("status-warning");
    expect(classes).not.toContain("status-error");
  });

  it('shows warning severity for CLEANUP_ALREADY_RUNNING conflict on per-host trigger', async () => {
    const { backupApiClient } = await import('../lib/BackupApiClient.ts');
    vi.mocked(backupApiClient.getStatus).mockResolvedValue({ "host1.example.com": "2024-01-01" });
    vi.mocked(backupApiClient.triggerCleanup).mockRejectedValue(mockCleanupAlreadyRunningError("host1.example.com"));

    const wrapper = mount(BackupStatusOverviewComponent);
    await flushPromises();

    await wrapper.findAll(".card-actions button")[1].trigger("click");
    await flushPromises();

    const messageEl = wrapper.get(".action-message");
    const text = messageEl.text();
    const classes = messageEl.classes();

    expect(text).toContain("host1.example.com");
    expect(text).toContain("already running");
    expect(classes).toContain("status-warning");
    expect(classes).not.toContain("status-error");
  });

  it('shows error severity for INTERNAL_ERROR on per-host backup trigger', async () => {
    const { backupApiClient } = await import('../lib/BackupApiClient.ts');
    vi.mocked(backupApiClient.getStatus).mockResolvedValue({ "host1.example.com": "2024-01-01" });
    vi.mocked(backupApiClient.triggerBackup).mockRejectedValue(
      mockInternalError("backup test-target failed: disk on fire")
    );

    const wrapper = mount(BackupStatusOverviewComponent);
    await flushPromises();

    await wrapper.findAll(".card-actions button")[0].trigger("click");
    await flushPromises();

    const messageEl = wrapper.get(".action-message");
    const text = messageEl.text();
    const classes = messageEl.classes();

    expect(text).toBe("backup test-target failed: disk on fire");
    expect(classes).toContain("status-error");
    expect(classes).not.toContain("status-warning");
  });

  it('shows error severity and transport message for text/plain failure on per-host backup trigger', async () => {
    const { backupApiClient } = await import('../lib/BackupApiClient.ts');
    vi.mocked(backupApiClient.getStatus).mockResolvedValue({ "host1.example.com": "2024-01-01" });
    vi.mocked(backupApiClient.triggerBackup).mockRejectedValue(mockTextPlainError());

    const wrapper = mount(BackupStatusOverviewComponent);
    await flushPromises();

    await wrapper.findAll(".card-actions button")[0].trigger("click");
    await flushPromises();

    const messageEl = wrapper.get(".action-message");
    const text = messageEl.text();
    const classes = messageEl.classes();

    expect(text).toBe("Request failed with status code 500");
    expect(text).not.toContain("undefined");
    expect(text.length).toBeGreaterThan(0);
    expect(classes).toContain("status-error");
  });

  it('shows warning severity for BACKUP_ALREADY_RUNNING conflict on trigger-all backup', async () => {
    const { backupApiClient } = await import('../lib/BackupApiClient.ts');
    vi.mocked(backupApiClient.getStatus).mockResolvedValue({ "host1.example.com": "2024-01-01" });
    vi.mocked(backupApiClient.triggerBackupAll).mockRejectedValue(mockBackupAlreadyRunningError("host1.example.com"));

    const wrapper = mount(BackupStatusOverviewComponent);
    await flushPromises();

    await (wrapper.vm as unknown as { triggerBackupAll: () => Promise<void> }).triggerBackupAll();
    await flushPromises();

    const messageEl = wrapper.get(".action-message");
    const text = messageEl.text();
    const classes = messageEl.classes();

    expect(text).toContain("host1.example.com");
    expect(text).toContain("already running");
    expect(classes).toContain("status-warning");
    expect(classes).not.toContain("status-error");
  });

  it('shows warning severity for CLEANUP_ALREADY_RUNNING conflict on trigger-all cleanup', async () => {
    const { backupApiClient } = await import('../lib/BackupApiClient.ts');
    vi.mocked(backupApiClient.getStatus).mockResolvedValue({ "host1.example.com": "2024-01-01" });
    vi.mocked(backupApiClient.triggerCleanupAll).mockRejectedValue(mockCleanupAlreadyRunningError("host1.example.com"));

    const wrapper = mount(BackupStatusOverviewComponent);
    await flushPromises();

    await (wrapper.vm as unknown as { triggerCleanupAll: () => Promise<void> }).triggerCleanupAll();
    await flushPromises();

    const messageEl = wrapper.get(".action-message");
    const text = messageEl.text();
    const classes = messageEl.classes();

    expect(text).toContain("host1.example.com");
    expect(text).toContain("already running");
    expect(classes).toContain("status-warning");
    expect(classes).not.toContain("status-error");
  });

  it('keeps success presentation unchanged with neutral styling', async () => {
    const { backupApiClient } = await import('../lib/BackupApiClient.ts');
    vi.mocked(backupApiClient.getStatus).mockResolvedValue({ "host1.example.com": "2024-01-01" });
    vi.mocked(backupApiClient.triggerBackup).mockResolvedValue({
      success: true,
      message: "backup host1.example.com completed",
    });

    const wrapper = mount(BackupStatusOverviewComponent);
    await flushPromises();

    await wrapper.findAll(".card-actions button")[0].trigger("click");
    await flushPromises();

    const messageEl = wrapper.get(".action-message");
    const text = messageEl.text();
    const classes = messageEl.classes();

    expect(text).toBe("backup host1.example.com completed");
    expect(classes).not.toContain("status-warning");
    expect(classes).not.toContain("status-error");
  });
});
