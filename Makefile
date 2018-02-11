GO     = go
GOFMT  = go fmt
GOLINT = golint

GOVERALLS      = goveralls
GOVERALLS_ARGS = -service=travis-ci

.PHONY: all
all: devel-lint devel-test build

.PHONY: build
build: deps
	$(GO) build

.PHONY: deps
deps:
	$(GO) get -d -v -t ./...

.PHONY: devel-deps
devel-deps: deps
	$(GO) get github.com/golang/lint/golint
	$(GO) get github.com/mattn/goveralls

.PHONY: devel-lint
devel-lint: devel-deps
	$(GO) vet ./...
	$(GOLINT) ./...

.PHONY: devel-test
devel-test: devel-deps
	$(GO) test -v ./...

.PHONY: devel-coverage
devel-coverage: devel-deps
	$(GOVERALLS) $(GOVERALLS_ARGS)