REGISTRY ?= docker.io
IMAGE ?= bborbe/backup
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
DIRS += $(shell find */* -maxdepth 0 -name Makefile -exec dirname "{}" \;)

build:
	docker build --no-cache --rm=true --platform=linux/amd64 -t $(REGISTRY)/$(IMAGE):$(BRANCH) -f Dockerfile .

upload:
	docker push $(REGISTRY)/$(IMAGE):$(BRANCH)

clean:
	docker rmi $(REGISTRY)/$(IMAGE):$(BRANCH) || true

precommit: ensure format generate test check
	@echo "ready to commit"

ensure:
	go mod tidy
	go mod verify
	go mod vendor

format:
	go run -mod=vendor github.com/incu6us/goimports-reviser/v3 -project-name github.com/bborbe/backup -format -excludes vendor ./...

generate:
	rm -rf mocks avro
	go generate -mod=vendor ./...

test:
	go test -mod=vendor -p=$${GO_TEST_PARALLEL:-1} -cover -race $(shell go list -mod=vendor ./... | grep -v /vendor/)

check: vet errcheck vulncheck

vet:
	go vet -mod=vendor $(shell go list -mod=vendor ./... | grep -v /vendor/)

errcheck:
	go run -mod=vendor github.com/kisielk/errcheck -ignore '(Close|Write|Fprint)' $(shell go list -mod=vendor ./... | grep -v /vendor/)

apply:
	@for i in $(DIRS); do \
		cd $$i; \
		echo "apply $${i}"; \
		make apply; \
		cd ..; \
	done

vulncheck:
	go run -mod=vendor golang.org/x/vuln/cmd/govulncheck $(shell go list -mod=vendor ./... | grep -v /vendor/)

generatek8s:
	rm -rf k8s/client ${GOPATH}/src/github.com/bborbe/backup
	chmod a+x vendor/k8s.io/code-generator/*.sh
	bash vendor/k8s.io/code-generator/generate-groups.sh applyconfiguration,client,deepcopy,informer,lister \
	github.com/bborbe/backup/k8s/client github.com/bborbe/backup/k8s/apis \
	backup.benjamin-borbe.de:v1
	cp -R ${GOPATH}/src/github.com/bborbe/backup/k8s .

deps:
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-config-parser@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-file@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-url@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-username@latest
	go install github.com/bborbe/teamvault-utils/cmd/teamvault-password@latest
	go install github.com/onsi/ginkgo/v2/ginkgo@latest
