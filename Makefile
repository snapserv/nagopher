GO             = go
GOFMT          = go fmt
REVIVE         = revive
GOVERALLS      = goveralls
GOVERALLS_ARGS = -service=travis-ci

.PHONY: all
all: lint test build

.PHONY: build
build:
	$(GO) build

.PHONY: devel-deps
devel-deps:
	$(GO) get -u github.com/mgechev/revive
	$(GO) get -u github.com/mattn/goveralls

.PHONY: lint
lint: devel-deps
	$(GO) vet ./...
	$(REVIVE) $(addprefix -exclude ,$(wildcard optional_*.go)) ./...

.PHONY: test
test: devel-deps
	$(GO) test -v ./...

.PHONY: coverage
coverage: devel-deps
	$(GOVERALLS) $(GOVERALLS_ARGS)
