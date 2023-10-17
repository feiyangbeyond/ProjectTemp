NOW=`date "+%Y-%m-%d %H:%M:%S"`
GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
PACKAGES=$(shell go list ./... | grep -v /vendor/)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	PROTO_FILES=$(shell $(Git_Bash) -c "find . -name *.proto")
	TEST_DIRS=$(shell $(Git_Bash) -c "find . -name '*_test.go' | awk -F '/[^/]*$$' '{print $$1}' | sort -u")
else
	PROTO_FILES=$(shell find . -name *.proto)
	TEST_DIRS=$(shell find . -name '*_test.go' | awk -F '/[^/]*$$' '{print $$1}' | sort -u)
endif

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/google/wire/cmd/wire@latest
	go install mvdan.cc/gofumpt@latest

.PHONY: lint
# 代码lint静态检测
lint:
	@lint=$${LINT:-golangci-lint}; \
	if ! command -v $$lint &> /dev/null; then \
		echo "$$lint is not installed. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.53.3; \
	fi; \
	$$lint run

.PHONY: proto
# 编译proto
proto:
	protoc --proto_path=./pkg/proto \
 	       --go_out=paths=source_relative:./pkg/proto \
 	       --go-grpc_out=paths=source_relative:./pkg/proto \
	       $(PROTO_FILES)

.PHONY: fmt
# 格式化代码
fmt:
	@echo "${NOW} Starting to format ..."
	@gofumpt -l -w ./
	@echo "${NOW} Format done!"

.PHONY: fmt-check
# 格式化检查
fmt-check:
	@echo "${NOW} Starting to format check ..."
	@UNFORMATTED=$$(gofumpt -l ./); \
	if [ -z "$$UNFORMATTED" ]; then \
		echo "Format check done!"; \
	else \
		echo "Below files need to format:"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	fi;

.PHONY: wire
# 依赖注入，织入
wire:
	wire ./cmd

.PHONY: test
# 运行单元测试
test:
	@go clean -testcache && go test -cover -v ${TEST_DIRS} -gcflags="all=-N -l"

.PHONY: vet
# 代码vet静态检查
vet:
	@echo "${NOW} Starting to vet check ..."
	@go vet --unsafeptr=false $(PACKAGES)
	@echo "${NOW} Vet done!"

.PHONY: build
# 编译
build:
	@echo "${NOW} Starting to build ..."
	rm -rf ./bin && mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...
	@echo "${NOW} Build done! files in ./bin as follow:"
	@ls -l bin

.PHONY: vuln
# 代码CVE漏洞监测
vuln:
	@vuln=$${LINT:-govulncheck}; \
	if ! command -v $$vuln &> /dev/null; then \
		echo "$$vuln is not installed. Installing..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi; \
	$$vuln ./...

.PHONY: all
# 编译测试
all:
	make proto;
	make fmt-check;
	make vet;
	#make test;
	make build;

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
