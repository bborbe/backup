---
status: completed
spec: [001-conflict-response-for-already-running-triggers]
summary: Created ApiError.ts with extractApiError/toActionFailure, updated BackupStatusOverviewComponent.vue to render server error messages with warning/error severity classes, added ApiError.test.ts and BackupStatusOverviewComponent.test.ts, updated CHANGELOG.md
execution_id: backup-409-json-error-exec-009-spec-001-conflict-dashboard
dark-factory-version: v0.192.9
created: "2026-08-09T09:10:00Z"
queued: "2026-08-09T09:33:00Z"
started: "2026-08-09T09:33:02Z"
completed: "2026-08-09T09:37:15Z"
---

# Show the server's already-running explanation as a warning in the dashboard

<summary>
- The dashboard currently shows "Request failed with status code 500" whenever a per-host trigger fails, no matter what the server actually said.
- After this change it reads the server's own explanation out of the failure response and shows that sentence instead.
- An already-running conflict is presented as a warning ("this host is busy, try later"), visually distinct from a genuine failure.
- Every other trigger failure is presented as a red error, and successful triggers keep exactly the neutral presentation they have today.
- The reading of the server's explanation happens in one shared place rather than being duplicated per button.
- When the failure has no structured explanation — network down, proxy rewrote the body, non-JSON answer — the dashboard falls back to the generic transport message it shows today.
- The fallback never leaves a blank bar, never shows "undefined", and never crashes the page.
- The host name coming back from the server is rendered as plain text, never as markup.
- The bulk "Backup All" / "Cleanup All" buttons keep their current failure behavior unchanged.
- Tests mount the real dashboard component and drive the conflict path, the genuine-error path, the no-structured-body fallback path, and the success path end to end.
</summary>

<objective>
Make the live dashboard render the server-supplied reason for a failed per-host backup or cleanup trigger — with warning severity for an already-running conflict and error severity for every other failure — so an operator can tell "this host is busy" apart from "the backup service is broken" without reading HTTP status codes.
</objective>

<context>
**Dependency: this prompt depends on the preceding prompt `008-spec-001-conflict-api` (same spec, `001-conflict-response-for-already-running-triggers`), which is approved and executes first.** That prompt defines the wire contract this one parses; do not re-derive or change it here. Its end state is: the per-host trigger endpoints `/backup/{name}` and `/cleanup/{name}` answer an already-running trigger with HTTP status 409, `Content-Type: application/json`, and the body

```json
{"error":{"code":"BACKUP_ALREADY_RUNNING","message":"backup for host1.example.com is already running"}}
```

(`CLEANUP_ALREADY_RUNNING` and `cleanup for … is already running` for the cleanup endpoint). The envelope shape comes from `github.com/bborbe/http`: an object with a single `error` member holding `code` (string), `message` (string), and optional `details` (object). Genuine server errors on the same two endpoints use the same envelope with code `INTERNAL_ERROR` and status 500.

This repo has no in-tree `CLAUDE.md`. Read `README.md`, `docs/dod.md`, and `docs/design-system.md` (§ Color Palette → Status Colors) before changing anything.

Read these files first — they are the real, verified structure this prompt builds on:

- `frontend/src/router.ts` — the app has exactly one route, `/` → `frontend/src/pages/DashboardPage.vue`.
- `frontend/src/pages/DashboardPage.vue` — mounts `BackupStatusOverviewComponent` (via `ref="backupOverviewRef"`) and `ActionButtonComponent`. It mounts nothing else. Its "Backup All" / "Cleanup All" buttons call `backupOverviewRef?.triggerBackupAll()` / `triggerCleanupAll()` through `defineExpose`.
- `frontend/src/components/BackupStatusOverviewComponent.vue` — **the only file with live per-host trigger logic, and the sole component this prompt edits.** Verified structure:
  - `const actionState = ref({ isBackingUp: false, isCleaningUp: false, message: null as string | null });` — a single shared message state, not per-host.
  - `const individualActionState = ref<Record<string, { isBackingUp: boolean; isCleaningUp: boolean }>>({});` — per-host button-disable booleans only; it carries no message and no result.
  - Four functions all write the same `actionState.value.message`: `triggerBackupAll()` and `triggerCleanupAll()` (bulk, hit `/backup/all` and `/cleanup/all`) and `triggerBackup(host)` and `triggerCleanup(host)` (per-host, hit `/backup/{name}` and `/cleanup/{name}`). Each of the four sets `actionState.value.message = null;` before its `try`, sets it again in both the success branch and the `catch` branch, and schedules `setTimeout(() => { actionState.value.message = null; }, 5000);` in both branches.
  - Template: a single `<div v-if="actionState.message" class="action-message">{{ actionState.message }}</div>` at the top of `.status-overview`. Its class today is always exactly `action-message`, with no success/failure distinction.
  - Per-host buttons live in `.card-actions` inside each `.status-card`: first the `label="Backup"` button (`@click="triggerBackup(String(host))"`), then the `label="Cleanup"` button (`@click="triggerCleanup(String(host))"`). `ActionButtonComponent` renders a plain `<button>`.
  - Scoped `<style>` already uses the warning tokens through a base-class + modifier-class convention: `.status-pill.status-warning { background-color: var(--status-warning-bg); color: var(--status-warning); border: 1px solid var(--status-warning-border); }`, plus `.status-pill.status-error`, `.detail-value.status-warning`, `.detail-value.status-error`. Follow that same `<base>.<status-*>` convention for the new message-bar severity classes.
  - Pre-existing `console.log` / `console.warn` / `console.error` calls in `loadStatus` are out of scope — do not remove them and do not add new ones.
- `frontend/src/lib/BackupApiClient.ts` — `triggerBackup(name)`, `triggerCleanup(name)`, `triggerBackupAll()`, `triggerCleanupAll()`, `getStatus()`. Each awaits `apiService.post(...)` / `apiService.get(...)` and returns an `ActionResult` (triggers) or `BackupStatus` (`getStatus`) on success; failures propagate as a thrown `AxiosError`. Exported as `export const backupApiClient = new BackupApiClient();` plus `export default backupApiClient;`.
- `frontend/src/lib/ApiService.ts` — thin axios wrapper; `post` returns `response.data` and rethrows axios errors untouched.
- `frontend/src/lib/types.ts` — declares `export interface ActionResult { success: boolean; message: string; }`. **`BackupStatusOverviewComponent.vue` never constructs or imports `ActionResult`** — it only reads `result.message` off the client's return value. `types.ts` is therefore not touched by this prompt (see requirement 1).
- `frontend/src/App.vue` — defines `--status-success*`, `--status-warning*` (`#ff9500` / `#2d1a00` / `#4a2c00`), and `--status-error*` on `:root`. These already exist; add no new colors.
- `frontend/src/tests/` — contains `ActionButtonComponent.test.ts`, `ApiService.test.ts`, `simple.test.ts`, `types.test.ts`. Naming convention is `<Name>.test.ts`; there is no test for `BackupStatusOverviewComponent.vue` yet. `ActionButtonComponent.test.ts` is the reference for the `mount` style, and every existing test file imports its helpers explicitly (`import { describe, it, expect } from 'vitest';`) rather than relying on globals — do the same.

Verified conventions and library facts for this frontend:

- `frontend/tsconfig.json` sets `"allowImportingTsExtensions": true`, `"moduleResolution": "bundler"`, `"strict": false`. Both extension-ful and extensionless specifiers resolve.
- **Import convention for value imports of `.ts` modules from a `.vue` file**: a `// @ts-ignore` comment line directly above an import that carries the explicit `.ts` extension — `BackupStatusOverviewComponent.vue` does exactly this for `import backupApiClient from "../lib/BackupApiClient.ts";`, and `BackupApiClient.ts` does the same for `import { apiService } from "./ApiService.ts";`. Type-only imports use the extensionless form (`import type { BackupStatus, LoadingState } from "../lib/types";`). Follow both forms exactly. Test files under `frontend/src/tests/` use extensionless value imports (`import { ApiService } from '../lib/ApiService';`).
- Tests run with vitest 4 (`environment: 'jsdom'`, `globals: true`) via `make -C frontend test`; `make -C frontend precommit` runs `lint build test`.
- `frontend/eslint.config.mjs` ignores `src/**/*.js`. The checked-in `*.vue.js` and `*.js` sibling files are stale build artifacts — never hand-write or hand-edit them.
- `AxiosError` is a real exported class from `axios` (installed version 1.18.0; constructor signature has been stable across the 1.x line). Verified constructor: `constructor(message?: string, code?: string, config?: InternalAxiosRequestConfig, request?: any, response?: AxiosResponse)`. Its `response` carries `{ status, statusText, headers, config, data }`.
- `flushPromises` is exported from `@vue/test-utils`. Verified implementation: `const scheduler = typeof setImmediate === 'function' ? setImmediate : setTimeout;` — it resolves via a timer. **Do NOT call `vi.useFakeTimers()` in the component test**: with fake timers installed, `flushPromises()` never resolves and the test hangs. Real timers are safe here because the component's 5-second auto-clear cannot fire within an assertion that runs immediately after `flushPromises()`.

Also read `/home/node/.claude/plugins/marketplaces/coding/docs/vue3-typescript-frontend-guide.md`.
</context>

<requirements>
1. Do NOT modify `frontend/src/lib/types.ts`. `ActionResult` stays exactly as it is (`{ success: boolean; message: string; }`), because `BackupStatusOverviewComponent.vue` does not use it and `frontend/src/tests/types.test.ts` asserts its current shape. The new severity type belongs in the new module created in requirement 2.

2. Create `frontend/src/lib/ApiError.ts` — the single shared place where a failing trigger response is turned into a displayable message plus severity. Give every exported symbol a doc comment (`docs/dod.md`). It must export exactly:

   ```ts
   export const BACKUP_ALREADY_RUNNING = "BACKUP_ALREADY_RUNNING";
   export const CLEANUP_ALREADY_RUNNING = "CLEANUP_ALREADY_RUNNING";

   export type ActionSeverity = "warning" | "error";

   export interface ApiErrorEnvelope {
     error: {
       code: string;
       message: string;
       details?: Record<string, unknown>;
     };
   }

   export interface ActionFailure {
     message: string;
     severity: ActionSeverity;
   }

   export function extractApiError(err: unknown): { code: string | null; message: string | null };

   export function toActionFailure(err: unknown, fallbackMessage: string): ActionFailure;
   ```

   The two code strings are frozen contracts spelled exactly as the Go handlers emit them. Do not import axios into this module — narrow `unknown` structurally instead. Avoid `any` (`@typescript-eslint/no-explicit-any` is configured as `'warn'`, so it won't fail `lint:analyse`, but using `any` here would still be sloppy narrowing — don't).

3. `extractApiError` behavior:

   a. Read the candidate body at `err?.response?.data`, reached by structural narrowing of `unknown` (no axios types, no type assertions to `any`).
   b. Treat it as the structured envelope only when it is a non-null, non-array object whose `error` member is a non-null, non-array object whose `message` is a non-empty string. In that case return that `message`, plus `error.code` when `code` is a non-empty string and `null` otherwise.
   c. For everything else — a `text/plain` string body, an array, `null`, a missing `response`, a non-object `err`, `undefined` — return `{ code: null, message: null }`. Do NOT attempt `JSON.parse` on a string body.
   d. The function must never throw, for any input, including `null`, `undefined`, primitives, and objects missing every expected member.

4. `toActionFailure(err, fallbackMessage)` behavior:

   a. Call `extractApiError(err)`.
   b. `message` is the extracted structured message when it is a non-empty string; otherwise the transport-level message `err.message` when `err instanceof Error` and that message is a non-empty string; otherwise `fallbackMessage`. The returned `message` must never be empty, `undefined`, or the literal string `"undefined"`.
   c. `severity` is `"warning"` when the extracted code equals `BACKUP_ALREADY_RUNNING` or `CLEANUP_ALREADY_RUNNING`, and `"error"` in every other case, including when no code was extracted. Key the severity on the error code, never on the raw HTTP status.
   d. It is only ever called from failure paths, so it always returns a concrete `ActionSeverity` — never `null`.

5. In `frontend/src/components/BackupStatusOverviewComponent.vue`, add the imports next to the existing ones at the top of `<script setup lang="ts">`, following the file's own two import forms:

   ```ts
   // @ts-ignore
   import { toActionFailure } from "../lib/ApiError.ts";
   import type { ActionSeverity } from "../lib/ApiError";
   ```

6. In the same file, extend the shared `actionState` ref with a severity member, leaving the existing three members unchanged:

   ```ts
   const actionState = ref({
     isBackingUp: false,
     isCleaningUp: false,
     message: null as string | null,
     severity: null as ActionSeverity | null,
   });
   ```

   `severity === null` means "no severity styling" and is the state a successful trigger leaves behind — that is what keeps the success presentation byte-identical to today's.

7. Reset the shared severity wherever the shared message is already reset at the start of a trigger. Because all four trigger functions write the one `actionState.value.message`, a stale severity from an earlier per-host warning would otherwise bleed into a later bulk-trigger render. Add exactly one line, `actionState.value.severity = null;`, immediately after the existing `actionState.value.message = null;` line at the top of **all four** functions: `triggerBackupAll`, `triggerCleanupAll`, `triggerBackup`, and `triggerCleanup`. This one-line reset is the only change permitted in `triggerBackupAll` and `triggerCleanupAll` — their `try`, success, `catch`, `finally`, and `setTimeout` bodies stay exactly as they are today.

8. In `triggerBackup(host)`, replace the `catch (err)` body so it goes through the shared helper, and clear the severity alongside the message in that branch's `setTimeout`:

   ```ts
   } catch (err) {
     const failure = toActionFailure(err, `Failed to trigger backup for ${host}`);
     actionState.value.message = failure.message;
     actionState.value.severity = failure.severity;
     setTimeout(() => {
       actionState.value.message = null;
       actionState.value.severity = null;
     }, 5000);
   } finally {
   ```

   Leave the success branch, the `individualActionState` handling, and the `finally` block exactly as they are.

9. Apply the identical change in `triggerCleanup(host)`, with fallback message `` `Failed to trigger cleanup for ${host}` ``. Both fallback strings match the strings those catch blocks produce today, so behavior is unchanged when no structured body and no transport message are available.

10. In the same file's `<template>`, give the message bar a severity modifier class, keeping Vue's escaping text interpolation for the message:

    ```html
    <div
      v-if="actionState.message"
      :class="['action-message', actionState.severity ? `status-${actionState.severity}` : '']"
    >
      {{ actionState.message }}
    </div>
    ```

    This yields `action-message status-warning` for a conflict, `action-message status-error` for any other failure, and plain `action-message` for a success — matching the file's existing `.status-pill.status-warning` / `.detail-value.status-error` convention. Do NOT introduce `v-html` anywhere: the message carries a server-echoed host name and must never be able to inject markup or script into the dashboard.

11. In the same file's scoped `<style>`, add two rules directly after the existing `.action-message` rule, reusing only the existing tokens and mirroring the `.status-pill.status-warning` / `.status-pill.status-error` declarations already in this file:

    ```css
    .action-message.status-warning {
      background-color: var(--status-warning-bg);
      color: var(--status-warning);
      border: 1px solid var(--status-warning-border);
    }

    .action-message.status-error {
      background-color: var(--status-error-bg);
      color: var(--status-error);
      border: 1px solid var(--status-error-border);
    }
    ```

    Leave the base `.action-message` rule untouched so the success presentation is unchanged.

12. Create `frontend/src/tests/ApiError.test.ts` (vitest, `import { describe, it, expect } from 'vitest';`, value import `from '../lib/ApiError'`). Build real `AxiosError` instances with `import { AxiosError } from 'axios';` using the verified constructor:

    ```ts
    const err = new AxiosError(
      "Request failed with status code 409",
      "ERR_BAD_REQUEST",
      undefined,
      undefined,
      { status: 409, statusText: "Conflict", headers: {}, config: {}, data } as never,
    );
    ```

    Cover `toActionFailure` for:

    a. 409 + `{ error: { code: "BACKUP_ALREADY_RUNNING", message: "backup for host1.example.com is already running" } }` → `severity === "warning"`; `message` equals the server message; `message` contains `host1.example.com` and `already running`; `message` contains neither `Request failed with status code` nor `500`.
    b. 409 + `{ error: { code: "CLEANUP_ALREADY_RUNNING", message: "cleanup for host1.example.com is already running" } }` → `severity === "warning"`; `message` equals the server message.
    c. 500 + `{ error: { code: "INTERNAL_ERROR", message: "backup test-target failed: disk on fire" } }` → `severity === "error"`; `message` equals the server message.
    d. 500 + a `text/plain` string body such as `"request failed: something broke\n"` → `severity === "error"`; `message` equals `err.message`; `message` is non-empty and does not contain `undefined`.
    e. An `AxiosError` with no `response` at all (network failure) → `severity === "error"`; `message` equals `err.message`.
    f. Non-`Error` inputs `null`, `undefined`, `"boom"`, `42`, and `{}` → each returns `severity === "error"` with `message` equal to the supplied `fallbackMessage`, and no call throws.
    g. Envelope present but `error.message` empty (`{ error: { code: "BACKUP_ALREADY_RUNNING", message: "" } }`) → falls back to the transport message, which is non-empty. Assert the envelope-shape expectations explicitly here — this is a runtime regression guard, not a build-time one (`vite build` does not type-check; there is no `vue-tsc`/`build:check` step and `tsconfig.json` sets `strict: false`), so this test is the only thing that would catch a future envelope-shape drift.
    h. Envelope-adjacent shapes that must NOT be treated as structured: `{ error: "boom" }`, `{ error: null }`, `{ error: [] }`, and an array body `[{ error: { code: "X", message: "y" } }]` → each falls back to the transport message at `"error"` severity.

13. Create `frontend/src/tests/BackupStatusOverviewComponent.test.ts` (vitest + `@vue/test-utils`; `import { describe, it, expect, vi, beforeEach } from 'vitest';`, `import { mount, flushPromises } from '@vue/test-utils';`). This is the test that proves the live render path, so it must mount the real component and let the component's own catch blocks and the real shared helper run — mock only the API client module:

    ```ts
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
    ```

    The specifier `"../lib/BackupApiClient.ts"` resolves from `frontend/src/tests/` to the same module the component imports from `frontend/src/components/`, so the mock applies. Do NOT mock `../lib/ApiError` — the real extractor must run.

    Per test: `vi.clearAllMocks()` in `beforeEach`; stub `getStatus` to resolve `{ "host1.example.com": "2024-01-01" }`; `const wrapper = mount(BackupStatusOverviewComponent);` (the `filters` prop defaults to all-true, so the single host renders); `await flushPromises();`; click the per-host button — `await wrapper.findAll(".card-actions button")[0].trigger("click")` for Backup and `[1]` for Cleanup; `await flushPromises();`; then read `wrapper.get(".action-message")`. Do not install fake timers. Cases:

    a. `triggerBackup` rejects with the 409 `BACKUP_ALREADY_RUNNING` `AxiosError` from case 12a → the `.action-message` text contains `host1.example.com` and `already running`, and contains neither `Request failed with status code` nor `500`; its `classes()` contain `status-warning` and do NOT contain `status-error`.
    b. `triggerCleanup` rejects with the 409 `CLEANUP_ALREADY_RUNNING` `AxiosError` from case 12b, clicking the second `.card-actions` button → text contains `host1.example.com` and `already running`; `classes()` contain `status-warning` and not `status-error`.
    c. `triggerBackup` rejects with the 500 structured `INTERNAL_ERROR` `AxiosError` from case 12c → text equals the server message; `classes()` contain `status-error` and do NOT contain `status-warning`.
    d. `triggerBackup` rejects with the 500 `text/plain` `AxiosError` from case 12d → the rendered text is non-empty, equals the axios transport message, does not contain `undefined`, and `classes()` contain `status-error`. The test must complete without an unhandled rejection or thrown exception.
    e. `triggerBackup` resolves with `{ success: true, message: "backup host1.example.com completed" }` → the `.action-message` text equals that message and its `classes()` contain neither `status-warning` nor `status-error`, proving the success presentation is unchanged.

14. Update `CHANGELOG.md`. Prompt `008-spec-001-conflict-api` creates a `## Unreleased` section directly above `## v3.9.25`. If that section is already present, append a bullet to it — do NOT create a second one. If it is absent, create it above `## v3.9.25`. The bullet states that the dashboard now shows the server-supplied reason for a failed per-host backup or cleanup trigger and renders an already-running conflict as a warning instead of an error.

15. Add a direct unit test in `frontend/src/tests/ApiError.test.ts` for `extractApiError` itself (not only through `toActionFailure`): assert it returns `{ code: "BACKUP_ALREADY_RUNNING", message: "backup for host1.example.com is already running" }` for the structured 409 body from case 12a, and `{ code: null, message: null }` for the `text/plain` body from case 12d. `extractApiError` is exported specifically so it can be tested directly — do not leave it covered only transitively through `toActionFailure`.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git.
- This prompt is frontend-only. Do NOT change any Go code, `main.go`, `pkg/`, `mocks/`, or the API contract.
- `frontend/src/components/TargetListComponent.vue` and `frontend/src/components/ActionPanelComponent.vue` are dead code — referenced by no `.vue`/`.ts` source, only by a stale build artifact. They are entirely out of scope: do not read them for guidance, do not edit them, do not test them.
- `BACKUP_ALREADY_RUNNING` and `CLEANUP_ALREADY_RUNNING` are frozen string contracts and must be matched exactly, character for character.
- Do NOT change the error-handling logic of `triggerBackupAll` or `triggerCleanupAll`. The single `actionState.value.severity = null;` reset line from requirement 7 is the only edit permitted in those two functions.
- Do NOT change error handling for any endpoint other than the per-host backup trigger and the per-host cleanup trigger. `/status`, `/list`, `/backup/all`, and `/cleanup/all` keep their current behavior.
- Do NOT add retry or backoff logic for already-running conflicts — the operator retries manually.
- Do NOT redesign the dashboard's error/toast presentation beyond adding the warning-vs-error severity distinction on the existing `.action-message` bar.
- Do NOT introduce new colors or design tokens. The existing `--status-warning*` and `--status-error*` custom properties from `docs/design-system.md` (defined on `:root` in `frontend/src/App.vue`) are the palette.
- Do NOT add a config flag, env var, or query parameter to switch behavior — the structured body is an invariant of the two per-host endpoints.
- Reflected input: the host name arrives from the server inside the message and is rendered in the message bar. It must go through Vue's escaping `{{ }}` interpolation — never `v-html`.
- The fallback path must not produce a blank bar, the string `undefined`, an unhandled rejection, or a crash, for any input including a missing response body.
- Successful triggers keep their existing behavior: same message text, same neutral `.action-message` styling, same 5-second auto-clear.
- All existing frontend tests under `frontend/src/tests/` continue to pass, including `types.test.ts`.
- Do NOT hand-edit the checked-in `*.vue.js` / `*.js` build artifacts under `frontend/src/`; they are generated and eslint-ignored.
- `docs/dod.md` applies: a `## Unreleased` CHANGELOG entry, doc comments on exported symbols, no debug output. The existing `console.log` / `console.warn` / `console.error` calls in `loadStatus` are pre-existing and out of scope — leave them, and add no new ones.
- `docs/dod.md`'s Go-specific criteria (Ginkgo/Gomega, Counterfeiter mocks, `github.com/bborbe/errors` wrapping, factory purity) do NOT apply to this frontend-only prompt — the frontend convention is vitest + `@vue/test-utils`, already followed by every other test file in `frontend/src/tests/`. Do not attempt to satisfy Go-specific DoD items here.
- `README.md` needs no update: this change alters neither usage, configuration, nor setup — do not add a section about the response contract or error handling to it.
</constraints>

<verification>
Run `make -C frontend test` — must pass (exit code 0), including the pre-existing suites and both new test files.

Run `grep -n 'action-message.status-warning' frontend/src/components/BackupStatusOverviewComponent.vue` — must return exactly one line, proving the new severity rule was added (the base `status-warning` token already appears elsewhere in this file pre-change, so a bare `status-warning` grep is not a valid signal).

Run `grep -n 'action-message.status-error' frontend/src/components/BackupStatusOverviewComponent.vue` — must return exactly one line.

Run `grep -n 'ALREADY_RUNNING' frontend/src/lib/ApiError.ts` — must return at least two lines, one per frozen error code.

Run `grep -rn 'v-html' frontend/src/components/ frontend/src/pages/ frontend/src/lib/` — must print nothing.

Run `make precommit` — must pass (this already runs `frontend-precommit`, which is lint + build + test for the whole frontend, plus the repo's Go generate/vulncheck/osv-scanner/trivy checks; no need to run `make -C frontend precommit` separately first).
</verification>
