# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.5.2

- add SetImagePullSecrets to PodSpecBuilder 
- go mod update

## v1.5.1

- fix DeploymentBuilder and StatefulSetBuilder

## v1.5.0

- allow define imagePullSecrets
- go mod update

## v1.4.0

- add ResourceEventHandler
- add EventHandler

## v1.3.8

- improve error message
- go mod update

## v1.3.7

- add k8s_ prefix to all go files

## v1.3.6

- allow BuildName with number at the end
- go mod update
- remove deprecated golint

## v1.3.5

- add name from pod function
- go mod update

## v1.3.4

- go mod update

## v1.3.3

- set some defaults for jobs

## v1.3.2

- allow set affinity in podSpec

## v1.3.1

- allow set restartPolicy on containerBuilder
- go mod update

## v1.3.0

- update to k8s v0.31.0
- go mod update

## v1.2.0

- add Name type
- use Name in services
- move validation

## v1.1.0

- skip ingress update if equal
- go mod update

## v1.0.0

- Initial Version
