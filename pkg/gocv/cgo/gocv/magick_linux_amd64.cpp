#include <Magick++.h>

#include <opencv2/opencv.hpp>

#include "api/openapi-spec/gocv/gocv.magick.pb.h"
#include "magick.h"

void sdk_gocv_magick_initialize_magick(void* req_data, int req_data_len,
                                       char** resp_data, int* resp_data_len) {
    sdk::api::gocv::MagickInitializeMagickResponse resp;

    try {
        do {
            if (!resp_data) {
                resp.mutable_error()->set_error_code(
                    sdk::types::code::Code::InvalidArgument);
                resp.mutable_error()->set_error_message("resp_data is nullptr");
                break;
            }
            if (!resp_data_len) {
                resp.mutable_error()->set_error_code(
                    sdk::types::code::Code::InvalidArgument);
                resp.mutable_error()->set_error_message(
                    "pointer of resp_data_len is nullptr");
                break;
            }

            sdk::api::gocv::MagickInitializeMagickRequest req;
            if (!req.ParseFromArray((char*)req_data, req_data_len)) {
                resp.mutable_error()->set_error_code(
                    sdk::types::code::Code::InvalidArgument);
                resp.mutable_error()->set_error_message("ParseFromArray");
                break;
            }
            if (req.path().empty()) {
                Magick::InitializeMagick(nullptr);
            } else {
                Magick::InitializeMagick(req.path().c_str());
            }
        } while (0);
    } catch (const std::exception& e) {
        resp.mutable_error()->set_error_code(sdk::types::code::Code::Internal);
        resp.mutable_error()->set_error_message(
            "Magick::InitializeMagick exception:" + std::string(e.what()));
    }

    *resp_data_len = resp.ByteSize();
    *resp_data = new char[*resp_data_len];
    resp.SerializeToArray(*resp_data, *resp_data_len);
}

void sdk_gocv_magick_image_decode(void* req_data, int req_data_len,
                                  char** resp_data, int* resp_data_len) {
    sdk::api::gocv::MagickImageDecodeResponse resp;

    try {
        do {
            if (!resp_data) {
                resp.mutable_error()->set_error_code(
                    sdk::types::code::Code::InvalidArgument);
                resp.mutable_error()->set_error_message("resp_data is nullptr");
                break;
            }
            if (!resp_data_len) {
                resp.mutable_error()->set_error_code(
                    sdk::types::code::Code::InvalidArgument);
                resp.mutable_error()->set_error_message(
                    "pointer of resp_data_len is nullptr");
                break;
            }

            sdk::api::gocv::MagickImageDecodeRequest req;
            if (!req.ParseFromArray((char*)req_data, req_data_len)) {
                resp.mutable_error()->set_error_code(
                    sdk::types::code::Code::InvalidArgument);
                resp.mutable_error()->set_error_message("ParseFromArray");
                break;
            }

            Magick::Image image;
            try {
                Magick::Blob blob((void*)(req.image().data()),
                                  req.image().length());
                image.read(blob);
            } catch (Magick::Warning& w) {
                std::cout << "warn: " << w.what() << std::endl;
                // ignore warn
            } catch (Magick::Error& e) {
                std::cout << "a Magick++ error occurred: " << e.what()
                          << std::endl;
                resp.mutable_error()->set_error_code(
                    sdk::types::code::Code::Internal);
                resp.mutable_error()->set_error_message(
                    "Magick::Blob read exception:" + std::string(e.what()));
                break;
            } catch (...) {
                std::cout << "an unhandled error has occurred" << std::endl;
                resp.mutable_error()->set_error_code(
                    sdk::types::code::Code::Internal);
                resp.mutable_error()->set_error_message(
                    "Magick::Blob read exception: an unhandled error has "
                    "occurred");
                break;
            }

            int rows = image.rows();
            int columns = image.columns();
            if (rows <= 0 || columns <= 0) {
                resp.mutable_error()->set_error_code(
                    sdk::types::code::Code::InvalidArgument);
                resp.mutable_error()->set_error_message(
                    std::string("invalid image resolution [") +
                    std::to_string(columns) + std::string(" x ") +
                    std::to_string(rows) + std::string("]"));
                break;
            }

            // https://www.imagemagick.org/Magick++/Image++.html
            cv::Mat mat;
            std::string map = req.target_color_space();
            do {
                image.colorSpace(Magick::RGBColorspace);
                if (map == "BGR") {
                    mat = ::cv::Mat(rows, columns, CV_8UC3);
                    image.write(0, 0, columns, rows, "BGR", Magick::CharPixel,
                                mat.data);
                    break;
                }
                if (map == "BGRA") {
                    mat = ::cv::Mat(rows, columns, CV_8UC4);
                    image.write(0, 0, columns, rows, "BGRA", Magick::CharPixel,
                                mat.data);
                    break;
                }
                if (map == "GRAY") {
                    image.type(Magick::GrayscaleType);
                    mat = ::cv::Mat(rows, columns, CV_8UC3);
                    image.write(0, 0, columns, rows, "BGR", Magick::CharPixel,
                                mat.data);
                    break;
                }
                if (map == "GRAYA") {
                    image.type(Magick::GrayscaleMatteType);
                    mat = ::cv::Mat(rows, columns, CV_8UC4);
                    image.write(0, 0, columns, rows, "BGRA", Magick::CharPixel,
                                mat.data);
                    break;
                }
                mat = ::cv::Mat(rows, columns, CV_8UC4);
                image.write(0, 0, columns, rows, map, Magick::CharPixel,
                            mat.data);

            } while (false);

            // set response
            resp.set_cv_mat_pointer(reinterpret_cast<int64>(new cv::Mat(mat)));
            resp.set_rows(image.rows());
            resp.set_columns(image.columns());
            resp.set_magick(image.magick());
            resp.set_orientation_type(
                static_cast<sdk::api::gocv::OrientationType>(
                    image.orientation()));
            resp.set_colorspace_type(
                static_cast<sdk::api::gocv::ColorspaceType>(
                    image.colorSpace()));

        } while (0);
    } catch (const std::exception& e) {
        resp.mutable_error()->set_error_code(sdk::types::code::Code::Internal);
        resp.mutable_error()->set_error_message("Magick::Blob read exception:" +
                                                std::string(e.what()));
    }

    *resp_data_len = resp.ByteSize();
    *resp_data = new char[*resp_data_len];
    resp.SerializeToArray(*resp_data, *resp_data_len);
}
