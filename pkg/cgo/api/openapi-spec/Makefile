MAKEFILE_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
TARGET := $(shell basename ${MAKEFILE_DIR})
PROJECT_ROOT_PATH := ${MAKEFILE_DIR}/../../
$(info ${PROJECT_ROOT_PATH})

.PHONY: generate 
generate:
	@echo "  >  starting generate protocol buffers for target ${TARGET} in golang"
	@bash -c "go generate  ../api/openapi-spec/proto.gen.go"
	@echo "  >  compiling protocol buffers by cpp"
	@if [[ -f "${PROJECT_ROOT_PATH}/build" ]]; then rm -Rv "${PROJECT_ROOT_PATH}/build"; fi
	@mkdir -p ${PROJECT_ROOT_PATH}/build; cd ${PROJECT_ROOT_PATH}/build; cmake ..
	@cmake --build ${PROJECT_ROOT_PATH}/build  --target ${TARGET} -- -j `nproc`;
	@echo "  >  install generated protocol buffers by cpp"
	@cd ${PROJECT_ROOT_PATH}/build; cmake --install .
	


