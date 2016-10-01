install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_cleanup/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_keep/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_list/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_old/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_resume/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_status_client/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_status_server/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_rsync_client/*.go
test:
	GO15VENDOREXPERIMENT=1 go test -cover `glide novendor`
vet:
	go tool vet .
	go tool vet --shadow .
lint:
	golint -min_confidence 1 ./...
errcheck:
	errcheck -ignore '(Close|Write)' ./...
check: lint vet errcheck
runbackupstatusserver:
	backup_status_server \
	-logtostderr \
	-v=2 \
	-port=8080 \
	-rootdir=/tmp
runbackuprsyncclient:
	mkdir -p /tmp/data /tmp/backup
	echo "Hello World" > /tmp/data/README.txt
	echo "- /dev" > /tmp/excludes
	backup_rsync_client \
	-logtostderr \
	-v=4 \
	-lock=/tmp/backup_rsync_client.lock \
	-target=/tmp/backup \
	-user=bborbe \
	-host=localhost \
	-port=22 \
	-dir=/opt/go1.7.1/ \
	-exclude_from=/tmp/excludes \
	-one-time
open:
	open http://localhost:8080/
format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .
prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	glide install
	npm install
update:
	glide up
clean:
	rm -rf vendor target
