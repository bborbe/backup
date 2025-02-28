# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

All notable changes to this project will be documented in this file.

## 3.5.2

- go mod update

## 3.5.1

- add cleanup already running error

## 3.5.0

- backup cleanup cron
- update golang
- update alpine

## 3.4.1

- refactor
- go mod update

## 3.4.0

- prevent concurrent backups
- go mod update

## 3.3.3

- fix backup cron on sunday
- go mod update

## 3.3.2

- go mod update
- backup hourly on sunday

## 3.3.1

- go mod update

## 3.3.0

- print rsync output by default
- go mod update

## 3.2.2

- update golang

## 3.2.1

- go mod update

## 3.2.0

- Sentry alert on failed backups

## 3.1.0

- Add status endpoint

## 3.0.0

- Complete rewrite

## 2.0.0

- Rename commands

## 1.3.1

- Cleanup backups even if one backup fails

## 1.3.0

- Cleanup backups even if one host fails
