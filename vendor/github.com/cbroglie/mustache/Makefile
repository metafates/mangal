.PHONY: all
all: bin/mustache

.PHONY: clean
clean:
	rm -rf bin

.PHONY: test
test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./...

SOURCES     := $(shell find . -name '*.go')
BUILD_FLAGS ?= -v
LDFLAGS     ?= -w -s

bin/%: $(SOURCES)
	CGO_ENABLED=0 go build -o $@ $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" ./cmd/$(@F)
