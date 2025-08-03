import { describe, it, expect } from 'vitest';
import { ApiService } from '../lib/ApiService';

describe('ApiService', () => {
  it('should create ApiService instance', () => {
    const apiService = new ApiService({ timeout: 5000 });
    expect(apiService).toBeInstanceOf(ApiService);
  });

  it('should create ApiService with default config', () => {
    const apiService = new ApiService({});
    expect(apiService).toBeInstanceOf(ApiService);
  });
});