include(CMakeParseArguments)
#
# log_execute_process is same to execute_process but message cmd to log
#

macro(log_debug msg)
    get_filename_component(name ${CMAKE_CURRENT_LIST_FILE} NAME)
    string(TIMESTAMP TIME_NOW "%Y-%m-%d %H:%M:%S")
    message("${TIME_NOW} - ${name}:${CMAKE_CURRENT_LIST_LINE} - ${msg}")
endmacro(log_debug)

macro(log_warn msg)
    get_filename_component(name ${CMAKE_CURRENT_LIST_FILE} NAME)
    string(TIMESTAMP TIME_NOW "%Y-%m-%d %H:%M:%S")
    message(WARNING "${TIME_NOW} - ${name}:${CMAKE_CURRENT_LIST_LINE} - ${msg}")
endmacro(log_warn)
#
macro(log_error msg)
    get_filename_component(name ${CMAKE_CURRENT_LIST_FILE} NAME)
    string(TIMESTAMP TIME_NOW "%Y-%m-%d %H:%M:%S")
    message(FATAL_ERROR "${TIME_NOW} - ${name}:${CMAKE_CURRENT_LIST_LINE} - ${msg}")
endmacro(log_error)

# debug(K1 K2 K3)
#
macro(debug)
    FOREACH(A ${ARGN})
      log_debug("${A}:${${A}}")
    ENDFOREACH(A)
endmacro(debug)

###################
#no Special symbol in cmd
macro(log_execute_process)
    set(options "" )
    set(oneValueArgs WORKING_DIRECTORY )
    set(multiValueArgs COMMAND)
    cmake_parse_arguments(log_execute_process "${options}" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})
    string(REPLACE ";" " " CMD "${log_execute_process_COMMAND}")

    if(NOT log_execute_process_WORKING_DIRECTORY)
         log_error("bad cmd:cd ${log_execute_process_WORKING_DIRECTORY} && ${CMD}")
    endif()

    if(NOT EXISTS ${log_execute_process_WORKING_DIRECTORY})
         log_error("bad cmd:cd ${log_execute_process_WORKING_DIRECTORY} && ${CMD}")
    endif()

    log_debug("cd ${log_execute_process_WORKING_DIRECTORY} && ${CMD}")
    execute_process(COMMAND ${log_execute_process_COMMAND}
                    WORKING_DIRECTORY ${log_execute_process_WORKING_DIRECTORY}
                    RESULT_VARIABLE rv ERROR_VARIABLE er
                    OUTPUT_FILE ${CMAKE_BINARY_DIR}/cmake_cmd_out.log
                    ERROR_FILE ${CMAKE_BINARY_DIR}/cmake_error.log)
    if(rv)
       log_debug("RESULT:${rv}, CMD:${log_execute_process_COMMAND}")
       execute_process(COMMAND cat ${CMAKE_BINARY_DIR}/cmake_error.log WORKING_DIRECTORY ${log_execute_process_WORKING_DIRECTORY})
       log_error("fatal error exit!")
    endif()

endmacro(log_execute_process)

macro(backup_and_mv_file)
    file(MAKE_DIRECTORY ${CMAKE_SOURCE_DIR}/.backup)
    FOREACH(C ${ARGN})
        if(EXISTS ${C})
             log_execute_process(COMMAND mv -f --backup=t ${C} -t ${CMAKE_SOURCE_DIR}/.backup WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})
        endif()
    ENDFOREACH(C)
endmacro(backup_and_mv_file)



macro(uncompress_package package output_dirname)
    get_filename_component(package_name ${package} NAME)
    get_filename_component(package_path ${package} PATH)

    if(${package_name} MATCHES ".zip$")
           string(FIND ${package_name} ".zip" LAST_POS REVERSE)
           set(UMCOMPRESS_CMD unzip -o)
    elseif(${package_name} MATCHES ".tar.gz$")
           string(FIND ${package_name} ".tar.gz" LAST_POS REVERSE)
           set(UMCOMPRESS_CMD tar zxvf)
    elseif(${package_name} MATCHES ".tar$")
           string(FIND ${package_name} ".tar" LAST_POS REVERSE)
           set(UMCOMPRESS_CMD tar xvf)
    endif()

    string(SUBSTRING ${package_name} 0 ${LAST_POS} dirname)
    set(package_dirname ${dirname} PARENT_SCOPE)

    if(NOT EXISTS ${output_dirname})
           file(MAKE_DIRECTORY ${output_dirname})
    endif()

    if(NOT EXISTS ${output_dirname}/${dirname})
           log_execute_process(COMMAND ${UMCOMPRESS_CMD} ${package} WORKING_DIRECTORY ${output_dirname})
    elseif(${package} IS_NEWER_THAN ${output_dirname}/${dirname})
           #backup_and_mv_file(${output_dirname}/${dirname})
           log_execute_process(COMMAND rm -fr ${output_dirname}/${dirname} WORKING_DIRECTORY ${output_dirname})
           log_execute_process(COMMAND ${UMCOMPRESS_CMD} ${package} WORKING_DIRECTORY ${output_dirname})
    endif()
    
endmacro(uncompress_package)


#
# check_and_install(PACKAGE TARGET)
#
function(check_and_install)

    get_filename_component(install_dir_path ${CMAKE_CURRENT_LIST_FILE} PATH)
    get_filename_component(install_dir_name ${install_dir_path} NAME)

    set(options "" )
    set(oneValueArgs "")
    set(multiValueArgs PACKAGE TARGET)
    cmake_parse_arguments(C "${options}" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})

    #0 parse args
    if(NOT C_PACKAGE OR NOT C_TARGET)
         log_error("check_and_install PACKAGE xx\n      TARGET xx\n INSTALL xx\n  CLEAR xx\n")
    endif()
    #debug(C_PACKAGE C_TARGET)

    # #1 check is package
    FOREACH(P ${C_PACKAGE})
         if(NOT EXISTS ${P})
              log_error("check_and_install: ${P} NOT FOUND")
         endif()
    ENDFOREACH(P)

    #1 check is need rebuild
    SET(${install_dir_name}_NEED_INSTALL FALSE)
    FOREACH(P ${C_PACKAGE})
        FOREACH(T ${C_TARGET})
             #1.1 target not exists
             if(NOT EXISTS ${T})
               SET(${install_dir_name}_NEED_INSTALL TRUE)
               log_debug("check_and_install: ${T} NOT FOUND, INSTALL ${C_PACKAGE} to ${install_dir_path}, ${install_dir_name}_NEED_INSTALL:${${install_dir_name}_NEED_INSTALL}")
               
               break()
            endif()

            #1.2 target has been changed
            if(${P} IS_NEWER_THAN ${T})
               log_debug("check_and_install: ${P} newer than ${T}, INSTALL ${C_PACKAGE} to ${install_dir_path}, ${install_dir_name}_NEED_INSTALL:${${install_dir_name}_NEED_INSTALL}")
               SET(${install_dir_name}_NEED_INSTALL TRUE)
               break()
            endif()
         ENDFOREACH(T)
    ENDFOREACH(P)

    #return
    SET(${install_dir_name}_NEED_INSTALL ${${install_dir_name}_NEED_INSTALL} PARENT_SCOPE)
endfunction(check_and_install)

function(find_and_set_package install_dir_path)

    set(options "" )
    set(oneValueArgs INCLUDE LIB BIN)
    set(multiValueArgs "")
    cmake_parse_arguments(find_and_set_package "${options}" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})

    if(NOT  find_and_set_package_INCLUDE)
        SET(find_and_set_package_INCLUDE include)
    endif()
    if(NOT find_and_set_package_LIB)
        SET(find_and_set_package_LIB lib)
    endif()
    if(NOT  find_and_set_package_BIN)
        SET(find_and_set_package_BIN bin)
    endif()

    get_filename_component(install_dir_name ${install_dir_path} NAME)
    set(PRE ${install_dir_name})
    set(${PRE}_INSTALL_DIR ${install_dir_path})
    SET(${PRE}_INCLUDE ${${PRE}_INSTALL_DIR}/${find_and_set_package_INCLUDE})
    SET(${PRE}_LIB ${${PRE}_INSTALL_DIR}/${find_and_set_package_LIB})
    SET(${PRE}_BIN ${${PRE}_INSTALL_DIR}/${find_and_set_package_BIN})
    

    # include
    if(NOT EXISTS ${${PRE}_INCLUDE} )
         log_error("${PRE}_INCLUDE:${${PRE}_INCLUDE} NOT FOUND, Please check_and_install first!")
    endif()
    file(GLOB_RECURSE ${PRE}_INCLUDES  ${${PRE}_INCLUDE}/*.h)
    #if(NOT ${PRE}_INCLUDES)
    #    log_warn("${${PRE}_INCLUDE} NOT FOUND HEADES, Please check_and_install first!")
    #endif()

    # use lib64 instead 
    if(NOT EXISTS ${${PRE}_LIB})
        SET(${PRE}_LIB  ${${PRE}_INSTALL_DIR}/lib64)
    endif()    
    # libs
    if(NOT EXISTS ${${PRE}_LIB} )
         log_error("${${PRE}_LIB} NOT FOUND, Please check_and_install first!")
    endif()   


    #static lib
    file(GLOB ${PRE}_STATIC_LIBS ${${PRE}_LIB}/lib*.a)
    FOREACH(A ${${PRE}_STATIC_LIBS})
         get_filename_component(_lib_file_name ${A} NAME_WE)
         string(SUBSTRING  ${_lib_file_name} 3 -1  _lib_name)
         if(NOT _lib_name)
              log_error("Bad lib name: ${A} in ${${PRE}_LIB}")
         endif()

        SET(${_lib_name}_A ${A} PARENT_SCOPE)
        add_library(${_lib_name}-static STATIC IMPORTED)
        set_property(TARGET ${_lib_name}-static PROPERTY IMPORTED_LOCATION ${A})
        log_debug("SET ${_lib_name}-static: ${A}")
    ENDFOREACH(A)

    #sharelib
    file(GLOB ${PRE}_SHARED_LIBS ${${PRE}_LIB}/lib*.so)
    FOREACH(A ${${PRE}_SHARED_LIBS})
         get_filename_component(_lib_file_name ${A} NAME_WE)
         string(SUBSTRING  ${_lib_file_name} 3 -1  _lib_name)
         if(NOT _lib_name)
              log_error("Bad lib name: ${A} in ${${PRE}_LIB}")
         endif()

        SET(${_lib_name}_SO ${A})
        add_library(${_lib_name}-shared SHARED IMPORTED)
        set_property(TARGET ${_lib_name}-shared PROPERTY IMPORTED_LOCATION ${A})
        log_debug("SET ${_lib_name}-shared: ${A}")
    ENDFOREACH(A)

    # cmds
    file(GLOB ${PRE}_BINS ${${PRE}_BIN}/*)
    FOREACH(A ${${PRE}_BINS})
        get_filename_component(P ${A} NAME)
        SET(${P}_BIN ${A})
        log_debug("SET ${P}_BIN: ${${P}_BIN}")
        SET(${P}_BIN ${A} PARENT_SCOPE)
    ENDFOREACH(A)

    # FILES for ide
    #file(GLOB_RECURSE ${PRE}_FILES ${install_dir_path}/*)

    #source_group(${PRE}_IDE_VS FILES ${${PRE}_FILES})
    #add_custom_target(${PRE}_IDE_QT SOURCES  ${${PRE}_FILES})

    #export to parent
    #set(PRE ${PRE} PARENT_SCOPE)
    #set(${PRE}_INSTALL_DIR ${${PRE}_INSTALL_DIR}  PARENT_SCOPE)
    SET(${PRE}_INCLUDE ${${PRE}_INCLUDE} PARENT_SCOPE)
    SET(${PRE}_LIB ${${PRE}_LIB} PARENT_SCOPE)
    SET(${PRE}_BIN ${${PRE}_BIN} PARENT_SCOPE)

    #set(${PRE}_INCLUDES ${${PRE}_INCLUDES} PARENT_SCOPE)
    #set(${PRE}_STATIC_LIBS ${${PRE}_STATIC_LIBS} PARENT_SCOPE)
    #set(${PRE}_SHARED_LIBS ${${PRE}_SHARED_LIBS} PARENT_SCOPE)
    #set(${PRE}_BINS ${${PRE}_BINS} PARENT_SCOPE)

    #set(${PRE}_FILES ${${PRE}_FILES} PARENT_SCOPE)
    set(${PRE}_FOUND TRUE PARENT_SCOPE)

endfunction(find_and_set_package)

# include all subdirs
function(recurse_include_directories DIR)

    file(GLOB_RECURSE _header_files FOLLOW_SYMLINKS ${DIR}/*.h)
    
    set(files_dirs "")
    foreach(src_file ${_header_files})        
             GET_FILENAME_COMPONENT(dir_path ${src_file} PATH)
             list(APPEND files_dirs ${dir_path})
    endforeach()
    LIST(REMOVE_DUPLICATES files_dirs)
    
    include_directories(${DIR})
    foreach(H ${files_dirs})
         include_directories(${H})
    endforeach(H)
    
    set(DIR_HEADERS ${_header_files} PARENT_SCOPE)
endfunction(recurse_include_directories)


function(add_source_to_qtcreator DIR)
    get_filename_component(dir_name ${DIR} NAME)
    string(REGEX REPLACE "/" "_" source_target ${DIR})
    set(source_target "${dir_name}_${source_target}")

    log_debug("add_files_to_qtcreator: ${source_target}")
    
    # all header store in HEADERS_${DIR}
    #recurse_include_directories(${DIR}) 
    
    # all c files .c / .cc / .cpp to c_src_files
    file(GLOB_RECURSE c_src_files FOLLOW_SYMLINKS ${DIR}/*[c|p|e])
    
    set(other_files "")
    foreach(suffix ${ARGN})
          file(GLOB_RECURSE other_src_files FOLLOW_SYMLINKS ${DIR}/*.${suffix})
          list(APPEND other_files ${other_src_files})
    endforeach()
    LIST(REMOVE_DUPLICATES other_files)
    
    add_custom_target(${source_target}
         SOURCES ${DIR_HEADERS} ${c_src_files} ${other_files}
         WORKING_DIRECTORY ${DIR})
    
    set(DIR_ALLS  ${DIR_HEADERS} ${c_src_files} ${other_files} PARENT_SCOPE)
endfunction(add_source_to_qtcreator)


macro(add_test_item)

    set(options "" )
    set(oneValueArgs "")
    set(multiValueArgs SOURCES INCLUDES LIBS)
    cmake_parse_arguments(C "${options}" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})
    
    SET(TARGET "")
    FOREACH(P ${C_SOURCES})
         get_filename_component(FILE_NAME ${P} NAME_WE)
         SET(TARGET "${FILE_NAME}")
         
         if(${TARGET} MATCHES "test_*")
             break()
         endif()
    ENDFOREACH(P)
    
    add_executable(${TARGET} ${C_SOURCES})

    if(C_INCLUDES)
        target_include_directories(${TARGET} PUBLIC ${C_INCLUDES})
    endif()

    if(C_LIBS)
       target_link_libraries(${TARGET} ${C_LIBS})
    endif()  
    
    add_test(NAME ${TARGET} COMMAND ${TARGET} WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR})
    log_debug("add_test:${TARGET}, srcs: ${C_SOURCES}")
    
endmacro(add_test_item)


# ngx_add_dynamic_module(DIR xxx SOURCES XXXX LIBS XXXX)
macro(ngx_add_dynamic_module)

    set(options "" )
    set(oneValueArgs "DIR")
    set(multiValueArgs SOURCES INCLUDES LIBS)
    cmake_parse_arguments(C "${options}" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})
    
    if(NOT C_SOURCES)
       file(GLOB C_SRCS "${C_DIR}/*.c")
       file(GLOB C_HEADERS "${C_DIR}/*.h")
       file(GLOB CC_SRCS "${C_DIR}/*.cc")
       file(GLOB CPP_SRCS "${C_DIR}/*.cpp")
       file(GLOB HPP_SRCS "${C_DIR}/*.hpp")       
       SET(C_SOURCES ${C_SRCS} ${C_HEADERS} ${CC_SRCS} ${CPP_SRCS} ${HPP_SRCS})
       #debug(C_SOURCES C_SRCS CMAKE_CURRENT_SOURCE_DIR)
    endif()
    
    get_filename_component(module_name ${C_DIR} NAME)
    
    log_debug("add_dynamic_module:${module_name}, srcs: ${C_SOURCES}")
    add_library(${module_name} SHARED ${C_SOURCES})

    if(C_INCLUDES)
       target_include_directories(${module_name} PUBLIC ${C_INCLUDES})
    endif()
    
    if(C_LIBS)
       target_link_libraries(${module_name} PUBLIC ${C_LIBS})
    endif()

    INSTALL(TARGETS ${C_DIR} DESTINATION shared)

endmacro(ngx_add_dynamic_module)

# ngx_add_static_module(DIR xxx SOURCES XXXX INCLUDES XXXLIBS XXXX)
macro(ngx_add_static_module)

    set(options "" )
    set(oneValueArgs "DIR")
    set(multiValueArgs SOURCES INCLUDES LIBS)
    cmake_parse_arguments(C "${options}" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})
    
    if(NOT C_SOURCES)
       file(GLOB C_SRCS "${C_DIR}/*.c")
       file(GLOB C_HEADERS "${C_DIR}/*.h")
       file(GLOB CC_SRCS "${C_DIR}/*.cc")
       file(GLOB CPP_SRCS "${C_DIR}/*.cpp")
       file(GLOB HPP_SRCS "${C_DIR}/*.hpp")       
       SET(C_SOURCES ${C_SRCS} ${C_HEADERS} ${CC_SRCS} ${CPP_SRCS} ${HPP_SRCS})
       #debug(C_SOURCES C_SRCS CMAKE_CURRENT_SOURCE_DIR)
    endif()
    
    get_filename_component(module_name ${C_DIR} NAME)



    log_debug("add_static_module:${module_name}, srcs: ${C_SOURCES}")
    add_library(${module_name} STATIC ${C_SOURCES})

    if(C_INCLUDES)
       target_include_directories(${module_name} PUBLIC ${C_INCLUDES})
    endif()
    
    if(C_LIBS)
       # target is nginx static link
       target_link_libraries(nginx PUBLIC ${module_name} ${C_LIBS})
    else()
       target_link_libraries(nginx PUBLIC ${module_name})
    endif()

    #INSTALL(TARGETS ${C_DIR} DESTINATION shared)
endmacro(ngx_add_static_module)


function(grpc_gen_code SRCS HDRS)
  if(NOT ARGN)
    message(SEND_ERROR "Error: grpc_gen_code() called without any proto files")
    return()
  endif()

  set(PROTOBUF_GENERATE_CPP_APPEND_PATH TRUE)
  if(PROTOBUF_GENERATE_CPP_APPEND_PATH)
    # Create an include path for each file specified
    foreach(FIL ${ARGN})
      get_filename_component(ABS_FIL ${FIL} ABSOLUTE)
      get_filename_component(ABS_PATH ${ABS_FIL} PATH)
      list(FIND _protobuf_include_path ${ABS_PATH} _contains_already)
      if(${_contains_already} EQUAL -1)
          list(APPEND _protobuf_include_path -I ${ABS_PATH})
      endif()
    endforeach()
  else()
    set(_protobuf_include_path -I ${CMAKE_CURRENT_SOURCE_DIR})
  endif()

  if(DEFINED PROTOBUF_IMPORT_DIRS)
    foreach(DIR ${PROTOBUF_IMPORT_DIRS})
      get_filename_component(ABS_PATH ${DIR} ABSOLUTE)
      list(FIND _protobuf_include_path ${ABS_PATH} _contains_already)
      if(${_contains_already} EQUAL -1)
          list(APPEND _protobuf_include_path -I ${ABS_PATH})
      endif()
    endforeach()
  endif()

  set(${SRCS})
  set(${HDRS})
  
  foreach(FIL ${ARGN})
    get_filename_component(ABS_FIL ${FIL} ABSOLUTE)
    get_filename_component(FIL_WE ${FIL} NAME_WE)

    list(APPEND ${SRCS} ${CMAKE_CURRENT_BINARY_DIR}/${FIL_WE}.pb.cc ${CMAKE_CURRENT_BINARY_DIR}/${FIL_WE}.grpc.pb.cc)
    list(APPEND ${HDRS} ${CMAKE_CURRENT_BINARY_DIR}/${FIL_WE}.pb.h ${CMAKE_CURRENT_BINARY_DIR}/${FIL_WE}.grpc.pb.h)
    
    add_custom_command(
      OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/${FIL_WE}.pb.cc ${CMAKE_CURRENT_BINARY_DIR}/${FIL_WE}.grpc.pb.cc
             ${CMAKE_CURRENT_BINARY_DIR}/${FIL_WE}.pb.h ${CMAKE_CURRENT_BINARY_DIR}/${FIL_WE}.grpc.pb.h
             ${CMAKE_CURRENT_BINARY_DIR}/${FIL_WE}_mock.grpc.pb.h
      COMMAND ${protoc_BIN}
      ARGS --grpc_out=generate_mock_code=true:${CMAKE_CURRENT_BINARY_DIR}
           --cpp_out=${CMAKE_CURRENT_BINARY_DIR}
           --plugin=protoc-gen-grpc=${grpc_cpp_plugin_BIN}
           ${_protobuf_include_path}
           ${ABS_FIL}
      DEPENDS ${ABS_FIL} ${PROTOBUF_PROTOC_EXECUTABLE}
      WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
      COMMENT "Running gRPC C++ protocol buffer compiler on ${FIL}"
      VERBATIM)

  endforeach()

  set_source_files_properties(${${SRCS}} ${${HDRS}} PROPERTIES GENERATED TRUE)
  
  set(${SRCS} ${${SRCS}} PARENT_SCOPE)
  set(${HDRS} ${${HDRS}} PARENT_SCOPE)
endfunction()

macro(add_shared_library DIR)

    SET(LIB_DIR ${DIR}/lib)
    if(NOT EXISTS ${LIB_DIR})
        SET(LIB_DIR ${DIR}/lib64)
    endif()

    INCLUDE_DIRECTORIES(${DIR}/include)


    foreach(libname ${ARGN})
            add_library(${libname} SHARED IMPORTED)
            set_property(TARGET ${libname} PROPERTY IMPORTED_LOCATION ${LIB_DIR}/lib${libname}.so)            
    endforeach()
    log_debug("add_shared_library: ${ARGN} in: ${DIR}")
endmacro(add_shared_library)

macro(add_static_library DIR)
    set_property(GLOBAL APPEND PROPERTY GLOBAL_INCLUDE_DIRS "${DIR}/include")
    SET(LIB_DIR ${DIR}/lib)
    if(NOT EXISTS ${LIB_DIR})
        SET(LIB_DIR ${DIR}/lib64)
    endif()
    set_property(GLOBAL APPEND PROPERTY GLOBAL_LINK_DIRS ${LIB_DIR})
    INCLUDE_DIRECTORIES(${DIR}/include)
    #add_files_to_qtcreator(${DIR}/include .h .hpp)
    foreach(libname ${ARGN})
      add_library(${libname} STATIC IMPORTED)
      set_property(TARGET ${libname} PROPERTY IMPORTED_LOCATION ${DIR}/lib/lib${libname}.a)
    endforeach()
    log_debug("add_static_library: ${ARGN} in: ${DIR}")
endmacro(add_static_library)


function(protoc_gen_code SRCS HDRS)
  if(NOT ARGN)
    message(SEND_ERROR "Error: protoc_gen_code() called without any proto files")
    return()
  endif()

  set(PROTOBUF_GENERATE_CPP_APPEND_PATH TRUE)
  if(PROTOBUF_GENERATE_CPP_APPEND_PATH)
    # Create an include path for each file specified
    foreach(FIL ${ARGN})
      get_filename_component(ABS_FIL ${FIL} ABSOLUTE)
      get_filename_component(ABS_PATH ${ABS_FIL} PATH)
      list(FIND _protobuf_include_path ${ABS_PATH} _contains_already)
      if(${_contains_already} EQUAL -1)
          list(APPEND _protobuf_include_path -I ${ABS_PATH})
      endif()
    endforeach()
  else()
    set(_protobuf_include_path -I ${CMAKE_CURRENT_SOURCE_DIR})
  endif()

  if(DEFINED PROTOBUF_IMPORT_DIRS)
    foreach(DIR ${PROTOBUF_IMPORT_DIRS})
      get_filename_component(ABS_PATH ${DIR} ABSOLUTE)
      list(FIND _protobuf_include_path ${ABS_PATH} _contains_already)
      if(${_contains_already} EQUAL -1)
          list(APPEND _protobuf_include_path -I ${ABS_PATH})
      endif()
    endforeach()
  endif()

  set(${SRCS})
  set(${HDRS})
  
  foreach(FIL ${ARGN})
    get_filename_component(ABS_FIL ${FIL} ABSOLUTE)
    get_filename_component(FIL_WE ${FIL} NAME_WE)
    get_filename_component(FIL_PATH ${FIL} DIRECTORY)

    list(APPEND ${SRCS} ${FIL_PATH}/${FIL_WE}.pb.cc)
    list(APPEND ${HDRS} ${FIL_PATH}/${FIL_WE}.pb.h)
    
    #message("${${SRCS}} ${${HDRS}}")
    add_custom_command(
      OUTPUT ${FIL_PATH}/${FIL_WE}.pb.cc
             ${FIL_PATH}/${FIL_WE}.pb.h
             #${CMAKE_CURRENT_BINARY_DIR}/${FIL_WE}_mock.grpc.pb.h
      COMMAND ${protoc_BIN}
      ARGS --cpp_out=${FIL_PATH}
           ${_protobuf_include_path}
           ${ABS_FIL}
      DEPENDS ${ABS_FIL} ${PROTOBUF_PROTOC_EXECUTABLE}
      WORKING_DIRECTORY ${FIL_PATH}
      COMMENT "Running C++ protocol buffer compiler on ${FIL}"
      VERBATIM)

  endforeach()

  set_source_files_properties(${${SRCS}} ${${HDRS}} PROPERTIES GENERATED TRUE)
  
  set(${SRCS} ${${SRCS}} PARENT_SCOPE)
  set(${HDRS} ${${HDRS}} PARENT_SCOPE)
endfunction()


MACRO(SUBDIRLIST result curdir)
    FILE(GLOB children RELATIVE ${curdir} ${curdir}/[a-zA-Z0-9]*)
  SET(dirlist "")
  FOREACH(child ${children})
    IF(IS_DIRECTORY ${curdir}/${child})
      LIST(APPEND dirlist ${child})
    ENDIF()
  ENDFOREACH()
  SET(${result} ${dirlist})
ENDMACRO()

MACRO(IMPORT_ONE_LIB dep)
    # include every include
    SET(CURRENT_DIR ${dep}/include)
    if(IS_DIRECTORY ${CURRENT_DIR})
        INCLUDE_DIRECTORIES(${CURRENT_DIR})
    endif()

    set(mkl_regex ".*mkl.*")
    # #for bin file
    # FILE(GLOB children RELATIVE ${dep}/bin ${dep}/bin/*)
    # foreach(_bin_name ${children})
    #     if(TARGET ${_bin_name})
    #         get_property(_bin_dir TARGET ${_bin_name}  PROPERTY IMPORTED_LOCATION)
    #         log_warn("${_bin_name} already add at:${_bin_dir} skip ${dep}/bin/${_bin_name}")
    #         continue()
    #     endif()

    #     add_executable(${_bin_name} IMPORTED)
    #     set_property(TARGET ${_bin_name} PROPERTY IMPORTED_LOCATION ${dep}/bin/${_bin_name})
    #     log_debug("SET ${_bin_name}: ${dep}/bin/${_bin_name}")
    # endforeach()

    #for lib file
    foreach(CURRENT_DIR ${dep}/lib64 ${dep}/lib)

        FILE(GLOB children RELATIVE ${CURRENT_DIR} ${CURRENT_DIR}/lib*.a)
        foreach(libs ${children})
            string(REGEX REPLACE "^lib" "" _lib_name ${libs})
            string(REGEX REPLACE "\\.a" "" _lib_name ${_lib_name})
            #set(_lib_name ${_lib_name})

            if(TARGET ${_lib_name})
                get_property(_lib_name_dir TARGET ${_lib_name}  PROPERTY IMPORTED_LOCATION)
                log_warn("${_lib_name} already add at:${_lib_name_dir} skip ${CURRENT_DIR}/${libs}")
                continue()
            endif()

            add_library(${_lib_name} STATIC IMPORTED)
            set_property(TARGET ${_lib_name}  PROPERTY IMPORTED_LOCATION ${CURRENT_DIR}/${libs})
            log_debug("SET ${_lib_name}: ${CURRENT_DIR}/${libs}")
        endforeach()

        #for lib file
        FILE(GLOB children RELATIVE ${CURRENT_DIR} ${CURRENT_DIR}/lib*.so)
        foreach(libs ${children})
            string(REGEX REPLACE "^lib" "" _lib_name ${libs})
            string(REGEX REPLACE "\\.so" "" _lib_name ${_lib_name})
            set(_lib_name ${_lib_name}-so)

            if(TARGET ${_lib_name})
                get_property(_lib_name_dir TARGET ${_lib_name}  PROPERTY IMPORTED_LOCATION)
                log_warn("${_lib_name} already add at:${_lib_name_dir} skip ${CURRENT_DIR}/${libs}")
                continue()
            endif()

            add_library(${_lib_name} SHARED IMPORTED)
            set_property(TARGET ${_lib_name}  PROPERTY IMPORTED_LOCATION ${CURRENT_DIR}/${libs})
            log_debug("SET ${_lib_name}: ${CURRENT_DIR}/${libs}")

            if ("${_lib_name}" MATCHES "${mkl_regex}")
                log_debug("${_lib_name}: IMPORTED_NO_SONAME")
                set_property(TARGET ${_lib_name}  PROPERTY IMPORTED_NO_SONAME 1)
            endif()

        endforeach()
    endforeach()
ENDMACRO()


MACRO(ADD_THIRD_LIB THIRD_PATH)
  SUBDIRLIST(SUB_THIRD_DIRS ${THIRD_PATH})
  message(STATUS "SUB_THIRD_LIBS: ${SUB_THIRD_DIRS}")

  set(cuda8_regex ".*cuda8.*")
  set(cuda9_regex ".*cuda9.*")
  set(cuda10_regex ".*cuda10.*")

  foreach(dep ${SUB_THIRD_DIRS})
      message(STATUS "dep: ${dep}")
      if("${dep}" MATCHES "${cuda8_regex}")
          if(${CUDA8_0})
            message(STATUS "add ${THIRD_PATH}/${deps}")
            IMPORT_ONE_LIB(${THIRD_PATH}/${dep})
            set(CUDA_NVCC_FLAGS "${CUDA_NVCC_FLAGS} -Xcompiler -fPIC -O3 --compiler-options -fno-strict-aliasing -lineinfo -Xptxas -dlcm=cg -use_fast_math -gencode arch=compute_61,code=sm_61 -gencode arch=compute_60,code=sm_60 -gencode arch=compute_52,code=sm_52" CACHE STRING "cuda flags")
            message(STATUS "CUDA_NVCC_FLAGS:${CUDA_NVCC_FLAGS}")
          endif(${CUDA8_0})
      elseif("${dep}" MATCHES "${cuda9_regex}")
          if(${CUDA9_0})
            message(STATUS "add ${THIRD_PATH}/${deps}")
            IMPORT_ONE_LIB(${THIRD_PATH}/${dep})
            set(CUDA_NVCC_FLAGS "${CUDA_NVCC_FLAGS} -Xcompiler -fPIC -O3 --compiler-options -fno-strict-aliasing -lineinfo -Xptxas -dlcm=cg -use_fast_math -gencode arch=compute_70,code=sm_70 -gencode arch=compute_61,code=sm_61 -gencode arch=compute_60,code=sm_60 -gencode arch=compute_52,code=sm_52" CACHE STRING "cuda flags")
            message(STATUS "CUDA_NVCC_FLAGS:${CUDA_NVCC_FLAGS}")
          endif(${CUDA9_0})
      elseif("${dep}" MATCHES "${cuda10_regex}")
          if(${CUDA10_0})
            message(STATUS "add ${THIRD_PATH}/${deps}")
            IMPORT_ONE_LIB(${THIRD_PATH}/${dep})
            set(CUDA_NVCC_FLAGS "${CUDA_NVCC_FLAGS} -Xcompiler -fPIC -O3 --compiler-options -fno-strict-aliasing -lineinfo -Xptxas -dlcm=cg -use_fast_math -gencode arch=compute_75,code=sm_75 -gencode arch=compute_70,code=sm_70 -gencode arch=compute_61,code=sm_61 -gencode arch=compute_60,code=sm_60 -gencode arch=compute_52,code=sm_52" CACHE STRING "cuda flags")
            message(STATUS "CUDA_NVCC_FLAGS:${CUDA_NVCC_FLAGS}")
          endif(${CUDA10_0})
      else("${dep}" MATCHES "${cuda8_regex}")
        message(STATUS "not cuda8 or cuda9 lib. normal import")
        IMPORT_ONE_LIB(${THIRD_PATH}/${dep})
      endif("${dep}" MATCHES "${cuda8_regex}")
  endforeach()
ENDMACRO()


MACRO(INSTALL_SO LIB_DIR)
# todo copy by bin
    foreach(so_target ${ARGN})
       get_property(SO_FILE TARGET ${so_target} PROPERTY IMPORTED_LOCATION)

       if(SO_FILE MATCHES ".so$")
           log_debug("INSTALL LIBRARY ${SO_FILE}")
           file(GLOB ALL_SOS "${SO_FILE}*")
           INSTALL(FILES ${ALL_SOS} DESTINATION ${LIB_DIR})
       endif()
    endforeach()

ENDMACRO()

MACRO(INSTALL_BIN BIN_DIR)
    foreach(bin_target ${ARGN})
        get_property(BIN_FILE TARGET ${bin_target} PROPERTY IMPORTED_LOCATION)
        log_debug("INSTALL BINARY ${BIN_FILE}")
        INSTALL(FILES ${BIN_FILE} DESTINATION ${BIN_DIR})
    endforeach()
ENDMACRO()

#build type
#flag to import CUDA8.0/9.0/10.0
option(BUILD_CUDA8_0 "BUILD_CUDA8_0" OFF)
set(CUDA8_0 FALSE)
if(BUILD_CUDA8_0)
    set(CUDA8_0 TRUE)
endif(BUILD_CUDA8_0)

option(BUILD_CUDA9_0 "BUILD_CUDA9_0" OFF)
set(CUDA9_0 FALSE)
if(BUILD_CUDA9_0)
    set(CUDA9_0 TRUE)
endif(BUILD_CUDA9_0)

option(BUILD_CUDA10_0 "BUILD_CUDA10_0" OFF)
set(CUDA10_0 FALSE)
if(BUILD_CUDA10_0)
    set(CUDA10_0 TRUE)
endif(BUILD_CUDA10_0)

set(THIRD_PATH  ${CMAKE_SOURCE_DIR}/third_path)
if(NOT EXISTS ${THIRD_PATH}) 
    set(THIRD_PATH ${CMAKE_SOURCE_DIR}/../third_path)
endif()
message(STATUS "THIRD_PATH;${THIRD_PATH}")

ADD_THIRD_LIB(${THIRD_PATH})


IF(CUDA8_0)
    set(CUDA_TOOLKIT_ROOT_DIR ${THIRD_PATH}/cuda8_0)
    include(FindCUDA)
    message(STATUS "FindCUDA - CUDA_TOOLKIT_ROOT_DIR: ${THIRD_PATH}/cuda8_0; ")
ELSE(CUDA8_0)
    IF(IS_DIRECTORY ${THIRD_PATH}/cuda)
    set(CUDA_TOOLKIT_ROOT_DIR ${THIRD_PATH}/cuda)
    include(FindCUDA)
    ELSE(IS_DIRECTORY ${THIRD_PATH}/cuda)
    ENDIF()
ENDIF()

IF(CUDA9_0)
    set(CUDA_TOOLKIT_ROOT_DIR ${THIRD_PATH}/cuda9_0)
    include(FindCUDA)
    message(STATUS "FindCUDA - CUDA_TOOLKIT_ROOT_DIR: ${THIRD_PATH}/cuda9_0; ")
ELSE(CUDA9_0)
    IF(IS_DIRECTORY ${THIRD_PATH}/cuda)
    set(CUDA_TOOLKIT_ROOT_DIR ${THIRD_PATH}/cuda)
    include(FindCUDA)
    ELSE(IS_DIRECTORY ${THIRD_PATH}/cuda)
    ENDIF()
ENDIF()

IF(CUDA10_0)
    set(CUDA_TOOLKIT_ROOT_DIR ${THIRD_PATH}/cuda10_0)
    include(FindCUDA)
    message(STATUS "FindCUDA - CUDA_TOOLKIT_ROOT_DIR: ${THIRD_PATH}/cuda10_0; ")
ELSE(CUDA10_0)
    IF(IS_DIRECTORY ${THIRD_PATH}/cuda)
    set(CUDA_TOOLKIT_ROOT_DIR ${THIRD_PATH}/cuda)
    include(FindCUDA)
    ELSE(IS_DIRECTORY ${THIRD_PATH}/cuda)
    ENDIF()
ENDIF()
