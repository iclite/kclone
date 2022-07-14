OS ?= $(shell uname)

ARCH := $(or $(ARCH),$(ARCH),amd64)

GOLINTER ?= golangci-lint
GOFORMATTER ?= gofmt

BIN_NAME := kclone
PROG := build/$(BIN_NAME)
ifneq (,$(findstring indows,$(OS)))
	PROG=build/$(BIN_NAME).exe
	OS=windows
else ifneq (,$(findstring arwin,$(OS)))
	PROG=build/$(BIN_NAME)_darwin
	OS=darwin
else ifneq (,$(findstring inux,$(OS)))
	PROG=build/$(BIN_NAME)_linux
	OS=linux
endif

SOURCES := $(wildcard cmd/*.go) $(wildcard cmd/*/*.go)

all:
	@echo Pick one of:
	@echo $$ make $(PROG)
	@echo $$ make run
	@echo $$ make clean
	@echo
	@echo Build for different OS's and ARCH's by defining these variables. Ex:
	@echo $$ make OS=windows ARCH=amd64 build/$(BIN_NAME).exe
	@echo $$ make OS=darwin  ARCH=amd64 build/$(BIN_NAME)_darwin
	@echo $$ make OS=linux   ARCH=amd64 build/$(BIN_NAME)_linux
	@echo
	@echo Run tests
	@echo $$ make test ARGS="<test args>"
	@echo
	@echo Release a new version of $(BIN_NAME)
	@echo $$ make release
	@echo
	@echo Clean everything
	@echo $$ make clean
	@echo
	@echo Configure local environment
	@echo $$ make config
	@echo
	@echo Generate a report on code-coverage
	@echo $$ make coverage-report

$(PROG): $(SOURCES)
	@echo Building project
	GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags "-X main.version=`git describe 2>/dev/null || echo unknown`" -o $(PROG) ./cmd/

run: $(PROG)
	@./$(PROG) $(ARGS) || true

lint:
	$(GOLINTER) run --config=.golangci.yml

format:
	$(GOFORMATTER) -s -w .

format-check:
	$(GOFORMATTER) -d . | tee format-check.out
	test ! -s format-check.out

.PHONY: test release config
test: $(SOURCES)
	cd cmd && rm -rf gitworks && go test $(ARGS) ./... -coverprofile ../cover.out && rm -rf gitworks

test-all: format-check coverage-check lint

coverage-report: test
	go tool cover -html=cover.out

coverage-check: test
	@echo Checking if test coverage is above 90%
	test `go tool cover -func cover.out | tail -1 | awk '{print ($$3 + 0)*10}'` -gt 900

test-public-index:
	@./scripts/test-public-index

test-xmllint-localrepository: $(PROG)
	@./scripts/test-xmllint-localrepository

test-on-windows:
	@./scripts/test-on-windows

release: test-all $(PROG)
	@./scripts/release

clean:
	rm -rf build

