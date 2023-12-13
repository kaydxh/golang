#ifndef GOCV_CGO_SWIG_GOCV_GOCV_H_
#define GOCV_CGO_SWIG_GOCV_GOCV_H_

#include <opencv2/opencv.hpp>
#include <string>

namespace gocv {
struct MagickInitializeMagickRequest {
    std::string path;
};

struct MagickInitializeMagickResponse {};

struct MagickImageDecodeRequest {
    std::string image;
    std::string target_color_space;
};

enum OrientationType {
    UndefinedOrientation = 0,
    TopLeftOrientation = 1,
    TopRightOrientation = 2,
    BottomRightOrientation = 3,
    BottomLeftOrientation = 4,
    LeftTopOrientation = 5,
    RightTopOrientation = 6,
    RightBottomOrientation = 7,
    LeftBottomOrientation = 8
};

// graphics-magick/include/magick/colorspace.h
enum ColorspaceType {
    UndefinedColorspace = 0,
    RGBColorspace = 1,  /* Plain old RGB colorspace */
    GRAYColorspace = 2, /* Plain old full-range grayscale */
    TransparentColorspace =
        3, /* RGB but preserve matte channel during quantize */
    OHTAColorspace = 4,
    XYZColorspace = 5, /* CIE XYZ */
    YCCColorspace = 6, /* Kodak PhotoCD PhotoYCC */
    YIQColorspace = 7,
    YPbPrColorspace = 8,
    YUVColorspace = 9,
    CMYKColorspace = 10, /* Cyan, magenta, yellow, black, alpha */
    sRGBColorspace = 11, /* Kodak PhotoCD sRGB */
    HSLColorspace = 12,  /* Hue, saturation, luminosity */
    HWBColorspace = 13,  /* Hue, whiteness, blackness */
    LABColorspace =
        14, /* LAB colorspace not supported yet other than via lcms */
    CineonLogRGBColorspace =
        15, /* RGB data with Cineon Log scaling, 2.048 density range */
    Rec601LumaColorspace = 16,  /* Luma (Y) according to ITU-R 601 */
    Rec601YCbCrColorspace = 17, /* YCbCr according to ITU-R 601 */
    Rec709LumaColorspace = 18,  /* Luma (Y) according to ITU-R 709 */
    Rec709YCbCrColorspace = 19  /* YCbCr according to ITU-R 709 */
};

struct MagickImageDecodeResponse {
    cv::Mat mat;
    int64_t rows;        // height
    int64_t columns;     // width
    std::string magick;  // File type magick identifier (.e.g "GIF")
    OrientationType orientation_type;
    ColorspaceType colorspace_type;
};

#ifdef SWIG
class MagicImage {
#else
class __attribute__((visibility("hidden"))) MagicImage {
#endif
public:
MagicImage() = default;
~MagicImage();

void MagickImageDecode(const MagickImageDecodeRequest& req,
                       MagickImageDecodeResponse* resp);

};

}  // namespace gocv

#endif
