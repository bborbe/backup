import { describe, it, expect } from 'vitest';
import type { 
  ApiResponse, 
  BackupHost, 
  BackupPort, 
  BackupUser, 
  BackupDir, 
  BackupExclude,
  BackupSpec,
  Target,
  BackupStatus,
  PageMeta,
  LoadingState,
  ActionResult
} from '../lib/types';

describe('Type Definitions', () => {
  it('should validate ApiResponse interface', () => {
    const response: ApiResponse<string> = {
      data: 'test data',
      success: true,
      message: 'optional message',
    };

    expect(response.data).toBe('test data');
    expect(response.success).toBe(true);
    expect(response.message).toBe('optional message');
  });

  it('should validate BackupSpec interface', () => {
    const spec: BackupSpec = {
      host: 'example.com',
      port: 22,
      user: 'root',
      dirs: ['/data', '/home'],
      excludes: ['*.log', 'tmp/'],
    };

    expect(spec.host).toBe('example.com');
    expect(spec.port).toBe(22);
    expect(spec.user).toBe('root');
    expect(spec.dirs).toEqual(['/data', '/home']);
    expect(spec.excludes).toEqual(['*.log', 'tmp/']);
  });

  it('should validate Target interface', () => {
    const target: Target = {
      metadata: {
        name: 'test-target',
        namespace: 'default',
        createdAt: '2024-01-01T00:00:00Z',
      },
      spec: {
        host: 'example.com',
        port: 22,
        user: 'root',
        dirs: ['/data'],
        excludes: [],
      },
    };

    expect(target.metadata.name).toBe('test-target');
    expect(target.metadata.namespace).toBe('default');
    expect(target.spec.host).toBe('example.com');
  });

  it('should validate BackupStatus interface', () => {
    const status: BackupStatus = {
      'host1.example.com': '2024-01-01',
      'host2.example.com': '2024-01-02',
    };

    expect(status['host1.example.com']).toBe('2024-01-01');
    expect(status['host2.example.com']).toBe('2024-01-02');
  });

  it('should validate PageMeta interface', () => {
    const meta: PageMeta = {
      title: 'Dashboard',
      description: 'Backup service dashboard',
    };

    expect(meta.title).toBe('Dashboard');
    expect(meta.description).toBe('Backup service dashboard');

    // Test without optional description
    const metaMinimal: PageMeta = {
      title: 'Dashboard',
    };

    expect(metaMinimal.title).toBe('Dashboard');
    expect(metaMinimal.description).toBeUndefined();
  });

  it('should validate LoadingState interface', () => {
    const loadingState: LoadingState = {
      isLoading: true,
      error: null,
    };

    expect(loadingState.isLoading).toBe(true);
    expect(loadingState.error).toBe(null);

    const errorState: LoadingState = {
      isLoading: false,
      error: 'Something went wrong',
    };

    expect(errorState.isLoading).toBe(false);
    expect(errorState.error).toBe('Something went wrong');
  });

  it('should validate ActionResult interface', () => {
    const successResult: ActionResult = {
      success: true,
      message: 'Operation completed successfully',
    };

    expect(successResult.success).toBe(true);
    expect(successResult.message).toBe('Operation completed successfully');

    const errorResult: ActionResult = {
      success: false,
      message: 'Operation failed',
    };

    expect(errorResult.success).toBe(false);
    expect(errorResult.message).toBe('Operation failed');
  });

  it('should validate individual backup value types', () => {
    const host: BackupHost = { value: 'example.com' };
    const port: BackupPort = { value: 22 };
    const user: BackupUser = { value: 'root' };
    const dir: BackupDir = { value: '/data' };
    const exclude: BackupExclude = { value: '*.log' };

    expect(host.value).toBe('example.com');
    expect(port.value).toBe(22);
    expect(user.value).toBe('root');
    expect(dir.value).toBe('/data');
    expect(exclude.value).toBe('*.log');
  });
});