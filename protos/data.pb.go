// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: protos/data.proto

package data

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

type PacketType int32

const (
	PacketType_REQUEST   PacketType = 0
	PacketType_RESPONSE  PacketType = 1
	PacketType_KEEPALIVE PacketType = 2
)

// Enum value maps for PacketType.
var (
	PacketType_name = map[int32]string{
		0: "REQUEST",
		1: "RESPONSE",
		2: "KEEPALIVE",
	}
	PacketType_value = map[string]int32{
		"REQUEST":   0,
		"RESPONSE":  1,
		"KEEPALIVE": 2,
	}
)

func (x PacketType) Enum() *PacketType {
	p := new(PacketType)
	*p = x
	return p
}

func (x PacketType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PacketType) Descriptor() protoreflect.EnumDescriptor {
	return file_protos_data_proto_enumTypes[0].Descriptor()
}

func (PacketType) Type() protoreflect.EnumType {
	return &file_protos_data_proto_enumTypes[0]
}

func (x PacketType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PacketType.Descriptor instead.
func (PacketType) EnumDescriptor() ([]byte, []int) {
	return file_protos_data_proto_rawDescGZIP(), []int{0}
}

type DataPacket struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TunnelId      string                 `protobuf:"bytes,1,opt,name=tunnel_id,json=tunnelId,proto3" json:"tunnel_id,omitempty"`
	RequestId     string                 `protobuf:"bytes,2,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Type          PacketType             `protobuf:"varint,3,opt,name=type,proto3,enum=data.PacketType" json:"type,omitempty"`
	Method        string                 `protobuf:"bytes,4,opt,name=method,proto3" json:"method,omitempty"`
	Url           string                 `protobuf:"bytes,5,opt,name=url,proto3" json:"url,omitempty"`
	Headers       map[string]string      `protobuf:"bytes,6,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Data          []byte                 `protobuf:"bytes,7,opt,name=data,proto3" json:"data,omitempty"`
	Status        int32                  `protobuf:"varint,8,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DataPacket) Reset() {
	*x = DataPacket{}
	mi := &file_protos_data_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataPacket) ProtoMessage() {}

func (x *DataPacket) ProtoReflect() protoreflect.Message {
	mi := &file_protos_data_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataPacket.ProtoReflect.Descriptor instead.
func (*DataPacket) Descriptor() ([]byte, []int) {
	return file_protos_data_proto_rawDescGZIP(), []int{0}
}

func (x *DataPacket) GetTunnelId() string {
	if x != nil {
		return x.TunnelId
	}
	return ""
}

func (x *DataPacket) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *DataPacket) GetType() PacketType {
	if x != nil {
		return x.Type
	}
	return PacketType_REQUEST
}

func (x *DataPacket) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *DataPacket) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *DataPacket) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *DataPacket) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *DataPacket) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

var File_protos_data_proto protoreflect.FileDescriptor

var file_protos_data_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xb9, 0x02, 0x0a, 0x0a, 0x44, 0x61,
	0x74, 0x61, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x75, 0x6e, 0x6e,
	0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x75, 0x6e,
	0x6e, 0x65, 0x6c, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x10, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x75, 0x72, 0x6c, 0x12, 0x37, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0x3a, 0x0a, 0x0c, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x36, 0x0a, 0x0a, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0x00,
	0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10, 0x01, 0x12, 0x0d,
	0x0a, 0x09, 0x4b, 0x45, 0x45, 0x50, 0x41, 0x4c, 0x49, 0x56, 0x45, 0x10, 0x02, 0x32, 0x45, 0x0a,
	0x0a, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x12, 0x37, 0x0a, 0x0b, 0x46,
	0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x1a, 0x10, 0x2e, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x00,
	0x28, 0x01, 0x30, 0x01, 0x42, 0x13, 0x5a, 0x11, 0x2e, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_protos_data_proto_rawDescOnce sync.Once
	file_protos_data_proto_rawDescData []byte
)

func file_protos_data_proto_rawDescGZIP() []byte {
	file_protos_data_proto_rawDescOnce.Do(func() {
		file_protos_data_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_protos_data_proto_rawDesc), len(file_protos_data_proto_rawDesc)))
	})
	return file_protos_data_proto_rawDescData
}

var file_protos_data_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_protos_data_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protos_data_proto_goTypes = []any{
	(PacketType)(0),    // 0: data.PacketType
	(*DataPacket)(nil), // 1: data.DataPacket
	nil,                // 2: data.DataPacket.HeadersEntry
}
var file_protos_data_proto_depIdxs = []int32{
	0, // 0: data.DataPacket.type:type_name -> data.PacketType
	2, // 1: data.DataPacket.headers:type_name -> data.DataPacket.HeadersEntry
	1, // 2: data.TunnelData.ForwardData:input_type -> data.DataPacket
	1, // 3: data.TunnelData.ForwardData:output_type -> data.DataPacket
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_protos_data_proto_init() }
func file_protos_data_proto_init() {
	if File_protos_data_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_protos_data_proto_rawDesc), len(file_protos_data_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_data_proto_goTypes,
		DependencyIndexes: file_protos_data_proto_depIdxs,
		EnumInfos:         file_protos_data_proto_enumTypes,
		MessageInfos:      file_protos_data_proto_msgTypes,
	}.Build()
	File_protos_data_proto = out.File
	file_protos_data_proto_goTypes = nil
	file_protos_data_proto_depIdxs = nil
}
