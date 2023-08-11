#!/usr/bin/env just --justfile

go-mod := `go list`

install:
    go install .

build:
	go build .

test:
    go test ./...

generate:
	go generate ./...
	./web/generate.sh

update:
	go get -u
	go mod tidy -v

publish tag:
    GOPROXY=proxy.golang.org go list -m {{go-mod}}@{{tag}}
