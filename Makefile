GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
PACKAGES=$(shell go list ./... | grep -v /vendor/)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(shell which bash)
	PROTO_FILES=$(shell $(Git_Bash) -c "find ./ -name *.proto")
	TEST_DIRS=$(shell $(Git_Bash) -c "find . -name '*_test.go' | awk -F '/[^/]*$$' '{print $$1}' | sort -u")
	GO_FILES=$(shell $(Git_Bash) -c "find . -name '*.go' -type f -not -path './vendor/*'")
else
	PROTO_FILES=$(shell find ./ -name *.proto)
	TEST_DIRS=$(shell find . -name '*_test.go' | awk -F '/[^/]*$$' '{print $$1}' | sort -u)
    GO_FILES=$(shell find . -name '*.go' -type f -not -path './vendor/*')
endif

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: wire
# wire
wire:
	cd cmd/ && wire gen ./...

.PHONY: generate
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: proto
# generate proto
proto:
	protoc --proto_path=./pkg/proto \
 	       --go_out=paths=source_relative:./pkg/proto \
 	       --go-grpc_out=paths=source_relative:./pkg/proto \
	       $(PROTO_FILES)

.PHONY: test
# test
test:
	@go clean -testcache && go test -cover -v ${TEST_DIRS} -gcflags="all=-N -l"

.PHONY: vet
# vet
vet:
	@go vet --unsafeptr=false $(PACKAGES)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-s -w -X main.Version=$(VERSION)" -o ./bin/core template/cmd

.PHONY: all
# generate all
all:
	make init;
	make generate;
	#make proto;
	make vet;
	#make test;
	make build

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help