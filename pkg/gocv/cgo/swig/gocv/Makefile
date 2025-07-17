MAKEFILE_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PROJECT_ROOT_DIR := $(realpath ${MAKEFILE_DIR}/../..)
PKG_CONFIG_PATH := ${PROJECT_ROOT_DIR}/pkgconfig
SCRIPT_PATH := ${MAKEFILE_DIR}/../../script
TARGET := $(shell basename ${MAKEFILE_DIR})

# https://web.mit.edu/gnu/doc/html/make_8.html
_empty:=
# define space with $(_empty)
_space:= $(_empty) $(_empty)
define joinwith
  $(subst $(_space),$1,$(strip $2))
endef

.PHONY: all
all: deps generate test
	@echo "${MAKEFILE_DIR}"
	@echo "${PROJECT_ROOT_DIR}"
	@echo "${TARGET}"

.PHONY: deps
deps:
	@echo "  >  downloading deps library"

.PHONY: version
version:
	@bash -c "bash ${SCRIPT_PATH}/version.sh gitinfos"

.PHONY: test
test: version
	@$(eval THIRD_LIB_PATHS := $(shell find -L ${PROJECT_ROOT_DIR}/third_path/ -type d -iname "lib*" -print0 |xargs -0 -I {} bash -c 'echo {}'|grep -v "stubs"))
	@echo ${THIRD_LIB_PATHS}
	@$(eval JOINED_THIRD_LIB_PATHS := $(call joinwith,:,$(THIRD_LIB_PATHS)))
	 PKG_CONFIG_PATH="${PKG_CONFIG_PATH}" LD_LIBRARY_PATH="$(JOINED_THIRD_LIB_PATHS):${LD_LIBRARY_PATH}" LIBRARY_PATH="$(JOINED_THIRD_LIB_PATHS):${LIBRARY_PATH}" go test -a -v .

.PHONY: generate 
generate:
	
