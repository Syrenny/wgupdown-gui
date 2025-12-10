##@ General
PACKAGE_NAME := wgupdown-gui
VERSION := $(shell cat VERSION)
DEB_DIR := deb
DEB_OUTPUT := $(DEB_DIR)/$(PACKAGE_NAME)_$(VERSION)_amd64.deb
BIN_DIR := bin
BINARY := $(BIN_DIR)/$(PACKAGE_NAME)
BINARY_CTL := $(BIN_DIR)/wgupdown

.PHONY: help all build deb lint fmt clean

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

all: build deb ## Build binary and .deb package

##@ Build
build: fmt lint ## Build binary
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-X wgupdown/internal/version.Version=$(VERSION)" -o $(BINARY_CTL) ./cmd/wgupdown
	go build -ldflags="-X wgupdown/internal/version.Version=$(VERSION)" -o $(BINARY) ./cmd/$(PACKAGE_NAME)


##@ Debian package
deb: build ## Build .deb package
	# Ensure maintainer scripts have correct permissions
	chmod 755 $(DEB_DIR)/DEBIAN/postinst $(DEB_DIR)/DEBIAN/postrm
	sed -i "s/^Version:.*/Version: $(VERSION)/" $(DEB_DIR)/DEBIAN/control
	# Copy binaries to deb structure
	cp $(BINARY_CTL) $(DEB_DIR)/usr/local/bin/
	cp $(BINARY) $(DEB_DIR)/usr/local/bin/
	# Build package
	dpkg-deb --build $(DEB_DIR) $(DEB_OUTPUT)
	@echo "Built package: $(DEB_OUTPUT)"

##@ Lint & fmt
lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format Go code
	go fmt ./...

##@ Clean
clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
	rm -f $(DEB_OUTPUT)
