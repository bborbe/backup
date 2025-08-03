import { describe, it, expect } from 'vitest';

describe('Simple Tests', () => {
  it('should pass basic math test', () => {
    expect(2 + 2).toBe(4);
  });

  it('should test string operations', () => {
    const str = 'Hello World';
    expect(str.toLowerCase()).toBe('hello world');
    expect(str.includes('World')).toBe(true);
  });

  it('should test array operations', () => {
    const arr = [1, 2, 3, 4, 5];
    expect(arr.length).toBe(5);
    expect(arr.includes(3)).toBe(true);
    expect(arr.filter(x => x > 3)).toEqual([4, 5]);
  });
});