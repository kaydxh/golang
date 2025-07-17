#include <pybind11/embed.h>

#include <memory>
#include <mutex>
#include <sstream>
#include <stdexcept>

#include "pycv.h"

namespace py = pybind11;
using namespace pybind11::literals;  // to bring in the `_a` literal
#ifdef __cplusplus
#define UIKIT_EXTERN extern "C" __attribute__((visibility("default")))
#else
#define UIKIT_EXTERN extern __attribute__((visibility("default")))
#endif
static std::once_flag python_once;
static std::unique_ptr<py::scoped_interpreter> python_interpreter_guard;
static std::unique_ptr<py::gil_scoped_release>
    python_interpreter_gil_scoped_release;

pycv::PyImage::~PyImage() {
    py::gil_scoped_acquire acquire;
    sdk_py.dec_ref();
}

void pycv::PyImage::GlobalInit(const std::string& model_dir, int gpu_id) {
    std::call_once(python_once, [] {
        python_interpreter_guard = std::unique_ptr<py::scoped_interpreter>(
            new py::scoped_interpreter());
        { pybind11::gil_scoped_acquire acquire; }
        python_interpreter_gil_scoped_release =
            std::unique_ptr<py::gil_scoped_release>(
                new py::gil_scoped_release());
    });
}

void pycv::PyImage::GlobalRelease() {}

void pycv::PyImage::LocalInit(const LocalInitRequest& req,
                              LocalInitResponse& resp) {
    // acquiring GIL as toPyObject creates new py::object
    // without grabbing the GIL.
    py::gil_scoped_acquire acquire;

    // https://pybind11.readthedocs.io/en/stable/advanced/pycpp/object.html#calling-python-methods
    // https://github.com/pybind/pybind11/issues/1201
    auto py_cv = py::module::import("pycv");
    if (py_cv.is_none()) {
        throw std::runtime_error("module pycv not found");
    }

    if (!py::hasattr(py_cv, "CVSDK")) {
        throw std::runtime_error("python class pycv.CVSD not found");
    }

    auto CVSDK = py_cv.attr("CVSDK");
    sdk_py = CVSDK();
    if (!py::hasattr(sdk_py, "init")) {
        throw std::runtime_error("python method py.CVSDK::init not found");
    }
    if (!py::hasattr(sdk_py, "do")) {
        throw std::runtime_error("python method py.CVSDK::do not found");
    }
    {  // 模型初始化
        py::dict dict("model_dir"_a = req.model_dir, "gpu_id"_a = req.gpu_id);
        dict["sdk_dir"] = req.sdk_dir;
        auto init_resp_py = sdk_py.attr("init")(**dict);
        auto init_resp = py::dict(init_resp_py);
        if (init_resp.contains("err")) {
            std::stringstream ss;
            ss << "py.CVSDK::init returns (";
            ss << init_resp["err"].cast<std::string>();
            ss << ")";
            throw std::runtime_error(ss.str());
        }
    }
}

void pycv::PyImage::Do(const DoRequest& req, DoResponse& resp) {
    py::gil_scoped_acquire acquire;

    if (sdk_py.is_none()) {
        throw std::runtime_error(
            "method py.CVSDK not initialized, call init at first");
    }
    if (!py::hasattr(sdk_py, "do")) {
        throw std::runtime_error("method py.do not found");
    }
    {
        py::dict dict("arg1"_a = py::bytes(req.arg1));
        dict["arg2"] = req.arg2;
        auto do_resp_py = sdk_py.attr("do")(**dict);
        auto do_resp = py::dict(do_resp_py);
        if (do_resp.contains("err")) {
            std::stringstream ss;
            ss << "py.CVSDK::do do_resp returns (";
            ss << do_resp["err"].cast<std::string>();
            ss << ")";
            throw std::runtime_error(ss.str());
        }
    }
}

