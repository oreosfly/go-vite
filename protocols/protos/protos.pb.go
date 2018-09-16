// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos.proto

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type BlockID struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Height               []byte   `protobuf:"bytes,2,opt,name=Height,proto3" json:"Height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockID) Reset()         { *m = BlockID{} }
func (m *BlockID) String() string { return proto.CompactTextString(m) }
func (*BlockID) ProtoMessage()    {}
func (*BlockID) Descriptor() ([]byte, []int) {
	return fileDescriptor_protos_3319151b0e98edd8, []int{0}
}
func (m *BlockID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockID.Unmarshal(m, b)
}
func (m *BlockID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockID.Marshal(b, m, deterministic)
}
func (dst *BlockID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockID.Merge(dst, src)
}
func (m *BlockID) XXX_Size() int {
	return xxx_messageInfo_BlockID.Size(m)
}
func (m *BlockID) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockID.DiscardUnknown(m)
}

var xxx_messageInfo_BlockID proto.InternalMessageInfo

func (m *BlockID) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *BlockID) GetHeight() []byte {
	if m != nil {
		return m.Height
	}
	return nil
}

type Segment struct {
	From                 *BlockID `protobuf:"bytes,1,opt,name=From,proto3" json:"From,omitempty"`
	To                   *BlockID `protobuf:"bytes,2,opt,name=To,proto3" json:"To,omitempty"`
	Step                 uint64   `protobuf:"varint,3,opt,name=step,proto3" json:"step,omitempty"`
	Forward              bool     `protobuf:"varint,4,opt,name=Forward,proto3" json:"Forward,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Segment) Reset()         { *m = Segment{} }
func (m *Segment) String() string { return proto.CompactTextString(m) }
func (*Segment) ProtoMessage()    {}
func (*Segment) Descriptor() ([]byte, []int) {
	return fileDescriptor_protos_3319151b0e98edd8, []int{1}
}
func (m *Segment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Segment.Unmarshal(m, b)
}
func (m *Segment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Segment.Marshal(b, m, deterministic)
}
func (dst *Segment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Segment.Merge(dst, src)
}
func (m *Segment) XXX_Size() int {
	return xxx_messageInfo_Segment.Size(m)
}
func (m *Segment) XXX_DiscardUnknown() {
	xxx_messageInfo_Segment.DiscardUnknown(m)
}

var xxx_messageInfo_Segment proto.InternalMessageInfo

func (m *Segment) GetFrom() *BlockID {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Segment) GetTo() *BlockID {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Segment) GetStep() uint64 {
	if m != nil {
		return m.Step
	}
	return 0
}

func (m *Segment) GetForward() bool {
	if m != nil {
		return m.Forward
	}
	return false
}

type AccountSegment struct {
	Segment              map[string]*Segment `protobuf:"bytes,1,rep,name=Segment,proto3" json:"Segment,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *AccountSegment) Reset()         { *m = AccountSegment{} }
func (m *AccountSegment) String() string { return proto.CompactTextString(m) }
func (*AccountSegment) ProtoMessage()    {}
func (*AccountSegment) Descriptor() ([]byte, []int) {
	return fileDescriptor_protos_3319151b0e98edd8, []int{2}
}
func (m *AccountSegment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountSegment.Unmarshal(m, b)
}
func (m *AccountSegment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountSegment.Marshal(b, m, deterministic)
}
func (dst *AccountSegment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountSegment.Merge(dst, src)
}
func (m *AccountSegment) XXX_Size() int {
	return xxx_messageInfo_AccountSegment.Size(m)
}
func (m *AccountSegment) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountSegment.DiscardUnknown(m)
}

var xxx_messageInfo_AccountSegment proto.InternalMessageInfo

func (m *AccountSegment) GetSegment() map[string]*Segment {
	if m != nil {
		return m.Segment
	}
	return nil
}

func init() {
	proto.RegisterType((*BlockID)(nil), "protos.BlockID")
	proto.RegisterType((*Segment)(nil), "protos.Segment")
	proto.RegisterType((*AccountSegment)(nil), "protos.AccountSegment")
	proto.RegisterMapType((map[string]*Segment)(nil), "protos.AccountSegment.SegmentEntry")
}

func init() { proto.RegisterFile("protos.proto", fileDescriptor_protos_3319151b0e98edd8) }

var fileDescriptor_protos_3319151b0e98edd8 = []byte{
	// 246 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x28, 0xca, 0x2f,
	0xc9, 0x2f, 0xd6, 0x03, 0x53, 0x42, 0x6c, 0x10, 0x9e, 0x92, 0x29, 0x17, 0xbb, 0x53, 0x4e, 0x7e,
	0x72, 0xb6, 0xa7, 0x8b, 0x90, 0x10, 0x17, 0x8b, 0x47, 0x62, 0x71, 0x86, 0x04, 0xa3, 0x02, 0xa3,
	0x06, 0x4f, 0x10, 0x98, 0x2d, 0x24, 0xc6, 0xc5, 0xe6, 0x91, 0x9a, 0x99, 0x9e, 0x51, 0x22, 0xc1,
	0x04, 0x16, 0x85, 0xf2, 0x94, 0x6a, 0xb9, 0xd8, 0x83, 0x53, 0xd3, 0x73, 0x53, 0xf3, 0x4a, 0x84,
	0x94, 0xb9, 0x58, 0xdc, 0x8a, 0xf2, 0x73, 0xc1, 0xda, 0xb8, 0x8d, 0xf8, 0xf5, 0xa0, 0xd6, 0x40,
	0x4d, 0x0d, 0x02, 0x4b, 0x0a, 0xc9, 0x73, 0x31, 0x85, 0xe4, 0x83, 0xcd, 0xc0, 0xa2, 0x84, 0x29,
	0x24, 0x1f, 0x64, 0x79, 0x71, 0x49, 0x6a, 0x81, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0x4b, 0x10, 0x98,
	0x2d, 0x24, 0xc1, 0xc5, 0xee, 0x96, 0x5f, 0x54, 0x9e, 0x58, 0x94, 0x22, 0xc1, 0xa2, 0xc0, 0xa8,
	0xc1, 0x11, 0x04, 0xe3, 0x2a, 0xcd, 0x61, 0xe4, 0xe2, 0x73, 0x4c, 0x4e, 0xce, 0x2f, 0xcd, 0x2b,
	0x81, 0x39, 0xc3, 0x16, 0xee, 0x22, 0x09, 0x46, 0x05, 0x66, 0x0d, 0x6e, 0x23, 0x65, 0x98, 0x35,
	0xa8, 0x0a, 0xf5, 0xa0, 0xb4, 0x6b, 0x5e, 0x49, 0x51, 0x65, 0x10, 0x4c, 0x8f, 0x94, 0x37, 0x17,
	0x0f, 0xb2, 0x84, 0x90, 0x00, 0x17, 0x73, 0x76, 0x6a, 0x25, 0xd8, 0x53, 0x9c, 0x41, 0x20, 0xa6,
	0x90, 0x2a, 0x17, 0x6b, 0x59, 0x62, 0x4e, 0x69, 0x2a, 0xba, 0x2f, 0xa0, 0xda, 0x82, 0x20, 0xb2,
	0x56, 0x4c, 0x16, 0x8c, 0x49, 0x90, 0xc0, 0x35, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xd5, 0x13,
	0x62, 0x53, 0x73, 0x01, 0x00, 0x00,
}
