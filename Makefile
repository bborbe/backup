REGISTRY ?= docker.io
IMAGE ?= bborbe/backup
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
DIRS += $(shell find */* -maxdepth 0 -name Makefile -exec dirname "{}" \;)

default: precommit

precommit: ensure format generate test check addlicense frontend-precommit
	@echo "ready to commit"

ensure:
	go mod tidy
	go mod verify
	rm -rf vendor

format:
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	go run -mod=mod github.com/incu6us/goimports-reviser/v3 -project-name github.com/bborbe/backup -format -excludes vendor ./...
	find . -type d -name vendor -prune -o -type f -name '*.go' -print0 | xargs -0 -n 10 go run -mod=mod github.com/segmentio/golines --max-len=100 -w

generate:
	rm -rf mocks avro
	go generate -mod=mod ./...

.PHONY: test
test:
	# -race
	go test -mod=mod -p=$${GO_TEST_PARALLEL:-1} -cover $(shell go list -mod=mod ./... | grep -v /vendor/)

# TODO: enable lint
# check: lint vet errcheck vulncheck osv-scanner gosec trivy
check: vet errcheck vulncheck osv-scanner gosec trivy

vet:
	go vet -mod=mod $(shell go list -mod=mod ./... | grep -v /vendor/)

errcheck:
	go run -mod=mod github.com/kisielk/errcheck -ignore '(Close|Write|Fprint)' $(shell go list -mod=mod ./... | grep -v /vendor/ | grep -v k8s/client)

vulncheck:
	go run -mod=mod golang.org/x/vuln/cmd/govulncheck $(shell go list -mod=mod ./... | grep -v /vendor/)

osv-scanner:
	@if [ -f .osv-scanner.toml ]; then \
		echo "Using .osv-scanner.toml"; \
		go run -mod=mod github.com/google/osv-scanner/v2/cmd/osv-scanner --config .osv-scanner.toml --recursive .; \
	else \
		echo "No config found, running default scan"; \
		go run -mod=mod github.com/google/osv-scanner/v2/cmd/osv-scanner --recursive .; \
	fi

gosec:
	go run -mod=mod github.com/securego/gosec/v2/cmd/gosec -exclude=G104 ./...

trivy:
	trivy fs --scanners vuln,secret --quiet --no-progress --disable-telemetry --exit-code 1 .

lint:
	go run -mod=mod github.com/golangci/golangci-lint/cmd/golangci-lint run --config .golangci.yml ./...

addlicense:
	go run -mod=mod github.com/google/addlicense -c "Benjamin Borbe" -y $$(date +'%Y') -l bsd $$(find . -name "*.go" -not -path './vendor/*')

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
