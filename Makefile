.PHONY: all test build clean cover certs

GO ?= go
VERSION=$(shell git describe --tags --always)

all: clean build
ifeq (,$(wildcard certs))
	$(MAKE) certs
endif

build:
	${GO} build -ldflags "-s -w -X main.version=${VERSION}" -o marabunta cmd/marabunta/main.go;
	# env GOOS=freebsd GOARCH=amd64 ${GO} build -ldflags "-s -w -X main.version=${VERSION}"

test:
	${GO} test -v

clean:
	@rm -rf marabunta ant *.out

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out

certs: SHELL:=/bin/bash
certs:
	@mkdir -p certs
	# crate CA
	openssl req -x509 -nodes -days 3650 -newkey ec:<(openssl ecparam -name prime256v1) -keyout certs/CA.pem -out certs/CA.crt -subj "/CN=marabunta"
	openssl ec -in certs/CA.pem -outform DER -out certs/CA.key

	# create server certs
	openssl req -new -nodes -newkey ec:<(openssl ecparam -name prime256v1) -keyout certs/server.key -out certs/server.csr -subj "/CN=marabunta.host"
	openssl x509 -days 3065 -sha256 -req -in certs/server.csr -CA certs/CA.crt -CAkey certs/CA.pem -set_serial 01 -out certs/server.crt
