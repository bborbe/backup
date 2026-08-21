# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v3.10.6

- chore: Run `gofmt -w` last in the `format` target so golines' wrapping is normalized before the gofmt lint check
- chore: Bump `golangci-lint` to v2.13.1 and `errcheck` to v1.20.0 (fixes staticcheck `buildir` panic and missing `context` types on Go 1.27)

## v3.10.5

- chore: update Go to 1.27.0 and update dependencies

## v3.10.4

- fix: skip frontend-precommit when npm is not installed, so make precommit succeeds in Node-less CI/automation containers

## v3.10.3

- chore: update dependencies

## v3.10.2

- deps: update frontend npm dependencies to latest minor/patch (vue, vite, vitest, eslint, axios, and others); drop redundant `defineExpose` macro import that conflicts with the vue-tsc compiler macro

## v3.10.1

- chore: bump go directive and Dockerfile build stage from `1.26.5` to `1.26.6` — clears 5 stdlib advisories flagged by `osv-scanner`, all fixed in 1.26.6
- chore: drop `.osv-scanner.toml` entirely — all four suppressions (`GHSA-pxq6-2prw-chj9`, `GHSA-x744-4wpc-v9h2`, `GO-2026-4923`, `GHSA-6jwv-w5xf-7j27`) were reported unused, so the file only served to fail the scanner's stale-ignore check. Suppressions are re-added if the advisories resurface
- fix: skip the `BackupFinder List` read-permission-denied spec when running as root instead of failing. The existing guard keyed the skip on `os.Chmod` returning an error, but chmod succeeds as root — so the skip never fired, while root still read the `0000` directory and `List` returned no error, failing the assertion. Now guarded on `os.Geteuid() == 0`. Surfaced by the containerized Go-update agent, which runs as root; CI runs non-root and was unaffected
- deps: bump `golang.org/x/mod` `v0.37.0` → `v0.40.0` — clears GO-2026-6179 (transparency-log tile verification bypass in `sumdb/tlog`) and GO-2026-6180 (unrelated unauthenticated hashes accepted in `sumdb` Lookup), both flagged by `make vulncheck`. Pulled `golang.org/x/net` `v0.57.0` → `v0.58.0`, `x/text` `v0.40.0` → `v0.41.0` and `x/tools` `v0.47.0` → `v0.49.0` with it. `make precommit` was red on master before this change

## v3.10.0

- fix: stop backup-cleaner and status-handler loops promptly when context is cancelled instead of running to completion
- feat: register `/gc` admin endpoint for triggering garbage collection on demand
- refactor: replace local `boolPointer` helper with `collection.Ptr` from `github.com/bborbe/collection`

## v3.9.27

- fix: point the frontend npm registry at `verdaccio.prod.nuke.benjamin-borbe.de` — the previous `verdaccio.quant.benjamin-borbe.de` host was decommissioned when verdaccio migrated to the nuke clusters, breaking every image build with a 404 on `npm install -g npm@11.8.0`
- chore: make the npm registry overridable via the `NPM_REGISTRY` Docker build arg / Makefile variable, so a build can target the dev verdaccio or public npm without editing the Dockerfile

## v3.9.26

- fix: per-host backup and cleanup trigger endpoints now answer an already-running trigger with HTTP 409 and a structured JSON error body carrying `BACKUP_ALREADY_RUNNING` / `CLEANUP_ALREADY_RUNNING`
- fix(frontend): dashboard now shows the server-supplied reason for a failed per-host backup or cleanup trigger, rendering an already-running conflict as a warning instead of an error
- security(frontend): bump nanoid to 3.3.18 (GHSA-2v37-7h3g-55p8)

## v3.9.25

- security(frontend): bump undici to 7.29.0 (GHSA-4cwx-7wf7-3272 High, GHSA-v3r7-h72x-cjcm, GHSA-m8rv-5g2x-5cg5, GHSA-jr45-8vmc-qm54, GHSA-8xcm-r25x-g524 Moderate)

## v3.9.24

- security(frontend): bump js-yaml 4.2.0 -> 4.3.1 (GHSA-52cp-r559-cp3m, CVE-2026-59869, High)
- security(frontend): bump brace-expansion 5.0.6 -> 5.0.9 and nested 2.1.1 -> 2.1.4 (CVE-2026-13149, GHSA-3jxr-9vmj-r5cp, GHSA-mh99-v99m-4gvg, High)
- security(frontend): bump postcss 8.5.15 -> 8.5.25 (GHSA-r28c-9q8g-f849, High)

## v3.9.23

- fix(deps): bump Go toolchain to 1.26.5 (GO-2026-5856 stdlib vuln)

## v3.9.22

- Bump Go module dependencies (bborbe/*, k8s.io/*, ginkgo, gomega, etc.)
- Bump alpine base image 3.23 -> 3.24 in Dockerfile
- Rework vulncheck to surface real govulncheck errors, ignore known panic/stdlib CVEs
- Add opt-in TESTFLAGS_RACE for race-detector test runs
- Bump frontend npm dependencies

## v3.9.21

- security(frontend): override esbuild to ^0.28.1 (GHSA-gv7w-rqvm-qjhr, High)

## v3.9.20

- bump Go 1.26.3 → 1.26.4
- bump k8s.io/{api,apimachinery,client-go,apiextensions-apiserver} 0.36.0 → 0.36.1
- bump bborbe/{collection,cron,http,k8s,log,run,sentry,service,time,validation} deps
- bump ginkgo v2.29.0, gomega v1.41.0, golang.org/x/* tools
- exclude cloud.google.com/go v0.26.0

## v3.9.19

- security(frontend): bump vitest to ^4.1.0 (GHSA-5xrq-8626-4rwp, Critical)

## v3.9.18

- security(frontend): bump js-cookie to 3.0.7 (CVE / GHSA-qjx8-664m-686j, High) and brace-expansion via `npm audit fix`
- security: bump golang.org/x/net to v0.55.0 (GO-2026-5025..5030)

## v3.9.17

- security(frontend): bump axios to ^1.16.0 (CVE-2026-42035, CVE-2026-42033, CVE-2026-42040, CVE-2026-42036, CVE-2026-42044 + 3 more)

## v3.9.16

- chore: Migrate to tools.env + Makefile @version pattern; remove tools.go and obsolete replace block; update //go:generate counterfeiter directives to use @v6.12.2 hardcoded version; clean up stale CVE suppressions for removed transitive deps (blevesearch/bleve, jackc/pgx, aws-sdk-go-v2)

## v3.9.15

- Update k8s dependencies to v0.36.0 (api, apimachinery, client-go, code-generator)
- Update bborbe/* deps (collection, cron, k8s, run)
- Update getsentry/sentry-go to v0.46.1
- Add postcss v8.5.10 to frontend dependencies
- Update misc indirect deps (gosec, ginkgo, vuln, genai, protobuf)

## v3.9.14

- update bborbe/* dependencies to latest patch versions
- update k8s.io/* to v0.35.4
- update golang.org/x/* packages (crypto, net, sys, text, tools)
- update getsentry/sentry-go to v0.46.0
- add .dark-factory.log to .gitignore

## v3.9.13

- Updated project files and removed stale dark-factory prompt

## v3.9.12

- chore: add -coverprofile=coverage.out to Makefile test target for codecov upload
- chore: bump axios to 1.15.0 (fixes GHSA-3p68-rc4w-qgx5, GHSA-fvcv-3m26-pcqx)

## v3.9.11

- Update Go 1.26.2 and bump go.mod toolchain
- Update bborbe/run, bborbe/time, bborbe/parse dependencies
- Update getsentry/sentry-go to v0.45.0
- Update counterfeiter to v6.12.2
- Add CVE ignores for known unfixable vulnerabilities; improve vulncheck Makefile target

## v3.9.10

- Fix vite CVEs (GHSA-p9ff-h696-f583, GHSA-v2wj-q39q-566r, GHSA-4w7w-66w2-5vf9) via npm audit fix

## v3.9.9

- Update bborbe/* dependencies (collection, cron, errors, http, k8s, log, run, sentry, service, validation)
- Update google.golang.org/genai, genproto, grpc and related cloud deps
- Update moby/buildkit, containerd, docker/cli and container ecosystem deps
- Update charmbracelet, go-openapi, golang.org/x/* and misc indirect deps
- Add replace directives for anthropic-sdk-go, diskfs, ginkgolinter and others

## v3.9.8

- update Go dependencies

## v3.9.7

- Update dependencies to fix security vulnerabilities (go-git/v5 v5.17.2)
- Add .trivyignore for docker/docker CVEs

## v3.9.6

- allow parallel golangci-lint runners
- fix npm security vulnerabilities (brace-expansion, picomatch, yaml)
- add osv-scanner ignore for unfixable docker indirect CVEs

## v3.9.5

- remove containerd replace directives, keep only runtime-spec v1.2.1 pin

## v3.9.4

- pin runtime-spec to v1.2.1 to fix containerd/Go 1.26 build incompatibility

## v3.9.3

- exclude containerd v1.7.30 to fix Go 1.26 build incompatibility

## v3.9.2

- upgrade bborbe/* dependencies to latest versions
- upgrade getsentry/sentry-go to v0.44.1
- upgrade golangci-lint to v2.11.4
- upgrade google/osv-scanner to v2.3.5
- upgrade shoenig/go-modtool to v0.7.1

## v3.9.1

- upgrade golangci-lint from v1 to v2
- standardize Makefile: add .PHONY declarations, multiline trivy, mocks mkdir
- update .golangci.yml to v2 format
- remove deprecated golang.org/x/lint/golint from tools.go
- fix transitive dep incompatibilities (go-diskfs, go-header, anthropic-sdk)
- setup dark-factory config

## v3.9.0

- Upgrade k8s dependencies from v0.33.9 to v0.35.2
- Migrate structured-merge-diff from v4 to v6
- Add GetKind, GetAPIVersion, GetNamespace methods to TargetApplyConfiguration
- Add IsApplyConfiguration() to TargetApplyConfiguration
- Update bborbe/* libs (collection, errors, k8s, run, time) and other deps

## v3.8.10

- Update frontend npm dependencies
- Add missing ESLint 10 peer dependencies (globals, @eslint/js, vue-eslint-parser)
- fix: update flatted to 3.4.2 and undici to 7.24.4 to resolve npm vulnerabilities (GHSA-25h7-pfq9-p65f, GHSA-2mjp-6q6p-2qxm and others)
- chore: use ghcr.io/aquasecurity/trivy-db as trivy DB repository to fix mirror.gcr.io connectivity

## v3.8.9

- go mod update

## v3.8.8

- Update Go to 1.26.0 in Dockerfile
- Update Go dependencies
- Add gosec nosec annotation for trusted rsync command args

## v3.8.7

- Use go-version-file in CI workflow instead of hardcoded Go version
- Update github.com/go-git/go-git/v5 from 5.16.4 to 5.16.5 (security fix)
- Update axios from 1.13.4 to 1.13.5 (security fix)

## v3.8.6

- Update GitHub workflows to v1 plugin system
- Simplify Claude Code action with inline conditions
- Add ready_for_review and reopened triggers

## v3.8.5

- Updated Go dependencies including sentry, ginkgo, and gomega
- Updated frontend dependencies including Vue 3.5.27 and related tooling
- Updated Alpine base image from 3.22 to 3.23
- Updated npm from 11.6.0 to 11.8.0 in Docker build

## v3.8.4
- fix tar security vulnerability CVE-2026-23950 (update tar 7.5.3 → 7.5.6)
- add .mcp-* to gitignore

## v3.8.3
- update multiple dependencies to latest versions
- add k8s v0.34.2 to exclusion list

## v3.8.2
- update Go version to 1.25.4
- update multiple dependencies to latest versions

## v3.8.1
- add make frontend-precommit target to Makefile
- update Vite from v7.1.5 to v7.1.11

## v3.8.0
- add golangci-lint configuration (.golangci.yml)
- enhance Makefile with improved build tooling and quality checks
- update Go version to 1.25.2
- standardize code formatting across all Go files
- add new security scanning tools (gosec, trivy, osv-scanner)
- improve formatting tools integration (goimports-reviser, golines)
- update dependencies in go.mod

## v3.7.2
- update Go to version 1.25.1 in Dockerfile
- update npm to version 11.6.0 in Dockerfile
- update frontend dependencies and build system
- update go.mod dependencies

## v3.7.1

- implement interactive filtering system for dashboard metric cards
- add multi-select filtering with visual feedback (active/inactive states)
- implement master toggle functionality for Total hosts filter
- enhance user experience with clickable metric cards and hover effects

## v3.7.0

- add GitHub workflows for CI/CD automation
- major frontend UI improvements and enhancements
- implement target finder functionality with multiple strategies (by hostname, by name, combined, list)
- add comprehensive tests for backup cleaner and executor components
- enhance factory pattern for better dependency injection
- improve backup and cleanup handlers
- add design system documentation
- refactor frontend components with better error handling and user experience

## v3.6.1

- fix backup failure when /tmp directory doesn't exist in container
- change container base from scratch to alpine with rsync, openssh-client, and tzdata
- mount /tmp volume in Kubernetes deployment

## v3.6.0

- add UI

## v3.5.4

- go mod update
- update Dockerfile

## v3.5.3

- go mod update

## v3.5.2

- go mod update

## v3.5.1

- add cleanup already running error

## v3.5.0

- backup cleanup cron
- update golang
- update alpine

## v3.4.1

- refactor
- go mod update

## v3.4.0

- prevent concurrent backups
- go mod update

## v3.3.3

- fix backup cron on sunday
- go mod update

## v3.3.2

- go mod update
- backup hourly on sunday

## v3.3.1

- go mod update

## v3.3.0

- print rsync output by default
- go mod update

## v3.2.2

- update golang

## v3.2.1

- go mod update

## v3.2.0

- Sentry alert on failed backups

## v3.1.0

- Add status endpoint

## v3.0.0

- Complete rewrite

## v2.0.0

- Rename commands

## v1.3.1

- Cleanup backups even if one backup fails

## v1.3.0

- Cleanup backups even if one host fails
