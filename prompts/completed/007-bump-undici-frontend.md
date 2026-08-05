---
status: completed
summary: Bumped transitive undici from 7.28.0 to 7.29.0 in frontend/package-lock.json, resolving 5 Dependabot advisories (GHSA-4cwx-7wf7-3272 High + 4 Moderate)
execution_id: backup-bump-undici-exec-007-bump-undici-frontend
dark-factory-version: v0.192.9
created: "2026-08-05T19:37:09Z"
queued: "2026-08-05T19:37:09Z"
started: "2026-08-05T19:38:19Z"
completed: "2026-08-05T19:41:23Z"
---
<summary>
- Bumps transitive `undici` from `7.28.0` to `7.29.0` in `frontend/package-lock.json`
- Resolves 5 Dependabot alerts: GHSA-4cwx-7wf7-3272 (High) + GHSA-v3r7-h72x-cjcm, GHSA-m8rv-5g2x-5cg5, GHSA-jr45-8vmc-qm54, GHSA-8xcm-r25x-g524 (Moderate)
- Lockfile-only change — the parent `jsdom` already declares `undici: ^7.21.0`, which accepts `7.29.0`, so no `overrides` block and no `package.json` edit are needed
- `make precommit` exits 0 after the change (lints + builds + tests the Vue frontend, and runs osv-scanner + trivy)
- CHANGELOG `## Unreleased` documents the bump
</summary>

<objective>
Patch the five Dependabot advisories on `undici` reported for bborbe/backup. `undici` is a transitive dev dependency of `jsdom`, currently resolving to the vulnerable `7.28.0`. All five advisories are fixed in `7.29.0`.
</objective>

<context>
Read `README.md` for project conventions.

Affected file: `frontend/package-lock.json` — `node_modules/undici` currently resolves to `7.28.0`. It is pulled in by `jsdom` (dev dependency), which declares `"undici": "^7.21.0"`.

Because that caret range already accepts `7.29.0`, this is a plain lockfile refresh — NOT an `overrides` case. Do not add an `overrides` block and do not add `undici` as a direct dependency of `frontend/package.json`.

Advisories (all fixed in 7.29.0):
- https://github.com/advisories/GHSA-4cwx-7wf7-3272 (High)
- https://github.com/advisories/GHSA-v3r7-h72x-cjcm (Moderate)
- https://github.com/advisories/GHSA-m8rv-5g2x-5cg5 (Moderate)
- https://github.com/advisories/GHSA-jr45-8vmc-qm54 (Moderate)
- https://github.com/advisories/GHSA-8xcm-r25x-g524 (Moderate)

Dependabot alerts: https://github.com/bborbe/backup/security/dependabot — alerts 82, 83, 84, 85, 86.

The repo root `package.json` does NOT declare undici. Do NOT touch root `package.json` / `package-lock.json`.

Frontend uses Vue 3 + Vite + Vitest 4. Root `make precommit` runs the Go checks (including `osv-scanner` and `trivy`) AND `make -C frontend precommit`, whose targets run `npm install` first and then `npm run lint:analyse`, `npm run build`, and `npm run test -- --run`. Lockfile is `frontend/package-lock.json`.

Every command below that operates inside `frontend/` is written either as a `make -C frontend <target>` invocation or wrapped in a `(cd frontend && ...)` subshell. Keep it that way — a bare `cd frontend` leaks into later commands and would make the root `make precommit` resolve to `frontend/Makefile`'s frontend-only target, silently skipping the osv-scanner and trivy checks that actually prove these advisories are gone.
</context>

<requirements>
1. Refresh the transitive `undici` pin in the frontend lockfile:
   ```bash
   (cd frontend && npm update undici)
   ```
   This rewrites only `frontend/package-lock.json`. The `node_modules/undici` entry must afterwards show `"version": "7.29.0"` (or higher within 7.x).

   If that does not move the version (e.g. npm keeps the cached resolution), retry with:
   ```bash
   (cd frontend && npm update undici --package-lock-only)
   ```
   Never use `npm install undici@7.29.0` — `npm install <spec>` writes the package to `package.json` even with `--package-lock-only`, which would turn a transitive dep into a direct one and violate the constraints below.

2. Confirm `frontend/package.json` is unchanged and the lockfile blast radius is small:
   ```bash
   git diff --exit-code frontend/package.json
   git diff --stat frontend/package-lock.json
   ```
   The first must exit 0. If it reports a diff, restore the file and redo step 1:
   ```bash
   git checkout -- frontend/package.json
   (cd frontend && npm update undici --package-lock-only)
   ```
   The second is an inspection: expect a small diff confined to the `undici` entry (version / resolved / integrity lines). If unrelated packages moved, revert `frontend/package-lock.json` and redo step 1 with `--package-lock-only`.

3. Run the frontend tests to confirm jsdom still works against the bumped undici:
   ```bash
   make -C frontend test
   ```
   Must exit 0. This target runs `npm install` first (so it works even after a `--package-lock-only` refresh) and passes `--run` to vitest — without that flag vitest 4 defaults to watch mode and hangs in a non-TTY container.

4. Run the full precommit from the repo root:
   ```bash
   make precommit
   ```
   Must exit 0, and must be invoked from the repo root — NOT from inside `frontend/`. This is the primary gate: it runs `osv-scanner` and `trivy`, neither of which has an undici entry in `.osv-scanner.toml` or `.trivyignore`, so it fails while any of the five advisories remain.

5. Update `CHANGELOG.md` at the repo root under `## Unreleased` (create the `## Unreleased` section directly above `## v3.9.24` if it does not already exist):
   ```
   - security(frontend): bump undici to 7.29.0 (GHSA-4cwx-7wf7-3272 High, GHSA-v3r7-h72x-cjcm, GHSA-m8rv-5g2x-5cg5, GHSA-jr45-8vmc-qm54, GHSA-8xcm-r25x-g524 Moderate)
   ```
</requirements>

<constraints>
- Only edit: `frontend/package-lock.json`, `CHANGELOG.md`
- Do NOT edit `frontend/package.json` — the existing `jsdom` caret range already permits 7.29.0
- Do NOT add an `overrides` block — unnecessary here, it would add ongoing carry-cost
- Do NOT touch root `package.json` / `package-lock.json` (different deps)
- Do NOT bump unrelated deps
- Do NOT bump `jsdom` itself (parent major bump is out of scope)
- Do NOT use `npm install undici@<ver>` (writes to `package.json`) or `npm audit fix --force` (may downgrade unrelated deps or perform major bumps)
- If `make precommit` fails on something unrelated to undici, report `status: partial` with the failure — do NOT add `.osv-scanner.toml` / `.trivyignore` entries and do NOT bump other deps to make it pass
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
</constraints>

<verification>
```bash
git diff --exit-code frontend/package.json                          # must exit 0 (file untouched)
# at least one undici entry must exist, and every one of them must be at 7.29.0
[ "$(jq -r '[.packages | to_entries[] | select(.key|endswith("node_modules/undici"))] | length' frontend/package-lock.json)" -ge 1 ]
jq -e '[.packages | to_entries[] | select(.key|endswith("node_modules/undici")) | .value.version] | all(. == "7.29.0")' frontend/package-lock.json
make -C frontend test                                               # must exit 0
make precommit                                                      # must exit 0, run from repo root
grep -F 'GHSA-4cwx-7wf7-3272' CHANGELOG.md                          # must match
```
</verification>
