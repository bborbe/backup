---
status: approved
approved: "2026-08-09T08:57:23Z"
branch: dark-factory/conflict-response-for-already-running-triggers
---

## Summary

- Asking the service to back up a host whose backup is already running is a normal, expected situation — today it is reported as a server error.
- The service will answer such a request with "conflict — the resource is busy" instead of "the server broke".
- The answer will carry a stable, machine-readable reason code plus a human-readable sentence naming the host.
- The dashboard will show that sentence to the operator instead of the meaningless "Request failed with status code 500".
- The dashboard will present it as a warning ("this is busy, try later"), visually distinct from a genuine failure.

## Problem

When an operator triggers a backup for a host whose backup is already in flight, the service refuses the second run — correctly — but reports the refusal as an internal server error with a plain-text body. The dashboard has no way to read that body, so it falls back to the transport library's generic text and shows "Request failed with status code 500". The operator cannot distinguish "the backup you asked for is already running" (nothing is wrong) from "the backup service is broken" (something is very wrong), and the same false server-error signal reaches any monitoring that watches HTTP status codes. The identical problem exists for cleanup triggers.

## Goal

Triggering a backup or a cleanup for a host that is already backing up or cleaning up is reported end to end as a busy-resource conflict: the API answers with the conflict status and a structured, machine-readable body, and the dashboard renders the server's own explanation as a warning that names the host. Genuine failures continue to look like failures, and the underlying refusal to start a second run is unchanged.

## Non-goals

- Do NOT change error handling for any endpoint other than the per-host backup trigger and the per-host cleanup trigger. The status endpoint, the list endpoint, and the trigger-all endpoints keep their current error format; revisiting them is a separate spec.
- Do NOT add retry or backoff logic in the dashboard for already-running conflicts. The operator retries manually.
- Do NOT change the only-once locking mechanism that decides whether a backup or cleanup is already running. Only the reporting of its refusal changes.
- Do NOT redesign the dashboard's error/toast presentation beyond adding the severity distinction this spec needs.
- Do NOT add a config flag, env var, or query parameter to switch between the plain-text and the structured error body — the structured body is an invariant of these two endpoints. If a future consumer needs the old format, that's a separate spec.
- Do NOT introduce new colors or design tokens; the existing warning tokens in `docs/design-system.md` are the palette.

## Acceptance Criteria

- [ ] Triggering a backup for an already-running host returns HTTP 409 (not 500).
      Evidence: a Go test drives the per-host backup trigger handler through its real HTTP wiring with a stubbed executor that returns the already-running sentinel, and asserts the recorded response status is exactly `409`. Run via `make test`, exit code 0.
- [ ] Triggering a cleanup for an already-running host returns HTTP 409 (not 500).
      Evidence: the equivalent Go test for the per-host cleanup trigger handler asserts the recorded response status is exactly `409`. Run via `make test`, exit code 0.
- [ ] The 409 response body is JSON with a stable machine-readable error code (`BACKUP_ALREADY_RUNNING` / `CLEANUP_ALREADY_RUNNING`), not text/plain.
      Evidence: the same two Go tests assert (a) response header `Content-Type` contains `application/json`, (b) the body unmarshals into an object with an `error` member whose `code` field equals exactly `BACKUP_ALREADY_RUNNING` (backup) / `CLEANUP_ALREADY_RUNNING` (cleanup), and (c) the `error.message` field is non-empty and contains both the requested host name and the substring `already running`. Run via `make test`, exit code 0.
- [ ] The dashboard displays the human-readable reason (e.g. "Backup for X is already running") instead of "Request failed with status code 500".
      Evidence: a frontend test mocks the HTTP client to answer the trigger call with status 409 and the structured body above, then asserts the rendered message bar's text contains the host name and the substring `already running`, and contains neither `Request failed with status code` nor `500`. Run via `make -C frontend test`, exit code 0.
- [ ] The dashboard renders a 409 as a neutral/warning message, not a red error bar.
      Evidence: the same frontend test asserts the rendered message element carries the warning severity class and does NOT carry the error severity class; a companion assertion in the same file drives a 500 response through the same code path and asserts the inverse (error class present, warning class absent). Additionally `grep -n 'status-warning\|warning' frontend/src/components/BackupStatusOverviewComponent.vue` returns at least one line, proving the live per-host trigger UI (the component actually mounted by `DashboardPage.vue`, the sole route in `router.ts`) now uses the existing `--status-warning*` design token rather than a new color. `ActionPanelComponent.vue` and `TargetListComponent.vue` are dead code (referenced by no `.vue`/`.ts` source, only by a stale build artifact) and must NOT be the target of this AC's evidence.
- [ ] A failing trigger whose response body is not the structured envelope falls back to the transport-level message, with no blank bar and no crash.
      Evidence: a frontend test answers the trigger call with status 500 and a `text/plain` body (not the structured envelope), then asserts the rendered message bar text is non-empty, equals the axios transport-level message, does not contain `undefined` or render blank, and that the test itself completes without an unhandled rejection or thrown exception. Run via `make -C frontend test`, exit code 0.

**Scenario coverage: none.** Every criterion above is reachable with an in-process HTTP recorder on the Go side and a mocked HTTP client on the frontend side. No real deployment, cluster, or long-running backup is required, so no end-to-end scenario is justified.

## Verification

### Container-executable (runs inside the YOLO container at prompt time)

```
make test
make -C frontend test
make precommit
grep -rn 'BACKUP_ALREADY_RUNNING\|CLEANUP_ALREADY_RUNNING' pkg/
grep -n 'status-warning\|warning' frontend/src/components/BackupStatusOverviewComponent.vue
```

Expected: `make test`, `make -C frontend test` and `make precommit` all exit 0. Both greps return at least one line each.

No operator-executable rung: this change ships as merged code with no deployment-observable step that the container cannot already prove.

## Desired Behavior

1. A per-host backup trigger for a host whose backup is already running is answered with the HTTP conflict status, because the request is well-formed and the resource is merely busy.
2. A per-host cleanup trigger for a host whose cleanup is already running is answered the same way.
3. Both conflict answers carry a structured body containing a stable machine-readable code and a human-readable message that names the host and states that the work is already running.
4. The dashboard extracts the server-supplied code and message from the response body for every failing backup and cleanup trigger, in one shared place rather than per call site.
5. When no structured body is present (network failure, proxy-rewritten body, non-JSON response), the dashboard falls back to the transport-level message it shows today — no blank bar, no crash.
6. The dashboard tags an already-running conflict with warning severity and every other failure with error severity, and the message bar renders those two severities differently.
7. Successful triggers, genuine server errors, and network failures keep their existing status codes, bodies, and dashboard behavior.

## Constraints

- The HTTP method and paths of the per-host backup trigger and the per-host cleanup trigger are frozen — only the failure status code and body shape change.
- Success responses for both endpoints — status code and body — are unchanged.
- The error codes `BACKUP_ALREADY_RUNNING` and `CLEANUP_ALREADY_RUNNING` are frozen string contracts. The dashboard and any future client match on them exactly.
- The JSON envelope shape is dictated by `github.com/bborbe/http` v1.26.16 (`{"error":{"code":…,"message":…,"details":…}}`). Do not hand-roll a competing envelope; use the library's status-code/error-code wrapping and its JSON error handler.
- The only-once refusal behavior is unchanged: after this work, a conflicting trigger still starts no second run.
- All existing Go tests and all existing frontend tests under `frontend/src/tests/` continue to pass.
- `docs/dod.md` applies: Go errors wrapped via `github.com/bborbe/errors`, Go tests in Ginkgo v2 / Gomega with Counterfeiter mocks, a `## Unreleased` CHANGELOG entry, doc comments on exported symbols.
- Frontend styling uses the existing `--status-warning*` and `--status-error*` custom properties from `docs/design-system.md`.

## Assumptions

- The already-running condition is already detectable as a distinct sentinel error inside the service; this spec only changes how that sentinel is translated to HTTP and to the UI.
- `github.com/bborbe/http` v1.26.16 already provides both the error-with-status-code/error-code wrapping and the JSON error handler, so no new dependency is needed.
- The dashboard's HTTP client is the only path by which the UI learns about trigger failures.
- Operators read the message bar; there is no separate notification channel to update.

## Failure Modes

| Trigger | Expected behavior | Detection | Recovery | Concurrency note | Reversibility |
|---|---|---|---|---|---|
| Two trigger requests for the same host arrive concurrently | Exactly one run starts; every other request gets 409 with the already-running code | Server logs one start and N-1 conflict lines for the host | None needed — the second operator retries after the run finishes | This is the core case: the lock, not the handler, decides the winner; the handler only reports the loser's outcome | N/A — correct steady-state behavior |
| The backup finishes between the lock check and the response being written | The in-flight request still returns 409 | Operator sees the warning bar while the dashboard already shows the host as idle | Operator triggers again; the next request succeeds | Benign race — the answer is stale by milliseconds, never wrong about "a run was in progress" | N/A — correct steady-state behavior |
| Browser holds a cached pre-change dashboard bundle while the new API is live | Old bundle shows the transport-level text (now "…status code 409") | Message bar reads "Request failed with status code 409" | Hard reload; the served bundle ships with the API in the same image | — | Reversible |
| A proxy or gateway rewrites the error body to non-JSON | Dashboard falls back to the transport-level message and error severity; no crash, no blank bar | Message bar shows the generic transport text instead of the host-specific sentence | Fix the proxy; no service-side change needed | — | Reversible |
| The already-running sentinel is renamed or restructured upstream so the conflict is no longer recognised | The endpoint falls back to the pre-change 500 behavior | The Go tests for AC 1–3 fail at build time, before merge | Re-point the detection at the new sentinel | — | Reversible, caught pre-merge |
| The HTTP library changes its JSON error envelope in a future version | Dashboard falls back to the transport-level message | The frontend test asserting the envelope shape fails at build time | Pin or adapt to the new envelope | — | Reversible, caught pre-merge |
| The API is unreachable (service down, network failure) | Unchanged from today: dashboard shows the transport error at error severity | Red error bar, no code extracted | Operator retries once the service is reachable | — | Unaffected by this spec |
| An operator hammers the trigger button while a run is in flight | Every extra request is rejected by the existing lock and answered 409; no extra backup work is started and no extra load beyond one cheap HTTP round-trip each | Repeated conflict log lines for the same host | None needed | The lock is the throttle; this spec adds no new resource consumption per rejected request | N/A — correct steady-state behavior |

## Security / Abuse Cases

- **Attacker-controlled input:** the host name in the trigger path. It is looked up against the configured targets before any work happens; an unknown name never reaches the conflict path and its failure mode is unchanged by this spec.
- **Reflected input:** the host name now appears in the response body and is rendered in the dashboard message bar. It must be rendered as text through the framework's escaping interpolation — never as raw HTML — so a crafted target name cannot inject markup or script into the dashboard.
- **Information disclosure:** the structured message is the wrapped error text. It must not carry filesystem paths, credentials, or internal stack detail beyond the host name and the already-running statement. Reviewers check the emitted message string against this rule.
- **Resource exhaustion:** rejected triggers do no work — the lock rejects before any backup process starts — so repeated triggering cannot be amplified into load. The conflict path adds no new allocation-heavy work.
- **Trust boundary:** unchanged. This spec adds no new endpoint, no new authentication surface, and no new privilege.

## Suggested Decomposition

Prompts run in this order — each row is a single prompt with a clear scope.

| # | Prompt focus | Covers DBs | Covers ACs | Depends on |
|---|---|---|---|---|
| 1 | API: translate the already-running sentinel into a conflict status plus stable error code for both per-host trigger handlers, and switch those two handlers' wiring to the JSON error responder; Go tests asserting status, content type, code, and message | 1, 2, 3, 7 (API half) | 1, 2, 3 | — |
| 2 | Dashboard: central extraction of the structured code and message from failing trigger responses with fallback to the transport message; warning-vs-error severity on the message bar; frontend tests for the 409, 500, and non-structured-body fallback paths | 4, 5, 6, 7 (UI half) | 4, 5, 6 | prompt 1 (defines the codes and body shape the UI parses) |

Rationale: the API prompt establishes the wire contract — the two frozen error codes and the JSON envelope — that the dashboard prompt parses and asserts against. Splitting the other way would force the UI prompt to invent a contract that the API prompt might then contradict. Each prompt is verifiable on its own: prompt 1 by `make test`, prompt 2 by `make -C frontend test`.

## Do-Nothing Option

If nothing changes, operators keep seeing "Request failed with status code 500" whenever they trigger a backup or cleanup that is already running. The cost is real but bounded: every such event looks identical to a genuine service failure, so the operator either investigates a non-incident or — worse — learns to ignore 500s from this service and misses a real one. Any monitoring that counts server errors on these endpoints inflates its error rate on entirely healthy behavior. The workaround available today is for the operator to check the dashboard's per-host status before triggering, which is exactly the manual step the trigger button exists to avoid. This is not acceptable as a permanent state, and the fix is small and self-contained.
