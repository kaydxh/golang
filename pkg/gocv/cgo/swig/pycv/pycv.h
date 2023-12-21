#ifndef GOCV_CGO_SWIG_PYCV_PYCV_H_
#define GOCV_CGO_SWIG_PYCV_PYCV_H_

#include <pybind11/embed.h>

#include <stdexcept>
#include <string>

namespace pybind11 {
class object;
}

namespace pycv {

struct LocalInitRequest {
    int gpu_id = -1;
    std::string sdk_dir;

    std::string model_dir;
};
struct LocalInitResponse {};

struct DoRequest {
    std::string arg1;
    std::string arg2;
};

struct DoResponse {};

#ifdef SWIG
class PyImage {
#else
class __attribute__((visibility("hidden"))) PyImage {
#endif
pybind11::object sdk_py;

public:
PyImage() = default;
~PyImage();

static void GlobalInit(const std::string& model_dir, int gpu_id);

static void GlobalRelease();

void LocalInit(const LocalInitRequest& req, LocalInitResponse& resp);

void Do(const DoRequest& req, DoResponse& resp);

};

}  // namespace pycv

#endif
