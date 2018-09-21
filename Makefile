ARCH?=amd64
OUT_DIR?=./_output
TEMP_DIR:=$(shell mktemp -d)
DOCKER_REPO=wavefronthq
DOCKER_IMAGE=wavefront-hpa-adapter
VERSION=0.9

.PHONY: all test verify-gofmt gofmt verify

all: build
build: vendor
	CGO_ENABLED=0 GOARCH=$(ARCH) go build -a -tags netgo -o $(OUT_DIR)/$(ARCH)/wavefront-adapter ./cmd/wavefront-adapter/

# Main driver used for dev test purposes
build-query:
	CGO_ENABLED=0 GOARCH=$(ARCH) go build -a -tags netgo -o $(OUT_DIR)/$(ARCH)/wavefront-query ./cmd/wavefront-query/

# Build linux executable
build-linux: vendor
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build -a -tags netgo -o $(OUT_DIR)/$(ARCH)/wavefront-adapter-linux ./cmd/wavefront-adapter/

vendor: glide.lock
	glide install -v

test: vendor
	CGO_ENABLED=0 go test ./pkg/...

lint:
	go vet -composites=false ./...

verify-gofmt:
	./hack/gofmt-all.sh -v

gofmt:
	./hack/gofmt-all.sh

verify: verify-gofmt test

container: build-linux
	cp deploy/Dockerfile $(TEMP_DIR)
	cp $(OUT_DIR)/$(ARCH)/wavefront-adapter-linux $(TEMP_DIR)/wavefront-adapter
	cd $(TEMP_DIR)
	docker build -t $(DOCKER_REPO)/$(DOCKER_IMAGE)-$(ARCH):$(VERSION) $(TEMP_DIR)
	docker tag $(DOCKER_REPO)/$(DOCKER_IMAGE)-$(ARCH):$(VERSION) $(DOCKER_REPO)/$(DOCKER_IMAGE):latest
	rm -rf $(TEMP_DIR)

clean:
	rm -rf $(OUT_DIR)
