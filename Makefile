.PHONY: all test build cover certs

GO ?= go
VERSION=$(shell git describe --tags --always)

all: clean build
ifeq (,$(wildcard certs))
	$(MAKE) certs
endif

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

certs: SHELL:=/bin/bash
certs:
	@mkdir -p certs
	# crate CA
	openssl req -x509 -newkey rsa:2048 -sha256 -nodes -keyout certs/CA.key -out certs/CA.pem -subj "/CN=example.com" -days 365
	# create server certs
	openssl req -newkey rsa:2048 -sha256 -nodes -keyout certs/server.key -out certs/server.csr -subj "/CN=marabunta.example.com"
	openssl x509 -days 3065 -sha256 -req -in certs/server.csr -CA certs/CA.pem -CAkey certs/CA.key -set_serial 01 -out certs/server.crt -extfile <(printf "subjectAltName = DNS:localhost,DNS:marabunta.example.com,IP:127.0.0.1,IP:0.0.0.0")
	# create client certs
	openssl req -newkey rsa:2048 -sha256 -nodes -keyout certs/client.key -out certs/client.csr -subj "/CN=client.example.com"
	openssl x509 -days 3065 -sha256 -req -in certs/client.csr -CA certs/CA.pem -CAkey certs/CA.key -set_serial 01 -out certs/client.crt -extfile <(printf "subjectAltName = DNS:localhost,DNS:client.example.com,IP:127.0.0.1,IP:0.0.0.0")
