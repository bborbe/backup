---
status: completed
summary: Bumped vitest from ^4.0.18 to ^4.1.0 in frontend/package.json, regenerated package-lock.json, verified tests pass and advisory GHSA-5xrq-8626-4rwp is resolved.
container: service-vitest-bump-exec-005-bump-vitest-frontend
dark-factory-version: v0.174.0-14-g29eee19
created: "2026-06-02T14:31:38Z"
queued: "2026-06-02T14:31:38Z"
started: "2026-06-02T14:32:22Z"
completed: "2026-06-02T14:34:21Z"
---
<summary>
- Bumps `vitest` in `frontend/package.json` from `^4.0.18` to `^4.1.0` (or latest 4.1.x)
- Resolves GHSA-5xrq-8626-4rwp (Critical): "Vitest UI server allows arbitrary file read/execution when listening"
- Regenerates `frontend/package-lock.json` via `npm install`
- `make precommit` exits 0 after the change (lints + builds + tests the Vue frontend)
- `npm audit --audit-level=high` shows no remaining vitest advisory
- CHANGELOG `## Unreleased` documents the bump
</summary>

<objective>
Patch the critical vitest advisory GHSA-5xrq-8626-4rwp reported by Dependabot on bborbe/backup by upgrading the dev dependency in the Vue frontend to `^4.1.0`.
</objective>

<context>
Read `CLAUDE.md` for project conventions.

Affected file: `frontend/package.json` — currently:
```json
"vitest": "^4.0.18",
```

Advisory: https://github.com/advisories/GHSA-5xrq-8626-4rwp
Dependabot alert: https://github.com/bborbe/backup/security/dependabot/72

The repo root `package.json` does NOT declare vitest. Do NOT touch root `package.json` / `package-lock.json`.

Frontend uses Vue 3 + Vite + Vitest. `make precommit` invokes `make -C frontend precommit` which runs `npm run lint`, `npm run build`, and `npm test` (`vitest`). Lockfile is `frontend/package-lock.json`.
</context>

<requirements>
1. Update vitest in the frontend:
   ```bash
   cd frontend
   npm install vitest@^4.1.0 --save-dev
   ```
   This rewrites `frontend/package.json` and `frontend/package-lock.json`.

2. Verify no unrelated deps changed unexpectedly:
   ```bash
   git diff frontend/package.json
   ```
   Only the `vitest` line under `devDependencies` should differ. (Transitive lockfile updates in `package-lock.json` are fine.)

3. Run the frontend tests directly to confirm the new vitest runs the suite green:
   ```bash
   cd frontend && npm test -- --run
   ```
   Must exit 0. The `--run` flag forces single-run mode — without it vitest 4 defaults to watch mode and hangs in a non-TTY container.

4. Run the full precommit from the repo root:
   ```bash
   make precommit
   ```
   Must exit 0. Lint + build + tests must all pass.

5. Confirm the advisory is gone:
   ```bash
   cd frontend && npm audit --audit-level=high
   ```
   Must NOT list GHSA-5xrq-8626-4rwp or any vitest advisory at high+ severity. If other unrelated high/critical advisories exist, leave them alone — they are out of scope for this prompt.

6. Update `CHANGELOG.md` at the repo root under `## Unreleased` (create the `## Unreleased` section directly above the most recent released version if it does not already exist):
   ```
   - security(frontend): bump vitest to ^4.1.0 (GHSA-5xrq-8626-4rwp, Critical)
   ```
</requirements>

<constraints>
- Only edit: `frontend/package.json`, `frontend/package-lock.json`, `CHANGELOG.md`
- Do NOT touch root `package.json` / `package-lock.json` (different deps)
- Do NOT bump unrelated deps
- Do NOT use `npm audit fix --force` (may downgrade unrelated deps or perform major bumps)
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
</constraints>

<verification>
```bash
grep '"vitest":' frontend/package.json          # must show ^4.1.0 or higher 4.x
cd frontend && npm test -- --run                # must exit 0
cd frontend && npm audit --audit-level=high    # must not list GHSA-5xrq-8626-4rwp
make precommit                                  # must exit 0
grep -F 'GHSA-5xrq-8626-4rwp' CHANGELOG.md     # must match
```
</verification>
