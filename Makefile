CURRENT_DIR         := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN                 := metrics
PKG                 := github.com/radar-go/$(BIN)
VERSION             := $(shell git describe --tags --always)
GIT_BRANCH          := $(shell git name-rev --name-only HEAD | sed "s/~.*//")
GIT_COMMIT          := $(shell git rev-parse HEAD)
BUILD_CREATOR       := $(shell git log --format=format:%ae | head -n 1)
REPORTS_DIR         := $(CURRENT_DIR)/reports
GOMETALINTER_CONFIG := $(CURRENT_DIR)/.gometalinter.json
PORT                := 6000

all: build

build-dirs:
	@mkdir -p bin
	@if [ ! -d vendor ]; then $(MAKE) --no-print-directory update-vendors; fi

update-vendors:
	@dep ensure -v

build: $(CURRENT_DIR)/bin/$(BIN)

$(CURRENT_DIR)/bin/$(BIN): build-dirs
	@echo "building: $@"
	@BIN=$(BIN) PKG=$(PKG) VERSION=$(VERSION) $(CURRENT_DIR)/scripts/build.sh

start:
	@if [ ! -f $(CURRENT_DIR)/.$(BIN).pid ]; then \
		echo -n "Starting $(BIN)"; \
		$(CURRENT_DIR)/bin)/$(BIN) -logtostderr=true 1> $(BIN).log < /dev/null 2>&1 & \
		echo $$! > $(CURRENT_DIR)/.$(BIN).pid ; \
		while ! curl localhost:$(PORT)/healthcheck > /dev/null 2>&1; do \
			/bin/sleep 1; \
			echo -n "."; \
		done; \
		echo; \
	fi

stop:
	@if [ -f $(CURRENT_DIR)/.$(BIN).pid ]; then \
		echo -n "Stopping $(BIN)"; \
		kill -s 15 `cat $(CURRENT_DIR)/.$(BIN).pid`; \
		while curl localhost:$(PORT)/healthcheck > /dev/null 2>&1; do \
			/bin/sleep 1; \
			echo -n "."; \
		done; \
		rm -f $(CURRENT_DIR)/.$(BIN).pid; \
		echo; \
	fi

restart:
	@$(MAKE) stop
	@$(MAKE) start

tests:
	@if [ ! -d vendor ]; then $(MAKE) --no-print-directory update-vendors; fi
	@REPORTS_DIR=$(REPORTS_DIR) $(CURRENT_DIR)/scripts/tests.sh

lint:
	@REPORTS_DIR=$(REPORTS_DIR) GOMETALINTER_CONFIG=$(GOMETALINTER_CONFIG) $(CURRENT_DIR)/scripts/lint.sh

coverage:
	@REPORTS_DIR=$(REPORTS_DIR) $(CURRENT_DIR)/scripts/coverage.sh

clean:
	@rm -fr bin
	@rm -fr vendor
	@rm -fr Gopkg.lock
	@rm -fr reports
	@rm -fr profile.cov

GOLDEN_PKG ?= github.com/radar-go/$(BIN)
update-golden-files:
	@go test $(GOLDEN_PKG) -update

