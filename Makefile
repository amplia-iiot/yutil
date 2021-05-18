GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GORUN=$(GOCMD) run
BINARY=yutil
VERSION?=$(shell git describe --tags --always)
BUILD=`git rev-parse HEAD`
DOCKER_REGISTRY?= #if set it should finished by /
EXPORT_RESULT?=true # for CI please set EXPORT_RESULT to true

LDFLAGS ?= "-X 'main.version=$(VERSION)' -X 'main.commit=$(shell git rev-parse HEAD)' -X 'main.date=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")' -X 'main.builtBy=$(shell whoami)'"


RED    := $(shell tput -Txterm setaf 1)
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
CYAN   := $(shell tput -Txterm setaf 6)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build

ifndef VERBOSE
.SILENT:
endif

all: help

###############
##@ Development

deps: ## Download project dependencies
	GO111MODULE=on $(GOCMD) mod download

set-up: ## Set up development environment
	$(GOCMD) install github.com/axw/gocov/gocov@latest
	$(GOCMD) install github.com/AlekSi/gocov-xml@latest
	$(GOCMD) install github.com/jstemmer/go-junit-report@latest
	$(GOCMD) install github.com/cosmtrek/air@latest
	$(GOCMD) install github.com/spf13/cobra/cobra@latest
	$(GOCMD) install github.com/goreleaser/goreleaser@latest
	$(GOCMD) install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
	$(GOCMD) install github.com/caarlos0/svu@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

clean: ## Remove build related files
	rm -rf ./out ./tmp ./dist
	rm -f ${BINARY} ./junit-report.xml checkstyle-report.xml ./coverage.xml ./profile.cov yamllint-checkstyle.xml

test: ## Run the tests of the project
ifeq ($(EXPORT_RESULT), true)
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	GOFLAGS="-count=1" $(GOTEST) -v -race ./... $(OUTPUT_OPTIONS)

check: ## Check the source code (test & lint)
	( ( ($(MAKE) -s test) && printf '${GREEN}Tests - OK${RESET}\n' ) || printf '${RED}Tests - failed${RESET}\n' )
	( ( ($(MAKE) -s lint) && printf '${GREEN}Lint  - OK${RESET}\n' ) || printf '${RED}Lint  - failed${RESET}\n' )

watch: ## Run air to execute tests when a change is detected
	air

coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov
ifeq ($(EXPORT_RESULT), true)
	gocov convert profile.cov | gocov-xml > coverage.xml
endif

lint: ## Run linters
	golangci-lint run

#########
##@ Build

build: ## Build project for current arch
	GO111MODULE=on $(GOCMD) build -ldflags=$(LDFLAGS) -o ${BINARY}

###########
##@ Release

changelog: ## Generate changelog
	git-chglog --next-tag $(VERSION) -o CHANGELOG.md

version: VERSION=$(shell svu next || echo "v1.0.0")
version: changelog ## Generate version
	git add CHANGELOG.md
	git commit -m "chore: update changelog for $(VERSION)"
	git tag -a $(VERSION) -m "$(patsubst v%,Version %,$(VERSION))"

release-notes: ## Print release notes for current version
	git-chglog -c .chglog/release-notes.yml $(shell svu current || echo "--next-tag v1.0.0")

release: ## Build release
ifeq ($(EXPORT_RESULT), true)
	goreleaser --rm-dist --release-notes <($(MAKE) -s release-notes)
else
	YUTIL_NEXT="$(VERSION)" goreleaser --snapshot --skip-publish --rm-dist
endif

########
##@ Help

help: ## Show this help
	@awk ' \
			BEGIN { \
				FS = ":.*##" ; \
				printf "Usage:\n  make ${YELLOW}<target>${RESET}\n" \
			} \
			/^[a-zA-Z_-]+:.*?##/ { \
				printf "  ${YELLOW}%-16s${RESET}%s\n", $$1, $$2 \
			} \
			/^##@/ { \
				printf "\n${CYAN}%s:${RESET}\n", substr($$0, 5) \
			} \
		' $(MAKEFILE_LIST)
