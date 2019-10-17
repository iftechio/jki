VERSION ?= 0.0.0
NAME ?= jki
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

.PHONY: all
all:
	@go build -ldflags \
		"-s -X github.com/iftechio/jki/pkg/cmd/version.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		-X github.com/iftechio/jki/pkg/cmd/version.gitCommit=$(shell git rev-parse HEAD) \
		-X github.com/iftechio/jki/pkg/cmd/version.version=$(VERSION)" -o "$(NAME)_$(VERSION)_$(OS)_$(ARCH)/$(NAME)"
	@tar czf  "$(NAME)_$(VERSION)_$(OS)_$(ARCH).tar.gz" "$(NAME)_$(VERSION)_$(OS)_$(ARCH)"

.PHONY: install
install:
	go install

.PHONY: debug
debug:
	@go build -ldflags \
		"-X github.com/iftechio/jki/pkg/cmd/version.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		-X github.com/iftechio/jki/pkg/cmd/version.gitCommit=$(shell git rev-parse HEAD) \
		-X github.com/iftechio/jki/pkg/cmd/version.version=$(VERSION)" -o "$(NAME)_$(VERSION)_$(OS)_$(ARCH)"

.PHONY: clean
clean:
	rm -rf jki_*
