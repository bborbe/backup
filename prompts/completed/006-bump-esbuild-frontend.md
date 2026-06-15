---
status: completed
summary: Added esbuild ^0.28.1 override to frontend/package.json, regenerated lockfile (node_modules/esbuild now 0.28.1), and updated CHANGELOG. make precommit exits 0, vitest 22/22 pass, npm audit shows 0 vulnerabilities, GHSA-gv7w-rqvm-qjhr no longer reported.
container: backup-fix-esbuild-exec-006-bump-esbuild-frontend
dark-factory-version: v0.177.1
created: "2026-06-15T16:11:15Z"
queued: "2026-06-15T16:11:15Z"
started: "2026-06-15T16:12:10Z"
completed: "2026-06-15T16:15:12Z"
---
<summary>
- Force-patches transitive `esbuild` from `0.27.3` to `^0.28.1` in `frontend/` via an npm `overrides` entry
- Resolves GHSA-gv7w-rqvm-qjhr (High): "esbuild: Missing binary integrity verification in Deno module enables remote code execution via NPM_CONFIG_REGISTRY"
- Regenerates `frontend/package-lock.json` so `node_modules/esbuild` resolves to `>= 0.28.1`
- `make precommit` exits 0 after the change (lints + builds + tests the Vue frontend)
- `npm audit --audit-level=high` shows no remaining esbuild advisory
- CHANGELOG `## Unreleased` documents the bump
</summary>

<objective>
Patch the high-severity esbuild advisory GHSA-gv7w-rqvm-qjhr reported by Dependabot on bborbe/backup. `esbuild` is a transitive dep of `vite ^7.3.1`, which pins `esbuild ^0.27.0` and so resolves to the vulnerable `0.27.3`. Bumping vite would be a major upgrade (vite 8 replaces esbuild with rolldown); instead force the patched esbuild via an npm `overrides` entry on `frontend/package.json`.
</objective>

<context>
Read `README.md` for project conventions.

Affected file: `frontend/package-lock.json` — `node_modules/esbuild` currently resolves to `0.27.3` (vulnerable range `>= 0.17.0, < 0.28.1`). Pulled in by `vite ^7.3.1` which declares `"esbuild": "^0.27.0"`.

Advisory: https://github.com/advisories/GHSA-gv7w-rqvm-qjhr (fixed in 0.28.1)
Dependabot alert: https://github.com/bborbe/backup/security/dependabot/74

The repo root `package.json` does NOT declare esbuild. Do NOT touch root `package.json` / `package-lock.json`.

Frontend uses Vue 3 + Vite 7 + Vitest 4. `make precommit` invokes `make -C frontend precommit` which runs `npm run lint`, `npm run build`, and `npm test` (`vitest`). Lockfile is `frontend/package-lock.json`.
</context>

<requirements>
1. Add an npm `overrides` entry to `frontend/package.json` so transitive esbuild resolves to the patched version. Add a top-level `"overrides"` block (or extend it if present):
   ```json
   "overrides": {
     "esbuild": "^0.28.1"
   }
   ```

2. Regenerate the lockfile:
   ```bash
   cd frontend
   npm install
   ```
   This rewrites `frontend/package-lock.json`. The `node_modules/esbuild` entry must now show `"version": "0.28.1"` (or higher 0.28.x).

3. Confirm the override took effect mechanically:
   ```bash
   jq -r '.overrides.esbuild' frontend/package.json
   ```
   Must print `^0.28.1`. (`<constraints>` already pins the edit allowlist, so a manual `git diff` eyeball is not needed.)

4. Run the frontend tests directly to confirm vite still works with the overridden esbuild:
   ```bash
   cd frontend && npm test -- --run
   ```
   Must exit 0. The `--run` flag forces single-run mode — without it vitest 4 defaults to watch mode and hangs in a non-TTY container.

5. Run the full precommit from the repo root:
   ```bash
   make precommit
   ```
   Must exit 0. Lint + build + tests must all pass.

6. Confirm the specific advisory is gone (without failing on unrelated high+ advisories that may exist in the lockfile):
   ```bash
   cd frontend && (npm audit --audit-level=high || true) | grep -F 'GHSA-gv7w-rqvm-qjhr' && exit 1 || true
   ```
   Must exit 0 (i.e. `GHSA-gv7w-rqvm-qjhr` not found in audit output). If other unrelated high/critical advisories exist, leave them alone — they are out of scope for this prompt.

7. Update `CHANGELOG.md` at the repo root under `## Unreleased` (create the `## Unreleased` section directly above the most recent released version if it does not already exist):
   ```
   - security(frontend): override esbuild to ^0.28.1 (GHSA-gv7w-rqvm-qjhr, High)
   ```
</requirements>

<constraints>
- Only edit: `frontend/package.json`, `frontend/package-lock.json`, `CHANGELOG.md`
- Do NOT touch root `package.json` / `package-lock.json` (different deps)
- Do NOT bump unrelated deps
- Do NOT bump vite to 8.x (major version, out of scope)
- Do NOT use `npm audit fix --force` (may downgrade unrelated deps or perform major bumps)
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
</constraints>

<verification>
```bash
[ "$(jq -r '.overrides.esbuild' frontend/package.json)" = "^0.28.1" ]  # must be true
grep -B1 -A3 '"node_modules/esbuild":' frontend/package-lock.json | grep -E '"version": "0\.2[89]\.|"version": "0\.[3-9][0-9]\.'  # must show 0.28+ version
cd frontend && npm test -- --run                                       # must exit 0
cd frontend && (npm audit --audit-level=high || true) | grep -F 'GHSA-gv7w-rqvm-qjhr' && exit 1 || true  # must not find advisory
make precommit                                                         # must exit 0
grep -F 'GHSA-gv7w-rqvm-qjhr' CHANGELOG.md                            # must match
```
</verification>
