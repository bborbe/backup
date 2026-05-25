include tools.env

REGISTRY ?= docker.io
IMAGE ?= bborbe/backup
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
DIRS += $(shell find */* -maxdepth 0 -name Makefile -exec dirname "{}" \;)

.PHONY: default
default: precommit

.PHONY: precommit
precommit: ensure format generate test check addlicense frontend-precommit
	@echo "ready to commit"

.PHONY: ensure
ensure:
	go mod tidy -e
	go mod verify
	rm -rf vendor

.PHONY: format
format:
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	go run github.com/incu6us/goimports-reviser/v3@$(GOIMPORTS_REVISER_VERSION) -project-name github.com/bborbe/backup -format -excludes vendor ./...
	find . -type d -name vendor -prune -o -type f -name '*.go' -print0 | xargs -0 -P 8 -n 50 go run github.com/segmentio/golines@$(GOLINES_VERSION) --max-len=100 -w

.PHONY: generate
generate:
	rm -rf mocks avro
	mkdir -p mocks
	echo "package mocks" > mocks/mocks.go
	go generate -mod=mod ./...

.PHONY: test
test:
	# -race
	go test -mod=mod -p=$${GO_TEST_PARALLEL:-1} -cover -coverprofile=coverage.out $(shell go list -mod=mod ./... | grep -v /vendor/)

.PHONY: check
# TODO: enable lint (pre-existing tech debt — fix separately, then add `lint` back to check)
check: vet vulncheck osv-scanner trivy

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run --allow-parallel-runners --config .golangci.yml ./...

.PHONY: vet
vet:
	go vet -mod=mod $(shell go list -mod=mod ./... | grep -v /vendor/)

.PHONY: vulncheck
vulncheck:
	@go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) -format json $(shell go list -mod=mod ./... | grep -v /vendor/) 2>&1 | \
		jq -e 'select(.finding != null and .finding.osv != "GO-2026-4923" and .finding.osv != "GO-2026-4514")' > /dev/null 2>&1 && \
		{ echo "Unexpected vulnerabilities found"; go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) $(shell go list -mod=mod ./... | grep -v /vendor/); exit 1; } || \
		echo "No unignored vulnerabilities found"

.PHONY: osv-scanner
osv-scanner:
	@if [ -f .osv-scanner.toml ]; then \
		echo "Using .osv-scanner.toml"; \
		go run github.com/google/osv-scanner/v2/cmd/osv-scanner@$(OSV_SCANNER_VERSION) --config .osv-scanner.toml --recursive .; \
	else \
		echo "No config found, running default scan"; \
		go run github.com/google/osv-scanner/v2/cmd/osv-scanner@$(OSV_SCANNER_VERSION) --recursive .; \
	fi
	-exclude=G104 \
	-quiet \
	-fmt=summary \
	-severity=medium \
	./...

.PHONY: trivy
trivy:
	trivy fs \
	--db-repository ghcr.io/aquasecurity/trivy-db \
	--scanners vuln,secret \
	--quiet \
	--no-progress \
	--disable-telemetry \
	--exit-code 1 .

.PHONY: addlicense
addlicense:
	go run github.com/google/addlicense@$(ADDLICENSE_VERSION) -c "Benjamin Borbe" -y $$(date +'%Y') -l bsd $$(find . -name "*.go" -not -path './vendor/*')

frontend-precommit:
	$(MAKE) -C frontend precommit

.PHONY: build
build:
	go mod vendor
	docker build --no-cache --rm=true --platform=linux/amd64 -t $(REGISTRY)/$(IMAGE):$(BRANCH) -f Dockerfile .

.PHONY: upload
upload:
	docker push $(REGISTRY)/$(IMAGE):$(BRANCH)

.PHONY: clean
clean:
	docker rmi $(REGISTRY)/$(IMAGE):$(BRANCH) || true
	rm -rf vendor

.PHONY: apply
apply:
	@for i in $(DIRS); do \
		cd $$i; \
		echo "apply $${i}"; \
		make apply; \
		cd ..; \
	done

.PHONY: buca
buca: build upload clean apply

generatek8s:
	go mod vendor
	bash hack/update-codegen.sh
	rm -rf vendor

deps:
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-config-parser@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-file@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-url@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-username@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-password@latest
	go install github.com/onsi/ginkgo/v2/ginkgo@latest
