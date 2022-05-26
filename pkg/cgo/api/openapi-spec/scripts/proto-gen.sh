# Created by kayxhding on 2020-10-11 12:40:37
#!/usr/bin/env bash

# exit by command return non-zero exit code
set -o errexit
# Indicate an error when it encounters an undefined variable
set -o nounset
# Fail on any error.
set -o pipefail
# set -o xtrace

# if script called by source, $0 is the name of father script, not the name of source run script
SCRIPT_PATH=$(cd `dirname "${BASH_SOURCE[0]}"`;pwd)

<<'COMMENT'
SCRIPT=$(readlink -f "${BASH_SOURCE[0]}")
SCRIPT_PATH=$(dirname "$SCRIPT")
echo ${SCRIPT_PATH}
COMMENT

PROTOC_FILE_DIR=
PROTO_HEADERS=
# THIRD_PARTY_DIR=$(realpath "${2:-${SCRIPT_PATH}/../../third_party}")
THIRD_PARTY_DIR="${SCRIPT_PATH}/../../third_party}"
WITH_DOC=
WITH_CPP=

function die() {
  echo 1>&2 "$*"
  exit 1
}

function getopts() {
  local -a protodirs
  while test $# -ne 0
  do
    case "$1" in
       -I|--proto_path=PATH)
             protodirs+=(
             "-I $(realpath "$2")"
            )
            shift
            ;;
       -T|--third_party_path=PATH)
           THIRD_PARTY_DIR=$(realpath "$2")
            shift
            ;;
       -D|--with-doc)
            WITH_DOC=1
            ;;
       -P|--find_proto_file_path=PATH)
            PROTOC_FILE_DIR=$(realpath "$2")
            shift
            ;;
       --with-cpp)
           WITH_CPP=1
           ;;
     esac
     shift
 done

 PROTO_HEADERS="${protodirs[*]}"
 # echo "${protodirs[*]}"
}

<<'COMMENT'
# This will place three binaries in your $GOBIN
# Make sure that your $GOBIN is in your $PATH
# install protoc-gen-doc on mac=> https:
# github.com/pseudomuto/protoc-gen-doc/issues/20  (make build, cp bin/protoc-gen-doc ${GOBIN})
 go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger 
COMMENT

echo `pwd`

getopts $@

echo "==> Checking tools..."
#GEN_PROTO_TOOLS=(protoc protoc-gen-go protoc-gen-grpc-gateway protoc-gen-govalidators)
GEN_PROTO_TOOLS=(protoc protoc-gen-go protoc-gen-grpc-gateway)
for tool in ${GEN_PROTO_TOOLS[@]}; do
   q=$(command -v ${tool}) || die "didn't find ${tool}"
   echo 1>&2 "${tool}: ${q}"
done


echo "==> Generating proto..."
#proto_headers="-I ${SCRIPT_PATH}/../../third_party"
#proto_headers="-I .. -I ${THIRD_PARTY_DIR}"
# "-I ." need behind PROTO_HEADERS, or remove it
proto_headers="${PROTO_HEADERS} -I ."
proto_headers="${proto_headers} -I ${THIRD_PARTY_DIR}/github.com/grpc-ecosystem/grpc-gateway"
# proto_headers="${proto_headers} -I ${SCRIPT_PATH}/../../third_party/github.com/grpc-ecosystem/grpc-gateway"
source_relative_option="paths=source_relative:."
go_tag_option="--go-tag_out=${source_relative_option}"
go_grpc_option="--go-grpc_out=${source_relative_option}"
doc_option=""
doc_out_option=""
cpp_option=""
grpc_gateway_option=""
grpc_gateway_out_option="--grpc-gateway_out=logtostderr=true"
grpc_gateway_delete_option="--grpc-gateway_opt=allow_delete_body=true"

for proto in $(find ${PROTOC_FILE_DIR} -type f -name '*.proto' -print0 | xargs -0); do
  echo "Generating ${proto}"
  api_conf_yaml_base_name="$(basename ${proto} .proto).yaml"
  api_conf_yaml_dir="$(dirname ${proto})"
  api_conf_yaml="${api_conf_yaml_dir}/$api_conf_yaml_base_name"
  grpc_api_yaml_option=""

  if [[ -f "${api_conf_yaml}" ]];then
    grpc_api_yaml_option="grpc_api_configuration=${api_conf_yaml},${source_relative_option}"
    grpc_gateway_option="${grpc_gateway_out_option},${grpc_api_yaml_option} ${grpc_gateway_delete_option}"
  fi

  if [[ "${WITH_DOC}" -eq 1 ]]; then
    doc_option="--doc_opt=markdown,docs.md"
    doc_out_option="--doc_out=${SCRIPT_PATH}/../../doc"
  fi

  if [[ "${WITH_CPP}" -eq 1 ]]; then
    cpp_option="--cpp_out=."
  fi

  protoc ${proto_headers} ${go_tag_option} ${go_grpc_option} ${grpc_gateway_option} ${cpp_option} ${doc_option} ${doc_out_option} "${proto}"
  #protoc -I . ${proto_headers} --go-tag_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=logtostderr=true,grpc_api_configuration=${api_conf_yaml},paths=source_relative:. --grpc-gateway_opt=allow_delete_body=true ${f}
done
