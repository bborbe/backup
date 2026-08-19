---
status: completed
summary: Add npm existence guard to frontend-precommit target so make precommit succeeds in Node-less containers
execution_id: backup-npm-guard-exec-010-guard-frontend-precommit-when-npm-absent
dark-factory-version: dev
created: "2026-08-19T13:40:00Z"
queued: "2026-08-19T13:49:21Z"
started: "2026-08-19T14:06:18Z"
completed: "2026-08-19T14:20:58Z"
---
<summary>
- `make precommit` no longer fails on machines and containers that have no `npm`
- Where `npm` IS available, the frontend lint/build/test checks run exactly as before — no loss of coverage
- When `npm` is missing, the frontend stage is skipped with an explicit message instead of dying with exit 127
- Unblocks the automated Go/deps update agent, whose container has no `npm` and which currently fails this repo on every run
- CHANGELOG `## Unreleased` documents the change
</summary>

<objective>
Make the root `frontend-precommit` target tolerate an environment without `npm`, so `make precommit` succeeds in Node-less automation containers while still running the full frontend checks wherever `npm` exists.
</objective>

<context>
Read `docs/dod.md` for project conventions (it is also the daemon's `validationPrompt`).

Affected file: `Makefile` (repo root) — the `frontend-precommit` target. It currently reads:

```make
frontend-precommit:
	$(MAKE) -C frontend precommit
```

It is invoked by the root `precommit` target. `frontend/Makefile`'s `lint` target begins with `npm install`, so in an environment without `npm` the sub-make dies with:

```
npm install
make[1]: npm: No such file or directory
make[1]: *** [Makefile:25: lint] Error 127
make: *** [Makefile:122: frontend-precommit] Error 2
```

Why this matters operationally: the `github-update-go-agent` runs this repo's `make precommit` as a release gate inside a container that has Go but no Node/npm. The gate therefore always fails, the agent's task returns to its planning phase, and it re-spawns — observed 11 times on 2026-08-19, each run ~8-9 minutes, monopolising that agent's single concurrency slot while ~18 other repositories queued behind it.

Note the asymmetry: the dark-factory YOLO container that executes THIS prompt does have `npm` (see `prompts/completed/007-bump-undici-frontend.md`, which ran `npm update` successfully). So the happy path and the skip path must both be exercised deliberately — simply running `make precommit` here proves only the happy path.

Do NOT change `frontend/Makefile`. The guard belongs at the root, which is the boundary between "Go project checks" and "frontend checks".
</context>

<requirements>
1. Replace the `frontend-precommit` target in the root `Makefile` with a guarded form. Keep the target name and keep it wired into `precommit` exactly as today:

   ```make
   frontend-precommit:
   	@if command -v npm >/dev/null 2>&1; then \
   		$(MAKE) -C frontend precommit; \
   	else \
   		echo "frontend-precommit: skipped (npm not found in PATH)"; \
   	fi
   ```

   Use `command -v npm`, not `which npm` — `which` is not guaranteed present and its exit status is unreliable across shells.

   Preserve the target's existing `.PHONY` treatment: if the current target has no `.PHONY` entry, do not add one; if it does, keep it.

2. Add a CHANGELOG entry under `## Unreleased` (create the section if absent), wording it around the behaviour, e.g.:
   `- fix: skip frontend-precommit when npm is not installed, so make precommit succeeds in Node-less CI/automation containers`

3. Update `README.md` where it documents `make precommit` (around the build-commands section) with a one-line note that the frontend checks are conditional, e.g.:
   `Frontend lint/build/test run automatically when npm is installed, and are skipped otherwise (e.g. Node-less CI/automation containers).`
   This is required by `docs/dod.md` § Documentation, since the change alters what `make precommit` does depending on the environment.

4. Do not alter any other target, and do not alter `frontend/Makefile`.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git
- Do NOT run `go mod vendor`
- Existing tests must still pass
- The frontend checks MUST still run when `npm` is present — this change is a skip guard, never an unconditional disable
- Do not add npm/Node installation logic to any Makefile or Dockerfile
- Do not touch `frontend/Makefile`, `frontend/package.json`, or either lockfile
</constraints>

<verification>
Run `make precommit` -- must pass, and its output must still show the frontend stage running (npm is present in this container, so the guard must take the TRUE branch).

Then prove the skip branch, by invoking the target with `npm`'s directory removed from PATH:

```bash
MAKE_BIN="$(command -v make)"
env -i PATH="$(dirname "$MAKE_BIN"):/usr/bin:/bin" HOME="$HOME" "$MAKE_BIN" frontend-precommit
```

This must exit 0 and print `frontend-precommit: skipped (npm not found in PATH)`.

Build the minimal PATH from `make`'s own location rather than by subtracting `npm`'s directory from `$PATH`. Subtracting is unreliable: where `npm` is provided by a shell function (nvm and similar), `command -v npm` returns the bare word `npm` instead of a path, `dirname` yields `.`, and the filter silently becomes a no-op — so the skip branch is never actually exercised. Before concluding this step passed, confirm the skip message really appeared in the output; if instead the frontend stage ran, the PATH was not stripped and this check has NOT been satisfied.

Finally confirm the blast radius is exactly the intended files (`Makefile`, `CHANGELOG.md`, `README.md` — and NOT `frontend/Makefile`):

```bash
grep -n 'command -v npm' Makefile
grep -n 'npm' frontend/Makefile
```

The first must match inside `frontend-precommit`; the second must show `frontend/Makefile` still contains its original `npm install` line (i.e. it was not edited).
</verification>
