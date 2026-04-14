---
status: completed
container: service-002-fix-codecov-coverage-file
dark-factory-version: v0.108.0-dirty
created: "2026-04-13T18:41:44Z"
queued: "2026-04-13T18:41:44Z"
started: "2026-04-13T18:41:49Z"
completed: "2026-04-14T06:06:31Z"
lastFailReason: 'validate completion report: completion report status: partial'
---

<summary>
- Generate a coverage profile file during test runs for codecov upload
- CI codecov step currently fails because no coverage file exists
- Add coverprofile flag to Makefile test target
- Coverage file includes results from all tested packages
- No functional changes to test behavior
</summary>

<objective>
Fix the CI build failure where the codecov upload step fails with "No coverage reports found". The Makefile test target uses `-cover` (console output only) but does not generate a `-coverprofile` file that codecov can upload.
</objective>

<context>
Read CLAUDE.md for project conventions.
Read `Makefile` — find the `test` target. It currently runs `go test -mod=mod -p=... -cover $(shell go list ...)` without `-coverprofile`.
Read `.github/workflows/ci.yml` — the codecov step expects coverage files to exist after `make precommit`.
Read `.gitignore` — check if `coverage.out` is already listed.
Modern Go (1.20+) correctly merges multi-package coverage when `-coverprofile` is used, so a simple flag addition is sufficient.
</context>

<requirements>
1. In the `Makefile` `test` target, add `-coverprofile=coverage.out` alongside the existing `-cover` flag
2. Keep the existing `-cover` flag (console output is still useful)
3. Add `/coverage.out` to `.gitignore` if not already present
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git
- Do NOT modify `.github/workflows/ci.yml`
- Do NOT change any Go source code
- Existing tests must still pass
</constraints>

<verification>
Run `make precommit` -- must pass.
Run `make test && test -f coverage.out` -- coverage file must exist after tests.
</verification>
