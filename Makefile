# ------------------------------------------------------------------------------
# Makefile definition for the Terraform Utils Project
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# Makefile Macros
# ------------------------------------------------------------------------------

VERSION:=$(CI_COMMIT_TAG)

# List of files for use in the verious go commands. GOFMT_FILES is used when
# running format checks and formatting the codebase with "go fmt", GOVET_FILES
# contains the package names to issue "go vet" against, GOTEST_FILES lists
# the packages used for testing with "go test", and GODOC_FILES contains the
# packages used with "godoc" to generate documentation for the project.
#
# Files that are part of the vendor directory are not included as part of the
# format check, vetting, testing, etc.
GOFMT_FILES := $(shell find . -name '*.go' | grep -v vendor)
GOVET_FILES := $(shell go list ./... | grep -v vendor)
GOTEST_FILES := $(GOVET_FILES)
GODOC_FILES := $(GOVET_FILES)
GOBUILD_FILES := $(GOVET_FILES)
GOINSTALL_FILES := $(GOVET_FILES)

# Root directory for (auto)generated project documentation
#
# NOTE(ALL): DOCS_DIR should be kept in sync with docs_dir in mkdocs.yml
DOCS_DIR := docs
# Directory to output static HTML generated from the `godoc` tool
GODOC_OUT_DIR := $(DOCS_DIR)/godoc
# wget options.  wget is used in the 'doc' target to generate static site
# documentation for the project.
#
# -r, --recursive
#   Turn on recusrive retrieving. The default maximum depth is 15
# -np, --no-parent
#   Disallow the retrieval of the links that refer to the hierarchy above the
#   beginning directory
# -nH, --no-host-directories
#   Disable generation of host-prefixed directories. By default, with the "-r"
#   option, will create a structure of directories beginning with the hostname.
#   This disables this behavior
# -nv, --no-verbose
#   Turn off verbose without being completely quiet, error messages and basic
#   information get printed
# -N, --timestamping
#   Turn on timestamping (only download files that do not already exist
#   locally or the remote has a newer version)
# -E, --adjust-extension
#   Appends the correct file suffix to the local filename if the downloaded
#   file does not already have it
# -p, --page-requisites
#   Download all the files that are necessary to properly display a given
#   HTML page. This includes things like inlined images, sounds, and
#   referenced stylesheets
# -k, --convert-links
#   After the download is complete, convert the links in the document to make
#   them suitable for local viewing
# -e, --execute
#   Execute a command as if it were part of .wgetrc. The commands are AFTER
#   comamnds in wgetrc, thus taking precedence
# -P, --directory-prefix
#   Set the directory prefix. All files and sub-directories will be saved to
#   this location to form the top of the retreival tree
WGET_OPTIONS := -r -np -nH -nv -N -E -p -k -e 'robots=off' -P "$(GODOC_OUT_DIR)"

ifndef VERSION
	VERSION:=$(shell git describe --always 2>/dev/null)
endif

HOST_OS := $(shell go env GOHOSTOS)
HOST_ARCH := $(shell go env GOHOSTARCH)
GOPKG := $${GOPATH}/pkg/$(HOST_OS)_$(HOST_ARCH)/github.com/wayfair/terraform-provider-utils

# ------------------------------------------------------------------------------
# Makefile Targets
# ------------------------------------------------------------------------------

# All of the Makefile targets are not the names of files and therefore are
# phony targets
.PHONY: all build clean clean-godoc default ensure format formatcheck godoc install test vet

# Default target - build the project
# Use the special built-in target name and human conventions
.DEFAULT: build
default: build
all: build

# Compiles the codebase into the target binary.  The binary will be in the
# output directory
build: formatcheck
	@echo "Building packages..."
	@go build -v $(GOBUILD_FILES)

# Removes the compiled binaries (if they exist), log files, and documentation
clean: clean-godoc
	@echo 'Cleaning archive files...'
	@rm -rf "$(GOPKG)" 2>/dev/null || true
	@echo 'Cleaning log files...'
	@find . -type f -name '*.log' -delete 2>/dev/null || true

# Removes all godoc files
clean-godoc:
	@echo 'Cleaning godoc files...'
	@if [ -d $(GODOC_OUT_DIR) ]; then \
		rm -rf $(GODOC_OUT_DIR); \
	fi;

# Ensure the project dependencies are in sync and up-to-date.  This will read
# the dependencies and constraints in the Gopkg.toml file and update the
# /vendor directory and Gopkg.lock file to reflect the constraints.
ensure:
	@echo 'Ensuring project dependencies are up to date...'
	@dep ensure

# Runs "go fmt" on the codebase and writes the output back to the source files
format:
	@echo 'Formatting codebase...'
	@gofmt -w $(GOFMT_FILES)

# Runs "go fmt" on the codebase, but unlike the "format" target it does not
# write the results back to the source files.  It captures the output of the
# files that violate the formatting and displays them to the console.
formatcheck:
	@echo 'Validating format of codebase...'
	@badFiles=$$(gofmt -l $(GOFMT_FILES)); \
	if [ -n "$$badFiles" ]; \
	then \
		echo 'The following files violate go formatting:'; \
		echo ''; \
		echo "$$badFiles"; \
		echo ''; \
		echo 'Run "make format" to reformat the code.'; \
		exit 1; \
	else \
		echo 'All files pass format check.'; \
		exit 0; \
	fi

# Generates godoc for the project and saves the static assets to GODOC_OUT_DIR
# through recursive downloads with wget.  The godoc can be read locally through
# a web viewport by browsing the filesystem.  The documentation is also used in
# conjunction with the documentation stage of a CI/CD pipeline for creating
# project documentation.
godoc:
	@echo "Generating godoc to $(GODOC_OUT_DIR)..."
	@if [ ! -d "$(GODOC_OUT_DIR)" ]; then \
		echo "Creating $(GODOC_OUT_DIR)"; \
		mkdir -p "$(GODOC_OUT_DIR)"; \
	fi;
	@godocAddr="127.0.0.1:8000"; \
	godoc -http="$${godocAddr}" & \
	godocPID="$$!"; \
	echo "godoc PID: [$${godocPID}]"; \
	echo "Sleeping while godoc initializes..."; \
	sleep 5; \
	echo "Downloading pages..."; \
	echo ''; \
	for pkg in $$(go list ./...); do \
		echo "Downloading $${pkg}"; \
		wget $(WGET_OPTIONS) "http://$${godocAddr}/pkg/$${pkg}"; \
	done; \
	echo ''; \
	echo 'done.'; \
	echo "Killing godoc process [$${godocPID}]"; \
	kill "$${godocPID}";

# Compiles the codebase and moves the target binary into the terraform plugins
# directory for use
install: formatcheck
	@echo "Installing packages to $(GOPKG)..."
	@go install -v $(GOINSTALL_FILES)

# Runs the go unit and integration tests on the codebase
test:
	@echo 'Running unit tests...'
	@testOutput=$$(go test $(GOTEST_FILES) 2>&1); \
	exitStatus=$$?; \
	if [ "$$exitStatus" -eq 0 ]; \
	then \
		echo 'All files pass tests'; \
		exit 0; \
	else \
		echo 'Codebase failed tests:'; \
		echo ''; \
		echo "$${testOutput}"; \
		echo ''; \
		exit 1; \
	fi

# Runs "go vet" on the codebase and writes any errors or suspicious program
# behavior to the console
vet:
	@echo "Vetting the codebase for suspicious constructs..."
	@vetOutput=$$(go vet $(GOVET_FILES) 2>&1); \
	exitStatus=$$?; \
	if [ "$$exitStatus" -eq 0 ]; \
	then \
		echo 'All files pass vet check'; \
		exit 0; \
	else \
		echo 'Codebase failed vet check:'; \
		echo ''; \
		echo "$${vetOutput}"; \
		echo ''; \
		exit 1; \
	fi

