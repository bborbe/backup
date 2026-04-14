# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

All notable changes to this project will be documented in this file.

## Unreleased

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
