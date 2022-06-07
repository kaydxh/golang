#include "stdio.h"
#include "sys/time.h"
#include "unistd.h"
#include "vector"
#include "opencv2/core/core.hpp"
#include "opencv2/highgui/highgui.hpp"
#include "opencv2/imgproc/imgproc.hpp"

int main() {
    struct timeval tv;
    gettimeofday(&tv,NULL);
    printf("cv::imread beg  %ldms\n", tv.tv_sec*1000 + tv.tv_usec/1000);
    cv::Mat img = cv::imread("/home/test_turbojpeg.jpg");
    gettimeofday(&tv,NULL);
    printf("cv::imread end %ldms\n", tv.tv_sec*1000 + tv.tv_usec/1000);

    std::vector<int> params;
    params.push_back(cv::IMWRITE_JPEG_QUALITY);
    params.push_back(100);

    gettimeofday(&tv,NULL);
    printf("cv::imencode beg %ldms\n", tv.tv_sec*1000 + tv.tv_usec/1000);
    std::vector<uchar> buf;
    cv::imencode(".jpg", img, buf, params);
    gettimeofday(&tv,NULL);
    printf("cv::imencode end %ldms\n", tv.tv_sec*1000 + tv.tv_usec/1000);
    return 0;
}
