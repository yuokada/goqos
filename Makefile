USER=$(shell whoami)
GROUP=$(shell groups)
.PHONY: build build_linux rpm
GOQOS_VERSION=$(shell git tag -l --points-at HEAD)
GOQOS_REVISION=$(shell git rev-parse --verify HEAD)

build:
	go build \
		-ldflags="-X main.version=$(GOQOS_VERSION) -X main.revision=$(GOQOS_REVISION)"

build_linux:
	GOOS=linux GOARCH=amd64 go build

rpm:
	/bin/bash ./buildrpm.sh
