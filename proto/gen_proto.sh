#   Copyright 2020-2022 Exactpro (Exactpro Systems Limited)
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

#! /bin/bash

# Changing GOPATH
export TEMP_PATH=$GOPATH
export GOPATH=$PWD/dependencies

# Downloading required proto file dependencies: th2-grpc-common, th2-grpc-check1, th2-grpc-act-template
go get github.com/th2-net/th2-grpc-common
go get github.com/th2-net/th2-grpc-check1
go get github.com/th2-net/th2-grpc-act-template

# Moving proto files from dependencies directory to proto/proto_files directory
mkdir proto/proto_files
mv dependencies/pkg/mod/github.com/th2-net/**/src/main/proto/**/*.proto proto/proto_files
mkdir proto/proto_files/th2_grpc_common
mv proto/proto_files/common.proto proto/proto_files/th2_grpc_common/common.proto

# TEMPORARY - Adding approriate import paths and output paths
sed -i '26 i option go_package = "/proto/proto_files/th2_grpc_common";' proto/proto_files/th2_grpc_common/common.proto
sed -i '23d' proto/proto_files/th2_grpc_common/common.proto
sed -i '23d' proto/proto_files/th2_grpc_common/common.proto

sed -i '20d' proto/proto_files/act_template.proto
sed -i '20 i option go_package = "/proto";' proto/proto_files/act_template.proto
sed -i '19d' proto/proto_files/act_template.proto
sed -i '19 i import "proto/proto_files/th2_grpc_common/common.proto";' proto/proto_files/act_template.proto

sed -i '24d' proto/proto_files/check1.proto
sed -i '24 i option go_package = "/proto";' proto/proto_files/check1.proto
sed -i '20d' proto/proto_files/check1.proto
sed -i '20 i import "proto/proto_files/th2_grpc_common/common.proto";' proto/proto_files/check1.proto

# Generating go code from proto files
protoc --go_out=. proto/proto_files/*.proto

# Changing the GOPATH back
export GOPATH=TEMP_PATH
unset TEMP_PATH