# CMAKE generated file: DO NOT EDIT!
# Generated by "Unix Makefiles" Generator, CMake Version 3.20

# Delete rule output on recipe failure.
.DELETE_ON_ERROR:

#=============================================================================
# Special targets provided by cmake.

# Disable implicit rules so canonical targets will work.
.SUFFIXES:

# Disable VCS-based implicit rules.
% : %,v

# Disable VCS-based implicit rules.
% : RCS/%

# Disable VCS-based implicit rules.
% : RCS/%,v

# Disable VCS-based implicit rules.
% : SCCS/s.%

# Disable VCS-based implicit rules.
% : s.%

.SUFFIXES: .hpux_make_needs_suffix_list

# Command-line flag to silence nested $(MAKE).
$(VERBOSE)MAKESILENT = -s

#Suppress display of executed commands.
$(VERBOSE).SILENT:

# A target that is always out of date.
cmake_force:
.PHONY : cmake_force

#=============================================================================
# Set environment variables for the build.

# The shell in which to execute make rules.
SHELL = /bin/sh

# The CMake executable.
CMAKE_COMMAND = /usr/local/bin/cmake

# The command to remove a file.
RM = /usr/local/bin/cmake -E rm -f

# Escaping for special characters.
EQUALS = =

# The top-level source directory on which CMake was run.
CMAKE_SOURCE_DIR = /data/home/kayxhding/workspace/github.com/kaydxh/golang/pkg/cgo

# The top-level build directory on which CMake was run.
CMAKE_BINARY_DIR = /data/home/kayxhding/workspace/github.com/kaydxh/golang/pkg/cgo/build

# Utility rule file for openapi-spec.

# Include any custom commands dependencies for this target.
include api/openapi-spec/CMakeFiles/openapi-spec.dir/compiler_depend.make

# Include the progress variables for this target.
include api/openapi-spec/CMakeFiles/openapi-spec.dir/progress.make

openapi-spec: api/openapi-spec/CMakeFiles/openapi-spec.dir/build.make
.PHONY : openapi-spec

# Rule to build all files generated by this target.
api/openapi-spec/CMakeFiles/openapi-spec.dir/build: openapi-spec
.PHONY : api/openapi-spec/CMakeFiles/openapi-spec.dir/build

api/openapi-spec/CMakeFiles/openapi-spec.dir/clean:
	cd /data/home/kayxhding/workspace/github.com/kaydxh/golang/pkg/cgo/build/api/openapi-spec && $(CMAKE_COMMAND) -P CMakeFiles/openapi-spec.dir/cmake_clean.cmake
.PHONY : api/openapi-spec/CMakeFiles/openapi-spec.dir/clean

api/openapi-spec/CMakeFiles/openapi-spec.dir/depend:
	cd /data/home/kayxhding/workspace/github.com/kaydxh/golang/pkg/cgo/build && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" /data/home/kayxhding/workspace/github.com/kaydxh/golang/pkg/cgo /data/home/kayxhding/workspace/github.com/kaydxh/golang/pkg/cgo/api/openapi-spec /data/home/kayxhding/workspace/github.com/kaydxh/golang/pkg/cgo/build /data/home/kayxhding/workspace/github.com/kaydxh/golang/pkg/cgo/build/api/openapi-spec /data/home/kayxhding/workspace/github.com/kaydxh/golang/pkg/cgo/build/api/openapi-spec/CMakeFiles/openapi-spec.dir/DependInfo.cmake --color=$(COLOR)
.PHONY : api/openapi-spec/CMakeFiles/openapi-spec.dir/depend
