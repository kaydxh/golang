MAKEFILE_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PROJECT_DIR := $(realpath ${MAKEFILE_DIR}/..)
PKG_CONFIG_PATH := ${PROJECT_DIR}/pkgconfig
SCRIPT_PATH := ${MAKEFILE_DIR }/../api/openapi-spec/script
TARGET := $(shell basename ${MAKEFILE_DIR})

.PHONY: all
all: generate 
	@echo "${MAKEFILE_DIR}"
	@echo "${PROJECT_DIR}"
	@echo "${TARGET}"

.PHONY: deps
deps:
	@echo "  >  downloading deps library"

.PHONY: generate 
generate:
	@echo "  >  start to generate protocol buffers for target ${TARGET} in golang"
	@bash -c "go generate  ../api/openapi-spec/proto.gen.go"
	@echo "  >  start to generate protocol buffers by cpp"
	@mkdir -p ${PROJECT_DIR}/build; cd ${PROJECT_DIR}/build; cmake .. --log-level "ERROR"
	@cmake --build ${PROJECT_DIR}/build --log-level "ERROR" --target proto-${TARGET} -- -j `nproc`;
	@echo "  >  install generated protocol buffers by cpp"
	@cd ${PROJECT_DIR}/build; cmake --install .
	
