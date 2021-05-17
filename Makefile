GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GORUN=$(GOCMD) run
BINARY=yutil
VERSION?=$(shell git describe --tags --always)
BUILD=`git rev-parse HEAD`
DOCKER_REGISTRY?= #if set it should finished by /
EXPORT_RESULT?=true # for CI please set EXPORT_RESULT to true

YELLOW := $(shell tput -Txterm setaf 3)
CYAN   := $(shell tput -Txterm setaf 6)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build

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

clean: ## Remove build related files
	rm -rf ./out ./tmp ./dist
	rm -f ${BINARY} ./junit-report.xml checkstyle-report.xml ./coverage.xml ./profile.cov yamllint-checkstyle.xml

test: ## Run the tests of the project
ifeq ($(EXPORT_RESULT), true)
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	GOFLAGS="-count=1" $(GOTEST) -v -race ./... $(OUTPUT_OPTIONS)

watch: ## Run air to execute tests when a change is detected
	air

coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov
ifeq ($(EXPORT_RESULT), true)
	gocov convert profile.cov | gocov-xml > coverage.xml
endif

#########
##@ Build

build: ## Build project for current arch
	GO111MODULE=on $(GOCMD) build ${LDFLAGS} -o ${BINARY}

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
