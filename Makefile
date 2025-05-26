# Paths
OPENAPI_FILE = openapi.yaml
API_DIR = internal/api
TYPES_FILE = $(API_DIR)/types.gen.go
SERVER_FILE = $(API_DIR)/server.gen.go

# Tools
OAPI_GEN = oapi-codegen
AIR_BIN = $(GOPATH)/bin/air

# Target: Install all tools
.PHONY: install-tools
install-tools:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	go install github.com/cosmtrek/air@latest

# Target: Generate API code from OpenAPI spec
.PHONY: generate
generate:
	@echo "🔄 Generating Go types and server interface from OpenAPI spec..."
	$(OAPI_GEN) -generate types -o $(TYPES_FILE) -package api $(OPENAPI_FILE)
	$(OAPI_GEN) -generate server -o $(SERVER_FILE) -package api $(OPENAPI_FILE)
	@echo "✅ Code generation complete."

# Target: Run Go server directly
.PHONY: run
run:
	go run main.go

# Target: Dev server with hot reload using air
.PHONY: dev
dev:
	$(AIR_BIN)

# Target: Clean generated code
.PHONY: clean
clean:
	rm -f $(TYPES_FILE) $(SERVER_FILE)
	@echo "🧹 Cleaned generated files."

# Target: Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Target: Build
.PHONY: build
build:
	go build -o bin/git-analyzer main.go
