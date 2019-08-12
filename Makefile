.PHONY: all clean

VERSION ?= 0.0.0
NAME ?= jki
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

all:
	@go build -ldflags \
		"-s -X github.com/bario/jki/pkg/cmd/version.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		-X github.com/bario/jki/pkg/cmd/version.gitCommit=$(shell git rev-parse HEAD) \
		-X github.com/bario/jki/pkg/cmd/version.version=$(VERSION)" -o "$(NAME)_$(VERSION)_$(OS)_$(ARCH)"

debug:
	@go build -ldflags \
		"-X github.com/bario/jki/pkg/cmd/version.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		-X github.com/bario/jki/pkg/cmd/version.gitCommit=$(shell git rev-parse HEAD) \
		-X github.com/bario/jki/pkg/cmd/version.version=$(VERSION)" -o "$(NAME)_$(VERSION)_$(OS)_$(ARCH)"

clean:
	@rm -f jki_*
