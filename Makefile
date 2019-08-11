.PHONY: all clean

VERSION=0.1.0

all:
	@go build -ldflags \
		"-s -X github.com/iftechio/jki/pkg/cmd/version.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		-X github.com/iftechio/jki/pkg/cmd/version.gitCommit=$(shell git rev-parse HEAD) \
		-X github.com/iftechio/jki/pkg/cmd/version.version=$(VERSION)"

debug:
	@go build -ldflags \
		"-X github.com/iftechio/jki/pkg/cmd/version.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		-X github.com/iftechio/jki/pkg/cmd/version.gitCommit=$(shell git rev-parse HEAD) \
		-X github.com/iftechio/jki/pkg/cmd/version.version=$(VERSION)"

clean:
