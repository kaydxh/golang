#include <Magick++.h>

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
        } while (false);
    }
    catch (const std::exception& e) {
        resp.mutable_error()->set_error_code(sdk::types::code::Code::Internal);
        resp.mutable_error()->set_error_message(
            "Magick::InitializeMagick exception:" + std::string(e.what()));
    }

    *resp_data_len = resp.ByteSize();
    *resp_data = new char[*resp_data_len];
    resp.SerializeToArray(*resp_data, *resp_data_len);
}
