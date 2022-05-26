package openapispec

//go:generate mkdir -p ./scripts
//go:generate bash -c "curl -s -L -o ./scripts/proto-gen.sh https://raw.githubusercontent.com/kaydxh/sea/master/api/scripts/go_proto_gen.sh"
//go:generate bash scripts/proto-gen.sh -P . -I . -I ../../ -T ../../../../third_party
