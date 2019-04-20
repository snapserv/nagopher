GO             = go
GOFMT          = go fmt
GOLINT         = golint
GOVERALLS      = goveralls
GOVERALLS_ARGS = -service=travis-ci

.PHONY: all
all: lint test build

.PHONY: build
build:
	$(GO) build

.PHONY: devel-deps
devel-deps:
	$(GO) get -u golang.org/x/lint/golint
	$(GO) get -u github.com/mattn/goveralls

.PHONY: lint
lint: devel-deps
	$(GO) vet ./...
	$(GOLINT) -set_exit_status ./...

.PHONY: test
test: devel-deps
	$(GO) test -v ./...

.PHONY: coverage
coverage: devel-deps
	$(GOVERALLS) $(GOVERALLS_ARGS)
