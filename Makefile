REGISTRY = xxxxx.hhh.com
GROUP = xxx
NAME = file-uploader
VERSION = ${TAG_NAME}

ifeq ($(strip $(WORKSPACE)),)
    WORKSPACE := `pwd`
endif

ifeq ($(strip $(VERSION)),)
    VERSION := latest
endif

IMAGE_NAME = $(REGISTRY)/$(GROUP)/$(NAME)

.PHONY: all build push build-file-uploader run-docker-compose

all: build

build-binary:
	cd src && env GOCACHE=$(WORKSPACE)/tmp GOPATH=$(WORKSPACE)/go/pkg go build -o ../bin/file-uploader

test:
	cd src && go test -p 1 -v -cover ./...

build:
	docker build -t $(IMAGE_NAME):$(VERSION) --rm .

publish: build
	docker push $(IMAGE_NAME):$(VERSION)

build-file-uploader:
	cd src && env GOOS=linux go build -o ../bin/file-uploader

run-docker-compose:
	cd src && env GOOS=linux go build -o ../bin/file-uploader
	docker-compose up --build --remove-orphans
