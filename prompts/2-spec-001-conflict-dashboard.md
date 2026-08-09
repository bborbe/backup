---
spec: ["001-conflict-response-for-already-running-triggers"]
status: draft
created: "2026-08-09T09:10:00Z"
---

# Show the server's already-running explanation as a warning in the dashboard

<summary>
- The dashboard currently shows "Request failed with status code 500" whenever a trigger fails, no matter what the server actually said.
- After this change it reads the server's own explanation out of the failure response and shows that sentence instead.
- An already-running conflict is presented as a warning ("this host is busy, try later"), visually distinct from a genuine failure.
- Every other trigger failure keeps the red error presentation it has today.
- The reading of the server's explanation happens in one shared place, so all four trigger buttons behave identically.
- When the failure has no structured explanation — network down, proxy rewrote the body, non-JSON answer — the dashboard falls back to the generic transport message it shows today.
- The fallback never leaves a blank bar, never shows "undefined", and never crashes the page.
- Successful triggers are completely unchanged.
- The host name coming back from the server is rendered as plain text, never as markup.
- Tests cover the conflict path, the genuine-error path, and the no-structured-body fallback path.
</summary>

<objective>
Make the dashboard render the server-supplied reason for a failed backup or cleanup trigger — with warning severity for an already-running conflict and error severity for everything else — so an operator can tell "this host is busy" apart from "the backup service is broken" without reading HTTP status codes.
</objective>

<context>
**This prompt depends on the preceding prompt `1-spec-001-conflict-api` (spec `001-conflict-response-for-already-running-triggers`), which must be executed and completed first.** That prompt defines the wire contract this one parses: the per-host trigger endpoints `/backup/{name}` and `/cleanup/{name}` answer an already-running trigger with HTTP 409 and the JSON body

```json
{"error":{"code":"BACKUP_ALREADY_RUNNING","message":"backup for host1.example.com is already running"}}
```

(`CLEANUP_ALREADY_RUNNING` for the cleanup endpoint). The envelope shape comes from `github.com/bborbe/http`: an object with a single `error` member holding `code` (string), `message` (string), and optional `details` (object). Do not re-derive or change that contract here.

This repo has no in-tree `CLAUDE.md`. Read `README.md`, `docs/dod.md`, and `docs/design-system.md` (§ Color Palette → Status Colors) before changing anything.

Read these files first:

- `frontend/src/lib/types.ts` — `export interface ActionResult { success: boolean; message: string; }`.
- `frontend/src/lib/BackupApiClient.ts` — `triggerBackup(name)`, `triggerCleanup(name)`, `triggerBackupAll()`, `triggerCleanupAll()`. Each awaits `apiService.post(...)` and returns an `ActionResult` on success; failures propagate as a thrown `AxiosError`.
- `frontend/src/lib/ApiService.ts` — thin axios wrapper; `post` returns `response.data` and rethrows axios errors untouched.
- `frontend/src/components/TargetListComponent.vue` — the per-host trigger UI. `triggerBackup(host)` / `triggerCleanup(host)` each have a `catch (err)` block that currently does `message: err instanceof Error ? err.message : "Backup failed"`. The result is rendered in the `div.action-result` whose class list is `['action-result', ...success ? 'success' : 'error']`.
- `frontend/src/components/ActionPanelComponent.vue` — the bulk (all-targets) trigger UI with the same two catch blocks and the same `div.action-result` class-list pattern.
- `frontend/src/App.vue` — defines the `--status-success*`, `--status-warning*`, and `--status-error*` custom properties on `:root`. These already exist; do not add new colors.
- `frontend/src/components/BackupStatusOverviewComponent.vue` — existing consumer of the `status-warning` token (per-host status pill path). Reference only; do not modify.
- `frontend/src/tests/ActionButtonComponent.test.ts` — reference for the project's vitest + `@vue/test-utils` `mount` style (explicit `import { describe, it, expect } from 'vitest'`).

Notes on this frontend:

- Tests run with vitest (`environment: 'jsdom'`, `globals: true`) via `make -C frontend test`; `make -C frontend precommit` runs `lint build test`.
- `frontend/eslint.config.mjs` ignores `src/**/*.js`. The checked-in `*.vue.js` and `*.ts` → `*.js` sibling files are build artifacts — do NOT hand-write or hand-edit them.
- Imports of TypeScript modules from `.vue` files use an extensionless specifier when they are type-only (`from "../lib/types"`). Use extensionless specifiers for the new module too.
- `AxiosError` is a real exported class: `new AxiosError(message?, code?, config?, request?, response?)`, with `response` carrying `{ status, statusText, headers, config, data }`.
- `flushPromises` is exported from `@vue/test-utils`.

Also read `/home/node/.claude/plugins/marketplaces/coding/docs/vue3-typescript-frontend-guide.md`.
</context>

<requirements>
1. In `frontend/src/lib/types.ts`, add an exported severity type and extend `ActionResult` with an optional severity, keeping the existing members unchanged so all current success-path constructions still type-check:

   ```ts
   export type ActionSeverity = "warning" | "error";

   export interface ActionResult {
     success: boolean;
     message: string;
     severity?: ActionSeverity;
   }
   ```

2. Create `frontend/src/lib/ApiError.ts` — the single shared place where a failing trigger response is turned into a displayable result. It must export:

   - `export const BACKUP_ALREADY_RUNNING = "BACKUP_ALREADY_RUNNING";` and `export const CLEANUP_ALREADY_RUNNING = "CLEANUP_ALREADY_RUNNING";` — frozen contracts, spelled exactly as the API emits them.
   - `export interface ApiErrorEnvelope { error: { code: string; message: string; details?: Record<string, unknown> } }` — the structured body shape produced by `github.com/bborbe/http`.
   - `export function extractApiError(err: unknown): { code: string | null; message: string | null }` — pulls the structured code and message out of an error's response body.
   - `export function toActionResult(err: unknown, fallbackMessage: string): ActionResult` — the function the components call from their catch blocks.

   Import `ActionResult` and `ActionSeverity` from `./types` (do not redeclare them; `ApiError.ts` imports from `types.ts`, never the other way round). Give every exported symbol a doc comment.

3. Specify `extractApiError` behavior precisely:

   a. Read the candidate body structurally, without importing axios types: the value at `err?.response?.data`.
   b. Treat it as the structured envelope only when it is a non-null object whose `error` member is a non-null object with a `message` that is a non-empty string. In that case return that `message`, and return `error.code` when it is a non-empty string, otherwise `null` for the code.
   c. For anything else — a string body (`text/plain`), an array, `null`, a missing `response`, a non-object `err`, `undefined` — return `{ code: null, message: null }`. Do NOT attempt `JSON.parse` on a string body.
   d. The function must never throw for any input, including `null`, `undefined`, primitives, and objects with getters that are absent.

4. Specify `toActionResult(err, fallbackMessage)` behavior precisely:

   a. Call `extractApiError(err)`.
   b. `message` is the structured message when present and non-empty; otherwise the transport-level message `err instanceof Error ? err.message : ""`; otherwise `fallbackMessage`. The result's `message` must never be empty, `undefined`, or the literal string `"undefined"`.
   c. `severity` is `"warning"` when the extracted code equals `BACKUP_ALREADY_RUNNING` or `CLEANUP_ALREADY_RUNNING`; `"error"` in every other case, including when there is no code at all. Key the severity on the error code, never on the raw HTTP status.
   d. Always returns `{ success: false, message, severity }` — it is only called from failure paths.

5. In `frontend/src/components/TargetListComponent.vue`, replace the body of both `catch (err)` blocks (in `triggerBackup` and `triggerCleanup`) with a single call to the shared helper, keeping the surrounding assignment to `actionStates.value[host].result`:

   ```ts
   } catch (err) {
     actionStates.value[host].result = toActionResult(err, "Backup failed");
   } finally {
   ```

   (`"Cleanup failed"` for `triggerCleanup`.) Import with `import { toActionResult } from "../lib/ApiError";`. Do not change the success paths, the `setTimeout` clearing, the loading-state handling, or `loadTargets`.

6. In `frontend/src/components/TargetListComponent.vue`, extend the result element's class binding so a warning result carries the `warning` class and not the `error` class:

   ```
   :class="[
     'action-result',
     actionStates[target.host]?.result?.success
       ? 'success'
       : actionStates[target.host]?.result?.severity === 'warning'
         ? 'warning'
         : 'error'
   ]"
   ```

   Add a `.action-result.warning` rule to the component's scoped `<style>` alongside the existing `.success` / `.error` rules, using only the existing tokens: `background-color: var(--status-warning-bg); color: var(--status-warning); border: 1px solid var(--status-warning-border);`.

7. Apply the same two changes to `frontend/src/components/ActionPanelComponent.vue`: both catch blocks call `toActionResult(err, "Backup failed")` / `toActionResult(err, "Cleanup failed")`, the two result elements gain the same three-way class binding based on `backupResult.severity` / `cleanupResult.severity`, and the scoped style gains the same `.action-result.warning` rule built from the `--status-warning*` tokens. The bulk endpoints still return plain-text failures, so this path simply falls back to the transport message at error severity — that is the intended behavior, not a bug.

8. Keep the message rendered through Vue's escaping text interpolation (`{{ ... }}`). Do NOT introduce `v-html` anywhere: the message contains a server-echoed host name and must never be able to inject markup or script into the dashboard.

9. Create `frontend/src/tests/ApiError.test.ts` (vitest, explicit `import { describe, it, expect } from 'vitest'`) covering `toActionResult` against real `AxiosError` instances built with `import { AxiosError } from "axios";`. Build them like:

   ```ts
   const err = new AxiosError(
     "Request failed with status code 409",
     "ERR_BAD_REQUEST",
     undefined,
     undefined,
     { status: 409, statusText: "Conflict", headers: {}, config: {}, data } as never,
   );
   ```

   Cases:

   a. 409 + `{ error: { code: "BACKUP_ALREADY_RUNNING", message: "backup for host1.example.com is already running" } }` → `severity === "warning"`, message equals the server message, message contains `host1.example.com` and `already running`, message contains neither `Request failed with status code` nor `500`.
   b. 409 + `{ error: { code: "CLEANUP_ALREADY_RUNNING", message: "cleanup for host1.example.com is already running" } }` → `severity === "warning"`, message equals the server message.
   c. 500 + `{ error: { code: "INTERNAL_ERROR", message: "backup test-target failed: disk on fire" } }` → `severity === "error"`, message equals the server message.
   d. 500 + a `text/plain` string body such as `"request failed: something broke\n"` → `severity === "error"`, message equals the axios transport message (`err.message`), message is non-empty and does not contain `undefined`.
   e. Network failure: an `AxiosError` with no `response` at all → `severity === "error"`, message equals the axios transport message.
   f. Non-Error inputs `null`, `undefined`, `"boom"`, and `{}` → each returns `severity === "error"` with `message` equal to the supplied `fallbackMessage`, and the call does not throw.
   g. Envelope present but `error.message` empty (`{ error: { code: "BACKUP_ALREADY_RUNNING", message: "" } }`) → falls back to the transport message and is non-empty. Assert the envelope-shape expectations explicitly so a future change to the library's envelope fails here at build time.

10. Create `frontend/src/tests/TargetListComponent.test.ts` (vitest + `@vue/test-utils`). Mock the API client module so the component's own catch blocks and the real shared extractor run:

    ```ts
    vi.mock("../lib/BackupApiClient.ts", () => {
      const client = { listTargets: vi.fn(), triggerBackup: vi.fn(), triggerCleanup: vi.fn() };
      return { default: client, backupApiClient: client, BackupApiClient: class {} };
    });
    ```

    In each test stub `listTargets` to resolve to a single target (`{ host: "host1.example.com", port: 22, user: "root", dirs: ["/data"], excludes: [] }`), `mount(TargetListComponent)`, `await flushPromises()`, click the first button inside `.target-actions`, `await flushPromises()`, then assert on the `.action-result` element. Cases:

    a. `triggerBackup` rejects with the 409 structured `AxiosError` from case 9a → the rendered `.action-result` text contains `host1.example.com` and `already running`, and contains neither `Request failed with status code` nor `500`; its classes contain `warning` and do NOT contain `error`.
    b. `triggerBackup` rejects with the 500 structured `AxiosError` from case 9c → the element's classes contain `error` and do NOT contain `warning`.
    c. `triggerBackup` rejects with the 500 `text/plain` `AxiosError` from case 9d → the rendered text is non-empty, equals the axios transport message, does not contain `undefined`, and the element's classes contain `error`. The test must complete without an unhandled rejection or thrown exception.

    Use `beforeEach` to reset the mocks (`vi.clearAllMocks()`); import `describe, it, expect, vi, beforeEach` explicitly from `vitest`.

11. Add an entry under the existing `## Unreleased` section in `CHANGELOG.md` describing that the dashboard now shows the server-supplied reason for a failed backup/cleanup trigger and renders an already-running conflict as a warning instead of an error.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git.
- Do NOT change any Go code, `main.go`, `pkg/`, or the API contract. This prompt is frontend-only.
- The error codes `BACKUP_ALREADY_RUNNING` and `CLEANUP_ALREADY_RUNNING` are frozen string contracts and must be matched exactly.
- Do NOT add retry or backoff logic for already-running conflicts — the operator retries manually.
- Do NOT redesign the dashboard's error/toast presentation beyond adding the warning-vs-error severity distinction.
- Do NOT introduce new colors or design tokens; the existing `--status-warning*` and `--status-error*` custom properties from `docs/design-system.md` (defined on `:root` in `frontend/src/App.vue`) are the palette.
- Do NOT add a config flag, env var, or query parameter to switch behavior — the structured body is an invariant of the two per-host endpoints.
- Reflected input: the host name arrives from the server and is rendered in the message bar. It must go through Vue's escaping interpolation — never `v-html`.
- The fallback path must not produce a blank bar, the string `undefined`, an unhandled rejection, or a crash, for any input including a missing response body.
- Successful triggers keep their existing behavior: same `ActionResult`, same `success` class, same auto-clear timeout.
- All existing frontend tests under `frontend/src/tests/` continue to pass.
- Do NOT hand-edit the checked-in `*.vue.js` / `*.js` build artifacts under `frontend/src/`; they are generated and eslint-ignored.
- `docs/dod.md` applies: a `## Unreleased` CHANGELOG entry, doc comments on exported symbols, no debug output.
</constraints>

<verification>
Run `make -C frontend test` — must pass (exit code 0), including the pre-existing suites.

Run `grep -n 'status-warning\|warning' frontend/src/components/ActionPanelComponent.vue` — must return at least one line.

Run `grep -n 'status-warning\|warning' frontend/src/components/TargetListComponent.vue` — must return at least one line.

Run `grep -rn 'v-html' frontend/src/` — must return no lines.

Run `make -C frontend precommit` — must pass (lint, build, test).

Run `make precommit` — must pass.
</verification>
