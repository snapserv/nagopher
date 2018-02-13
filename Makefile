GO             = go
GOFMT          = go fmt
GOLINT         = golint
GOVERALLS      = goveralls
GOVERALLS_ARGS = -service=travis-ci

.PHONY: all
all: lint test build

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