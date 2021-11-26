.PHONY: build test clean prepare update docker

GO = CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/device-uart

.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_uart_go
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-uart.Version=$(VERSION)"

tidy:
	go mod tidy

build: $(MICROSERVICES)

cmd/device-uart:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd

test:
	$(GOCGO) test ./... -coverprofile=coverage.out
	$(GOCGO) vet ./...

clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

docker_device_uart_go:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/device-uart:$(GIT_SHA) \
		-t edgexfoundry/device-uart:$(VERSION)-dev \
		.
