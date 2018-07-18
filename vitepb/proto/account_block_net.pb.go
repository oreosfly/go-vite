// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/account_block_net.proto

package vitepb

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

type AccountBlockNet struct {
	Meta                 *AccountBlockMeta `protobuf:"bytes,1,opt,name=meta,proto3" json:"meta,omitempty"`
	To                   []byte            `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	PrevHash             []byte            `protobuf:"bytes,3,opt,name=prevHash,proto3" json:"prevHash,omitempty"`
	FromHash             []byte            `protobuf:"bytes,4,opt,name=fromHash,proto3" json:"fromHash,omitempty"`
	Hash                 []byte            `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
	TokenId              []byte            `protobuf:"bytes,6,opt,name=tokenId,proto3" json:"tokenId,omitempty"`
	Amount               []byte            `protobuf:"bytes,7,opt,name=amount,proto3" json:"amount,omitempty"`
	Balance              []byte            `protobuf:"bytes,8,opt,name=balance,proto3" json:"balance,omitempty"`
	Data                 string            `protobuf:"bytes,9,opt,name=data,proto3" json:"data,omitempty"`
	SnapshotTimestamp    []byte            `protobuf:"bytes,10,opt,name=snapshotTimestamp,proto3" json:"snapshotTimestamp,omitempty"`
	Timestamp            uint64            `protobuf:"varint,11,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Signature            []byte            `protobuf:"bytes,12,opt,name=signature,proto3" json:"signature,omitempty"`
	Nounce               []byte            `protobuf:"bytes,13,opt,name=nounce,proto3" json:"nounce,omitempty"`
	Difficulty           []byte            `protobuf:"bytes,14,opt,name=difficulty,proto3" json:"difficulty,omitempty"`
	FAmount              []byte            `protobuf:"bytes,15,opt,name=fAmount,proto3" json:"fAmount,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AccountBlockNet) Reset()         { *m = AccountBlockNet{} }
func (m *AccountBlockNet) String() string { return proto.CompactTextString(m) }
func (*AccountBlockNet) ProtoMessage()    {}
func (*AccountBlockNet) Descriptor() ([]byte, []int) {
	return fileDescriptor_account_block_net_fc109534c0ac795e, []int{0}
}
func (m *AccountBlockNet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountBlockNet.Unmarshal(m, b)
}
func (m *AccountBlockNet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountBlockNet.Marshal(b, m, deterministic)
}
func (dst *AccountBlockNet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountBlockNet.Merge(dst, src)
}
func (m *AccountBlockNet) XXX_Size() int {
	return xxx_messageInfo_AccountBlockNet.Size(m)
}
func (m *AccountBlockNet) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountBlockNet.DiscardUnknown(m)
}

var xxx_messageInfo_AccountBlockNet proto.InternalMessageInfo

func (m *AccountBlockNet) GetMeta() *AccountBlockMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *AccountBlockNet) GetTo() []byte {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *AccountBlockNet) GetPrevHash() []byte {
	if m != nil {
		return m.PrevHash
	}
	return nil
}

func (m *AccountBlockNet) GetFromHash() []byte {
	if m != nil {
		return m.FromHash
	}
	return nil
}

func (m *AccountBlockNet) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *AccountBlockNet) GetTokenId() []byte {
	if m != nil {
		return m.TokenId
	}
	return nil
}

func (m *AccountBlockNet) GetAmount() []byte {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *AccountBlockNet) GetBalance() []byte {
	if m != nil {
		return m.Balance
	}
	return nil
}

func (m *AccountBlockNet) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *AccountBlockNet) GetSnapshotTimestamp() []byte {
	if m != nil {
		return m.SnapshotTimestamp
	}
	return nil
}

func (m *AccountBlockNet) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *AccountBlockNet) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *AccountBlockNet) GetNounce() []byte {
	if m != nil {
		return m.Nounce
	}
	return nil
}

func (m *AccountBlockNet) GetDifficulty() []byte {
	if m != nil {
		return m.Difficulty
	}
	return nil
}

func (m *AccountBlockNet) GetFAmount() []byte {
	if m != nil {
		return m.FAmount
	}
	return nil
}

type AccountBlockListNet struct {
	Blocks               []*AccountBlockNet `protobuf:"bytes,1,rep,name=blocks,proto3" json:"blocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *AccountBlockListNet) Reset()         { *m = AccountBlockListNet{} }
func (m *AccountBlockListNet) String() string { return proto.CompactTextString(m) }
func (*AccountBlockListNet) ProtoMessage()    {}
func (*AccountBlockListNet) Descriptor() ([]byte, []int) {
	return fileDescriptor_account_block_net_fc109534c0ac795e, []int{1}
}
func (m *AccountBlockListNet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountBlockListNet.Unmarshal(m, b)
}
func (m *AccountBlockListNet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountBlockListNet.Marshal(b, m, deterministic)
}
func (dst *AccountBlockListNet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountBlockListNet.Merge(dst, src)
}
func (m *AccountBlockListNet) XXX_Size() int {
	return xxx_messageInfo_AccountBlockListNet.Size(m)
}
func (m *AccountBlockListNet) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountBlockListNet.DiscardUnknown(m)
}

var xxx_messageInfo_AccountBlockListNet proto.InternalMessageInfo

func (m *AccountBlockListNet) GetBlocks() []*AccountBlockNet {
	if m != nil {
		return m.Blocks
	}
	return nil
}

func init() {
	proto.RegisterType((*AccountBlockNet)(nil), "vitepb.AccountBlockNet")
	proto.RegisterType((*AccountBlockListNet)(nil), "vitepb.AccountBlockListNet")
}

func init() {
	proto.RegisterFile("proto/account_block_net.proto", fileDescriptor_account_block_net_fc109534c0ac795e)
}

var fileDescriptor_account_block_net_fc109534c0ac795e = []byte{
	// 344 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x4f, 0x6b, 0xf2, 0x40,
	0x10, 0xc6, 0x89, 0xe6, 0x8d, 0x3a, 0xfa, 0x2a, 0xdd, 0x42, 0xbb, 0x88, 0x2d, 0xc1, 0x53, 0x0e,
	0xa2, 0x60, 0x3f, 0x81, 0x3d, 0x94, 0x16, 0x5a, 0x0f, 0xa1, 0x77, 0xd9, 0x24, 0x9b, 0x1a, 0x34,
	0xd9, 0x90, 0x9d, 0x08, 0xfd, 0x64, 0xfd, 0x7a, 0x65, 0x67, 0xe3, 0x1f, 0xa8, 0xb7, 0x7d, 0x9e,
	0xdf, 0x33, 0xc3, 0x0c, 0xb3, 0xf0, 0x50, 0x56, 0x0a, 0xd5, 0x42, 0xc4, 0xb1, 0xaa, 0x0b, 0xdc,
	0x44, 0x7b, 0x15, 0xef, 0x36, 0x85, 0xc4, 0x39, 0xf9, 0xcc, 0x3b, 0x64, 0x28, 0xcb, 0x68, 0x3c,
	0xb9, 0x16, 0x4b, 0x22, 0x9b, 0x9a, 0xfe, 0xb4, 0x61, 0xb4, 0xb2, 0xe8, 0xd9, 0x90, 0xb5, 0x44,
	0x36, 0x03, 0x37, 0x97, 0x28, 0xb8, 0xe3, 0x3b, 0x41, 0x7f, 0xc9, 0xe7, 0xb6, 0xd1, 0xfc, 0x32,
	0xf6, 0x21, 0x51, 0x84, 0x94, 0x62, 0x43, 0x68, 0xa1, 0xe2, 0x2d, 0xdf, 0x09, 0x06, 0x61, 0x0b,
	0x15, 0x1b, 0x43, 0xb7, 0xac, 0xe4, 0xe1, 0x55, 0xe8, 0x2d, 0x6f, 0x93, 0x7b, 0xd2, 0x86, 0xa5,
	0x95, 0xca, 0x89, 0xb9, 0x96, 0x1d, 0x35, 0x63, 0xe0, 0x6e, 0x8d, 0xff, 0x8f, 0x7c, 0x7a, 0x33,
	0x0e, 0x1d, 0x54, 0x3b, 0x59, 0xbc, 0x25, 0xdc, 0x23, 0xfb, 0x28, 0xd9, 0x1d, 0x78, 0x22, 0x37,
	0xe3, 0xf0, 0x0e, 0x81, 0x46, 0x99, 0x8a, 0x48, 0xec, 0x45, 0x11, 0x4b, 0xde, 0xb5, 0x15, 0x8d,
	0x34, 0xfd, 0x13, 0x81, 0x82, 0xf7, 0x7c, 0x27, 0xe8, 0x85, 0xf4, 0x66, 0x33, 0xb8, 0xd1, 0x85,
	0x28, 0xf5, 0x56, 0xe1, 0x67, 0x96, 0x4b, 0x8d, 0x22, 0x2f, 0x39, 0x50, 0xdd, 0x5f, 0xc0, 0x26,
	0xd0, 0xc3, 0x53, 0xaa, 0xef, 0x3b, 0x81, 0x1b, 0x9e, 0x0d, 0x43, 0x75, 0xf6, 0x55, 0x08, 0xac,
	0x2b, 0xc9, 0x07, 0xd4, 0xe3, 0x6c, 0x98, 0x79, 0x0b, 0x55, 0x9b, 0xb1, 0xfe, 0xdb, 0x79, 0xad,
	0x62, 0x8f, 0x00, 0x49, 0x96, 0xa6, 0x59, 0x5c, 0xef, 0xf1, 0x9b, 0x0f, 0x89, 0x5d, 0x38, 0x66,
	0x9f, 0x74, 0x65, 0x17, 0x1d, 0xd9, 0x7d, 0x1a, 0x39, 0x7d, 0x81, 0xdb, 0xcb, 0x8b, 0xbc, 0x67,
	0x1a, 0xcd, 0xf1, 0x16, 0xe0, 0xd1, 0x89, 0x35, 0x77, 0xfc, 0x76, 0xd0, 0x5f, 0xde, 0x5f, 0x3b,
	0xdf, 0x5a, 0x62, 0xd8, 0xc4, 0x22, 0x8f, 0x3e, 0xc2, 0xd3, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x18, 0x70, 0x49, 0x6e, 0x4f, 0x02, 0x00, 0x00,
}
