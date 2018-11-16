.PHONY: all test build cover cert

GO ?= go
VERSION=$(shell git describe --tags --always)

all: clean build

build:
	${GO} build -ldflags "-s -w -X main.version=${VERSION}" -o marabunta cmd/marabunta/main.go;
	${GO} build -ldflags "-s -w -X main.version=${VERSION}" -o ant cmd/ant/main.go;
	# env GOOS=freebsd GOARCH=amd64 ${GO} build -ldflags "-s -w -X main.version=${VERSION}"

test:
	${GO} test -v

clean:
	@rm -rf marabunta ant *.out

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out

cert:
	openssl req -x509 -newkey rsa:2048 -sha256 -nodes -keyout server.key -out server.crt -subj "/CN=localhost" -days 3650
