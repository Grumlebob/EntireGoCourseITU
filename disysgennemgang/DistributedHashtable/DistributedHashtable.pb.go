// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.15.8
// source: DistributedHashtable/DistributedHashtable.proto

package go_DistributedHashtable

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

type PutMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   int32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value int32 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *PutMsg) Reset() {
	*x = PutMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_DistributedHashtable_DistributedHashtable_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutMsg) ProtoMessage() {}

func (x *PutMsg) ProtoReflect() protoreflect.Message {
	mi := &file_DistributedHashtable_DistributedHashtable_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutMsg.ProtoReflect.Descriptor instead.
func (*PutMsg) Descriptor() ([]byte, []int) {
	return file_DistributedHashtable_DistributedHashtable_proto_rawDescGZIP(), []int{0}
}

func (x *PutMsg) GetKey() int32 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *PutMsg) GetValue() int32 {
	if x != nil {
		return x.Value
	}
	return 0
}

type PutRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Succes bool `protobuf:"varint,1,opt,name=Succes,proto3" json:"Succes,omitempty"`
}

func (x *PutRsp) Reset() {
	*x = PutRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_DistributedHashtable_DistributedHashtable_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutRsp) ProtoMessage() {}

func (x *PutRsp) ProtoReflect() protoreflect.Message {
	mi := &file_DistributedHashtable_DistributedHashtable_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutRsp.ProtoReflect.Descriptor instead.
func (*PutRsp) Descriptor() ([]byte, []int) {
	return file_DistributedHashtable_DistributedHashtable_proto_rawDescGZIP(), []int{1}
}

func (x *PutRsp) GetSucces() bool {
	if x != nil {
		return x.Succes
	}
	return false
}

type GetMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key int32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (x *GetMsg) Reset() {
	*x = GetMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_DistributedHashtable_DistributedHashtable_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMsg) ProtoMessage() {}

func (x *GetMsg) ProtoReflect() protoreflect.Message {
	mi := &file_DistributedHashtable_DistributedHashtable_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMsg.ProtoReflect.Descriptor instead.
func (*GetMsg) Descriptor() ([]byte, []int) {
	return file_DistributedHashtable_DistributedHashtable_proto_rawDescGZIP(), []int{2}
}

func (x *GetMsg) GetKey() int32 {
	if x != nil {
		return x.Key
	}
	return 0
}

type GetRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int32 `protobuf:"varint,1,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *GetRsp) Reset() {
	*x = GetRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_DistributedHashtable_DistributedHashtable_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRsp) ProtoMessage() {}

func (x *GetRsp) ProtoReflect() protoreflect.Message {
	mi := &file_DistributedHashtable_DistributedHashtable_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRsp.ProtoReflect.Descriptor instead.
func (*GetRsp) Descriptor() ([]byte, []int) {
	return file_DistributedHashtable_DistributedHashtable_proto_rawDescGZIP(), []int{3}
}

func (x *GetRsp) GetValue() int32 {
	if x != nil {
		return x.Value
	}
	return 0
}

var File_DistributedHashtable_DistributedHashtable_proto protoreflect.FileDescriptor

var file_DistributedHashtable_DistributedHashtable_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x48, 0x61, 0x73,
	0x68, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2f, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x64, 0x48, 0x61, 0x73, 0x68, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x14, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x48, 0x61,
	0x73, 0x68, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x30, 0x0a, 0x06, 0x50, 0x75, 0x74, 0x4d, 0x73,
	0x67, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x20, 0x0a, 0x06, 0x50, 0x75, 0x74,
	0x52, 0x73, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x06, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x22, 0x1a, 0x0a, 0x06, 0x47,
	0x65, 0x74, 0x4d, 0x73, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x22, 0x1e, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x52, 0x73,
	0x70, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x32, 0x99, 0x01, 0x0a, 0x09, 0x48, 0x61, 0x73, 0x68,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x47, 0x0a, 0x03, 0x50, 0x75, 0x74, 0x12, 0x1c, 0x2e, 0x44,
	0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x48, 0x61, 0x73, 0x68, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x2e, 0x50, 0x75, 0x74, 0x4d, 0x73, 0x67, 0x1a, 0x1c, 0x2e, 0x44, 0x69, 0x73,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x48, 0x61, 0x73, 0x68, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x2e, 0x50, 0x75, 0x74, 0x52, 0x73, 0x70, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x43,
	0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1c, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x74, 0x65, 0x64, 0x48, 0x61, 0x73, 0x68, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74,
	0x4d, 0x73, 0x67, 0x1a, 0x1c, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x64, 0x48, 0x61, 0x73, 0x68, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x73,
	0x70, 0x22, 0x00, 0x42, 0x3d, 0x5a, 0x3b, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x2d, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x64, 0x48, 0x61, 0x73, 0x68, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x3b, 0x67, 0x6f, 0x5f, 0x44, 0x69,
	0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x48, 0x61, 0x73, 0x68, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_DistributedHashtable_DistributedHashtable_proto_rawDescOnce sync.Once
	file_DistributedHashtable_DistributedHashtable_proto_rawDescData = file_DistributedHashtable_DistributedHashtable_proto_rawDesc
)

func file_DistributedHashtable_DistributedHashtable_proto_rawDescGZIP() []byte {
	file_DistributedHashtable_DistributedHashtable_proto_rawDescOnce.Do(func() {
		file_DistributedHashtable_DistributedHashtable_proto_rawDescData = protoimpl.X.CompressGZIP(file_DistributedHashtable_DistributedHashtable_proto_rawDescData)
	})
	return file_DistributedHashtable_DistributedHashtable_proto_rawDescData
}

var file_DistributedHashtable_DistributedHashtable_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_DistributedHashtable_DistributedHashtable_proto_goTypes = []interface{}{
	(*PutMsg)(nil), // 0: DistributedHashtable.PutMsg
	(*PutRsp)(nil), // 1: DistributedHashtable.PutRsp
	(*GetMsg)(nil), // 2: DistributedHashtable.GetMsg
	(*GetRsp)(nil), // 3: DistributedHashtable.GetRsp
}
var file_DistributedHashtable_DistributedHashtable_proto_depIdxs = []int32{
	0, // 0: DistributedHashtable.Hashtable.Put:input_type -> DistributedHashtable.PutMsg
	2, // 1: DistributedHashtable.Hashtable.Get:input_type -> DistributedHashtable.GetMsg
	1, // 2: DistributedHashtable.Hashtable.Put:output_type -> DistributedHashtable.PutRsp
	3, // 3: DistributedHashtable.Hashtable.Get:output_type -> DistributedHashtable.GetRsp
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_DistributedHashtable_DistributedHashtable_proto_init() }
func file_DistributedHashtable_DistributedHashtable_proto_init() {
	if File_DistributedHashtable_DistributedHashtable_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_DistributedHashtable_DistributedHashtable_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutMsg); i {
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
		file_DistributedHashtable_DistributedHashtable_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutRsp); i {
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
		file_DistributedHashtable_DistributedHashtable_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMsg); i {
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
		file_DistributedHashtable_DistributedHashtable_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRsp); i {
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
			RawDescriptor: file_DistributedHashtable_DistributedHashtable_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_DistributedHashtable_DistributedHashtable_proto_goTypes,
		DependencyIndexes: file_DistributedHashtable_DistributedHashtable_proto_depIdxs,
		MessageInfos:      file_DistributedHashtable_DistributedHashtable_proto_msgTypes,
	}.Build()
	File_DistributedHashtable_DistributedHashtable_proto = out.File
	file_DistributedHashtable_DistributedHashtable_proto_rawDesc = nil
	file_DistributedHashtable_DistributedHashtable_proto_goTypes = nil
	file_DistributedHashtable_DistributedHashtable_proto_depIdxs = nil
}
