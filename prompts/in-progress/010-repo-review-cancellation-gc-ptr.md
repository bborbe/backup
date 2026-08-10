---
status: approved
summary: Add ctx.Done() guards to the two filesystem loops, register the /gc admin endpoint, and replace the boolPointer helper with collection.Ptr
created: "2026-08-10T09:31:24Z"
queued: "2026-08-10T09:31:24Z"
---

<summary>
- Two loops that touch the filesystem once per iteration now stop promptly when the context is cancelled, instead of running to completion.
- The one that matters most deletes backup directories; a cancelled run should stop deleting, not finish the list.
- The admin router gains the standard `/gc` endpoint it was missing; the other four canonical endpoints were already there.
- A local one-line `boolPointer` helper is replaced by the fleet's `collection.Ptr`.
- No behaviour changes beyond prompt cancellation; no signatures change.
</summary>

<objective>
Apply three independent code-review fixes in `/workspace`:
1. cancellation guards in two per-iteration filesystem loops (finding M1)
2. register the canonical `/gc` admin endpoint (finding M2)
3. replace the `boolPointer` helper with `collection.Ptr` (finding S1)
</objective>

<context>
Read `/workspace/docs/dod.md` for project conventions and the definition of done (Ginkgo v2 / Gomega, counterfeiter, `github.com/bborbe/errors` wrapping, CHANGELOG expectations). This repo has no committed `CLAUDE.md` — `docs/dod.md` is the conventions doc and is also the pipeline's `validationPrompt`.

Read these files IN FULL before editing:
- `/workspace/pkg/handler/handler_status.go` — the outer `for _, target := range targets` loop calls `backupFinder.List(ctx, target.Spec.Host)` each iteration, which reads the filesystem. `ctx` is already in scope from the enclosing `JsonHandlerFunc`.
- `/workspace/pkg/backup-cleaner.go` — `Clean(ctx, backupHost)` loops `for i, date := range dates` and calls `os.RemoveAll` on the delete path.
- `/workspace/main.go` — find the admin router registration; `/healthz`, `/readiness`, `/metrics` and `/setloglevel` are already registered, `/gc` is not.
- `/workspace/pkg/k8s-connector.go` — contains the `boolPointer` helper and its single call site.

For the guard, follow the wrapping convention already used in these files: `errors.Wrap(ctx, ctx.Err(), "context cancelled")`.
</context>

<requirements>
1. In `/workspace/pkg/handler/handler_status.go`, add a non-blocking cancellation check at the TOP of the **outer** `for _, target := range targets` loop body:
   ```go
   select {
   case <-ctx.Done():
       return nil, errors.Wrap(ctx, ctx.Err(), "context cancelled")
   default:
   }
   ```
   Do NOT add a guard to the inner `for _, backupTime := range dates` loop — it is pure in-memory comparison and is already covered by the outer guard.
2. In `/workspace/pkg/backup-cleaner.go`, add a cancellation check at the TOP of the `for i, date := range dates` loop body, returning `errors.Wrap(ctx, ctx.Err(), "context cancelled")`. It must come BEFORE the `i < b.backupKeepAmount` and `backupCleanEnabled` checks, so a cancelled context stops the loop even on iterations that would skip.
3. In `/workspace/main.go`, register `/gc` on the same router that serves `/healthz`, `/readiness`, `/metrics` and `/setloglevel`, using `libhttp.NewGarbageCollectorHandler()` from `github.com/bborbe/http` (already imported as `libhttp`; `/healthz` and `/readiness` come from the same package, while `/setloglevel` comes from `github.com/bborbe/log` — match the `router.Path(...).Handler(...)` shape of `/healthz`). Do not change any existing route.
4. In `/workspace/pkg/k8s-connector.go`, replace the single `boolPointer(true)` call with `collection.Ptr(true)` from `github.com/bborbe/collection`, then delete the now-unused `boolPointer` helper. Add the import if absent.
5. Do NOT change any function signature, return type, or interface definition. In particular do NOT modify the `K8sConnector` interface — that is deliberately deferred.
6. Do NOT touch anything under `/workspace/k8s/client/`, `/workspace/mocks/`, or any file marked `// Code generated ... DO NOT EDIT.`
7. Add a test to `/workspace/pkg/backup-cleaner_test.go` covering the new cancellation branch: call `Clean` with an already-cancelled context (`ctx, cancel := context.WithCancel(context.Background()); cancel()`) and assert it returns a non-nil error. Follow the existing Ginkgo v2 / Gomega style already in that file. `docs/dod.md` requires tests for changed behaviour, and this is a new code path.
8. Add a `## Unreleased` section to the TOP of `/workspace/CHANGELOG.md`, directly above the current top entry `## v3.9.27`, with three bullets in the existing style of the entries below it — one per change: the cancellation guards, the `/gc` endpoint, and the `collection.Ptr` replacement. `docs/dod.md` requires a `## Unreleased` entry; the file currently has none.
</requirements>

<verification>
- `cd /workspace && make precommit` exits 0.
- `grep -n 'ctx.Done()' /workspace/pkg/handler/handler_status.go` returns exactly one match.
- `grep -n 'ctx.Done()' /workspace/pkg/backup-cleaner.go` returns exactly one match.
- `grep -cn 'ctx.Done()' /workspace/pkg/handler/handler_status.go` is 1, NOT 2 — the inner loop must remain unguarded.
- `grep -n '"/gc"' /workspace/main.go` returns at least one match.
- `grep -c 'boolPointer' /workspace/pkg/k8s-connector.go` returns 0.
- `grep -n 'collection.Ptr' /workspace/pkg/k8s-connector.go` returns at least one match.
- `grep -n '## Unreleased' /workspace/CHANGELOG.md` returns exactly one match, and it appears before the `## v3.9.27` line.
- `grep -c 'ctx.Done()\|WithCancel' /workspace/pkg/backup-cleaner_test.go` is at least 1.
- `cd /workspace && git diff --name-only` lists ONLY: `CHANGELOG.md`, `main.go`, `pkg/backup-cleaner.go`, `pkg/backup-cleaner_test.go`, `pkg/handler/handler_status.go`, `pkg/k8s-connector.go`.
</verification>

<allowed_files>
- /workspace/main.go
- /workspace/pkg/backup-cleaner.go
- /workspace/pkg/handler/handler_status.go
- /workspace/pkg/k8s-connector.go
- /workspace/pkg/backup-cleaner_test.go
- /workspace/CHANGELOG.md
</allowed_files>
