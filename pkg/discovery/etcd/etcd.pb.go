// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.2
// source: pkg/discovery/etcd/etcd.proto

package etcd

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Etcd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enabled            bool     `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Addresses          []string `protobuf:"bytes,2,rep,name=addresses,proto3" json:"addresses,omitempty"`
	Username           string   `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	Password           string   `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	MaxCallSendMsgSize int32    `protobuf:"varint,5,opt,name=max_call_send_msg_size,json=maxCallSendMsgSize,proto3" json:"max_call_send_msg_size,omitempty"` //If 0, it defaults to 2.0 MiB (2 * 1024 * 1024)
	// MaxCallRecvMsgSize is the client-side response receive limit.
	// If 0, it defaults to "math.MaxInt32", because range response can
	// easily exceed request send limits.
	// Make sure that "MaxCallRecvMsgSize" >= server-side default
	// send/recv limit.
	// ("--max-request-bytes" flag to etcd or
	// "embed.Config.MaxRequestBytes").
	MaxCallRecvMsgSize int32 `protobuf:"varint,6,opt,name=max_call_recv_msg_size,json=maxCallRecvMsgSize,proto3" json:"max_call_recv_msg_size,omitempty"`
	// AutoSyncInterval is the interval to update endpoints with its latest members.
	// 0 disables auto-sync. By default auto-sync is disabled.
	AutoSyncInterval  *durationpb.Duration `protobuf:"bytes,7,opt,name=auto_sync_interval,json=autoSyncInterval,proto3" json:"auto_sync_interval,omitempty"`
	DialTimeout       *durationpb.Duration `protobuf:"bytes,8,opt,name=dial_timeout,json=dialTimeout,proto3" json:"dial_timeout,omitempty"`
	MaxWaitDuration   *durationpb.Duration `protobuf:"bytes,12,opt,name=max_wait_duration,json=maxWaitDuration,proto3" json:"max_wait_duration,omitempty"`
	FailAfterDuration *durationpb.Duration `protobuf:"bytes,13,opt,name=fail_after_duration,json=failAfterDuration,proto3" json:"fail_after_duration,omitempty"`
	WatchPaths        []string             `protobuf:"bytes,14,rep,name=watch_paths,json=watchPaths,proto3" json:"watch_paths,omitempty"`
}

func (x *Etcd) Reset() {
	*x = Etcd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_discovery_etcd_etcd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Etcd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Etcd) ProtoMessage() {}

func (x *Etcd) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_discovery_etcd_etcd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Etcd.ProtoReflect.Descriptor instead.
func (*Etcd) Descriptor() ([]byte, []int) {
	return file_pkg_discovery_etcd_etcd_proto_rawDescGZIP(), []int{0}
}

func (x *Etcd) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *Etcd) GetAddresses() []string {
	if x != nil {
		return x.Addresses
	}
	return nil
}

func (x *Etcd) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Etcd) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Etcd) GetMaxCallSendMsgSize() int32 {
	if x != nil {
		return x.MaxCallSendMsgSize
	}
	return 0
}

func (x *Etcd) GetMaxCallRecvMsgSize() int32 {
	if x != nil {
		return x.MaxCallRecvMsgSize
	}
	return 0
}

func (x *Etcd) GetAutoSyncInterval() *durationpb.Duration {
	if x != nil {
		return x.AutoSyncInterval
	}
	return nil
}

func (x *Etcd) GetDialTimeout() *durationpb.Duration {
	if x != nil {
		return x.DialTimeout
	}
	return nil
}

func (x *Etcd) GetMaxWaitDuration() *durationpb.Duration {
	if x != nil {
		return x.MaxWaitDuration
	}
	return nil
}

func (x *Etcd) GetFailAfterDuration() *durationpb.Duration {
	if x != nil {
		return x.FailAfterDuration
	}
	return nil
}

func (x *Etcd) GetWatchPaths() []string {
	if x != nil {
		return x.WatchPaths
	}
	return nil
}

var File_pkg_discovery_etcd_etcd_proto protoreflect.FileDescriptor

var file_pkg_discovery_etcd_etcd_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x6b, 0x67, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f,
	0x65, 0x74, 0x63, 0x64, 0x2f, 0x65, 0x74, 0x63, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x15, 0x67, 0x6f, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x2e, 0x65, 0x74, 0x63, 0x64, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x98, 0x04, 0x0a, 0x04, 0x45, 0x74, 0x63, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12,
	0x32, 0x0a, 0x16, 0x6d, 0x61, 0x78, 0x5f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x73, 0x65, 0x6e, 0x64,
	0x5f, 0x6d, 0x73, 0x67, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x12, 0x6d, 0x61, 0x78, 0x43, 0x61, 0x6c, 0x6c, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x53,
	0x69, 0x7a, 0x65, 0x12, 0x32, 0x0a, 0x16, 0x6d, 0x61, 0x78, 0x5f, 0x63, 0x61, 0x6c, 0x6c, 0x5f,
	0x72, 0x65, 0x63, 0x76, 0x5f, 0x6d, 0x73, 0x67, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x12, 0x6d, 0x61, 0x78, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x63, 0x76,
	0x4d, 0x73, 0x67, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x47, 0x0a, 0x12, 0x61, 0x75, 0x74, 0x6f, 0x5f,
	0x73, 0x79, 0x6e, 0x63, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10,
	0x61, 0x75, 0x74, 0x6f, 0x53, 0x79, 0x6e, 0x63, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c,
	0x12, 0x3c, 0x0a, 0x0c, 0x64, 0x69, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0b, 0x64, 0x69, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x45,
	0x0a, 0x11, 0x6d, 0x61, 0x78, 0x5f, 0x77, 0x61, 0x69, 0x74, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x6d, 0x61, 0x78, 0x57, 0x61, 0x69, 0x74, 0x44, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x49, 0x0a, 0x13, 0x66, 0x61, 0x69, 0x6c, 0x5f, 0x61, 0x66,
	0x74, 0x65, 0x72, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x11, 0x66,
	0x61, 0x69, 0x6c, 0x41, 0x66, 0x74, 0x65, 0x72, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x73, 0x18,
	0x0e, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x61, 0x74, 0x63, 0x68, 0x50, 0x61, 0x74, 0x68,
	0x73, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6b, 0x61, 0x79, 0x64, 0x78, 0x68, 0x2f, 0x67, 0x6f, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x65, 0x74, 0x63, 0x64, 0x3b, 0x65, 0x74, 0x63,
	0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_discovery_etcd_etcd_proto_rawDescOnce sync.Once
	file_pkg_discovery_etcd_etcd_proto_rawDescData = file_pkg_discovery_etcd_etcd_proto_rawDesc
)

func file_pkg_discovery_etcd_etcd_proto_rawDescGZIP() []byte {
	file_pkg_discovery_etcd_etcd_proto_rawDescOnce.Do(func() {
		file_pkg_discovery_etcd_etcd_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_discovery_etcd_etcd_proto_rawDescData)
	})
	return file_pkg_discovery_etcd_etcd_proto_rawDescData
}

var file_pkg_discovery_etcd_etcd_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_discovery_etcd_etcd_proto_goTypes = []interface{}{
	(*Etcd)(nil),                // 0: go.pkg.discovery.etcd.Etcd
	(*durationpb.Duration)(nil), // 1: google.protobuf.Duration
}
var file_pkg_discovery_etcd_etcd_proto_depIdxs = []int32{
	1, // 0: go.pkg.discovery.etcd.Etcd.auto_sync_interval:type_name -> google.protobuf.Duration
	1, // 1: go.pkg.discovery.etcd.Etcd.dial_timeout:type_name -> google.protobuf.Duration
	1, // 2: go.pkg.discovery.etcd.Etcd.max_wait_duration:type_name -> google.protobuf.Duration
	1, // 3: go.pkg.discovery.etcd.Etcd.fail_after_duration:type_name -> google.protobuf.Duration
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pkg_discovery_etcd_etcd_proto_init() }
func file_pkg_discovery_etcd_etcd_proto_init() {
	if File_pkg_discovery_etcd_etcd_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_discovery_etcd_etcd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Etcd); i {
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
			RawDescriptor: file_pkg_discovery_etcd_etcd_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_discovery_etcd_etcd_proto_goTypes,
		DependencyIndexes: file_pkg_discovery_etcd_etcd_proto_depIdxs,
		MessageInfos:      file_pkg_discovery_etcd_etcd_proto_msgTypes,
	}.Build()
	File_pkg_discovery_etcd_etcd_proto = out.File
	file_pkg_discovery_etcd_etcd_proto_rawDesc = nil
	file_pkg_discovery_etcd_etcd_proto_goTypes = nil
	file_pkg_discovery_etcd_etcd_proto_depIdxs = nil
}
