# Created by kayxhding on 2020-10-11 12:40:37
#!/usr/bin/env bash

# Fail on any error.
set -euo pipefail
# set -o xtrace

SCRIPT_PATH=$(cd `dirname $0`;pwd)

<<'COMMENT'
SCRIPT=$(readlink -f "$0")
SCRIPT_PATH=$(dirname "$SCRIPT")
echo ${SCRIPT_PATH}
COMMENT

PROTOC_FILE_DIR="$1"

function die() {
  echo 1>&2 "$*"
  exit 1
}

<<'COMMENT'
# This will place three binaries in your $GOBIN
# Make sure that your $GOBIN is in your $PATH
 go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger 
COMMENT

echo `pwd`

echo "==> Checking tools..."
#GEN_PROTO_TOOLS=(protoc protoc-gen-go protoc-gen-grpc-gateway protoc-gen-govalidators)
GEN_PROTO_TOOLS=(protoc protoc-gen-go protoc-gen-grpc-gateway)
for tool in ${GEN_PROTO_TOOLS[@]}; do
   q=$(command -v ${tool}) || die "didn't find ${tool}"
   echo 1>&2 "${tool}: ${q}"
done


echo "==> Generating proto..."
proto_headers="-I ${SCRIPT_PATH}/../../third_party"
proto_headers="${proto_headers} -I ${SCRIPT_PATH}/../../third_party/github.com/grpc-ecosystem/grpc-gateway"
source_relative_option="paths=source_relative:."
go_tag_option="--go-tag_out=${source_relative_option}"
go_grpc_option="--go-grpc_out=${source_relative_option}"
grpc_gateway_out_option="--grpc-gateway_out=logtostderr=true"
grpc_gateway_delete_option="--grpc-gateway_opt=allow_delete_body=true"

for proto in $(find ${PROTOC_FILE_DIR} -type f -name '*.proto' -print0 | xargs -0); do
  echo "Generating ${proto}"
  api_conf_yaml_base_name="$(basename ${proto} .proto).yaml"
  api_conf_yaml_dir="$(dirname ${proto})"
  api_conf_yaml="${api_conf_yaml_dir}/$api_conf_yaml_base_name"
  grpc_api_yaml_option="grpc_api_configuration=${api_conf_yaml},${source_relative_option}"

  grpc_option="${grpc_gateway_out_option},${grpc_api_yaml_option} ${grpc_gateway_delete_option}"
  protoc -I . ${proto_headers} ${go_tag_option} ${go_grpc_option} ${grpc_option} "${proto}"
  #protoc -I . ${proto_headers} --go-tag_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=logtostderr=true,grpc_api_configuration=${api_conf_yaml},paths=source_relative:. --grpc-gateway_opt=allow_delete_body=true ${f}
done
