install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_cleanup/backup_cleanup.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_keep/backup_keep.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_list/backup_list.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_old/backup_old.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_resume/backup_resume.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_status_client/backup_status_client.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_status_server/backup_status_server.go
test:
	GO15VENDOREXPERIMENT=1 go test `glide novendor`
check:
	golint ./...
	errcheck ./...
run:
	backup_status_server \
	-loglevel=DEBUG \
	-port=8080 \
	-rootdir=/tmp
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
