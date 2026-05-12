---
status: completed
summary: Bumped axios in frontend/package.json from ^1.15.0 to ^1.16.0, regenerated package-lock.json, and added CHANGELOG entry covering the CVE fixes.
container: service-004-bump-axios-frontend
dark-factory-version: v0.156.1-1-g04f3863-dirty
created: "2026-05-12T13:00:00Z"
queued: "2026-05-12T17:14:40Z"
started: "2026-05-12T19:42:00Z"
completed: "2026-05-12T19:45:12Z"
---

<summary>
- Bumps `axios` in `frontend/package.json` from `^1.15.0` to `^1.15.1` (or latest 1.16.x)
- Resolves Dependabot axios advisories for bborbe/backup: CVE-2026-42035 (High), CVE-2026-42033 (High), CVE-2026-42040 (Low), CVE-2026-42036 (Moderate), CVE-2026-42044 (Moderate), plus 3 more
- Regenerates `frontend/package-lock.json` via `npm install`
- `make precommit` exits 0 after the change (runs `frontend-precommit` which lints + builds + tests the Vue frontend)
- CHANGELOG `## Unreleased` documents the bump
</summary>

<objective>
Patch the axios CVEs reported by Dependabot on bborbe/backup by upgrading the dependency in the Vue frontend to `^1.15.1`.
</objective>

<context>
Read `CLAUDE.md` for project conventions.

Affected file: `frontend/package.json` — currently:
```json
"axios": "^1.15.0",
```

The repo's root `package.json` declares `follow-redirects` only (no axios). Do NOT touch root `package.json`.

Frontend uses Vue 3 + Vite. `make precommit` invokes `make -C frontend precommit` which runs `npm run lint`, `npm run build`, and `npm test`. Lockfile is `frontend/package-lock.json`.
</context>

<requirements>
1. Update axios in the frontend:
   ```bash
   cd frontend
   npm install axios@^1.15.1 --save
   ```
   This rewrites `frontend/package.json` and `frontend/package-lock.json`.

2. Verify no other deps changed unexpectedly:
   ```bash
   git diff frontend/package.json
   ```
   Only the axios line should differ.

3. Run the full precommit:
   ```bash
   make precommit
   ```
   Must exit 0. The frontend build + lint + tests must all pass.

4. Update `CHANGELOG.md` at the repo root under `## Unreleased`:
   ```
   - security(frontend): bump axios to ^1.15.1 (CVE-2026-42035, CVE-2026-42033, CVE-2026-42040, CVE-2026-42036, CVE-2026-42044 + 3 more)
   ```
</requirements>

<constraints>
- Only edit: `frontend/package.json`, `frontend/package-lock.json`, `CHANGELOG.md`
- Do NOT touch root `package.json` / `package-lock.json` (different deps)
- Do NOT bump unrelated deps
- Do NOT use `npm audit fix --force` (may downgrade unrelated deps)
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
</constraints>

<verification>
```bash
grep '"axios"' frontend/package.json       # must show ^1.15.1 or ^1.16.x
make precommit                              # must exit 0
```
</verification>
