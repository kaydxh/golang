package gocv

/*
#cgo linux,amd64 pkg-config: opencv4 graphics-magick
#cgo !linux !amd64 CXXFLAGS: -D__CGO_UNKNOWN_PLATFORM__
#cgo CXXFLAGS: -std=c++11
#include <stdlib.h>
*/
import "C"
