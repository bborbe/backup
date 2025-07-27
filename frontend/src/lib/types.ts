// API Response Types
export interface ApiResponse<T> {
  data: T;
  success: boolean;
  message?: string;
}

// Backup Service Types (based on Go service types)
export interface BackupHost {
  value: string;
}

export interface BackupPort {
  value: number;
}

export interface BackupUser {
  value: string;
}

export interface BackupDir {
  value: string;
}

export interface BackupExclude {
  value: string;
}

export interface BackupSpec {
  host: string;
  port: number;
  user: string;
  dirs: string[];
  excludes: string[];
}

export interface Target {
  metadata: {
    name: string;
    namespace: string;
    createdAt: string;
  };
  spec: BackupSpec;
}

// Status response: map of host to last backup date
export interface BackupStatus {
  [host: string]: string; // date in YYYY-MM-DD format
}

// Component Props Types
export interface PageMeta {
  title: string;
  description?: string;
}

// Loading states
export interface LoadingState {
  isLoading: boolean;
  error: string | null;
}

// Action results
export interface ActionResult {
  success: boolean;
  message: string;
}