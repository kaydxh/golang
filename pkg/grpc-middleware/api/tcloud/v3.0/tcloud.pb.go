// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        v3.13.0
// source: pkg/grpc-middleware/api/tcloud/v3.0/tcloud.proto

package tcloud

import (
	_ "github.com/golang/protobuf/ptypes/duration"
	_struct "github.com/golang/protobuf/ptypes/struct"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http/request.go
type TCloudBaseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HttpMethod string            `protobuf:"bytes,1,opt,name=http_method,json=httpMethod,proto3" json:"http_method,omitempty"`
	Scheme     string            `protobuf:"bytes,2,opt,name=scheme,proto3" json:"scheme,omitempty"`
	RootDomain string            `protobuf:"bytes,3,opt,name=root_domain,json=rootDomain,proto3" json:"root_domain,omitempty"`
	Domain     string            `protobuf:"bytes,4,opt,name=domain,proto3" json:"domain,omitempty"`
	Path       string            `protobuf:"bytes,5,opt,name=path,proto3" json:"path,omitempty"`
	Params     map[string]string `protobuf:"bytes,6,rep,name=params,proto3" json:"params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	FormParams map[string]string `protobuf:"bytes,7,rep,name=form_params,json=formParams,proto3" json:"form_params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Service    string            `protobuf:"bytes,8,opt,name=service,proto3" json:"service,omitempty"`
	Version    string            `protobuf:"bytes,9,opt,name=version,json=Version,proto3" json:"version,omitempty"`
	Action     string            `protobuf:"bytes,10,opt,name=action,json=Action,proto3" json:"action,omitempty"`
}

func (x *TCloudBaseRequest) Reset() {
	*x = TCloudBaseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TCloudBaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TCloudBaseRequest) ProtoMessage() {}

func (x *TCloudBaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TCloudBaseRequest.ProtoReflect.Descriptor instead.
func (*TCloudBaseRequest) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescGZIP(), []int{0}
}

func (x *TCloudBaseRequest) GetHttpMethod() string {
	if x != nil {
		return x.HttpMethod
	}
	return ""
}

func (x *TCloudBaseRequest) GetScheme() string {
	if x != nil {
		return x.Scheme
	}
	return ""
}

func (x *TCloudBaseRequest) GetRootDomain() string {
	if x != nil {
		return x.RootDomain
	}
	return ""
}

func (x *TCloudBaseRequest) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *TCloudBaseRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *TCloudBaseRequest) GetParams() map[string]string {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *TCloudBaseRequest) GetFormParams() map[string]string {
	if x != nil {
		return x.FormParams
	}
	return nil
}

func (x *TCloudBaseRequest) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *TCloudBaseRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *TCloudBaseRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

// https://github.com/TencentCloud/tencentcloud-sdk-go/blob/master/tencentcloud/common/http/response.go
type TCloudResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response *_struct.Struct `protobuf:"bytes,1,opt,name=response,json=Response,proto3" json:"response,omitempty"`
}

func (x *TCloudResponse) Reset() {
	*x = TCloudResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TCloudResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TCloudResponse) ProtoMessage() {}

func (x *TCloudResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TCloudResponse.ProtoReflect.Descriptor instead.
func (*TCloudResponse) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescGZIP(), []int{1}
}

func (x *TCloudResponse) GetResponse() *_struct.Struct {
	if x != nil {
		return x.Response
	}
	return nil
}

type TCloudErrorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ErrorResponse
	Response *ErrorResponse `protobuf:"bytes,1,opt,name=response,json=Response,proto3" json:"response,omitempty"`
}

func (x *TCloudErrorResponse) Reset() {
	*x = TCloudErrorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TCloudErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TCloudErrorResponse) ProtoMessage() {}

func (x *TCloudErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TCloudErrorResponse.ProtoReflect.Descriptor instead.
func (*TCloudErrorResponse) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescGZIP(), []int{2}
}

func (x *TCloudErrorResponse) GetResponse() *ErrorResponse {
	if x != nil {
		return x.Response
	}
	return nil
}

// ErrorResponse
type ErrorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error     *TCloudError `protobuf:"bytes,1,opt,name=error,json=Error,proto3" json:"error,omitempty"`
	RequestId string       `protobuf:"bytes,2,opt,name=request_id,json=RequestId,proto3" json:"request_id,omitempty"`
}

func (x *ErrorResponse) Reset() {
	*x = ErrorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorResponse) ProtoMessage() {}

func (x *ErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorResponse.ProtoReflect.Descriptor instead.
func (*ErrorResponse) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescGZIP(), []int{3}
}

func (x *ErrorResponse) GetError() *TCloudError {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *ErrorResponse) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

type TCloudError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string `protobuf:"bytes,1,opt,name=code,json=Code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,json=Message,proto3" json:"message,omitempty"`
}

func (x *TCloudError) Reset() {
	*x = TCloudError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TCloudError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TCloudError) ProtoMessage() {}

func (x *TCloudError) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TCloudError.ProtoReflect.Descriptor instead.
func (*TCloudError) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescGZIP(), []int{4}
}

func (x *TCloudError) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *TCloudError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto protoreflect.FileDescriptor

var file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDesc = []byte{
	0x0a, 0x30, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x6d, 0x69, 0x64, 0x64, 0x6c,
	0x65, 0x77, 0x61, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x2f, 0x76, 0x33, 0x2e, 0x30, 0x2f, 0x74, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0d, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76,
	0x33, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xf8, 0x03, 0x0a, 0x11, 0x54, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x42, 0x61, 0x73, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x6d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x68, 0x74, 0x74, 0x70,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x12, 0x1f,
	0x0a, 0x0b, 0x72, 0x6f, 0x6f, 0x74, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x6f, 0x6f, 0x74, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12,
	0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x44, 0x0a, 0x06, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x74, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x33, 0x2e, 0x54, 0x43, 0x6c, 0x6f,
	0x75, 0x64, 0x42, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x12, 0x51, 0x0a, 0x0b, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x2e, 0x76, 0x33, 0x2e, 0x54, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x42, 0x61, 0x73,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x72, 0x6d, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x66, 0x6f, 0x72, 0x6d, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x1a, 0x39, 0x0a, 0x0b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3d, 0x0a, 0x0f, 0x46,
	0x6f, 0x72, 0x6d, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x45, 0x0a, 0x0e, 0x54, 0x43,
	0x6c, 0x6f, 0x75, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x4f, 0x0a, 0x13, 0x54, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x74, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x33, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x60, 0x0a, 0x0d, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e,
	0x76, 0x33, 0x2e, 0x54, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x64, 0x22, 0x3b, 0x0a, 0x0b, 0x54, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6b, 0x61, 0x79, 0x64, 0x78, 0x68, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x76, 0x33, 0x2e,
	0x30, 0x3b, 0x74, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescOnce sync.Once
	file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescData = file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDesc
)

func file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescGZIP() []byte {
	file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescOnce.Do(func() {
		file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescData)
	})
	return file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDescData
}

var file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_goTypes = []interface{}{
	(*TCloudBaseRequest)(nil),   // 0: api.tcloud.v3.TCloudBaseRequest
	(*TCloudResponse)(nil),      // 1: api.tcloud.v3.TCloudResponse
	(*TCloudErrorResponse)(nil), // 2: api.tcloud.v3.TCloudErrorResponse
	(*ErrorResponse)(nil),       // 3: api.tcloud.v3.ErrorResponse
	(*TCloudError)(nil),         // 4: api.tcloud.v3.TCloudError
	nil,                         // 5: api.tcloud.v3.TCloudBaseRequest.ParamsEntry
	nil,                         // 6: api.tcloud.v3.TCloudBaseRequest.FormParamsEntry
	(*_struct.Struct)(nil),      // 7: google.protobuf.Struct
}
var file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_depIdxs = []int32{
	5, // 0: api.tcloud.v3.TCloudBaseRequest.params:type_name -> api.tcloud.v3.TCloudBaseRequest.ParamsEntry
	6, // 1: api.tcloud.v3.TCloudBaseRequest.form_params:type_name -> api.tcloud.v3.TCloudBaseRequest.FormParamsEntry
	7, // 2: api.tcloud.v3.TCloudResponse.response:type_name -> google.protobuf.Struct
	3, // 3: api.tcloud.v3.TCloudErrorResponse.response:type_name -> api.tcloud.v3.ErrorResponse
	4, // 4: api.tcloud.v3.ErrorResponse.error:type_name -> api.tcloud.v3.TCloudError
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_init() }
func file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_init() {
	if File_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TCloudBaseRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TCloudResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TCloudErrorResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrorResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TCloudError); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_goTypes,
		DependencyIndexes: file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_depIdxs,
		MessageInfos:      file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_msgTypes,
	}.Build()
	File_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto = out.File
	file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_rawDesc = nil
	file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_goTypes = nil
	file_pkg_grpc_middleware_api_tcloud_v3_0_tcloud_proto_depIdxs = nil
}