MAKEFLAGS += --silent

ldflags := -X 'github.com/metafates/mangal/constant.BuiltAt=$(shell date -u)'
ldflags += -X 'github.com/metafates/mangal/constant.BuiltBy=$(shell whoami)@$(shell hostname)'
ldflags += -X 'github.com/metafates/mangal/constant.Revision=$(shell git rev-parse HEAD)'
ldflags += -s
ldlags  += -w

build_flags := -ldflags=${ldflags}

all: help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build the mangal binary"
	@echo "  install      Install the mangal binary"
	@echo "  test         Run the tests"
	@echo "  help         Show this help message"
	@echo ""

install:
	@go install "$(build_flags)"


build:
	@go build "$(build_flags)"

test:
	@go test -v ./...

uninstall:
	@rm -f $(shell which mangal)