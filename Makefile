ARCH?=amd64
OUT_DIR?=./_output
DOCKER_REPO=wavefronthq
DOCKER_IMAGE=wavefront-hpa-adapter

VERSION=0.9.3
GOLANG_VERSION?=1.13
BINARY_NAME=wavefront-adapter
GIT_COMMIT:=$(shell git rev-parse --short HEAD)

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

build: vendor
	CGO_ENABLED=0 GOARCH=$(ARCH) go build -ldflags "$(LDFLAGS)" -a -tags netgo -o $(OUT_DIR)/$(ARCH)/$(BINARY_NAME) ./cmd/wavefront-adapter/

# Build linux executable
build-linux: vendor
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build -ldflags "$(LDFLAGS)" -a -tags netgo -o $(OUT_DIR)/$(ARCH)/$(BINARY_NAME)-linux ./cmd/wavefront-adapter/

vendor: glide.lock
	glide install -v

test: vendor
	CGO_ENABLED=0 go test ./pkg/...

lint:
	go vet -composites=false ./...

container:
	# Run build in a container in order to have reproducible builds
	docker run --rm -v $(TEMP_DIR):/build -v $(REPO_DIR):/go/src/github.com/wavefronthq/wavefront-kubernetes-adapter -w /go/src/github.com/wavefronthq/wavefront-kubernetes-adapter golang:$(GOLANG_VERSION) /bin/bash -c "\
		cp /etc/ssl/certs/ca-certificates.crt /build \
		&& GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags \"$(LDFLAGS)\" -a -tags netgo -o /build/$(BINARY_NAME) github.com/wavefronthq/wavefront-kubernetes-adapter/cmd/wavefront-adapter/"

	cp deploy/Dockerfile $(TEMP_DIR)
	docker build --pull -t $(DOCKER_REPO)/$(DOCKER_IMAGE):$(VERSION) $(TEMP_DIR)
	rm -rf $(TEMP_DIR)
ifneq ($(OVERRIDE_IMAGE_NAME),)
	docker tag $(DOCKER_REPO)/$(DOCKER_IMAGE):$(VERSION) $(OVERRIDE_IMAGE_NAME)
endif

clean:
	rm -rf $(OUT_DIR)
