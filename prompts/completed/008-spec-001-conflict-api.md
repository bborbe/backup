---
status: completed
spec: [001-conflict-response-for-already-running-triggers]
summary: Per-host backup and cleanup trigger endpoints now return HTTP 409 with structured JSON error bodies (BACKUP_ALREADY_RUNNING/CLEANUP_ALREADY_RUNNING) when a run is already in flight
execution_id: backup-409-json-error-exec-008-spec-001-conflict-api
dark-factory-version: v0.192.9
created: "2026-08-09T09:10:00Z"
queued: "2026-08-09T09:12:56Z"
started: "2026-08-09T09:13:29Z"
completed: "2026-08-09T09:18:31Z"
---

# Answer already-running per-host triggers with HTTP 409 and a structured JSON error

<summary>
- Asking the service to back up a host whose backup is already running is a normal situation, but today it is reported as an internal server error.
- After this change the service answers such a request with "conflict — the resource is busy" instead of "the server broke".
- The same treatment applies to a cleanup trigger for a host whose cleanup is already running.
- The refusal now carries a stable, machine-readable reason code so any client can react to it without parsing English text.
- It also carries a human-readable sentence that names the affected host and states that the work is already running.
- The two per-host trigger endpoints now answer failures as structured JSON instead of plain text.
- Successful triggers are completely unchanged — same status, same body.
- Genuine failures on those endpoints still report a server error, only now in the structured shape.
- The mechanism that decides whether a run is already in progress is untouched — only how its refusal is reported changes.
- Automated tests drive both endpoints through their real HTTP wiring and prove the status, content type, code, and message.
</summary>

<objective>
Translate the already-running sentinel errors into an HTTP conflict response with a stable machine-readable error code on the per-host backup and cleanup trigger endpoints, so operators and monitoring can tell "this host is busy" apart from "the backup service is broken". End state: both endpoints answer an already-running trigger with status 409 and a JSON body `{"error":{"code":"BACKUP_ALREADY_RUNNING"|"CLEANUP_ALREADY_RUNNING","message":"…"}}`.
</objective>

<context>
This repo has no in-tree `CLAUDE.md`. Read `README.md` and `docs/dod.md` for project conventions before changing anything.

Read these files first:

- `pkg/handler/handler_backup.go` — `NewBackupHandler(targetFinder pkg.TargetFinder, backupExectuor pkg.BackupExectuor) libhttp.WithError`. The body already wraps executor failures with `errors.Wrapf(ctx, err, "backup %s failed", target.Name)`.
- `pkg/handler/handler_cleanup.go` — `NewCleanupHandler(targetFinder pkg.TargetFinder, cleanupExectuor pkg.BackupCleaner) libhttp.WithError`.
- `pkg/backup-executor-only-once.go` — declares `var BackupAlreadyRunningError = stderrors.New("backup already running")` and returns it from `(*backupExectuorOnlyOnce).Backup` when a run is in flight.
- `pkg/backup-cleaner-only-once.go` — declares `var CleanupAlreadyRunningError = stderrors.New("cleanup already running")` and returns it from `(*backupCleanerOnlyOnce).Clean`.
- `pkg/factory/factory.go` — `CreateBackupHandler` and `CreateCleanupHandler` currently wrap their handler in `libhttp.NewErrorHandler(...)`. `CreateListHandler` and `CreateStatusHandler` also use `libhttp.NewErrorHandler`; `CreateBackupActionHandler` and `CreateCleanActionHandler` use `libhttp.NewBackgroundRunHandler`. All four are OUT OF SCOPE.
- `main.go`, `createHttpServer` — routes `/backup/{name}` to `factory.CreateBackupHandler(...)` and `/cleanup/{name}` to `factory.CreateCleanupHandler(...)`. This is the only entry point that mounts these two handlers; `cmd/backup-rsync` does not use them, so no other call site needs updating.
- `pkg/handler/handler_suite_test.go` — the Ginkgo v2 suite bootstrap for package `handler_test`. There are currently no tests in `pkg/handler/`; the first test files created by this prompt join this suite.
- `pkg/action-backup_test.go` — reference for the project's Ginkgo v2 / Gomega / Counterfeiter test style and for constructing a `v1.Target` with `metav1.ObjectMeta{Name: …}` plus a `v1.BackupSpec{Host: …, Port: …, User: …, Dirs: …, Excludes: …}`.
- `mocks/target-finder.go` (`mocks.TargetFinder`, stub setter `TargetReturns(*v1.Target, error)`), `mocks/backup-executor.go` (`mocks.BackupExecutor`, `BackupReturns(error)`), `mocks/backup-cleaner.go` (`mocks.BackupCleaner`, `CleanReturns(error)`).

Library API of `github.com/bborbe/http` v1.26.16 (already a direct dependency in `go.mod`; no new dependency is needed) — these are the exact signatures, do not invent others:

```go
// wraps an error with both an error code and an HTTP status code
func WrapWithCode(err error, code string, statusCode int) error

// converts a returned error into a JSON body and status code
func NewJSONErrorHandler(withError WithError) http.Handler

type ErrorResponse struct {
    Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
    Code    string         `json:"code"`
    Message string         `json:"message"`
    Details map[string]any `json:"details,omitempty"`
}
```

`NewJSONErrorHandler` reads the status code from an error implementing `StatusCode() int`, the code from an error implementing `Code() string` (defaulting to `INTERNAL_ERROR` / 500), sets `Content-Type: application/json`, and encodes `ErrorResponse` with `Message: err.Error()`.

Error helpers from `github.com/bborbe/errors` (v1.5.15, imported as `errors` in both handler files today):

```go
func Errorf(ctx context.Context, format string, args ...interface{}) error
func Wrapf(ctx context.Context, err error, format string, args ...interface{}) error
func Is(err, target error) bool
```

Also read `/home/node/.claude/plugins/marketplaces/coding/docs/go-json-error-handler-guide.md` and `/home/node/.claude/plugins/marketplaces/coding/docs/go-testing-guide.md`. Note: that guide's "Good" example calls `WrapWithDetails(err, http.StatusBadRequest, libhttp.ErrorCodeValidation, …)` with code and status transposed relative to the real signature. The signatures pinned above in this `<context>` block are authoritative — use them, not the guide's example ordering.
</context>

<requirements>
1. In `pkg/handler/handler_backup.go`, declare an exported constant with a doc comment:

   ```go
   // ErrorCodeBackupAlreadyRunning is the stable, machine-readable error code returned when a
   // backup trigger is refused because a backup for the same host is already in flight.
   const ErrorCodeBackupAlreadyRunning = "BACKUP_ALREADY_RUNNING"
   ```

2. In `pkg/handler/handler_cleanup.go`, declare the matching constant with a doc comment:

   ```go
   // ErrorCodeCleanupAlreadyRunning is the stable, machine-readable error code returned when a
   // cleanup trigger is refused because a cleanup for the same host is already in flight.
   const ErrorCodeCleanupAlreadyRunning = "CLEANUP_ALREADY_RUNNING"
   ```

   Both literal strings are frozen contracts — spell them in SCREAMING_SNAKE_CASE exactly as above.

3. In `NewBackupHandler` (`pkg/handler/handler_backup.go`), change only the executor-failure branch. When `backupExectuor.Backup(ctx, target.Spec)` returns an error, first test it with `errors.Is(err, pkg.BackupAlreadyRunningError)`. On a match, return `libhttp.WrapWithCode(...)` around a freshly built message error, using `http.StatusConflict`:

   ```go
   if err := backupExectuor.Backup(ctx, target.Spec); err != nil {
       if errors.Is(err, pkg.BackupAlreadyRunningError) {
           return libhttp.WrapWithCode(
               errors.Errorf(ctx, "backup for %s is already running", target.Spec.Host),
               ErrorCodeBackupAlreadyRunning,
               http.StatusConflict,
           )
       }
       return errors.Wrapf(ctx, err, "backup %s failed", target.Name)
   }
   ```

   Build a fresh message with `errors.Errorf` rather than wrapping the sentinel, so the emitted message is exactly the host plus the already-running statement and carries no internal chain text. The non-matching branch keeps its current `errors.Wrapf` behavior verbatim.

4. Apply the identical change in `NewCleanupHandler` (`pkg/handler/handler_cleanup.go`): match `errors.Is(err, pkg.CleanupAlreadyRunningError)`, return `libhttp.WrapWithCode(errors.Errorf(ctx, "cleanup for %s is already running", target.Spec.Host), ErrorCodeCleanupAlreadyRunning, http.StatusConflict)`, and leave the generic branch's `errors.Wrapf(ctx, err, "cleanup %s failed", target.Name)` unchanged.

5. Leave the `targetFinder.Target` failure branch in both handlers exactly as it is — it still returns the raw error and therefore still produces a 500. Do NOT map an unknown target to 404; that is explicitly out of scope for this spec.

6. Leave both success paths untouched: `libhttp.WriteAndGlog(resp, "backup %s completed", target.Name)` and `libhttp.WriteAndGlog(resp, "cleanup %s completed", target.Name)`, returning `nil`.

7. In `pkg/factory/factory.go`, change `CreateBackupHandler` and `CreateCleanupHandler` to wrap their handler in `libhttp.NewJSONErrorHandler(...)` instead of `libhttp.NewErrorHandler(...)`. Change nothing else in either function. Leave `CreateListHandler` and `CreateStatusHandler` on `libhttp.NewErrorHandler`, and leave `CreateBackupActionHandler` / `CreateCleanActionHandler` on `libhttp.NewBackgroundRunHandler`. This step is mandatory, not optional: the handler tests in requirements 8–9 construct `libhttp.NewJSONErrorHandler(...)` directly around the handler under test, so they pass even if this factory wiring change is skipped — the two `<verification>` greps below are the only guard that production code actually serves JSON instead of text/plain, so do not treat them as a formality.

8. Create `pkg/handler/handler_backup_test.go` in package `handler_test`, Ginkgo v2 + Gomega, following the style of `pkg/action-backup_test.go`. Import aliases: import the stdlib `errors` package plainly as `errors` (matching `pkg/action-backup_test.go:9`), and import `github.com/bborbe/errors` as `bberrors` (matching `pkg/backup-executor-only-once_test.go:13` / `pkg/backup-cleaner-only-once_test.go:13` / `pkg/target-finder_test.go:11`) — do NOT use `stderrors` as the alias, it does not match this repo's existing test files. Drive the handler through its real HTTP wiring rather than calling the closure directly: build a `mux.NewRouter()` (import `github.com/gorilla/mux`), register `router.Path("/backup/{name}").Handler(libhttp.NewJSONErrorHandler(handler.NewBackupHandler(mockTargetFinder, mockBackupExecutor)))`, then serve a `httptest.NewRequest(http.MethodPost, "/backup/<name>", nil)` into a `httptest.NewRecorder()`. Stub `mockTargetFinder.TargetReturns(target, nil)` with a target whose `ObjectMeta.Name` and `Spec.Host` are distinguishable (for example name `test-target` and host `host1.example.com`). Cover these cases:

   a. `mockBackupExecutor.BackupReturns(pkg.BackupAlreadyRunningError)` — assert `recorder.Code` is exactly `http.StatusConflict` (409); assert the `Content-Type` response header contains `application/json`; unmarshal the body into `libhttp.ErrorResponse` and assert `Error.Code` equals exactly `"BACKUP_ALREADY_RUNNING"`; assert `Error.Message` is non-empty, contains the target's host string, and contains the substring `already running`; assert `Error.Details` is empty (nil or zero-length map) — this locks the information-disclosure constraint so a future context-data call cannot silently leak internals into the 409 body.
   b. `mockBackupExecutor.BackupReturns(bberrors.Wrapf(ctx, pkg.BackupAlreadyRunningError, "wrapped"))` (or any error that wraps the sentinel) — same assertions as (a), proving the detection survives an error chain.
   c. `mockBackupExecutor.BackupReturns(errors.New("disk on fire"))` — assert status is exactly `http.StatusInternalServerError` (500) and the decoded `Error.Code` equals `"INTERNAL_ERROR"`.
   d. `mockTargetFinder.TargetReturns(nil, pkg.TargetNotFoundError)` — assert status is exactly `http.StatusInternalServerError` (500); the executor must not have been called (`mockBackupExecutor.BackupCallCount()` is 0).
   e. Success: `mockBackupExecutor.BackupReturns(nil)` — assert status is exactly `http.StatusOK` (200) and the body contains `backup test-target completed`, proving the success response is unchanged.

9. Create `pkg/handler/handler_cleanup_test.go` mirroring requirement 8 for `handler.NewCleanupHandler` on route `/cleanup/{name}`, with `mocks.BackupCleaner` (`CleanReturns`), the sentinel `pkg.CleanupAlreadyRunningError`, the expected code `"CLEANUP_ALREADY_RUNNING"`, and success body `cleanup test-target completed`. Use the same import aliases as requirement 8 (`errors` for stdlib, `bberrors` for `github.com/bborbe/errors`). Include the same five cases (a–e), including the `Error.Details` empty assertion in case (a)/(b), using `mockBackupCleaner.CleanCallCount()` for case (d).

10. Add a `## Unreleased` section at the top of the version list in `CHANGELOG.md` (directly above `## v3.9.25`; the section does not exist yet) with an entry describing that the per-host backup and cleanup trigger endpoints now answer an already-running trigger with HTTP 409 and a structured JSON error body carrying `BACKUP_ALREADY_RUNNING` / `CLEANUP_ALREADY_RUNNING`.

11. `README.md:242-253` documents `/backup/{name}` and `/cleanup/{name}` in an "API Endpoints" table but describes no response/error contracts for any endpoint today (only method + one-line description). This change does not add or remove an endpoint and does not change the documented request shape, so no README edit is required — do NOT add a response-contract subsection there; leave `README.md` untouched by this prompt.

12. Do not change `pkg/backup-executor-only-once.go` or `pkg/backup-cleaner-only-once.go`. The lock decides the winner; the handler only reports the loser's outcome. In particular, do not re-check whether a run is still in flight before writing the response — if the run finishes between the lock check and the response, the in-flight request still correctly returns 409.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git.
- The HTTP method and paths of the per-host backup trigger (`/backup/{name}`) and the per-host cleanup trigger (`/cleanup/{name}`) are frozen — only the failure status code and body shape change.
- Success responses for both endpoints — status code and body — are unchanged.
- `BACKUP_ALREADY_RUNNING` and `CLEANUP_ALREADY_RUNNING` are frozen string contracts; the dashboard and any future client match on them exactly.
- The JSON envelope shape is dictated by `github.com/bborbe/http` v1.26.16 (`{"error":{"code":…,"message":…,"details":…}}`). Do NOT hand-roll a competing envelope; use `libhttp.WrapWithCode` and `libhttp.NewJSONErrorHandler`.
- Do NOT change error handling for any endpoint other than the per-host backup trigger and the per-host cleanup trigger. `/status`, `/list`, `/backup/all`, and `/cleanup/all` keep their current error format.
- Do NOT change the only-once locking mechanism that decides whether a backup or cleanup is already running. Only the reporting of its refusal changes.
- Do NOT add a config flag, env var, or query parameter to switch between the plain-text and the structured error body — the structured body is an invariant of these two endpoints.
- The conflict path must add no new allocation-heavy work or per-request state: an operator hammering the trigger button costs one cheap HTTP round-trip per rejected request and nothing more.
- Information disclosure: the emitted message must contain only the host name and the already-running statement — no filesystem paths, no credentials, no internal stack or chain detail.
- All existing Go tests continue to pass. Do NOT touch `frontend/` in this prompt.
- `docs/dod.md` applies: Go errors wrapped via `github.com/bborbe/errors`, tests in Ginkgo v2 / Gomega with Counterfeiter mocks, a `## Unreleased` CHANGELOG entry, doc comments on exported symbols.
- Do NOT run `go mod vendor` and do NOT add `-mod=vendor` to any command; this repo does not commit `vendor/`.
</constraints>

<verification>
Run `make test` — must pass (exit code 0).

Run `grep -rn 'BACKUP_ALREADY_RUNNING\|CLEANUP_ALREADY_RUNNING' pkg/` — must return at least one line for each code.

Run `grep -n 'NewJSONErrorHandler' pkg/factory/factory.go` — must return exactly two lines (`CreateBackupHandler` and `CreateCleanupHandler`).

Run `grep -n 'libhttp.NewErrorHandler' pkg/factory/factory.go` — must still return exactly two lines (`CreateStatusHandler` and `CreateListHandler`).

Run `make precommit` — must pass.
</verification>
