TARGET_DIR?=$(shell pwd)
PROTO_DIR=proto
SRC_MAIN_PROTO_DIR=src/main/proto
GITHUB_GROUP=github.com/th2-net

TH2_GRPC_COMMON=th2-grpc-common
TH2_GRPC_COMMON_URL=$(GITHUB_GROUP)/$(TH2_GRPC_COMMON)@makefile

TH2_GRPC_CHECK1=th2-grpc-check1
TH2_GRPC_CHECK1_URL=$(GITHUB_GROUP)/$(TH2_GRPC_CHECK1)@makefile

TH2_GRPC_ACT_TEMPLATE=th2-grpc-act-template
TH2_GRPC_ACT_TEMPLATE_URL=$(GITHUB_GROUP)/$(TH2_GRPC_ACT_TEMPLATE)@makefile

MODULE_NAME=th2-grpc
MODULE_DIR=$(TARGET_DIR)/$(MODULE_NAME)

configure-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

clean-dependencies:
	-rm -rf $(PROTO_DIR)

prepare-dependencies: clean-dependencies
	mkdir $(PROTO_DIR)

	go get -u -t $(TH2_GRPC_COMMON_URL)
	go get -u -t $(TH2_GRPC_CHECK1_URL)
	go get -u -t $(TH2_GRPC_ACT_TEMPLATE_URL)
	sleep 1
	cp -r --no-preserve=mode,ownership $(shell go list -m -f '{{.Dir}}' $(TH2_GRPC_COMMON_URL))/$(SRC_MAIN_PROTO_DIR)/* $(PROTO_DIR)
	sleep 1
	cp -r --no-preserve=mode,ownership $(shell go list -m -f '{{.Dir}}' $(TH2_GRPC_CHECK1_URL))/$(SRC_MAIN_PROTO_DIR)/* $(PROTO_DIR)
	sleep 1
	cp -r --no-preserve=mode,ownership $(shell go list -m -f '{{.Dir}}' $(TH2_GRPC_ACT_TEMPLATE_URL))/$(SRC_MAIN_PROTO_DIR)/* $(PROTO_DIR)

clean-module:
	-rm -rf $(MODULE_DIR)
	-rm go.work go.work.sum

generate-module: clean-module prepare-dependencies configure-go
	mkdir $(MODULE_DIR)
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(MODULE_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(MODULE_DIR) --go-grpc_opt=paths=source_relative \
		$(shell find $(PROTO_DIR) -name '*.proto')
	cd $(MODULE_DIR) && go mod init $(MODULE_NAME) && go get github.com/golang/protobuf && go get google.golang.org/grpc
	cd $(TARGET_DIR) ; go work init ; go work use ./$(MODULE_NAME) 