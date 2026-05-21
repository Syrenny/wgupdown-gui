##@ General
PACKAGE_NAME := wgupdown-gui
VERSION := $(shell cat VERSION)
DEB_DIR := deb
DEB_OUTPUT := $(DEB_DIR)/$(PACKAGE_NAME)_$(VERSION)_amd64.deb
BIN_DIR := bin
BINARY := $(BIN_DIR)/$(PACKAGE_NAME)
HELPER_BINARY := $(BIN_DIR)/wgupdown
GO ?= $(shell command -v go 2>/dev/null || echo /usr/local/go/bin/go)
GOFMT ?= $(shell command -v gofmt 2>/dev/null || echo /usr/local/go/bin/gofmt)
GOLANGCI_LINT ?= $(shell command -v golangci-lint 2>/dev/null)

.PHONY: help all build deb lint fmt clean

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

all: build deb ## Build binary and .deb package

##@ Build
build: fmt lint ## Build binary
	@mkdir -p $(BIN_DIR)
	$(GO) build -ldflags="-X github.com/Syrenny/wgupdown-gui/internal/version.Version=$(VERSION)" -o $(HELPER_BINARY) ./cmd/wgupdown
	$(GO) build -ldflags="-X github.com/Syrenny/wgupdown-gui/internal/version.Version=$(VERSION)" -o $(BINARY) ./cmd/$(PACKAGE_NAME)


##@ Debian package
deb: build ## Build .deb package
	# Ensure maintainer scripts have correct permissions
	chmod 755 $(DEB_DIR)/DEBIAN/postinst $(DEB_DIR)/DEBIAN/postrm
	find $(DEB_DIR)/etc -type d -exec chmod 755 {} +
	find $(DEB_DIR)/etc -type f -exec chmod 644 {} +
	find $(DEB_DIR)/etc/sudoers.d -type f -exec chmod 440 {} +
	sed -i "s/^Version:.*/Version: $(VERSION)/" $(DEB_DIR)/DEBIAN/control
	# Copy binaries to deb structure
	rm -f $(DEB_DIR)/usr/local/bin/wgupdown $(DEB_DIR)/usr/local/bin/wgupdown-gui
	cp $(HELPER_BINARY) $(DEB_DIR)/usr/local/bin/
	cp $(BINARY) $(DEB_DIR)/usr/local/bin/
	rm -f $(DEB_DIR)/usr/local/bin/.gitkeep
	# Build package
	dpkg-deb --root-owner-group --build $(DEB_DIR) $(DEB_OUTPUT)
	@echo "Built package: $(DEB_OUTPUT)"

##@ Lint & fmt
lint: ## Run golangci-lint
ifneq ($(strip $(GOLANGCI_LINT)),)
	$(GOLANGCI_LINT) run ./...
else
	@echo "golangci-lint not found; skipping lint"
endif

fmt: ## Format Go code
	$(GOFMT) -w $$(find . -name '*.go' -not -path './vendor/*')

##@ Clean
clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
	rm -f $(DEB_OUTPUT)
