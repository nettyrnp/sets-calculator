# Basic go commands
GOCMD      = go
GOBUILD    = $(GOCMD) build
GORUN      = $(GOCMD) run
GOCLEAN    = $(GOCMD) clean
GOTEST     = $(GOCMD) test
GOGET      = $(GOCMD) get
GOFMT      = $(GOCMD) fmt
GOLINT     = $(GOCMD)lint

#
BINARY = scalc
PKGS = ./...

# Texts
TEST_STRING = "TEST"

.PHONY: all help clean test lint format build run

all: clean format lint test build

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    clean              Remove binary.'
	@echo '    test               Run all tests.'
	@echo '    lint               Run golint on package sources.'
	@echo '    format             Run gofmt on package sources.'
	@echo '    build              Compile packages and dependencies.'
	@echo '    run                Compile and run Go program.'
	@echo ''
	@echo 'Targets run by default are: clean format lint test build.'
	@echo ''

clean:
	@echo "[CLEAN]"
	@$(GOCLEAN)

test:
	@echo "[$(TEST_STRING)]"
	@$(GOTEST) -v $(PKGS)

lint:
	@echo "[LINT]"
	@-$(GOLINT) -min_confidence 0.8 $(PKGS)

format:
	@echo "[FORMAT]"
	@$(GOFMT) $(PKGS)

build:
	@echo "[BUILD]"
	@$(GOBUILD) -o $(BINARY) ./*.go

run: build
	@echo "[RUN]"
	@./$(BINARY) -task='[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]'


