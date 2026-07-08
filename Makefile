BINARY := traceling
PKG    := github.com/chinny/traceling/cmd/traceling
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(DATE)

.PHONY: all build test vet fmt clean install

all: build

# go-fitz statically links its bundled MuPDF via cgo, so cross-compiling
# requires a matching C cross-toolchain; build natively on each target
# platform (see .github/workflows/release.yml). Building with CGO_ENABLED=0
# falls back to purego, which dlopens a system libmupdf at runtime instead.
build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) $(PKG)

install:
	go install -ldflags "$(LDFLAGS)" $(PKG)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -s -w .

clean:
	rm -rf bin dist
