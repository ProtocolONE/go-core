ifndef VERBOSE
.SILENT:
endif

override ROOT_DIR = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

TAG ?= unknown
CACHE_TAG ?= unknown_cache
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0

define go_docker
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	docker run --rm \
		-v /${ROOT_DIR}:/${ROOT_DIR} \
		-v /$${GO_PATH}/pkg/mod:/$${GO_PATH}/pkg/mod \
		-w /${ROOT_DIR} \
		-e GOPATH=/$${GO_PATH}:/go \
		$${GO_IMAGE}:$${GO_IMAGE_TAG} \
		sh -c 'GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} TAG=${TAG} $(subst ",,${1})'
endef

clean: ## remove generated files, tidy vendor dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make clean") ;\
    else \
        export GO111MODULE=on ;\
        go mod tidy ;\
    	rm -rf profile.out vendor ;\
    fi;
.PHONY: clean

dev-test: test lint ## test application in dev env with race and lint
.PHONY: dev-test

dind: ## useful for windows
	if [ "${DIND}" = "1" ]; then \
		echo "Already in DIND" ;\
    else \
	    . ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	    docker run -it --rm --name dind --privileged \
            -v //var/run/docker.sock://var/run/docker.sock \
            -v /${ROOT_DIR}:/${ROOT_DIR} \
            -v /$${GO_PATH}/pkg/mod:/$${GO_PATH}/pkg/mod \
            -w /${ROOT_DIR} \
            nerufa/docker-dind:19 ;\
    fi;
.PHONY: dind

generate: go-generate ## execute all generators
.PHONY: generate

github-test: vendor test-with-coverage ## test application in GitLab CI
.PHONY: github-test

go-depends: ## view final versions that will be used in a build for all direct and indirect dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make go-depends") ;\
    else \
        cd $(ROOT_DIR) ;\
        GO111MODULE=on go list -m all ;\
    fi;
.PHONY: go-depends

go-generate: ## go generate
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make go-generate") ;\
    else \
        cd $(ROOT_DIR) ;\
        GO111MODULE=on go generate $$(go list ./...) || exit 1 ;\
        $(MAKE) vendor  ;\
    fi;
.PHONY: go-generate

go-update-all: ## view available minor and patch upgrades for all direct and indirect
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make go-update-all") ;\
    else \
        cd $(ROOT_DIR) ;\
    	GO111MODULE=on go list -u -m all ;\
    fi;
.PHONY: go-update-all

lint: ## execute linter
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make lint") ;\
    else \
        golangci-lint run --no-config --issues-exit-code=0 --deadline=30m \
          --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
          --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
          --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck ./... ;\
    fi;
.PHONY: lint

test-with-coverage: ## test application with race and total coverage
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make test-with-coverage") ;\
    else \
        GO111MODULE=on CGO_ENABLED=1 \
        go test -mod vendor -v -race -covermode atomic -coverprofile profile.out ./... || exit 1 ;\
        go tool cover -func=profile.out && rm -rf profile.out ;\
    fi;
.PHONY: test-with-coverage

test: ## test application with race
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make test") ;\
    else \
        GO111MODULE=on CGO_ENABLED=1 \
        go test -mod vendor  -race -v ./... ;\
    fi;
.PHONY: test

vendor: ## update vendor dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make vendor") ;\
    else \
        rm -rf $(ROOT_DIR)/vendor ;\
    	GO111MODULE=on \
    	go mod vendor ;\
    fi;
.PHONY: vendor

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

.DEFAULT_GOAL := help