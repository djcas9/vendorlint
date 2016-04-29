export GO15VENDOREXPERIMENT=1

NAME="vendorlint"

# Get the git commit
SHA=$(shell git rev-parse --short HEAD)
BUILD_COUNT=$(shell git rev-list --count HEAD)

build: lint
	@echo "Building ${NAME}..."
	@mkdir -p bin/
	@go build \
    -o bin/${NAME}

install: lint
	@echo "Installing ${NAME}..."
	@go install

lint:
	@go vet  $$(go list ./... | grep -v /vendor/)
	@for pkg in $$(go list ./... |grep -v /vendor/ |grep -v /kuber/) ; do \
		golint -min_confidence=1 $$pkg ; \
		done
test: deps
	go list ./... | xargs -n1 go test

clean:
	rm -rf bin/

.PHONY: build
