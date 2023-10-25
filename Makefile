# Copied from https://code.8labs.io/platform/cicd/templates/go/raw/master/Makefile
#
# Same as $(PACKAGES) except we get directory paths. We exclude the first line
# because it contains the top level directory which contains /vendor/
PACKAGE_DIRS = $(shell go list -f '{{ .Dir }}' ./...)
PACKAGES = $(shell go list ./...)
SOURCES = $(shell for f in $(PACKAGES); do ls $$GOPATH/src/$$f/*.go; done)

# Dependencies which are only required for testing.
# We use gomodules for this project, so I'm not sure we ever really want to install things out-of-band like this;
# they'll just get added to the go.mod file anyway.
TEST_DEPENDENCIES =

# Controls the docker image name this repository publishes.
IMAGE ?= gcr.io/mission-e/dnm-searcher
IMAGE_VERSION ?= $(shell git describe --tags --always)
NAMESPACE ?= secureworks


check: lint vet test coverage benchmark build

deps:
	go get $(TEST_DEPENDENCIES)

fmt: deps
	gofmt -w -s ./...
	goimports -w -l $(SOURCES)

lint: deps
	[[ -f .golangci.yaml || -f .golangci.yml ]] || wget -O .golangci.yaml https://dev.8labs.io:8443/repo/file/platform/cicd/templates/go?path=.golangci.default.yml
	golangci-lint run -c .golangci.yaml

vet: deps
	go vet $(PACKAGES) || (go clean $(PACKAGES); go vet $(PACKAGES))

# coverage runs the tests to collect coverage but does not attempt to look
# for race conditions.
coverage: deps $(patsubst %,%.coverage,$(PACKAGES))
	@rm -f .gocoverage/cover.txt
	gocovmerge .gocoverage/*.out > coverage.txt
	go tool cover -html=coverage.txt -o .gocoverage/index.html
	go tool cover -func=coverage.txt

%.coverage:
	@[ -d .gocoverage ] || mkdir .gocoverage
	go test -covermode=count -coverprofile=.gocoverage/$(subst /,-,$*).out $* -v

test: deps
	go test -race ./... -v

benchmark: deps
	go test -bench=. -benchmem -v ./...

build: deps
	go build .

docker: build
	docker build -t gcr.io/mission-e/dnm-searcher:latest .

install_local:
	helm install dnm-searcher charts/dnm-searcher

uninstall_local:
	helm uninstall dnm-searcher
