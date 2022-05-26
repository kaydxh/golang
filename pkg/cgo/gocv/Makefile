MAKEFILE_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PROJECT_DIR := $(realpath ${MAKEFILE_DIR}/..)
PKG_CONFIG_PATH := ${PROJECT_DIR}/pkgconfig
SCRIPT_PATH := ${MAKEFILE_DIR}/../script
TARGET := $(shell basename ${MAKEFILE_DIR})

define join
join = $(subst $(_space),$(1),$(strip $(2)))
endef

.PHONY: all
all: deps generate  test
	@echo "${MAKEFILE_DIR}"
	@echo "${PROJECT_DIR}"
	@echo "${TARGET}"

.PHONY: deps
deps:
	@echo "  >  downloading deps library"

.PHONY: version
version:
	@bash -c "bash ${SCRIPT_PATH}/version.sh gitinfos"

.PHONY: test
test: version
	@$(eval THIRD_LIB_PATHS := $(shell find -L ${PROJECT_DIR}/third_path/ -type d -iname "lib*" -print0 |xargs -0 -I {} bash -c 'echo {}'|grep -v "stubs"))
	@$(eval JOINED_THIRD_LIB_PATHS := $(call join,:,$(THIRD_LIB_PATHS)))
	 PKG_CONFIG_PATH="${PKG_CONFIG_PATH}" LD_LIBRARY_PATH="$(JOINED_THIRD_LIB_PATHS):${LD_LIBRARY_PATH}" LIBRARY_PATH="$(JOINED_THIRD_LIB_PATHS):${LIBRARY_PATH}" go test -a -v .

.PHONY: generate 
generate:
	@echo "  >  start to generate protocol buffers for target ${TARGET} in golang"
	@bash -c "go generate  ../api/openapi-spec/proto.gen.go"
	@echo "  >  start to generate protocol buffers by cpp"
	@mkdir -p ${PROJECT_DIR}/build; cd ${PROJECT_DIR}/build; cmake .. --log-level "ERROR"
	@cmake --build ${PROJECT_DIR}/build --log-level "ERROR" --target proto-${TARGET} -- -j `nproc`;
	@echo "  >  install generated protocol buffers by cpp"
	@cd ${PROJECT_DIR}/build; cmake --install .
	
