# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

All notable changes to this project will be documented in this file.

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
