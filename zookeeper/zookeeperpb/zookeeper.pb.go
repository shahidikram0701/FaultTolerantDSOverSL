// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: zookeeper.proto

package zookeeper

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

type ACL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scheme string `protobuf:"bytes,1,opt,name=scheme,proto3" json:"scheme,omitempty"`
	Id     string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Perms  uint32 `protobuf:"varint,3,opt,name=perms,proto3" json:"perms,omitempty"`
}

func (x *ACL) Reset() {
	*x = ACL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zookeeper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ACL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ACL) ProtoMessage() {}

func (x *ACL) ProtoReflect() protoreflect.Message {
	mi := &file_zookeeper_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ACL.ProtoReflect.Descriptor instead.
func (*ACL) Descriptor() ([]byte, []int) {
	return file_zookeeper_proto_rawDescGZIP(), []int{0}
}

func (x *ACL) GetScheme() string {
	if x != nil {
		return x.Scheme
	}
	return ""
}

func (x *ACL) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ACL) GetPerms() uint32 {
	if x != nil {
		return x.Perms
	}
	return 0
}

type Stat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Czxid          int64 `protobuf:"varint,1,opt,name=czxid,proto3" json:"czxid,omitempty"`
	Mzxid          int64 `protobuf:"varint,2,opt,name=mzxid,proto3" json:"mzxid,omitempty"`
	Ctime          int64 `protobuf:"varint,3,opt,name=ctime,proto3" json:"ctime,omitempty"`
	Mtime          int64 `protobuf:"varint,4,opt,name=mtime,proto3" json:"mtime,omitempty"`
	Version        int32 `protobuf:"varint,5,opt,name=version,proto3" json:"version,omitempty"`
	Cversion       int32 `protobuf:"varint,6,opt,name=cversion,proto3" json:"cversion,omitempty"`
	EphemeralOwner int32 `protobuf:"varint,7,opt,name=ephemeralOwner,proto3" json:"ephemeralOwner,omitempty"`
	DataLength     int32 `protobuf:"varint,8,opt,name=dataLength,proto3" json:"dataLength,omitempty"`
	NumChildren    int32 `protobuf:"varint,9,opt,name=numChildren,proto3" json:"numChildren,omitempty"`
	Pzxid          int32 `protobuf:"varint,10,opt,name=pzxid,proto3" json:"pzxid,omitempty"`
}

func (x *Stat) Reset() {
	*x = Stat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zookeeper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stat) ProtoMessage() {}

func (x *Stat) ProtoReflect() protoreflect.Message {
	mi := &file_zookeeper_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stat.ProtoReflect.Descriptor instead.
func (*Stat) Descriptor() ([]byte, []int) {
	return file_zookeeper_proto_rawDescGZIP(), []int{1}
}

func (x *Stat) GetCzxid() int64 {
	if x != nil {
		return x.Czxid
	}
	return 0
}

func (x *Stat) GetMzxid() int64 {
	if x != nil {
		return x.Mzxid
	}
	return 0
}

func (x *Stat) GetCtime() int64 {
	if x != nil {
		return x.Ctime
	}
	return 0
}

func (x *Stat) GetMtime() int64 {
	if x != nil {
		return x.Mtime
	}
	return 0
}

func (x *Stat) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Stat) GetCversion() int32 {
	if x != nil {
		return x.Cversion
	}
	return 0
}

func (x *Stat) GetEphemeralOwner() int32 {
	if x != nil {
		return x.EphemeralOwner
	}
	return 0
}

func (x *Stat) GetDataLength() int32 {
	if x != nil {
		return x.DataLength
	}
	return 0
}

func (x *Stat) GetNumChildren() int32 {
	if x != nil {
		return x.NumChildren
	}
	return 0
}

func (x *Stat) GetPzxid() int32 {
	if x != nil {
		return x.Pzxid
	}
	return 0
}

type ZNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Acl  []*ACL `protobuf:"bytes,3,rep,name=acl,proto3" json:"acl,omitempty"`
	Stat *Stat  `protobuf:"bytes,4,opt,name=stat,proto3" json:"stat,omitempty"`
}

func (x *ZNode) Reset() {
	*x = ZNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zookeeper_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ZNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ZNode) ProtoMessage() {}

func (x *ZNode) ProtoReflect() protoreflect.Message {
	mi := &file_zookeeper_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ZNode.ProtoReflect.Descriptor instead.
func (*ZNode) Descriptor() ([]byte, []int) {
	return file_zookeeper_proto_rawDescGZIP(), []int{2}
}

func (x *ZNode) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ZNode) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ZNode) GetAcl() []*ACL {
	if x != nil {
		return x.Acl
	}
	return nil
}

func (x *ZNode) GetStat() *Stat {
	if x != nil {
		return x.Stat
	}
	return nil
}

type Path struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *Path) Reset() {
	*x = Path{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zookeeper_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Path) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Path) ProtoMessage() {}

func (x *Path) ProtoReflect() protoreflect.Message {
	mi := &file_zookeeper_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Path.ProtoReflect.Descriptor instead.
func (*Path) Descriptor() ([]byte, []int) {
	return file_zookeeper_proto_rawDescGZIP(), []int{3}
}

func (x *Path) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zookeeper_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_zookeeper_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_zookeeper_proto_rawDescGZIP(), []int{4}
}

type SetZNodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path    string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Data    []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Version int32  `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *SetZNodeRequest) Reset() {
	*x = SetZNodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zookeeper_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetZNodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetZNodeRequest) ProtoMessage() {}

func (x *SetZNodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_zookeeper_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetZNodeRequest.ProtoReflect.Descriptor instead.
func (*SetZNodeRequest) Descriptor() ([]byte, []int) {
	return file_zookeeper_proto_rawDescGZIP(), []int{5}
}

func (x *SetZNodeRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *SetZNodeRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *SetZNodeRequest) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

type GetZNodeChildrenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Children []string `protobuf:"bytes,1,rep,name=children,proto3" json:"children,omitempty"`
}

func (x *GetZNodeChildrenResponse) Reset() {
	*x = GetZNodeChildrenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zookeeper_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetZNodeChildrenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetZNodeChildrenResponse) ProtoMessage() {}

func (x *GetZNodeChildrenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_zookeeper_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetZNodeChildrenResponse.ProtoReflect.Descriptor instead.
func (*GetZNodeChildrenResponse) Descriptor() ([]byte, []int) {
	return file_zookeeper_proto_rawDescGZIP(), []int{6}
}

func (x *GetZNodeChildrenResponse) GetChildren() []string {
	if x != nil {
		return x.Children
	}
	return nil
}

var File_zookeeper_proto protoreflect.FileDescriptor

var file_zookeeper_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x22, 0x43, 0x0a, 0x03,
	0x41, 0x43, 0x4c, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x70,
	0x65, 0x72, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x70, 0x65, 0x72, 0x6d,
	0x73, 0x22, 0x94, 0x02, 0x0a, 0x04, 0x53, 0x74, 0x61, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x7a,
	0x78, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x7a, 0x78, 0x69, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x6d, 0x7a, 0x78, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x6d, 0x7a, 0x78, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x6d, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6d, 0x74, 0x69,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08,
	0x63, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x63, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x0e, 0x65, 0x70, 0x68, 0x65,
	0x6d, 0x65, 0x72, 0x61, 0x6c, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0e, 0x65, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x4f, 0x77, 0x6e, 0x65, 0x72,
	0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68,
	0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x75, 0x6d, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x6e, 0x75, 0x6d, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72,
	0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x7a, 0x78, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x70, 0x7a, 0x78, 0x69, 0x64, 0x22, 0x76, 0x0a, 0x05, 0x5a, 0x4e, 0x6f, 0x64,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x20, 0x0a, 0x03, 0x61, 0x63, 0x6c,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70,
	0x65, 0x72, 0x2e, 0x41, 0x43, 0x4c, 0x52, 0x03, 0x61, 0x63, 0x6c, 0x12, 0x23, 0x0a, 0x04, 0x73,
	0x74, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b,
	0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x52, 0x04, 0x73, 0x74, 0x61, 0x74,
	0x22, 0x1a, 0x0a, 0x04, 0x50, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x07, 0x0a, 0x05,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x53, 0x0a, 0x0f, 0x53, 0x65, 0x74, 0x5a, 0x4e, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x36, 0x0a, 0x18, 0x47, 0x65,
	0x74, 0x5a, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72,
	0x65, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72,
	0x65, 0x6e, 0x32, 0xd2, 0x02, 0x0a, 0x09, 0x5a, 0x6f, 0x6f, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x12, 0x30, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5a, 0x4e, 0x6f, 0x64, 0x65, 0x12,
	0x10, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x5a, 0x4e, 0x6f, 0x64,
	0x65, 0x1a, 0x0f, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x50, 0x61,
	0x74, 0x68, 0x12, 0x30, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5a, 0x4e, 0x6f, 0x64,
	0x65, 0x12, 0x0f, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x50, 0x61,
	0x74, 0x68, 0x1a, 0x10, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x12, 0x2f, 0x0a, 0x0b, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x5a, 0x4e,
	0x6f, 0x64, 0x65, 0x12, 0x0f, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e,
	0x50, 0x61, 0x74, 0x68, 0x1a, 0x0f, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x12, 0x2d, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x5a, 0x4e, 0x6f, 0x64,
	0x65, 0x12, 0x0f, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x50, 0x61,
	0x74, 0x68, 0x1a, 0x10, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x5a,
	0x4e, 0x6f, 0x64, 0x65, 0x12, 0x37, 0x0a, 0x08, 0x53, 0x65, 0x74, 0x5a, 0x4e, 0x6f, 0x64, 0x65,
	0x12, 0x1a, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x74,
	0x5a, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x7a,
	0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x12, 0x48, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x5a, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65,
	0x6e, 0x12, 0x0f, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x50, 0x61,
	0x74, 0x68, 0x1a, 0x23, 0x2e, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x47,
	0x65, 0x74, 0x5a, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x24, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x63, 0x61, 0x6c, 0x6f, 0x67, 0x2f, 0x73, 0x63, 0x61,
	0x6c, 0x6f, 0x67, 0x2f, 0x7a, 0x6f, 0x6f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_zookeeper_proto_rawDescOnce sync.Once
	file_zookeeper_proto_rawDescData = file_zookeeper_proto_rawDesc
)

func file_zookeeper_proto_rawDescGZIP() []byte {
	file_zookeeper_proto_rawDescOnce.Do(func() {
		file_zookeeper_proto_rawDescData = protoimpl.X.CompressGZIP(file_zookeeper_proto_rawDescData)
	})
	return file_zookeeper_proto_rawDescData
}

var file_zookeeper_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_zookeeper_proto_goTypes = []interface{}{
	(*ACL)(nil),                      // 0: zookeeper.ACL
	(*Stat)(nil),                     // 1: zookeeper.Stat
	(*ZNode)(nil),                    // 2: zookeeper.ZNode
	(*Path)(nil),                     // 3: zookeeper.Path
	(*Empty)(nil),                    // 4: zookeeper.Empty
	(*SetZNodeRequest)(nil),          // 5: zookeeper.SetZNodeRequest
	(*GetZNodeChildrenResponse)(nil), // 6: zookeeper.GetZNodeChildrenResponse
}
var file_zookeeper_proto_depIdxs = []int32{
	0, // 0: zookeeper.ZNode.acl:type_name -> zookeeper.ACL
	1, // 1: zookeeper.ZNode.stat:type_name -> zookeeper.Stat
	2, // 2: zookeeper.ZooKeeper.CreateZNode:input_type -> zookeeper.ZNode
	3, // 3: zookeeper.ZooKeeper.DeleteZNode:input_type -> zookeeper.Path
	3, // 4: zookeeper.ZooKeeper.ExistsZNode:input_type -> zookeeper.Path
	3, // 5: zookeeper.ZooKeeper.GetZNode:input_type -> zookeeper.Path
	5, // 6: zookeeper.ZooKeeper.SetZNode:input_type -> zookeeper.SetZNodeRequest
	3, // 7: zookeeper.ZooKeeper.GetZNodeChildren:input_type -> zookeeper.Path
	3, // 8: zookeeper.ZooKeeper.CreateZNode:output_type -> zookeeper.Path
	4, // 9: zookeeper.ZooKeeper.DeleteZNode:output_type -> zookeeper.Empty
	1, // 10: zookeeper.ZooKeeper.ExistsZNode:output_type -> zookeeper.Stat
	2, // 11: zookeeper.ZooKeeper.GetZNode:output_type -> zookeeper.ZNode
	1, // 12: zookeeper.ZooKeeper.SetZNode:output_type -> zookeeper.Stat
	6, // 13: zookeeper.ZooKeeper.GetZNodeChildren:output_type -> zookeeper.GetZNodeChildrenResponse
	8, // [8:14] is the sub-list for method output_type
	2, // [2:8] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_zookeeper_proto_init() }
func file_zookeeper_proto_init() {
	if File_zookeeper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_zookeeper_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ACL); i {
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
		file_zookeeper_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stat); i {
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
		file_zookeeper_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ZNode); i {
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
		file_zookeeper_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Path); i {
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
		file_zookeeper_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_zookeeper_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetZNodeRequest); i {
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
		file_zookeeper_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetZNodeChildrenResponse); i {
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
			RawDescriptor: file_zookeeper_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_zookeeper_proto_goTypes,
		DependencyIndexes: file_zookeeper_proto_depIdxs,
		MessageInfos:      file_zookeeper_proto_msgTypes,
	}.Build()
	File_zookeeper_proto = out.File
	file_zookeeper_proto_rawDesc = nil
	file_zookeeper_proto_goTypes = nil
	file_zookeeper_proto_depIdxs = nil
}
