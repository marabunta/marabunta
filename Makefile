.PHONY: all deps test build cover

GO ?= go
VERSION=$(shell git describe --tags --always)

all: build

deps:
	${GO} get github.com/nbari/violetear
	${GO} get github.com/nbari/violetear/middleware

build: deps
	# ${GO} build -ldflags "-X main.version=${VERSION}"
	env GOOS=freebsd GOARCH=amd64 ${GO} build -ldflags "-s -w -X main.version=${VERSION}"

test: deps
	${GO} test -v

clean:
	@rm -rf comments

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
