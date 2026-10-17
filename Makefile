PROTO_SRC := api
PROTO_OUT := pkg/api
PROTO_FILES := $(shell find $(PROTO_SRC) -name '*.proto')

LOCAL_BIN := $(CURDIR)/bin
VENDOR_PROTO := $(CURDIR)/vendor.protogen

PROTOC_GEN_GO_VERSION := v1.36.11
PROTOC_GEN_GO_GRPC_VERSION := v1.6.2
PROTOC_GEN_GRPC_GATEWAY_VERSION := v2.31.0
GOOGLEAPIS_COMMIT := b29d17037e889281b718bfffdeea534aed00d36f
GOOGLEAPIS_URL := https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_COMMIT)

.PHONY: build run proto proto-deps

check:
	go build -o /dev/null ./cmd/factory_analyst/main.go

build:
	go build ./cmd/factory_analyst/main.go

run: build
	./main

proto-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@$(PROTOC_GEN_GRPC_GATEWAY_VERSION)
	mkdir -p $(VENDOR_PROTO)/google/api
	curl -sSfL -o $(VENDOR_PROTO)/google/api/annotations.proto $(GOOGLEAPIS_URL)/google/api/annotations.proto
	curl -sSfL -o $(VENDOR_PROTO)/google/api/http.proto $(GOOGLEAPIS_URL)/google/api/http.proto

proto:
	rm -rf $(PROTO_OUT)
	mkdir -p $(PROTO_OUT)
	PATH=$(LOCAL_BIN):$$PATH protoc -I $(PROTO_SRC) -I $(VENDOR_PROTO) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PROTO_OUT) --grpc-gateway_opt=paths=source_relative \
		$(PROTO_FILES)
