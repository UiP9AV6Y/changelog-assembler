.DEFAULT_GOAL: default

GOPATH ?= $(shell pwd)/_workspace
GOBASE := $(firstword $(subst :, ,$(GOPATH)))
GOBIN := $(GOBASE)/bin
GOFMT ?= $(GOBIN)/goimports
GOLINT ?= $(GOBIN)/golangci-lint
CHLOG ?= $(GOBIN)/changelog-assembler
CSUM ?= sha256sum
TAR ?= tar
ZIP ?= zip
GIT ?= git
GO ?= go
INSTALL ?= install
INSTALL_PROGRAM = $(INSTALL)
INSTALL_DATA = $(INSTALL) -m 644

VERSION ?= $(shell $(GIT) describe --abbrev=0 --tags 2>/dev/null || echo 0.0.0)
COMMIT ?= $(shell $(GIT) rev-parse --short HEAD 2>/dev/null || echo HEAD)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SOURCE_DATE_EPOCH ?= $(shell $(GIT) log -1 --format='%ct' 2>/dev/null || echo 0)
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)

export SOURCE_DATE_EPOCH
export CGO_ENABLED = 0
export GO111MODULE = on
export GOARCH
export GOOS
export GOBIN

GO_MODULE := $(shell $(GO) list -m)
GO_PACKAGES := $(shell $(GO) list ./... | grep -vE '/(vendor|_workspace)')
GO_SOURCES := $(shell find . -type f -name '*.go' | grep -vE '/(vendor|_workspace)/')
GO_CMDS := $(notdir $(wildcard ./cmd/*))

PROJECT_NAME ?= $(notdir $(GO_MODULE))
DIST_NAME := $(PROJECT_NAME)_$(VERSION)
BINARY_DIST_NAME := $(DIST_NAME)-$(GOOS)-$(GOARCH)
CHANGELOGS_DIR ?= changelogs/unreleased
BUILD_DIR ?= out

GO_LDFLAGS := '-extldflags "-static"
GO_LDFLAGS += -X $(GO_MODULE)/pkg/version.version=$(VERSION)
GO_LDFLAGS += -X $(GO_MODULE)/pkg/version.commit=$(COMMIT)
GO_LDFLAGS += -X $(GO_MODULE)/pkg/version.date=$(BUILD_DATE)
GO_LDFLAGS += -w -s # Drop debugging symbols.
GO_LDFLAGS += -buildid= # Reproducable builds
GO_LDFLAGS += '

ifeq ($(GOOS),windows)
	PROGRAMS := $(addsuffix .exe,$(GO_CMDS))
	BINARY_DIST_EXT := zip
else
	PROGRAMS := $(GO_CMDS)
	BINARY_DIST_EXT := tar.gz
endif

.PHONY: default
default: all

.PHONY: all
all: lint test build

.PHONY: clean
clean:
	-$(RM) *.gz *.xz *.tar *.zip *.tgz *.txz
	-$(RM) CHANGELOG.*.md *.sha256
	-$(RM) -r $(BUILD_DIR)
	@$(GO) clean -x $(GO_PACKAGES)

.PHONY: lint
lint: $(GO_SOURCES) $(GOLINT)
	@$(GOLINT) run ./...

.PHONY: format
format: $(GO_SOURCES) $(GOFMT)
	@$(GOFMT) -w $(GO_SOURCES)

.PHONY: test
test: $(GO_SOURCES)
	@$(GO) test $(GO_PACKAGES)

.PHONY: build
build: $(addprefix $(BUILD_DIR)/,$(PROGRAMS))

.PHONY: binary-dist
binary-dist: $(BINARY_DIST_NAME).$(BINARY_DIST_EXT) $(BINARY_DIST_NAME).$(BINARY_DIST_EXT).sha256

.PHONY: release
release: CHANGELOG.md $(wildcard ./changelogs/unreleased/*.yml) $(CHLOG)
	@if $(GIT) rev-parse $(VERSION) >/dev/null 2>&1; then \
		$(error Version $(VERSION) has already been released previously); \
	fi
	$(CHLOG) release -k --file $< $(VERSION)
	$(GIT) rm changelogs/unreleased/*.yml
	$(GIT) add $<
	$(GIT) commit -m "Release of $(VERSION)"
	$(GIT) tag -a -m "Release of $(VERSION)" $(VERSION)

$(GOBIN)/%:
	# go install -v -tags tools ./...
	$(GO) list -tags tools -f '{{range .Imports}}{{ . }} {{end}}' ./tools | \
		xargs -n1 $(GO) install -v

$(BUILD_DIR)/%: $(GO_SOURCES)
	$(GO) build \
		-o $@ \
		-trimpath \
		-ldflags $(GO_LDFLAGS) \
		$(GO_MODULE)/cmd/$*

CHANGELOG.tip.md: CHANGELOG.md $(CHLOG)
	$(CHLOG) parse --file $< $(VERSION) > $@

CHANGELOG.%.md: CHANGELOG.md $(CHLOG)
	$(CHLOG) parse --file $< $* > $@

docs/%.md: $(GO_SOURCES)
	$(GO) run docs/generator.go -m -o $(dir $@)

docs/%.1: $(GO_SOURCES)
	$(GO) run docs/generator.go -o $(dir $@)

$(CHANGELOGS_DIR):
	$(INSTALL) -d $@

$(CHANGELOGS_DIR)/%.yml: $(CHANGELOGS_DIR) $(CHLOG)
	$(CHLOG) create -d $(dir $@)

%.zip: build
	$(ZIP) -jr $@ $(BUILD_DIR)

%.tar.gz %.tgz: build
	$(TAR) -czf $@ \
		-C $(BUILD_DIR) \
		$(notdir $(wildcard $(BUILD_DIR)/*))

%.tar.xz %.txz: build
	$(TAR) -cJf $@ \
		-C $(BUILD_DIR) \
		$(notdir $(wildcard $(BUILD_DIR)/*))

%.sha256: %
	$(CSUM) $< > $@
