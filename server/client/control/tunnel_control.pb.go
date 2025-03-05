// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: proto/control.proto

package control

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TunnelRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ClientId      string                 `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	TargetHost    string                 `protobuf:"bytes,2,opt,name=target_host,json=targetHost,proto3" json:"target_host,omitempty"`
	TargetPort    int32                  `protobuf:"varint,3,opt,name=target_port,json=targetPort,proto3" json:"target_port,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TunnelRequest) Reset() {
	*x = TunnelRequest{}
	mi := &file_proto_tunnel_control_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TunnelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TunnelRequest) ProtoMessage() {}

func (x *TunnelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tunnel_control_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TunnelRequest.ProtoReflect.Descriptor instead.
func (*TunnelRequest) Descriptor() ([]byte, []int) {
	return file_proto_tunnel_control_proto_rawDescGZIP(), []int{0}
}

func (x *TunnelRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *TunnelRequest) GetTargetHost() string {
	if x != nil {
		return x.TargetHost
	}
	return ""
}

func (x *TunnelRequest) GetTargetPort() int32 {
	if x != nil {
		return x.TargetPort
	}
	return 0
}

type TunnelResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TunnelId      string                 `protobuf:"bytes,1,opt,name=tunnel_id,json=tunnelId,proto3" json:"tunnel_id,omitempty"`
	PublicUrl     string                 `protobuf:"bytes,2,opt,name=public_url,json=publicUrl,proto3" json:"public_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TunnelResponse) Reset() {
	*x = TunnelResponse{}
	mi := &file_proto_tunnel_control_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TunnelResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TunnelResponse) ProtoMessage() {}

func (x *TunnelResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tunnel_control_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TunnelResponse.ProtoReflect.Descriptor instead.
func (*TunnelResponse) Descriptor() ([]byte, []int) {
	return file_proto_tunnel_control_proto_rawDescGZIP(), []int{1}
}

func (x *TunnelResponse) GetTunnelId() string {
	if x != nil {
		return x.TunnelId
	}
	return ""
}

func (x *TunnelResponse) GetPublicUrl() string {
	if x != nil {
		return x.PublicUrl
	}
	return ""
}

var File_proto_tunnel_control_proto protoreflect.FileDescriptor

var file_proto_tunnel_control_proto_rawDesc = string([]byte{
	0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x6e, 0x0a, 0x0d, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x68, 0x6f, 0x73, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x48, 0x6f,
	0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50,
	0x6f, 0x72, 0x74, 0x22, 0x4c, 0x0a, 0x0e, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c,
	0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x55, 0x72,
	0x6c, 0x32, 0x4c, 0x0a, 0x0d, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x12, 0x3b, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x75, 0x6e, 0x6e,
	0x65, 0x6c, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x75, 0x6e, 0x6e, 0x65,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x10, 0x5a, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_tunnel_control_proto_rawDescOnce sync.Once
	file_proto_tunnel_control_proto_rawDescData []byte
)

func file_proto_tunnel_control_proto_rawDescGZIP() []byte {
	file_proto_tunnel_control_proto_rawDescOnce.Do(func() {
		file_proto_tunnel_control_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_tunnel_control_proto_rawDesc), len(file_proto_tunnel_control_proto_rawDesc)))
	})
	return file_proto_tunnel_control_proto_rawDescData
}

var file_proto_tunnel_control_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_tunnel_control_proto_goTypes = []any{
	(*TunnelRequest)(nil),  // 0: proto.TunnelRequest
	(*TunnelResponse)(nil), // 1: proto.TunnelResponse
}
var file_proto_tunnel_control_proto_depIdxs = []int32{
	0, // 0: proto.TunnelControl.CreateTunnel:input_type -> proto.TunnelRequest
	1, // 1: proto.TunnelControl.CreateTunnel:output_type -> proto.TunnelResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_tunnel_control_proto_init() }
func file_proto_tunnel_control_proto_init() {
	if File_proto_tunnel_control_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_tunnel_control_proto_rawDesc), len(file_proto_tunnel_control_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_tunnel_control_proto_goTypes,
		DependencyIndexes: file_proto_tunnel_control_proto_depIdxs,
		MessageInfos:      file_proto_tunnel_control_proto_msgTypes,
	}.Build()
	File_proto_tunnel_control_proto = out.File
	file_proto_tunnel_control_proto_goTypes = nil
	file_proto_tunnel_control_proto_depIdxs = nil
}
