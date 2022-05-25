add_custom_target(openapi-spec)

find_package(Protobuf)

if (EXISTS ${CMAKE_CURRENT_SOURCE_DIR}/thirdparty.srv.cmake)
   include(${CMAKE_CURRENT_SOURCE_DIR}/thirdparty.srv.cmake)
endif ()

GET_MODEL_DIRS(${CMAKE_CURRENT_SOURCE_DIR} _module_dirs)
FOREACH (_module_dir ${_module_dirs})
    if (EXISTS ${_module_dir})
      file(RELATIVE_PATH _rel_module_dir ${CMAKE_CURRENT_SOURCE_DIR} ${_module_dir})
      message(STATUS "add_subdirectory=${_rel_module_dir}")
      add_subdirectory(${_rel_module_dir})
      # add_dependencies(openapi-spec proto-${_rel_module_dir})
    endif()
ENDFOREACH ()
