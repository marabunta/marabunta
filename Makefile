.PHONY: all test build cover certs

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
	openssl ecparam -genkey -name prime256v1 -out certs/CA.key
	openssl req -x509 -new -SHA256 -nodes -key certs/CA.key -out certs/CA.crt -subj "/CN=marabunta.host" -days 3650
	# create server cert and sign it with the CA
	openssl ecparam -genkey -name prime256v1 -out certs/server.key
	openssl req -new -SHA256 -key certs/server.key -nodes -out certs/server.csr -subj "/CN=marabunta.host"
	openssl x509 -days 3065 -sha256 -req -in certs/server.csr -CA certs/CA.crt -CAkey certs/CA.key -set_serial 01 -out certs/server.crt
	# create client cert and sign it with the CA
	openssl ecparam -genkey -name prime256v1 -out certs/client.key
	openssl req -new -SHA256 -key certs/server.key -nodes -out certs/client.csr -subj "/CN=ant.marabunta.host"
	openssl x509 -days 3065 -sha256 -req -in certs/client.csr -CA certs/CA.crt -CAkey certs/CA.key -set_serial 01 -out certs/client.crt
