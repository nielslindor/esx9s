BINARY := esx9s
GO ?= go
VERSION ?= 0.1.0-dev
VERSION_PKG := github.com/nielslindor/esx9s/internal/version
LDFLAGS := -X $(VERSION_PKG).Version=$(VERSION)

.PHONY: build run test fmt clean

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/esx9s

run:
	$(GO) run ./cmd/esx9s --mock

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -rf bin
