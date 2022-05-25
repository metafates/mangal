BINARY=mangal

VERSION=`git describe --always --long`
BUILD=`date +%FT%T%z`
BINARY_PATH=`go list -f '{{.Target}}'`

LDFLAGS=-ldflags="-X main.version=${VERSION} -X main.build=${BUILD}"

build:
	@go build ${LDFLAGS} -o ${BINARY}

install:
	@go install ${LDFLAGS}

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

uninstall:
	@if [ -f ${BINARY_PATH} ] ; then rm ${BINARY_PATH} ; fi

.PHONY: install uninstall
