# build:
# 	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.buildVersion=novu" -o health_check ./health_check.go

# Makefile template borrowed from https://sohlich.github.io/post/go_makefile/
PROJECT=one-wire-temp
PROJECT_VERSION=`cat VERSION.txt`
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVERSION=1.15
GOFLAGS="-X main.buildVersion=$(PROJECT_VERSION)"
BINARY_NAME=$(PROJECT)

all: test build
build:
	$(GOBUILD) -ldflags $(GOFLAGS) -o "$(BINARY_NAME)" -v "cmd/$(PROJECT).go"
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
# deps:
# 	GO111MODULE=on $(GOGET) github.com/docker/docker/client@master
build-all: build-linux build-arm build-darwin

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o "$(BINARY_NAME)_linux" -v "cmd/$(PROJECT).go"
build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o "$(BINARY_NAME)_arm" -v "cmd/$(PROJECT).go"
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o "$(BINARY_NAME)_darwin" -v "cmd/$(PROJECT).go"
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w "/go/src/github.com/novu/$(PROJECT)" golang:$(GOVERSION) $(GOBUILD) -o "$(BINARY_NAME)" -v "cmd/$(PROJECT).go"
