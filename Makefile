VERSION ?= 0.0.0
NAME ?= jki
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT := $(shell git rev-parse HEAD)

.PHONY: all
all:
	@go build -ldflags \
		"-s -X github.com/iftechio/jki/pkg/info.BuildDate=$(BUILD_DATE) \
		-X github.com/iftechio/jki/pkg/info.GitCommit=$(GIT_COMMIT) \
		-X github.com/iftechio/jki/pkg/info.Version=$(VERSION)" -o "$(NAME)_$(VERSION)_$(OS)_$(ARCH)/$(NAME)"
	@tar czf  "$(NAME)_$(VERSION)_$(OS)_$(ARCH).tar.gz" "$(NAME)_$(VERSION)_$(OS)_$(ARCH)"

.PHONY: install
install:
	go install

.PHONY: debug
debug:
	@go build -ldflags \
		"-X github.com/iftechio/jki/pkg/info.BuildDate=$(BUILD_DATE) \
		-X github.com/iftechio/jki/pkg/info.GitCommit=$(GIT_COMMIT) \
		-X github.com/iftechio/jki/pkg/info.Version=$(VERSION)" -o "$(NAME)_$(VERSION)_$(OS)_$(ARCH)"

.PHONY: clean
clean:
	rm -rf jki_*
