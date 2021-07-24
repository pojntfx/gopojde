// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: instances.proto

package v1

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type InstancesMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instances []*InstanceMessage `protobuf:"bytes,1,rep,name=instances,proto3" json:"instances,omitempty"`
}

func (x *InstancesMessage) Reset() {
	*x = InstancesMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_instances_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstancesMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstancesMessage) ProtoMessage() {}

func (x *InstancesMessage) ProtoReflect() protoreflect.Message {
	mi := &file_instances_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstancesMessage.ProtoReflect.Descriptor instead.
func (*InstancesMessage) Descriptor() ([]byte, []int) {
	return file_instances_proto_rawDescGZIP(), []int{0}
}

func (x *InstancesMessage) GetInstances() []*InstanceMessage {
	if x != nil {
		return x.Instances
	}
	return nil
}

type InstanceMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string  `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Ports  []int32 `protobuf:"varint,2,rep,packed,name=Ports,proto3" json:"Ports,omitempty"`
	Status string  `protobuf:"bytes,3,opt,name=Status,proto3" json:"Status,omitempty"`
}

func (x *InstanceMessage) Reset() {
	*x = InstanceMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_instances_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceMessage) ProtoMessage() {}

func (x *InstanceMessage) ProtoReflect() protoreflect.Message {
	mi := &file_instances_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceMessage.ProtoReflect.Descriptor instead.
func (*InstanceMessage) Descriptor() ([]byte, []int) {
	return file_instances_proto_rawDescGZIP(), []int{1}
}

func (x *InstanceMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InstanceMessage) GetPorts() []int32 {
	if x != nil {
		return x.Ports
	}
	return nil
}

func (x *InstanceMessage) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type InstanceReferenceMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *InstanceReferenceMessage) Reset() {
	*x = InstanceReferenceMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_instances_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceReferenceMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceReferenceMessage) ProtoMessage() {}

func (x *InstanceReferenceMessage) ProtoReflect() protoreflect.Message {
	mi := &file_instances_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceReferenceMessage.ProtoReflect.Descriptor instead.
func (*InstanceReferenceMessage) Descriptor() ([]byte, []int) {
	return file_instances_proto_rawDescGZIP(), []int{2}
}

func (x *InstanceReferenceMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type LogMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chunk string `protobuf:"bytes,1,opt,name=Chunk,proto3" json:"Chunk,omitempty"`
}

func (x *LogMessage) Reset() {
	*x = LogMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_instances_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogMessage) ProtoMessage() {}

func (x *LogMessage) ProtoReflect() protoreflect.Message {
	mi := &file_instances_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogMessage.ProtoReflect.Descriptor instead.
func (*LogMessage) Descriptor() ([]byte, []int) {
	return file_instances_proto_rawDescGZIP(), []int{3}
}

func (x *LogMessage) GetChunk() string {
	if x != nil {
		return x.Chunk
	}
	return ""
}

type InstanceRemovalOptionsMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name           string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Customizations bool   `protobuf:"varint,2,opt,name=Customizations,proto3" json:"Customizations,omitempty"`
	DEBCache       bool   `protobuf:"varint,3,opt,name=DEBCache,proto3" json:"DEBCache,omitempty"`
	Preferences    bool   `protobuf:"varint,4,opt,name=Preferences,proto3" json:"Preferences,omitempty"`
	UserData       bool   `protobuf:"varint,5,opt,name=UserData,proto3" json:"UserData,omitempty"`
	Transfer       bool   `protobuf:"varint,6,opt,name=Transfer,proto3" json:"Transfer,omitempty"`
}

func (x *InstanceRemovalOptionsMessage) Reset() {
	*x = InstanceRemovalOptionsMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_instances_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceRemovalOptionsMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceRemovalOptionsMessage) ProtoMessage() {}

func (x *InstanceRemovalOptionsMessage) ProtoReflect() protoreflect.Message {
	mi := &file_instances_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceRemovalOptionsMessage.ProtoReflect.Descriptor instead.
func (*InstanceRemovalOptionsMessage) Descriptor() ([]byte, []int) {
	return file_instances_proto_rawDescGZIP(), []int{4}
}

func (x *InstanceRemovalOptionsMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InstanceRemovalOptionsMessage) GetCustomizations() bool {
	if x != nil {
		return x.Customizations
	}
	return false
}

func (x *InstanceRemovalOptionsMessage) GetDEBCache() bool {
	if x != nil {
		return x.DEBCache
	}
	return false
}

func (x *InstanceRemovalOptionsMessage) GetPreferences() bool {
	if x != nil {
		return x.Preferences
	}
	return false
}

func (x *InstanceRemovalOptionsMessage) GetUserData() bool {
	if x != nil {
		return x.UserData
	}
	return false
}

func (x *InstanceRemovalOptionsMessage) GetTransfer() bool {
	if x != nil {
		return x.Transfer
	}
	return false
}

var File_instances_proto protoreflect.FileDescriptor

var file_instances_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6f, 0x6a, 0x74, 0x69, 0x6e, 0x67, 0x65, 0x72,
	0x2e, 0x66, 0x65, 0x6c, 0x69, 0x78, 0x2e, 0x70, 0x6f, 0x6a, 0x64, 0x65, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x10, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x48, 0x0a,
	0x09, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2a, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6f, 0x6a, 0x74, 0x69, 0x6e, 0x67, 0x65, 0x72,
	0x2e, 0x66, 0x65, 0x6c, 0x69, 0x78, 0x2e, 0x70, 0x6f, 0x6a, 0x64, 0x65, 0x2e, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x09, 0x69, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x22, 0x53, 0x0a, 0x0f, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x05, 0x52, 0x05, 0x50,
	0x6f, 0x72, 0x74, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x2e, 0x0a, 0x18,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x22, 0x0a, 0x0a,
	0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x68,
	0x75, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x22, 0xd1, 0x01, 0x0a, 0x1d, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x6d,
	0x6f, 0x76, 0x61, 0x6c, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e,
	0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1a,
	0x0a, 0x08, 0x44, 0x45, 0x42, 0x43, 0x61, 0x63, 0x68, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x44, 0x45, 0x42, 0x43, 0x61, 0x63, 0x68, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0b, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x65, 0x72, 0x32, 0xcf, 0x04, 0x0a, 0x10, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x53, 0x0a, 0x0c, 0x47, 0x65, 0x74,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x2b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6f, 0x6a, 0x74, 0x69, 0x6e, 0x67, 0x65,
	0x72, 0x2e, 0x66, 0x65, 0x6c, 0x69, 0x78, 0x2e, 0x70, 0x6f, 0x6a, 0x64, 0x65, 0x2e, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x67,
	0x0a, 0x07, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x67, 0x73, 0x12, 0x33, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x70, 0x6f, 0x6a, 0x74, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x6c, 0x69, 0x78, 0x2e,
	0x70, 0x6f, 0x6a, 0x64, 0x65, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x25,
	0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6f, 0x6a, 0x74, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x66,
	0x65, 0x6c, 0x69, 0x78, 0x2e, 0x70, 0x6f, 0x6a, 0x64, 0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x30, 0x01, 0x12, 0x5c, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x33, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70,
	0x6f, 0x6a, 0x74, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x6c, 0x69, 0x78, 0x2e, 0x70,
	0x6f, 0x6a, 0x64, 0x65, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x5b, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x70, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x33, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6f, 0x6a, 0x74,
	0x69, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x6c, 0x69, 0x78, 0x2e, 0x70, 0x6f, 0x6a, 0x64,
	0x65, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x5e, 0x0a, 0x0f, 0x52, 0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x33, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6f, 0x6a, 0x74,
	0x69, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x6c, 0x69, 0x78, 0x2e, 0x70, 0x6f, 0x6a, 0x64,
	0x65, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x62, 0x0a, 0x0e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x12, 0x38, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6f, 0x6a, 0x74, 0x69,
	0x6e, 0x67, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x6c, 0x69, 0x78, 0x2e, 0x70, 0x6f, 0x6a, 0x64, 0x65,
	0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x61, 0x6c,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6f, 0x6a, 0x6e, 0x74, 0x66, 0x78, 0x2f, 0x70, 0x6f, 0x6a,
	0x64, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_instances_proto_rawDescOnce sync.Once
	file_instances_proto_rawDescData = file_instances_proto_rawDesc
)

func file_instances_proto_rawDescGZIP() []byte {
	file_instances_proto_rawDescOnce.Do(func() {
		file_instances_proto_rawDescData = protoimpl.X.CompressGZIP(file_instances_proto_rawDescData)
	})
	return file_instances_proto_rawDescData
}

var file_instances_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_instances_proto_goTypes = []interface{}{
	(*InstancesMessage)(nil),              // 0: com.pojtinger.felicitas.pojde.InstancesMessage
	(*InstanceMessage)(nil),               // 1: com.pojtinger.felicitas.pojde.InstanceMessage
	(*InstanceReferenceMessage)(nil),      // 2: com.pojtinger.felicitas.pojde.InstanceReferenceMessage
	(*LogMessage)(nil),                    // 3: com.pojtinger.felicitas.pojde.LogMessage
	(*InstanceRemovalOptionsMessage)(nil), // 4: com.pojtinger.felicitas.pojde.InstanceRemovalOptionsMessage
	(*empty.Empty)(nil),                   // 5: google.protobuf.Empty
}
var file_instances_proto_depIdxs = []int32{
	1, // 0: com.pojtinger.felicitas.pojde.InstancesMessage.instances:type_name -> com.pojtinger.felicitas.pojde.InstanceMessage
	5, // 1: com.pojtinger.felicitas.pojde.InstancesService.GetInstances:input_type -> google.protobuf.Empty
	2, // 2: com.pojtinger.felicitas.pojde.InstancesService.GetLogs:input_type -> com.pojtinger.felicitas.pojde.InstanceReferenceMessage
	2, // 3: com.pojtinger.felicitas.pojde.InstancesService.StartInstance:input_type -> com.pojtinger.felicitas.pojde.InstanceReferenceMessage
	2, // 4: com.pojtinger.felicitas.pojde.InstancesService.StopInstance:input_type -> com.pojtinger.felicitas.pojde.InstanceReferenceMessage
	2, // 5: com.pojtinger.felicitas.pojde.InstancesService.RestartInstance:input_type -> com.pojtinger.felicitas.pojde.InstanceReferenceMessage
	4, // 6: com.pojtinger.felicitas.pojde.InstancesService.RemoveInstance:input_type -> com.pojtinger.felicitas.pojde.InstanceRemovalOptionsMessage
	0, // 7: com.pojtinger.felicitas.pojde.InstancesService.GetInstances:output_type -> com.pojtinger.felicitas.pojde.InstancesMessage
	3, // 8: com.pojtinger.felicitas.pojde.InstancesService.GetLogs:output_type -> com.pojtinger.felicitas.pojde.LogMessage
	5, // 9: com.pojtinger.felicitas.pojde.InstancesService.StartInstance:output_type -> google.protobuf.Empty
	5, // 10: com.pojtinger.felicitas.pojde.InstancesService.StopInstance:output_type -> google.protobuf.Empty
	5, // 11: com.pojtinger.felicitas.pojde.InstancesService.RestartInstance:output_type -> google.protobuf.Empty
	5, // 12: com.pojtinger.felicitas.pojde.InstancesService.RemoveInstance:output_type -> google.protobuf.Empty
	7, // [7:13] is the sub-list for method output_type
	1, // [1:7] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_instances_proto_init() }
func file_instances_proto_init() {
	if File_instances_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_instances_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstancesMessage); i {
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
		file_instances_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceMessage); i {
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
		file_instances_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceReferenceMessage); i {
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
		file_instances_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogMessage); i {
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
		file_instances_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceRemovalOptionsMessage); i {
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
			RawDescriptor: file_instances_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_instances_proto_goTypes,
		DependencyIndexes: file_instances_proto_depIdxs,
		MessageInfos:      file_instances_proto_msgTypes,
	}.Build()
	File_instances_proto = out.File
	file_instances_proto_rawDesc = nil
	file_instances_proto_goTypes = nil
	file_instances_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// InstancesServiceClient is the client API for InstancesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type InstancesServiceClient interface {
	GetInstances(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*InstancesMessage, error)
	GetLogs(ctx context.Context, in *InstanceReferenceMessage, opts ...grpc.CallOption) (InstancesService_GetLogsClient, error)
	StartInstance(ctx context.Context, in *InstanceReferenceMessage, opts ...grpc.CallOption) (*empty.Empty, error)
	StopInstance(ctx context.Context, in *InstanceReferenceMessage, opts ...grpc.CallOption) (*empty.Empty, error)
	RestartInstance(ctx context.Context, in *InstanceReferenceMessage, opts ...grpc.CallOption) (*empty.Empty, error)
	RemoveInstance(ctx context.Context, in *InstanceRemovalOptionsMessage, opts ...grpc.CallOption) (*empty.Empty, error)
}

type instancesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInstancesServiceClient(cc grpc.ClientConnInterface) InstancesServiceClient {
	return &instancesServiceClient{cc}
}

func (c *instancesServiceClient) GetInstances(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*InstancesMessage, error) {
	out := new(InstancesMessage)
	err := c.cc.Invoke(ctx, "/com.pojtinger.felicitas.pojde.InstancesService/GetInstances", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instancesServiceClient) GetLogs(ctx context.Context, in *InstanceReferenceMessage, opts ...grpc.CallOption) (InstancesService_GetLogsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_InstancesService_serviceDesc.Streams[0], "/com.pojtinger.felicitas.pojde.InstancesService/GetLogs", opts...)
	if err != nil {
		return nil, err
	}
	x := &instancesServiceGetLogsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type InstancesService_GetLogsClient interface {
	Recv() (*LogMessage, error)
	grpc.ClientStream
}

type instancesServiceGetLogsClient struct {
	grpc.ClientStream
}

func (x *instancesServiceGetLogsClient) Recv() (*LogMessage, error) {
	m := new(LogMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *instancesServiceClient) StartInstance(ctx context.Context, in *InstanceReferenceMessage, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/com.pojtinger.felicitas.pojde.InstancesService/StartInstance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instancesServiceClient) StopInstance(ctx context.Context, in *InstanceReferenceMessage, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/com.pojtinger.felicitas.pojde.InstancesService/StopInstance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instancesServiceClient) RestartInstance(ctx context.Context, in *InstanceReferenceMessage, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/com.pojtinger.felicitas.pojde.InstancesService/RestartInstance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instancesServiceClient) RemoveInstance(ctx context.Context, in *InstanceRemovalOptionsMessage, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/com.pojtinger.felicitas.pojde.InstancesService/RemoveInstance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InstancesServiceServer is the server API for InstancesService service.
type InstancesServiceServer interface {
	GetInstances(context.Context, *empty.Empty) (*InstancesMessage, error)
	GetLogs(*InstanceReferenceMessage, InstancesService_GetLogsServer) error
	StartInstance(context.Context, *InstanceReferenceMessage) (*empty.Empty, error)
	StopInstance(context.Context, *InstanceReferenceMessage) (*empty.Empty, error)
	RestartInstance(context.Context, *InstanceReferenceMessage) (*empty.Empty, error)
	RemoveInstance(context.Context, *InstanceRemovalOptionsMessage) (*empty.Empty, error)
}

// UnimplementedInstancesServiceServer can be embedded to have forward compatible implementations.
type UnimplementedInstancesServiceServer struct {
}

func (*UnimplementedInstancesServiceServer) GetInstances(context.Context, *empty.Empty) (*InstancesMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInstances not implemented")
}
func (*UnimplementedInstancesServiceServer) GetLogs(*InstanceReferenceMessage, InstancesService_GetLogsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetLogs not implemented")
}
func (*UnimplementedInstancesServiceServer) StartInstance(context.Context, *InstanceReferenceMessage) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartInstance not implemented")
}
func (*UnimplementedInstancesServiceServer) StopInstance(context.Context, *InstanceReferenceMessage) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopInstance not implemented")
}
func (*UnimplementedInstancesServiceServer) RestartInstance(context.Context, *InstanceReferenceMessage) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestartInstance not implemented")
}
func (*UnimplementedInstancesServiceServer) RemoveInstance(context.Context, *InstanceRemovalOptionsMessage) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveInstance not implemented")
}

func RegisterInstancesServiceServer(s *grpc.Server, srv InstancesServiceServer) {
	s.RegisterService(&_InstancesService_serviceDesc, srv)
}

func _InstancesService_GetInstances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstancesServiceServer).GetInstances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.pojtinger.felicitas.pojde.InstancesService/GetInstances",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstancesServiceServer).GetInstances(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstancesService_GetLogs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(InstanceReferenceMessage)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(InstancesServiceServer).GetLogs(m, &instancesServiceGetLogsServer{stream})
}

type InstancesService_GetLogsServer interface {
	Send(*LogMessage) error
	grpc.ServerStream
}

type instancesServiceGetLogsServer struct {
	grpc.ServerStream
}

func (x *instancesServiceGetLogsServer) Send(m *LogMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _InstancesService_StartInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstanceReferenceMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstancesServiceServer).StartInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.pojtinger.felicitas.pojde.InstancesService/StartInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstancesServiceServer).StartInstance(ctx, req.(*InstanceReferenceMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstancesService_StopInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstanceReferenceMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstancesServiceServer).StopInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.pojtinger.felicitas.pojde.InstancesService/StopInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstancesServiceServer).StopInstance(ctx, req.(*InstanceReferenceMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstancesService_RestartInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstanceReferenceMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstancesServiceServer).RestartInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.pojtinger.felicitas.pojde.InstancesService/RestartInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstancesServiceServer).RestartInstance(ctx, req.(*InstanceReferenceMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstancesService_RemoveInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstanceRemovalOptionsMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstancesServiceServer).RemoveInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.pojtinger.felicitas.pojde.InstancesService/RemoveInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstancesServiceServer).RemoveInstance(ctx, req.(*InstanceRemovalOptionsMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _InstancesService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "com.pojtinger.felicitas.pojde.InstancesService",
	HandlerType: (*InstancesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInstances",
			Handler:    _InstancesService_GetInstances_Handler,
		},
		{
			MethodName: "StartInstance",
			Handler:    _InstancesService_StartInstance_Handler,
		},
		{
			MethodName: "StopInstance",
			Handler:    _InstancesService_StopInstance_Handler,
		},
		{
			MethodName: "RestartInstance",
			Handler:    _InstancesService_RestartInstance_Handler,
		},
		{
			MethodName: "RemoveInstance",
			Handler:    _InstancesService_RemoveInstance_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetLogs",
			Handler:       _InstancesService_GetLogs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "instances.proto",
}
