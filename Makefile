ARCH?=amd64
OUT_DIR?=./_output
DOCKER_REPO?=wavefronthq
DOCKER_IMAGE?=wavefront-hpa-adapter

VERSION?=0.9.10

BINARY_NAME=wavefront-adapter
GIT_COMMIT:=$(shell git rev-parse --short HEAD)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)


REPO_DIR:=$(shell pwd)
ifndef TEMP_DIR
TEMP_DIR:=$(shell mktemp -d /tmp/wavefront.XXXXXX)
endif

# for testing, the built image will also be tagged with this name
OVERRIDE_IMAGE_NAME?=$(ADAPTER_TEST_IMAGE)

LDFLAGS=-w -X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT)

.PHONY: all test verify-gofmt gofmt verify

all: build

fmt:
	find . -type f -name "*.go" | grep -v "./vendor*" | xargs gofmt -s -w

.PHONY: build
build:
	CGO_ENABLED=0 GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -a -tags netgo -o build/$(GOOS)/$(GOARCH)/$(BINARY_NAME) ./cmd/wavefront-adapter/


test:
	CGO_ENABLED=0 go test ./pkg/...

lint:
	go vet -composites=false ./...


BUILDER_SUFFIX=$(shell echo $(PREFIX) | cut -d '/' -f1)

.PHONY: publish
publish:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build -o fmt -o vet
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 make build -o fmt -o vet
	docker buildx create --use --node wavefront_k8s_adapter_builder_$(BUILDER_SUFFIX)
	docker buildx build --platform linux/amd64,linux/arm64 --push --pull -t $(DOCKER_REPO)/$(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_REPO)/$(DOCKER_IMAGE):latest -f Dockerfile build
clean:
	rm -rf $(OUT_DIR)