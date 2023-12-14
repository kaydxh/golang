#include <Magick++.h>

#include <stdexcept>

#include "gocv.h"

namespace gocv {

MagicImage::~MagicImage() {}

void MagicImage::MagickInitializeMagick(
    const MagickInitializeMagickRequest& req,
    MagickInitializeMagickResponse& resp) {
    if (req.path.empty()) {
        Magick::InitializeMagick(nullptr);
        return;
    }
    Magick::InitializeMagick(req.path.c_str());
}

void MagicImage::MagickImageDecode(const MagickImageDecodeRequest& req,
                                   MagickImageDecodeResponse& resp) {
    if (req.image.empty()) {
        throw std::invalid_argument("image is empty");
    }

    Magick::Image image;
    try {
        Magick::Blob blob((void*)(req.image.data()), req.image.length());
        image.read(blob);
    } catch (Magick::Warning& w) {
        std::cout << "warn: " << w.what() << std::endl;
        // ignore warn
    } catch (Magick::Error& e) {
        std::cout << "a Magick++ error occurred: " << e.what() << std::endl;
        throw;
    } catch (...) {
        std::cout << "an unhandled error has occurred" << std::endl;
        throw;
    }

    int rows = image.rows();
    int columns = image.columns();
    if (rows <= 0 || columns <= 0) {
        throw std::invalid_argument(std::string("invalid image resolution [") +
                                    std::to_string(columns) +
                                    std::string(" x ") + std::to_string(rows) +
                                    std::string("]"));
    }

    cv::Mat mat;
    std::string map = req.target_color_space;
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
        image.write(0, 0, columns, rows, map, Magick::CharPixel, mat.data);
    } while (false);

    resp.mat = mat.clone();
    resp.magick = image.magick();
    resp.orientation_type =
        static_cast<gocv::OrientationType>(image.orientation());
    resp.colorspace_type =
        static_cast<gocv::ColorspaceType>(image.colorSpace());

    resp.rows = rows;
    resp.columns = columns;
    return;
}

}  // namespace gocv
