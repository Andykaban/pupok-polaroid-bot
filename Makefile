GOOS ?= linux
GOARCH ?= amd64

ifeq ($(OS), Windows_NT)
	CURRENTDIR := $(shell cmd /c cd)
else
	CURRENTDIR := $(shell pwd)
endif

BUILD_IMAGE := pupok-polaroid-bot:latest

all: image build

image:
	docker build -t $(BUILD_IMAGE) .

build: image
	docker run --rm -i \
		-v $(CURRENTDIR):/build/src/github.com/Andykaban/pupok-polaroid-bot \
		-e GO111MODULE=off \
		-e GOOS=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		-w /build/src/github.com/Andykaban/pupok-polaroid-bot \
		$(BUILD_IMAGE) go build

