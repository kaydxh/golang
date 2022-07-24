macro(GENERATE_PROTOBUF_LIB IMPORT_DIRS DEPEND_TARGETS) 
   if (IMPORT_DIRS)
        list(APPEND Protobuf_IMPORT_DIRS "${IMPORT_DIRS}")
   endif ()

  get_filename_component(CURRENT_FOLDER ${CMAKE_CURRENT_SOURCE_DIR} NAME)
  set(target ${CURRENT_FOLDER})

  if (EXISTS ${CMAKE_CURRENT_SOURCE_DIR}/third_path/protobuf/bin/protoc)
    set(Protobuf_PROTOC_EXECUTABLE ${CMAKE_CURRENT_SOURCE_DIR}/third_path/protobuf/bin/protoc)
  elseif (EXISTS ${PROJECT_SOURCE_DIR}/third_path/protobuf/bin/protoc)
    set(Protobuf_PROTOC_EXECUTABLE ${PROJECT_SOURCE_DIR}/third_path/protobuf/bin/protoc)
  elseif (EXISTS ${PROJECT_SOURCE_DIR}/third_party/protobuf/bin/protoc)
    set(Protobuf_PROTOC_EXECUTABLE ${PROJECT_SOURCE_DIR}/third_party/protobuf/bin/protoc)
  elseif (EXISTS ${PROJECT_SOURCE_DIR}/api/openapi-spec/third_path/protobuf/bin/protoc)
    set(Protobuf_PROTOC_EXECUTABLE ${PROJECT_SOURCE_DIR}/api/openapi-spec/third_path/protobuf/bin/protoc)
  else ()
    message(SEND_ERROR "Error: protoc not found")
    return()
  endif ()

  message(STATUS Protobuf_PROTOC_EXECUTABLE=${Protobuf_PROTOC_EXECUTABLE})
  if (NOT TARGET protobuf::protoc)
     add_executable(protobuf::protoc IMPORTED)
     if (EXISTS "${Protobuf_PROTOC_EXECUTABLE}")
         set_target_properties(protobuf::protoc PROPERTIES
         IMPORTED_LOCATION "${Protobuf_PROTOC_EXECUTABLE}")
     endif ()
  endif ()

  set(protobuf_generate_LANGUAGE cpp)
  set(protobuf_generate_PROTOC_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")
  set(Protobuf_USE_STATIC_LIBS ON)
  file(GLOB_RECURSE protofiles RELATIVE ${CMAKE_CURRENT_SOURCE_DIR} "*.proto") 

  # PROTO_SRCS 存储.pb.cc文件的变量名称
  # PROTO_HDRS 存储.pb.h文件的变量名称
  # PROTOBUF_GENERATE_CPP generate *.h *.cxx files
  # https://github.com/protocolbuffers/protobuf/blob/main/cmake/protobuf-module.cmake.in
  PROTOBUF_GENERATE_CPP(PROTO_SRCS PROTO_HDRS ${protofiles})

  add_library(proto-${target} STATIC ${PROTO_SRCS} ${PROTO_HDRS})
  target_include_directories(proto-${target} PUBLIC ${CMAKE_CURRENT_BINARY_DIR})
  target_link_libraries(proto-${target} PUBLIC protobuf ${DEPEND_TARGETS})

  install(TARGETS proto-${target}
        OPTIONAL
        LIBRARY DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}/
        ARCHIVE DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}/
        RUNTIME DESTINATION DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}/
        )
   
  foreach (_abs_file ${PROTO_HDRS})
          file(RELATIVE_PATH _rel_file ${CMAKE_CURRENT_BINARY_DIR} ${_abs_file})
          get_filename_component(_rel_dir ${_rel_file} DIRECTORY)
          install(FILES ${_abs_file} DESTINATION ${PROJECT_SOURCE_DIR}/${_rel_dir} OPTIONAL)
  endforeach ()
endmacro()

macro(GET_MODEL_DIRS direction list_return)
  FILE(GLOB_RECURSE dirs LIST_DIRECTORIES true RELATIVE ${direction} ${direction}/*)
  SET(submodules_dirs)
  FOREACH (dir ${dirs}) 
    IF(EXISTS ${direction}/${dir}/CMakeLists.txt)
      LIST(APPEND submodules_dirs ${direction}/${dir})
    endif()
  ENDFOREACH()
  SET(${list_return} ${submodules_dirs})
endmacro()


# This file contains backwards compatibility patches for various legacy functions and variables
# Functions
function(PROTOBUF_GENERATE_CPP SRCS HDRS)
  cmake_parse_arguments(protobuf_generate_cpp "" "EXPORT_MACRO" "" ${ARGN})

  set(_proto_files "${protobuf_generate_cpp_UNPARSED_ARGUMENTS}")
  if(NOT _proto_files)
    message(SEND_ERROR "Error: PROTOBUF_GENERATE_CPP() called without any proto files")
    return()
  endif()

  if(PROTOBUF_GENERATE_CPP_APPEND_PATH)
    set(_append_arg APPEND_PATH)
  endif()

  if(DEFINED Protobuf_IMPORT_DIRS)
    set(_import_arg IMPORT_DIRS ${Protobuf_IMPORT_DIRS})
  endif()

  set(_outvar)
  protobuf_generate(${_append_arg} LANGUAGE cpp EXPORT_MACRO ${protobuf_generate_cpp_EXPORT_MACRO} OUT_VAR _outvar ${_import_arg} PROTOS ${_proto_files})

  set(${SRCS})
  set(${HDRS})
  foreach(_file ${_outvar})
    if(_file MATCHES "cc$")
      list(APPEND ${SRCS} ${_file})
    else()
      list(APPEND ${HDRS} ${_file})
    endif()
  endforeach()
  set(${SRCS} ${${SRCS}} PARENT_SCOPE)
  set(${HDRS} ${${HDRS}} PARENT_SCOPE)
endfunction()


function(protobuf_generate)
    set(_options APPEND_PATH DESCRIPTORS)
    set(_singleargs LANGUAGE OUT_VAR EXPORT_MACRO PROTOC_OUT_DIR PLUGIN)
    if (COMMAND target_sources)
        list(APPEND _singleargs TARGET)
    endif ()
    set(_multiargs PROTOS IMPORT_DIRS GENERATE_EXTENSIONS)

    cmake_parse_arguments(protobuf_generate "${_options}" "${_singleargs}" "${_multiargs}" "${ARGN}")

    if (NOT protobuf_generate_PROTOS AND NOT protobuf_generate_TARGET)
        message(SEND_ERROR "Error: protobuf_generate called without any targets or source files")
        return()
    endif ()

    if (NOT protobuf_generate_OUT_VAR AND NOT protobuf_generate_TARGET)
        message(SEND_ERROR "Error: protobuf_generate called without a target or output variable")
        return()
    endif ()

    if (NOT protobuf_generate_LANGUAGE)
        set(protobuf_generate_LANGUAGE cpp)
    endif ()
    string(TOLOWER ${protobuf_generate_LANGUAGE} protobuf_generate_LANGUAGE)

    if (NOT protobuf_generate_PROTOC_OUT_DIR)
        set(protobuf_generate_PROTOC_OUT_DIR ${CMAKE_CURRENT_BINARY_DIR})
    endif ()

    if (protobuf_generate_EXPORT_MACRO AND protobuf_generate_LANGUAGE STREQUAL cpp)
        set(_dll_export_decl "dllexport_decl=${protobuf_generate_EXPORT_MACRO}:")
    endif ()

    if (protobuf_generate_PLUGIN)
        set(_plugin "--plugin=${protobuf_generate_PLUGIN}")
    endif ()

    if (NOT protobuf_generate_GENERATE_EXTENSIONS)
        if (protobuf_generate_LANGUAGE STREQUAL cpp)
            set(protobuf_generate_GENERATE_EXTENSIONS .pb.h .pb.cc)
        elseif (protobuf_generate_LANGUAGE STREQUAL python)
            set(protobuf_generate_GENERATE_EXTENSIONS _pb2.py)
        else ()
            message(SEND_ERROR "Error: protobuf_generate given unknown Language ${LANGUAGE}, please provide a value for GENERATE_EXTENSIONS")
            return()
        endif ()
    endif ()

    if (protobuf_generate_TARGET)
        get_target_property(_source_list ${protobuf_generate_TARGET} SOURCES)
        foreach (_file ${_source_list})
            if (_file MATCHES "proto$")
                list(APPEND protobuf_generate_PROTOS ${_file})
            endif ()
        endforeach ()
    endif ()

    if (NOT protobuf_generate_PROTOS)
        message(SEND_ERROR "Error: protobuf_generate could not find any .proto files")
        return()
    endif ()

    if (protobuf_generate_APPEND_PATH)
        # Create an include path for each file specified
        foreach (_file ${protobuf_generate_PROTOS})
            get_filename_component(_abs_file ${_file} ABSOLUTE)
            get_filename_component(_abs_path ${_abs_file} PATH)
            list(FIND _protobuf_include_path ${_abs_path} _contains_already)
            if (${_contains_already} EQUAL -1)
                list(APPEND _protobuf_include_path -I ${_abs_path})
            endif ()
        endforeach ()
    else ()
        set(_protobuf_include_path -I ${CMAKE_CURRENT_SOURCE_DIR})
    endif ()

    foreach (DIR ${protobuf_generate_IMPORT_DIRS})
        get_filename_component(ABS_PATH ${DIR} ABSOLUTE)
        list(FIND _protobuf_include_path ${ABS_PATH} _contains_already)
        if (${_contains_already} EQUAL -1)
            list(APPEND _protobuf_include_path -I ${ABS_PATH})
        endif ()
    endforeach ()

    set(_generated_srcs_all)
    foreach (_proto ${protobuf_generate_PROTOS})
        get_filename_component(_abs_file ${_proto} ABSOLUTE)
        get_filename_component(_abs_dir ${_abs_file} DIRECTORY)
        get_filename_component(_basename ${_proto} NAME_WLE)
        file(RELATIVE_PATH _rel_dir ${PROJECT_SOURCE_DIR} ${_abs_dir})

        set(_possible_rel_dir)
        if (NOT protobuf_generate_APPEND_PATH)
            set(_possible_rel_dir ${_rel_dir}/)
        endif ()

        set(_generated_srcs)
        foreach (_ext ${protobuf_generate_GENERATE_EXTENSIONS})
            list(APPEND _generated_srcs "${protobuf_generate_PROTOC_OUT_DIR}/${_possible_rel_dir}${_basename}${_ext}")
        endforeach ()

        if (protobuf_generate_DESCRIPTORS AND protobuf_generate_LANGUAGE STREQUAL cpp)
            set(_descriptor_file "${CMAKE_CURRENT_BINARY_DIR}/${_basename}.desc")
            set(_dll_desc_out "--descriptor_set_out=${_descriptor_file}")
            list(APPEND _generated_srcs ${_descriptor_file})
        endif ()
        list(APPEND _generated_srcs_all ${_generated_srcs})

        file(RELATIVE_PATH _rel_proto_file ${PROJECT_SOURCE_DIR} ${_abs_file})
        add_custom_command(
                OUTPUT ${_generated_srcs}
                COMMAND protobuf::protoc
                ARGS --${protobuf_generate_LANGUAGE}_out ${_dll_export_decl}${protobuf_generate_PROTOC_OUT_DIR} ${_plugin} ${_dll_desc_out} ${_protobuf_include_path} ${_rel_proto_file}
                WORKING_DIRECTORY ${PROJECT_SOURCE_DIR}
                DEPENDS ${_abs_file} protobuf::protoc
                COMMENT "Running ${protobuf_generate_LANGUAGE} protocol buffer compiler on ${_proto}"
                VERBATIM)
    endforeach ()

    set_source_files_properties(${_generated_srcs_all} PROPERTIES GENERATED TRUE)
    if (protobuf_generate_OUT_VAR)
        set(${protobuf_generate_OUT_VAR} ${_generated_srcs_all} PARENT_SCOPE)
    endif ()
    if (protobuf_generate_TARGET)
        target_sources(${protobuf_generate_TARGET} PRIVATE ${_generated_srcs_all})
    endif ()
endfunction()
