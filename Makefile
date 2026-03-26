# Makefile for Vorpalstacks

BINARY := vorpalstacks
GO := go

.PHONY: all build clean test vet fmt tidy lint run help

all: build

build:
	$(GO) build -o $(BINARY) .

clean:
	rm -f $(BINARY)
	$(GO) clean

test:
	$(GO) test -v ./...

vet:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

lint: vet fmt
	@echo "Lint passed."

run: build
	SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./data ./$(BINARY)

help:
	@echo "Vorpalstacks Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build   Build the server binary"
	@echo "  test    Run unit tests"
	@echo "  vet     Run go vet"
	@echo "  fmt     Format code"
	@echo "  tidy    Tidy dependencies"
	@echo "  lint    Run vet + fmt"
	@echo "  run     Build and run server (dev mode)"
	@echo "  clean   Remove build artifacts"
	@echo "  help    Show this help"
