/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        v3.13.0
// source: pkg/middleware/api/trivial/v1/api.proto

package v1

import (
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

type TrivialResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId string `protobuf:"bytes,1,opt,name=session_id,proto3" json:"session_id,omitempty"`
	ErrorCode int32  `protobuf:"varint,2,opt,name=error_code,json=errorcode,proto3" json:"error_code,omitempty"`
	ErrorMsg  string `protobuf:"bytes,3,opt,name=error_msg,json=errormsg,proto3" json:"error_msg,omitempty"`
}

func (x *TrivialResponse) Reset() {
	*x = TrivialResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_middleware_api_trivial_v1_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrivialResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrivialResponse) ProtoMessage() {}

func (x *TrivialResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_middleware_api_trivial_v1_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrivialResponse.ProtoReflect.Descriptor instead.
func (*TrivialResponse) Descriptor() ([]byte, []int) {
	return file_pkg_middleware_api_trivial_v1_api_proto_rawDescGZIP(), []int{0}
}

func (x *TrivialResponse) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *TrivialResponse) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *TrivialResponse) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

type TrivialErrorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId string `protobuf:"bytes,1,opt,name=session_id,proto3" json:"session_id,omitempty"`
	ErrorCode int32  `protobuf:"varint,2,opt,name=error_code,json=errorcode,proto3" json:"error_code,omitempty"`
	ErrorMsg  string `protobuf:"bytes,3,opt,name=error_msg,json=errormsg,proto3" json:"error_msg,omitempty"`
}

func (x *TrivialErrorResponse) Reset() {
	*x = TrivialErrorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_middleware_api_trivial_v1_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrivialErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrivialErrorResponse) ProtoMessage() {}

func (x *TrivialErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_middleware_api_trivial_v1_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrivialErrorResponse.ProtoReflect.Descriptor instead.
func (*TrivialErrorResponse) Descriptor() ([]byte, []int) {
	return file_pkg_middleware_api_trivial_v1_api_proto_rawDescGZIP(), []int{1}
}

func (x *TrivialErrorResponse) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *TrivialErrorResponse) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *TrivialErrorResponse) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

var File_pkg_middleware_api_trivial_v1_api_proto protoreflect.FileDescriptor

var file_pkg_middleware_api_trivial_v1_api_proto_rawDesc = []byte{
	0x0a, 0x27, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x72, 0x69, 0x76, 0x69, 0x61, 0x6c, 0x2f, 0x76, 0x31, 0x2f,
	0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x61, 0x70, 0x69, 0x2e, 0x74,
	0x72, 0x69, 0x76, 0x69, 0x61, 0x6c, 0x2e, 0x76, 0x31, 0x22, 0x6d, 0x0a, 0x0f, 0x54, 0x72, 0x69,
	0x76, 0x69, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x6d, 0x73, 0x67, 0x22, 0x72, 0x0a, 0x14, 0x54, 0x72, 0x69, 0x76,
	0x69, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64,
	0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x6d, 0x73, 0x67, 0x42, 0x3b, 0x5a, 0x39,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x61, 0x79, 0x64, 0x78,
	0x68, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x69, 0x64,
	0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x72, 0x69, 0x76,
	0x69, 0x61, 0x6c, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_pkg_middleware_api_trivial_v1_api_proto_rawDescOnce sync.Once
	file_pkg_middleware_api_trivial_v1_api_proto_rawDescData = file_pkg_middleware_api_trivial_v1_api_proto_rawDesc
)

func file_pkg_middleware_api_trivial_v1_api_proto_rawDescGZIP() []byte {
	file_pkg_middleware_api_trivial_v1_api_proto_rawDescOnce.Do(func() {
		file_pkg_middleware_api_trivial_v1_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_middleware_api_trivial_v1_api_proto_rawDescData)
	})
	return file_pkg_middleware_api_trivial_v1_api_proto_rawDescData
}

var file_pkg_middleware_api_trivial_v1_api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_middleware_api_trivial_v1_api_proto_goTypes = []interface{}{
	(*TrivialResponse)(nil),      // 0: api.trivial.v1.TrivialResponse
	(*TrivialErrorResponse)(nil), // 1: api.trivial.v1.TrivialErrorResponse
}
var file_pkg_middleware_api_trivial_v1_api_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_middleware_api_trivial_v1_api_proto_init() }
func file_pkg_middleware_api_trivial_v1_api_proto_init() {
	if File_pkg_middleware_api_trivial_v1_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_middleware_api_trivial_v1_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrivialResponse); i {
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
		file_pkg_middleware_api_trivial_v1_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrivialErrorResponse); i {
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
			RawDescriptor: file_pkg_middleware_api_trivial_v1_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_middleware_api_trivial_v1_api_proto_goTypes,
		DependencyIndexes: file_pkg_middleware_api_trivial_v1_api_proto_depIdxs,
		MessageInfos:      file_pkg_middleware_api_trivial_v1_api_proto_msgTypes,
	}.Build()
	File_pkg_middleware_api_trivial_v1_api_proto = out.File
	file_pkg_middleware_api_trivial_v1_api_proto_rawDesc = nil
	file_pkg_middleware_api_trivial_v1_api_proto_goTypes = nil
	file_pkg_middleware_api_trivial_v1_api_proto_depIdxs = nil
}
