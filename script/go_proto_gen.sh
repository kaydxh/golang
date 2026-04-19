# Created by kayxhding on 2020-10-11 12:40:37
#!/usr/bin/env bash

# exit by command return non-zero exit code
set -o errexit
# Indicate an error when it encounters an undefined variable
set -o nounset
# Fail on any error.
set -o pipefail
#set -o xtrace

# example, generate golang proto files
# bash go_proto_gen.sh -I . --proto_file_path pkg/webserver/webserver.proto --with-go

# if script called by source, $0 is the name of father script, not the name of source run script
SCRIPT_PATH=$(cd `dirname "${BASH_SOURCE[0]}"`;pwd)

<<'COMMENT'
SCRIPT=$(readlink -f "${BASH_SOURCE[0]}")
SCRIPT_PATH=$(dirname "$SCRIPT")
echo ${SCRIPT_PATH}
COMMENT

PROTOC_FILE_DIR=
PROTO_HEADERS=
# 从 Go Module 缓存中获取 grpc-gateway v1 的路径，其 third_party/googleapis 包含 google/api/annotations.proto
GRPC_GATEWAY_DIR=$(go list -m -f '{{.Dir}}' github.com/grpc-ecosystem/grpc-gateway 2>/dev/null || echo "")
# 兼容本地 third_party 目录（当 go module 不可用时作为回退）
THIRD_PARTY_DIR="${SCRIPT_PATH}/third_party"
WITH_DOC=
WITH_CPP=
WITH_GO=

function die() {
  echo 1>&2 "$*"
  exit 1
}

function getopts() {
  local -a protodirs
  while test $# -ne 0
  do
    case "$1" in
       -I|--proto_path)
             protodirs+=(
             "-I $(realpath "$2")"
            )
            shift
            ;;
       --third_party_path)
           THIRD_PARTY_DIR=$(realpath "$2")
           GRPC_GATEWAY_DIR=""
            shift
            ;;
       --with-doc)
            WITH_DOC=1
            ;;
       --proto_file_path)
            PROTOC_FILE_DIR=$(realpath "$2")
            shift
            ;;
       --with-cpp)
           WITH_CPP=1
           ;;
       --with-go)
           WITH_GO=1
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
# "-I ." need behind PROTO_HEADERS, or remove it
proto_headers="${PROTO_HEADERS} -I `pwd`"
# 优先从 Go Module 缓存获取 google/api proto 依赖，无需在项目内维护 third_party 副本
if [[ -n "${GRPC_GATEWAY_DIR}" ]]; then
  proto_headers="${proto_headers} -I ${GRPC_GATEWAY_DIR}/third_party/googleapis"
elif [[ -d "${THIRD_PARTY_DIR}/github.com/grpc-ecosystem/grpc-gateway" ]]; then
  proto_headers="${proto_headers} -I ${THIRD_PARTY_DIR}/github.com/grpc-ecosystem/grpc-gateway"
else
  die "无法找到 google/api proto 依赖，请确保 go.mod 中包含 github.com/grpc-ecosystem/grpc-gateway 或通过 --third_party_path 指定"
fi
source_relative_option="paths=source_relative:."
go_opt_option=""
go_out_option=""
go_tag_option=""
go_grpc_option=""
doc_option=""
doc_out_option=""
cpp_option=""
cpp_out_option=""
cpp_grpc_option=""
# HTTP 路由优先从 proto 里的 google.api.http 注解读取，同时支持 yaml 配置文件兜底。
grpc_gateway_out_option="--grpc-gateway_out=logtostderr=true"
grpc_gateway_delete_option="--grpc-gateway_opt=allow_delete_body=true"
grpc_gateway_option=""

for proto in $(find ${PROTOC_FILE_DIR} -type f -name '*.proto' -print0 | xargs -0); do
  echo "Generating ${proto}"
  proto_base_name="$(basename ${proto} .proto)"
  api_conf_yaml_base_name="${proto_base_name}.yaml"
  api_conf_yaml_dir="$(dirname ${proto})"
  api_conf_yaml="${api_conf_yaml_dir}/$api_conf_yaml_base_name"
  grpc_api_yaml_option=""
  grpc_gateway_option=""

  # 如果存在同名 yaml 配置文件，则使用 grpc_api_configuration 指定 HTTP 路由映射
  if [[ -f "${api_conf_yaml}" ]];then
    grpc_api_yaml_option="grpc_api_configuration=${api_conf_yaml},${source_relative_option}"
    grpc_gateway_option="${grpc_gateway_out_option},${grpc_api_yaml_option} ${grpc_gateway_delete_option}"
  else
    grpc_gateway_option="${grpc_gateway_out_option},${source_relative_option} ${grpc_gateway_delete_option}"
  fi

  if [[ "${WITH_DOC}" -eq 1 ]]; then
    # output file name
    doc_option="--doc_opt=markdown,${proto_base_name}.md"
    # docs output directory is created on demand so downstream projects do not need to pre-create it
    doc_out_dir="${SCRIPT_PATH}/../docs"
    mkdir -p "${doc_out_dir}"
    doc_out_option="--doc_out=${doc_out_dir}"
  fi

  if [[ "${WITH_CPP}" -eq 1 ]]; then
    cpp_option="--cpp_out=."
    cpp_out_option="--grpc_out=."
    cpp_grpc_option="--plugin=protoc-gen-grpc=`which grpc_cpp_plugin`"
  fi

  if [[ "${WITH_GO}" -eq 1 ]]; then
    # go_tag_option="--go-tag_out=${source_relative_option}"
    go_out_option="--go_out=."
    go_opt_option="--go_opt=paths=source_relative"
    go_grpc_option="--go-grpc_out=${source_relative_option}"
  fi

  # 打印最终 protoc 调用命令，便于排查 proto 生成问题。
  echo "+ protoc ${proto_headers} ${go_out_option} ${go_tag_option} ${go_opt_option} ${go_grpc_option} ${grpc_gateway_option} ${cpp_out_option} ${cpp_option} ${doc_option} ${doc_out_option} ${cpp_grpc_option} ${proto}"
  protoc ${proto_headers} ${go_out_option} ${go_tag_option} ${go_opt_option} ${go_grpc_option} ${grpc_gateway_option} ${cpp_out_option} ${cpp_option} ${doc_option} ${doc_out_option} ${cpp_grpc_option} "${proto}"
  #protoc -I . ${proto_headers} --go-tag_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=logtostderr=true,grpc_api_configuration=${api_conf_yaml},paths=source_relative:. --grpc-gateway_opt=allow_delete_body=true ${f}
done
