#!/usr/bin/env just --justfile

go-mod := `go list`

# install mangal to the ~/go/bin
install:
    go install .

# generate and install mangal
full: generate && install

# build
build:
	go build .

# run tests
test:
    go test ./...

# generate assets
generate:
	go generate ./...
	./web/generate.sh

# update deps
update:
	go get -u
	go mod tidy -v

# publish
publish tag:
    GOPROXY=proxy.golang.org go list -m {{go-mod}}@{{tag}}
