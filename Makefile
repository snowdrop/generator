PROJECT     := github.com/snowdrop/generator
GITCOMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_FLAGS := -ldflags="-w -X main.GITCOMMIT=$(GITCOMMIT) -X main.VERSION=$(VERSION)"
GO          ?= go
GOFMT       ?= $(GO)fmt

VFSGENDEV   := $(GOPATH)/bin/vfsgendev
VFSGENDEV_SRC := $(GOPATH)/src/github.com/shurcooL/vfsgen
PREFIX      ?= $(shell pwd)

all: clean build

clean:
	@echo "> Remove build dir"
	@rm -rf ./build

build: clean
	@echo "> Build go application"
	go build ${BUILD_FLAGS} -o generator main.go

cross: clean
	gox -osarch="darwin/amd64 linux/amd64" -output="dist/bin/{{.OS}}-{{.Arch}}/generator" $(BUILD_FLAGS)

image:
	imagebuilder -t spring-boot-generator:$(VERSION) -f docker/Dockerfile_generator .

	$(eval TAG_ID = $(shell docker images -q spring-boot-generator:$(VERSION)))

	docker tag $(TAG_ID) quay.io/snowdrop/spring-boot-generator:$(VERSION)
	# docker login quai.io
	docker push quay.io/snowdrop/spring-boot-generator:$(VERSION)

prepare-release: cross
	./scripts/prepare_release.sh

upload: prepare-release
	./scripts/upload_assets.sh

assets: $(VFSGENDEV)
	@echo ">> writing assets"
	cd $(PREFIX)/pkg/template && go generate

$(VFSGENDEV): $(VFSGENDEV_SRC)
	go get -u github.com/shurcooL/vfsgen/cmd/vfsgendev

$(VFSGENDEV_SRC):
	go get -u github.com/shurcooL/vfsgen

gofmt:
	@echo ">> checking code style"
	@fmtRes=$$($(GOFMT) -d $$(find . -path ./vendor -prune -o -name '*.go' -print)); \
	if [ -n "$${fmtRes}" ]; then \
		echo "gofmt checking failed!"; echo "$${fmtRes}"; echo; \
		echo "Please ensure you are using $$($(GO) version) for formatting code."; \
		exit 1; \
	fi

version:
	@echo $(VERSION)