GOPATH=$(CURDIR)/.go
USER=$(shell whoami)
GROUP=$(shell groups)
.PHONY: rpm

goget:
	mkdir $(GOPATH)
	@GOPATH=$(GOPATH) go get github.com/pkg/errors

build: goget
	GOPATH=$(GOPATH) go build

build_linux: goget
	GOOS=linux GOARCH=amd64 GOPATH=$(GOPATH)  go build

rpm:
		/bin/bash ./buildrpm.sh
