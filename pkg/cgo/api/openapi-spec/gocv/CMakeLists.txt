find_package(Protobuf)
list(APPEND IMPORT_DIRS "${PROJECT_SOURCE_DIR}")

GENERATE_PROTOBUF_LIB("${IMPORT_DIRS}" "proto-types")
