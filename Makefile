VERSION ?= $(shell git describe --tags || echo "unknown")
NAME ?= jki
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT := $(shell git rev-parse HEAD)
GO_LDFLAGS ?= "-w -s -X github.com/iftechio/jki/pkg/info.BuildDate=$(BUILD_DATE) \
	-X github.com/iftechio/jki/pkg/info.GitCommit=$(GIT_COMMIT) \
	-X github.com/iftechio/jki/pkg/info.Version=$(VERSION)"
GOBUILD = CGO_ENABLED=0 go build -trimpath -ldflags $(GO_LDFLAGS)

.PHONY: all
all:
	$(GOBUILD) -o "$(NAME)_$(VERSION)_$(OS)_$(ARCH)/$(NAME)"
	@tar czf  "$(NAME)_$(VERSION)_$(OS)_$(ARCH).tar.gz" "$(NAME)_$(VERSION)_$(OS)_$(ARCH)"

.PHONY: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64 releases
linux_amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(NAME)_$(VERSION)_$@/$(NAME)
	tar czf "$(NAME)_$(VERSION)_$@.tar.gz" "$(NAME)_$(VERSION)_$@"

linux_arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(NAME)_$(VERSION)_$@/$(NAME)
	tar czf "$(NAME)_$(VERSION)_$@.tar.gz" "$(NAME)_$(VERSION)_$@"

darwin_amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(NAME)_$(VERSION)_$@/$(NAME)
	tar czf "$(NAME)_$(VERSION)_$@.tar.gz" "$(NAME)_$(VERSION)_$@"

darwin_arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(NAME)_$(VERSION)_$@/$(NAME)
	tar czf "$(NAME)_$(VERSION)_$@.tar.gz" "$(NAME)_$(VERSION)_$@"

releases: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64

.PHONY: install
install:
	@go install -trimpath -ldflags $(GO_LDFLAGS)

.PHONY: debug
debug:
	@go build -trimpath -ldflags $(GO_LDFLAGS) -o jki

.PHONY: clean
clean:
	rm -rf jki_*
