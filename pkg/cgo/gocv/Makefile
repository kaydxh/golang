MAKEFILE_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PROJECT_DIR := $(realpath ${MAKEFILE_DIR}/..)
PKG_CONFIG_PATH := ${PROJECT_ROOT_PATH}/pkgconfig
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
	@echo "  >  starting generate protocol buffers for target ${TARGET} in golang"
	@bash -c "go generate  ../api/openapi-spec/proto.gen.go"
	
