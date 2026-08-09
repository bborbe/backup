import { describe, it, expect } from 'vitest';
import { AxiosError } from 'axios';
import {
  extractApiError,
  toActionFailure,
} from '../lib/ApiError';

describe('ApiError', () => {
  describe('extractApiError', () => {
    it('extracts code and message from structured 409 BACKUP_ALREADY_RUNNING response', () => {
      const err = new AxiosError(
        "Request failed with status code 409",
        "ERR_BAD_REQUEST",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: { code: "BACKUP_ALREADY_RUNNING", message: "backup for host1.example.com is already running" } } } as never,
      );

      const result = extractApiError(err);

      expect(result.code).toBe("BACKUP_ALREADY_RUNNING");
      expect(result.message).toBe("backup for host1.example.com is already running");
    });

    it('returns nulls for text/plain body', () => {
      const err = new AxiosError(
        "Request failed with status code 500",
        "ERR_INTERNAL_SERVER_ERROR",
        undefined,
        undefined,
        { status: 500, statusText: "Internal Server Error", headers: {}, config: {}, data: "request failed: something broke\n" } as never,
      );

      const result = extractApiError(err);

      expect(result.code).toBeNull();
      expect(result.message).toBeNull();
    });

    it('returns nulls for null input', () => {
      const result = extractApiError(null);
      expect(result).toEqual({ code: null, message: null });
    });

    it('returns nulls for undefined input', () => {
      const result = extractApiError(undefined);
      expect(result).toEqual({ code: null, message: null });
    });

    it('returns nulls for primitive string input', () => {
      const result = extractApiError("boom");
      expect(result).toEqual({ code: null, message: null });
    });

    it('returns nulls for numeric input', () => {
      const result = extractApiError(42);
      expect(result).toEqual({ code: null, message: null });
    });

    it('returns nulls for empty object input', () => {
      const result = extractApiError({});
      expect(result).toEqual({ code: null, message: null });
    });

    it('returns nulls for array body', () => {
      const err = new AxiosError(
        "Request failed",
        "ERR",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: [{ error: { code: "X", message: "y" } }] } as never,
      );
      const result = extractApiError(err);
      expect(result).toEqual({ code: null, message: null });
    });

    it('returns nulls for error field being a string', () => {
      const err = new AxiosError(
        "Request failed",
        "ERR",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: "boom" } } as never,
      );
      const result = extractApiError(err);
      expect(result).toEqual({ code: null, message: null });
    });

    it('returns nulls for error field being null', () => {
      const err = new AxiosError(
        "Request failed",
        "ERR",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: null } } as never,
      );
      const result = extractApiError(err);
      expect(result).toEqual({ code: null, message: null });
    });

    it('returns nulls for error field being an array', () => {
      const err = new AxiosError(
        "Request failed",
        "ERR",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: [] } } as never,
      );
      const result = extractApiError(err);
      expect(result).toEqual({ code: null, message: null });
    });

    it('returns nulls for error.message being empty string', () => {
      const err = new AxiosError(
        "Request failed",
        "ERR",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: { code: "BACKUP_ALREADY_RUNNING", message: "" } } } as never,
      );
      const result = extractApiError(err);
      expect(result).toEqual({ code: null, message: null });
    });
  });

  describe('toActionFailure', () => {
    const fallbackMessage = "fallback message";

    it('returns warning severity for BACKUP_ALREADY_RUNNING 409', () => {
      const err = new AxiosError(
        "Request failed with status code 409",
        "ERR_BAD_REQUEST",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: { code: "BACKUP_ALREADY_RUNNING", message: "backup for host1.example.com is already running" } } } as never,
      );

      const failure = toActionFailure(err, fallbackMessage);

      expect(failure.severity).toBe("warning");
      expect(failure.message).toBe("backup for host1.example.com is already running");
      expect(failure.message).toContain("host1.example.com");
      expect(failure.message).toContain("already running");
      expect(failure.message).not.toContain("Request failed with status code");
      expect(failure.message).not.toContain("500");
    });

    it('returns warning severity for CLEANUP_ALREADY_RUNNING 409', () => {
      const err = new AxiosError(
        "Request failed with status code 409",
        "ERR_BAD_REQUEST",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: { code: "CLEANUP_ALREADY_RUNNING", message: "cleanup for host1.example.com is already running" } } } as never,
      );

      const failure = toActionFailure(err, fallbackMessage);

      expect(failure.severity).toBe("warning");
      expect(failure.message).toBe("cleanup for host1.example.com is already running");
    });

    it('returns error severity for INTERNAL_ERROR 500', () => {
      const err = new AxiosError(
        "Request failed with status code 500",
        "ERR_INTERNAL_SERVER_ERROR",
        undefined,
        undefined,
        { status: 500, statusText: "Internal Server Error", headers: {}, config: {}, data: { error: { code: "INTERNAL_ERROR", message: "backup test-target failed: disk on fire" } } } as never,
      );

      const failure = toActionFailure(err, fallbackMessage);

      expect(failure.severity).toBe("error");
      expect(failure.message).toBe("backup test-target failed: disk on fire");
    });

    it('falls back to transport message for text/plain body at 500', () => {
      const err = new AxiosError(
        "Request failed with status code 500",
        "ERR_INTERNAL_SERVER_ERROR",
        undefined,
        undefined,
        { status: 500, statusText: "Internal Server Error", headers: {}, config: {}, data: "request failed: something broke\n" } as never,
      );

      const failure = toActionFailure(err, fallbackMessage);

      expect(failure.severity).toBe("error");
      expect(failure.message).toBe("Request failed with status code 500");
      expect(failure.message).not.toContain("undefined");
      expect(failure.message.length).toBeGreaterThan(0);
    });

    it('falls back to err.message for network failure (no response)', () => {
      const err = new AxiosError(
        "Network Error",
        "ERR_NETWORK",
        undefined,
        undefined,
        undefined,
      );

      const failure = toActionFailure(err, fallbackMessage);

      expect(failure.severity).toBe("error");
      expect(failure.message).toBe("Network Error");
    });

    it('returns error severity with fallbackMessage for null input', () => {
      const failure = toActionFailure(null, fallbackMessage);
      expect(failure.severity).toBe("error");
      expect(failure.message).toBe(fallbackMessage);
    });

    it('returns error severity with fallbackMessage for undefined input', () => {
      const failure = toActionFailure(undefined, fallbackMessage);
      expect(failure.severity).toBe("error");
      expect(failure.message).toBe(fallbackMessage);
    });

    it('returns error severity with fallbackMessage for string primitive input', () => {
      const failure = toActionFailure("boom", fallbackMessage);
      expect(failure.severity).toBe("error");
      expect(failure.message).toBe(fallbackMessage);
    });

    it('returns error severity with fallbackMessage for numeric input', () => {
      const failure = toActionFailure(42 as unknown, fallbackMessage);
      expect(failure.severity).toBe("error");
      expect(failure.message).toBe(fallbackMessage);
    });

    it('returns error severity with fallbackMessage for empty object input', () => {
      const failure = toActionFailure({} as unknown, fallbackMessage);
      expect(failure.severity).toBe("error");
      expect(failure.message).toBe(fallbackMessage);
    });

    it('falls back to transport message for envelope with empty error.message', () => {
      const err = new AxiosError(
        "Request failed with status code 409",
        "ERR_BAD_REQUEST",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: { code: "BACKUP_ALREADY_RUNNING", message: "" } } } as never,
      );

      const failure = toActionFailure(err, fallbackMessage);

      expect(failure.severity).toBe("error"); // code present but message empty → falls back to transport, so severity is error
      expect(failure.message).toBe("Request failed with status code 409");
      expect(failure.message.length).toBeGreaterThan(0);
    });

    it('falls back to transport message for error field being a string', () => {
      const err = new AxiosError(
        "Request failed",
        "ERR",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: "boom" } } as never,
      );
      const failure = toActionFailure(err, fallbackMessage);
      expect(failure.severity).toBe("error");
      expect(failure.message).toBe("Request failed");
    });

    it('falls back to transport message for error field being null', () => {
      const err = new AxiosError(
        "Request failed",
        "ERR",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: null } } as never,
      );
      const failure = toActionFailure(err, fallbackMessage);
      expect(failure.severity).toBe("error");
      expect(failure.message).toBe("Request failed");
    });

    it('falls back to transport message for error field being an array', () => {
      const err = new AxiosError(
        "Request failed",
        "ERR",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: { error: [] } } as never,
      );
      const failure = toActionFailure(err, fallbackMessage);
      expect(failure.severity).toBe("error");
      expect(failure.message).toBe("Request failed");
    });

    it('falls back to transport message for array body', () => {
      const err = new AxiosError(
        "Request failed",
        "ERR",
        undefined,
        undefined,
        { status: 409, statusText: "Conflict", headers: {}, config: {}, data: [{ error: { code: "X", message: "y" } }] } as never,
      );
      const failure = toActionFailure(err, fallbackMessage);
      expect(failure.severity).toBe("error");
      expect(failure.message).toBe("Request failed");
    });

    it('never throws for any input including weird objects', () => {
      expect(() => toActionFailure(Symbol("test") as unknown, fallbackMessage)).not.toThrow();
      expect(() => toActionFailure({ deep: { nested: { value: true } } } as unknown, fallbackMessage)).not.toThrow();
      expect(() => toActionFailure(() => {}, fallbackMessage)).not.toThrow();
    });
  });
});
