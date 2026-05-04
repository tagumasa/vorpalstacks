# Makefile for VorpalStacks
#
# Build targets:
#   make build          - Build webconsole frontend + Go server binary
#   make build-console  - Build webconsole frontend only (npm run build)
#   make build-tools    - Build proto_generator tool
#   make clean          - Remove build artifacts and generated files
#   make deps           - Install protoc plugins and webconsole dependencies
#
# Proto targets:
#   make proto-generate - Generate .proto files from Smithy JSON models
#   make proto          - Generate Go code from .proto files
#   make proto-ts       - Generate TypeScript code for webconsole (buf + protoc-gen-es)
#   make proto-dart     - Generate Dart code for Flutter UI
#   make proto-all      - Generate .proto, Go, and TypeScript code
#   make proto-<svc>    - Generate for a specific service (e.g., proto-s3)
#   make proto-storage  - Generate Go code from storage proto files
#   make pebble-load    - Load Smithy models into Pebble DB
#
# Server targets:
#   make run            - Run server in dev mode (no signature verification)
#   make start/stop     - Start/stop server on port 8080
#   make start-test/stop-test - Start/stop test server on port 8081
#
# Test targets:
#   make test           - Run unit tests
#   make test-cli       - Run CLI integration tests
#   make test-integration - Run all integration tests
#
# Other:
#   make fmt            - Format Go code
#   make tidy           - Tidy Go dependencies
#   make help           - Show help

PROTO_DIR := proto/aws
PB_DIR := internal/pb/aws
STORAGE_PROTO_DIR := proto/storage
STORAGE_PB_DIR := internal/pb/storage
PROTOC := protoc
GENERATOR := ./proto_generator
SMITHY_DIR := third_party/api-models-aws/models
WEBCONSOLE_DIR := webconsole

.PHONY: all build build-console build-tools clean deps
.PHONY: proto-generate proto proto-ts proto-% pebble-load
.PHONY: run start stop start-test stop-test test test-cli test-integration
.PHONY: fmt tidy help

all: build

build: build-console
	go build -o tmp/vorpalstacks .

# Build webconsole frontend
build-console:
	cd $(WEBCONSOLE_DIR) && npm run build

build-tools: build-proto-generator

build-proto-generator:
	go build -o proto_generator ./cmd/proto_generator

clean:
	rm -f proto_generator tmp/vorpalstacks
	rm -rf $(PB_DIR)
	rm -rf $(WEBCONSOLE_DIR)/dist

deps:
	@echo "Installing protoc plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	@echo "Installing webconsole dependencies..."
	cd $(WEBCONSOLE_DIR) && npm install
	@echo "Dependencies installed."

deps-dart:
	@echo "Installing Dart protoc plugin..."
	dart pub global activate protoc_plugin
	@echo "Dart protoc plugin installed. Add ~/.pub-cache/bin to PATH if needed."

# Generate proto files from Smithy JSON models
proto-generate: build-proto-generator
	@echo "Generating proto files from Smithy models..."
	@mkdir -p $(PROTO_DIR)
	@for svc in $(SERVICES); do \
		json=$$(find $(SMITHY_DIR)/$$svc/service -name "*.json" 2>/dev/null | head -1); \
		if [ -n "$$json" ]; then \
			echo "Processing $$svc..."; \
			$(GENERATOR) --smithy "$$json" --output $(PROTO_DIR) || exit 1; \
		else \
			echo "Warning: No Smithy JSON found for $$svc"; \
		fi; \
	done
	@echo "Proto files generated."

# Generate Go code from proto files
proto:
	@echo "Generating Protocol Buffer Go code..."
	@mkdir -p $(PB_DIR)
	@for proto_file in $(wildcard $(PROTO_DIR)/*.proto); do \
		echo "Compiling $$(basename $$proto_file)..."; \
		$(PROTOC) \
			--proto_path=$(PROTO_DIR) \
			--go_out=. \
			--go_opt=module=vorpalstacks \
			--go-grpc_out=. \
			--go-grpc_opt=module=vorpalstacks \
			--connect-go_out=. \
			--connect-go_opt=module=vorpalstacks \
			"$$proto_file" || exit 1; \
	done
	@echo "Protocol Buffer code generated."

# Generate TypeScript code from proto files for webconsole (buf + protoc-gen-es)
proto-ts:
	@echo "Generating TypeScript Protocol Buffer code for webconsole..."
	PATH="$(WEBCONSOLE_DIR)/node_modules/.bin:$$PATH" && buf generate --template $(WEBCONSOLE_DIR)/buf.gen.yaml
	@echo "TypeScript Protocol Buffer code generated in $(WEBCONSOLE_DIR)/src/gen/"

# Generate Dart code from proto files for Flutter UI
DART_OUT := web_ui/lib/src/generated
proto-dart:
	@echo "Generating Protocol Buffer Dart code..."
	@mkdir -p $(DART_OUT)
	@for proto_file in $(wildcard $(PROTO_DIR)/*.proto); do \
		echo "Compiling $$(basename $$proto_file) to Dart..."; \
		$(PROTOC) \
			--proto_path=$(PROTO_DIR) \
			--dart_out=grpc:$(DART_OUT) \
			"$$proto_file" || exit 1; \
	done
	@echo "Dart Protocol Buffer code generated in $(DART_OUT)"

# Generate both proto files and Go code
proto-all: proto-generate proto proto-ts

# Generate Go code from storage proto files
proto-storage:
	@echo "Generating Storage Protocol Buffer Go code..."
	@mkdir -p $(STORAGE_PB_DIR)
	@for proto_file in $(wildcard $(STORAGE_PROTO_DIR)/*.proto); do \
		echo "Compiling $$(basename $$proto_file)..."; \
		$(PROTOC) \
			--proto_path=$(STORAGE_PROTO_DIR) \
			--proto_path=third_party/protobuf \
			--go_out=. \
			--go_opt=module=vorpalstacks \
			"$$proto_file" || exit 1; \
	done
	@echo "Storage Protocol Buffer code generated."

# Load Smithy models into Pebble DB with auto-detection of implementations
pebble-load:
	@echo "Loading Smithy models into Pebble DB..."
	go run ./cmd/pebble-loader --model $(SMITHY_DIR) --data ./data --proto $(PROTO_DIR)
	@echo "Pebble DB loaded."

# Generate for a specific service
proto-%: build-proto-generator
	@echo "Generating proto for $*..."
	@json=$$(find $(SMITHY_DIR)/$*/service -name "*.json" 2>/dev/null | head -1); \
	if [ -z "$$json" ]; then \
		echo "Error: No Smithy JSON found for $*"; \
		exit 1; \
	fi; \
	$(GENERATOR) --smithy "$$json" --output $(PROTO_DIR); \
	$(PROTOC) \
		--proto_path=$(PROTO_DIR) \
		--go_out=. \
		--go_opt=module=vorpalstacks \
		--go-grpc_out=. \
		--go-grpc_opt=module=vorpalstacks \
		--connect-go_out=. \
		--connect-go_opt=module=vorpalstacks \
		$(PROTO_DIR)/$*.proto

run:
	SIGNATURE_VERIFICATION_ENABLED=false go run main.go

start:
	@bash scripts/start_server.sh 8080

stop:
	@bash scripts/stop_server.sh 8080

start-test:
	@bash scripts/start_server.sh 8081

stop-test:
	@bash scripts/stop_server.sh 8081

fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Formatting completed."

tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "Dependencies tidied."

test:
	@echo "Running unit tests..."
	go test -v ./...

test-cli:
	@bash scripts/run_cli_tests.sh

test-integration:
	@bash scripts/run_all_tests.sh

help:
	@echo "VorpalStacks Makefile"
	@echo ""
	@echo "Build targets:"
	@echo "  make build                Build webconsole + server binary"
	@echo "  make build-console        Build webconsole frontend (npm run build)"
	@echo "  make build-tools          Build all tools (proto_generator)"
	@echo "  make clean                Remove build artifacts, pb files, and webconsole/dist"
	@echo ""
	@echo "Proto targets:"
	@echo "  make deps                 Install protoc plugins"
	@echo "  make proto-generate       Generate .proto files from Smithy JSON"
	@echo "  make proto                Generate Go code from .proto files"
	@echo "  make proto-all            Generate .proto, Go code, and TypeScript code"
	@echo "  make proto-ts             Generate TypeScript code for webconsole (buf)"
	@echo "  make proto-<service>      Generate for specific service (e.g., proto-s3)"
	@echo "  make pebble-load          Load Smithy models into Pebble DB"
	@echo ""
	@echo "Server targets:"
	@echo "  make run                  Run server (dev mode, no signature verification)"
	@echo "  make start                Start server in background (port 8080)"
	@echo "  make stop                 Stop server (port 8080)"
	@echo "  make start-test           Start test server (port 8081)"
	@echo "  make stop-test            Stop test server (port 8081)"
	@echo ""
	@echo "Test targets:"
	@echo "  make test                 Run unit tests"
	@echo "  make test-cli             Run CLI integration tests"
	@echo "  make test-integration     Run all integration tests"
	@echo ""
	@echo "Other targets:"
	@echo "  make fmt                  Format code"
	@echo "  make tidy                 Tidy dependencies"
	@echo ""
	@echo "Available services (27):"
	@for svc in $(SERVICES); do echo "  - $$svc"; done
