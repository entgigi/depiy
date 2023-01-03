VERSION ?= 0.0.1

IMAGE_TAG_BASE ?= gigiozzz/bundle-rest-app

# Image URL to use all building/pushing image targets
IMG ?= $(IMAGE_TAG_BASE):$(VERSION)

PROJECT?=github.com/entgigi/bundle-rest-app/
RELEASE?=$(VERSION)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications

build:
	go build \
			-v -ldflags="-s -w -X '${PROJECT}version.Release=${RELEASE}' \
			-X '${PROJECT}version.Commit=${COMMIT}' -X '${PROJECT}version.BuildTime=${BUILD_TIME}'" \
			-o bin/bundle-rest-app

run: build
	bin/bundle-rest-app
docker-build:
	docker build -t ${IMG} .

docker-push:
	docker push ${IMG} .