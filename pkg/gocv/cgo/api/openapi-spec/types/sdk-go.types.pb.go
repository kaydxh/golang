// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: types/sdk-go.types.proto

package types

import (
	code "github.com/kaydxh/golang/pkg/cgo/api/openapi-spec/types/code"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Symbols defined in public import of api/openapi-spec/types/code/sdk-go.code.proto.

type Code = code.Code

const Code_OK = code.Code_OK
const Code_Canceled = code.Code_Canceled
const Code_Unknown = code.Code_Unknown
const Code_InvalidArgument = code.Code_InvalidArgument
const Code_DeadlineExceeded = code.Code_DeadlineExceeded
const Code_NotFound = code.Code_NotFound
const Code_AlreadyExists = code.Code_AlreadyExists
const Code_PermissionDenied = code.Code_PermissionDenied
const Code_ResourceExhausted = code.Code_ResourceExhausted
const Code_FailedPrecondition = code.Code_FailedPrecondition
const Code_Aborted = code.Code_Aborted
const Code_OutOfRange = code.Code_OutOfRange
const Code_Unimplemented = code.Code_Unimplemented
const Code_Internal = code.Code_Internal
const Code_Unavailable = code.Code_Unavailable
const Code_DataLoss = code.Code_DataLoss
const Code_Unauthenticated = code.Code_Unauthenticated

var Code_name = code.Code_name
var Code_value = code.Code_value

type CgoError = code.CgoError

var File_types_sdk_go_types_proto protoreflect.FileDescriptor

var file_types_sdk_go_types_proto_rawDesc = []byte{
	0x0a, 0x18, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x2e, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x73, 0x64, 0x6b, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x1a, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x6f,
	0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x2d, 0x73, 0x70, 0x65, 0x63, 0x2f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x2e, 0x63, 0x6f,
	0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x61, 0x79, 0x64, 0x78, 0x68, 0x2f, 0x67, 0x6f,
	0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x2d, 0x73, 0x70, 0x65, 0x63, 0x2f, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x73, 0x50, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var file_types_sdk_go_types_proto_goTypes = []interface{}{}
var file_types_sdk_go_types_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_types_sdk_go_types_proto_init() }
func file_types_sdk_go_types_proto_init() {
	if File_types_sdk_go_types_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_types_sdk_go_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_types_sdk_go_types_proto_goTypes,
		DependencyIndexes: file_types_sdk_go_types_proto_depIdxs,
	}.Build()
	File_types_sdk_go_types_proto = out.File
	file_types_sdk_go_types_proto_rawDesc = nil
	file_types_sdk_go_types_proto_goTypes = nil
	file_types_sdk_go_types_proto_depIdxs = nil
}