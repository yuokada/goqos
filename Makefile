GOPATH=$(CURDIR)/.go
USER=$(shell whoami)
GROUP=$(shell groups)
.PHONY: rpm
GOQOS_VERSION=$(shell git tag -l --points-at HEAD)
GOQOS_REVISION=$(shell git rev-parse --verify HEAD)

goget:
	if [ ! -e $(GOPATH) ]; then \
		mkdir $(GOPATH) ;\
		GOPATH=$(GOPATH) go get github.com/pkg/errors ; \
	fi

build: goget
	GOPATH=$(GOPATH) go build \
		-ldflags="-X main.version=$(GOQOS_VERSION) -X main.revision=$(GOQOS_REVISION)"

build_linux: goget
	GOOS=linux GOARCH=amd64 GOPATH=$(GOPATH)  go build

rpm:
		/bin/bash ./buildrpm.sh
