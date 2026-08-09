/**
 * API Error parsing utilities for structured error responses.
 *
 * Provides extraction of error details from failed API responses and
 * conversion to displayable action failure information with severity levels.
 */

/**
 * Error code returned when a backup trigger is rejected because
 * a backup is already running for the target host.
 */
export const BACKUP_ALREADY_RUNNING = "BACKUP_ALREADY_RUNNING";

/**
 * Error code returned when a cleanup trigger is rejected because
 * a cleanup is already running for the target host.
 */
export const CLEANUP_ALREADY_RUNNING = "CLEANUP_ALREADY_RUNNING";

/**
 * Severity level for an action failure displayed in the dashboard.
 * "warning" is used for already-running conflicts; "error" for all other failures.
 */
export type ActionSeverity = "warning" | "error";

/**
 * The structured error envelope returned by the backup service
 * on per-host trigger endpoints.
 */
export interface ApiErrorEnvelope {
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}

/**
 * A displayable action failure with message text and severity.
 */
export interface ActionFailure {
  message: string;
  severity: ActionSeverity;
}

/**
 * Extracts the structured error code and message from an unknown error value.
 *
 * This function narrows the shape of an unknown error by checking for
 * an AxiosError-like structure with a response.data body. It only treats
 * a body as the structured envelope when it is a non-null, non-array object
 * with a non-null, non-array `error` member whose `message` is a non-empty string.
 *
 * @param err - Any value, typically an AxiosError thrown by ApiService.
 * @returns An object with the extracted code (non-empty string or null) and
 *          message (non-empty string or null). Returns nulls for any input
 *          that does not match the expected envelope shape.
 */
export function extractApiError(err: unknown): { code: string | null; message: string | null } {
  // Attempt to access the response body via structural narrowing
  const responseData = (err as { response?: { data?: unknown } })?.response?.data;

  // Only treat responseData as the envelope when it is a plain object
  if (
    responseData !== null &&
    typeof responseData === "object" &&
    !Array.isArray(responseData)
  ) {
    const envelope = responseData as Record<string, unknown>;
    const errorField = envelope.error;

    // error must also be a plain non-null object
    if (
      errorField !== null &&
      typeof errorField === "object" &&
      !Array.isArray(errorField)
    ) {
      const errorObj = errorField as Record<string, unknown>;
      const message = errorObj.message;

      // message must be a non-empty string to be considered structured
      if (typeof message === "string" && message.length > 0) {
        const code = errorObj.code;
        return {
          code: typeof code === "string" && code.length > 0 ? code : null,
          message: message,
        };
      }
    }
  }

  // Fall through for: null, undefined, primitives, arrays, missing response,
  // non-object err, object without error member, or error.message not a
  // non-empty string.
  return { code: null, message: null };
}

/**
 * Converts an unknown error from a trigger action into a displayable ActionFailure.
 *
 * When the error carries the structured API error envelope, its message is preferred
 * over the transport-level message. Severity is "warning" only for the explicit
 * already-running codes; every other case (including no code) is "error".
 *
 * @param err - The error caught in a trigger function, typically an AxiosError.
 * @param fallbackMessage - The message to use when neither a structured body
 *                          nor a transport message is available.
 * @returns An ActionFailure with a non-empty, non-undefined message and a
 *          concrete severity value.
 */
export function toActionFailure(err: unknown, fallbackMessage: string): ActionFailure {
  const extracted = extractApiError(err);

  let message: string;

  if (extracted.message && extracted.message.length > 0) {
    message = extracted.message;
  } else if (err instanceof Error && err.message && err.message.length > 0) {
    message = err.message;
  } else {
    message = fallbackMessage;
  }

  const isAlreadyRunningConflict =
    extracted.code === BACKUP_ALREADY_RUNNING ||
    extracted.code === CLEANUP_ALREADY_RUNNING;

  return {
    message,
    severity: isAlreadyRunningConflict ? "warning" : "error",
  };
}
