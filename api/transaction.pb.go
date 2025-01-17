// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/aipadad/aipa/api/transaction.proto

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	github.com/aipadad/aipa/api/transaction.proto
	github.com/aipadad/aipa/api/basic-transaction.proto
	github.com/aipadad/aipa/api/chain.proto

It has these top-level messages:
	Transaction
	BasicTransaction
	SendTransactionResponse
	GetTransactionRequest
	GetTransactionResponse
	GetBlockRequest
	GetBlockResponse
	GetInfoRequest
	GetInfoResponse
	GetAccountRequest
	GetAccountResponse
	GetKeyValueRequest
	GetKeyValueResponse
	GetAbiRequest
	GetAbiResponse
	GetTransferCreditRequest
	GetTransferCreditResponse
*/
package api

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

// Transaction definition for gprc interface
type Transaction struct {
	Version     uint32 `protobuf:"varint,1,opt,name=version" json:"version"`
	CursorNum   uint64 `protobuf:"varint,2,opt,name=cursor_num,json=cursorNum" json:"cursor_num"`
	CursorLabel uint32 `protobuf:"varint,3,opt,name=cursor_label,json=cursorLabel" json:"cursor_label"`
	Lifetime    uint64 `protobuf:"varint,4,opt,name=lifetime" json:"lifetime"`
	Sender      string `protobuf:"bytes,5,opt,name=sender" json:"sender"`
	Contract    string `protobuf:"bytes,6,opt,name=contract" json:"contract"`
	Method      string `protobuf:"bytes,7,opt,name=method" json:"method"`
	Param       string `protobuf:"bytes,8,opt,name=param" json:"param"`
	SigAlg      uint32 `protobuf:"varint,9,opt,name=sig_alg,json=sigAlg" json:"sig_alg"`
	Signature   string `protobuf:"bytes,10,opt,name=signature" json:"signature"`
}

func (m *Transaction) Reset()                    { *m = Transaction{} }
func (m *Transaction) String() string            { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()               {}
func (*Transaction) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Transaction) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Transaction) GetCursorNum() uint64 {
	if m != nil {
		return m.CursorNum
	}
	return 0
}

func (m *Transaction) GetCursorLabel() uint32 {
	if m != nil {
		return m.CursorLabel
	}
	return 0
}

func (m *Transaction) GetLifetime() uint64 {
	if m != nil {
		return m.Lifetime
	}
	return 0
}

func (m *Transaction) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *Transaction) GetContract() string {
	if m != nil {
		return m.Contract
	}
	return ""
}

func (m *Transaction) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *Transaction) GetParam() string {
	if m != nil {
		return m.Param
	}
	return ""
}

func (m *Transaction) GetSigAlg() uint32 {
	if m != nil {
		return m.SigAlg
	}
	return 0
}

func (m *Transaction) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func init() {
	proto.RegisterType((*Transaction)(nil), "api.Transaction")
}

func init() {
	proto.RegisterFile("github.com/aipadad/aipa/api/transaction.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 265 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0x4f, 0x4b, 0xc4, 0x30,
	0x10, 0x47, 0xe9, 0xfe, 0x69, 0xb7, 0xb3, 0x7a, 0x09, 0xa2, 0x41, 0x14, 0xaa, 0x88, 0xf4, 0xe2,
	0xf6, 0x20, 0x78, 0xd7, 0xb3, 0x78, 0x28, 0x9e, 0xbc, 0x2c, 0x69, 0x36, 0x66, 0x23, 0x4d, 0xa6,
	0x24, 0x53, 0x3f, 0x8e, 0x9f, 0x55, 0x9a, 0xd6, 0xdd, 0xab, 0xc7, 0xf7, 0x7e, 0x3c, 0x18, 0x06,
	0x9e, 0xb4, 0xa1, 0x7d, 0xdf, 0x6c, 0x24, 0xda, 0xaa, 0x41, 0x22, 0x0c, 0x0f, 0x9d, 0xc7, 0x2f,
	0x25, 0x69, 0xc2, 0x4a, 0x74, 0xa6, 0x22, 0x2f, 0x5c, 0x10, 0x92, 0x0c, 0xba, 0x4d, 0xe7, 0x91,
	0x90, 0xcd, 0x45, 0x67, 0x6e, 0x7f, 0x66, 0xb0, 0x7e, 0x3f, 0x4e, 0x8c, 0x43, 0xf6, 0xad, 0x7c,
	0x30, 0xe8, 0x78, 0x52, 0x24, 0xe5, 0x69, 0xfd, 0x87, 0xec, 0x1a, 0x40, 0xf6, 0x3e, 0xa0, 0xdf,
	0xba, 0xde, 0xf2, 0x59, 0x91, 0x94, 0x8b, 0x3a, 0x1f, 0xcd, 0x5b, 0x6f, 0xd9, 0x0d, 0x9c, 0x4c,
	0x73, 0x2b, 0x1a, 0xd5, 0xf2, 0x79, 0xac, 0xd7, 0xa3, 0x7b, 0x1d, 0x14, 0xbb, 0x84, 0x55, 0x6b,
	0x3e, 0x15, 0x19, 0xab, 0xf8, 0x22, 0xf6, 0x07, 0x66, 0xe7, 0x90, 0x06, 0xe5, 0x76, 0xca, 0xf3,
	0x65, 0x91, 0x94, 0x79, 0x3d, 0xd1, 0xd0, 0x48, 0x74, 0xe4, 0x85, 0x24, 0x9e, 0xc6, 0xe5, 0xc0,
	0x43, 0x63, 0x15, 0xed, 0x71, 0xc7, 0xb3, 0xb1, 0x19, 0x89, 0x9d, 0xc1, 0xb2, 0x13, 0x5e, 0x58,
	0xbe, 0x8a, 0x7a, 0x04, 0x76, 0x01, 0x59, 0x30, 0x7a, 0x2b, 0x5a, 0xcd, 0xf3, 0x78, 0x5b, 0x1a,
	0x8c, 0x7e, 0x6e, 0x35, 0xbb, 0x82, 0x3c, 0x18, 0xed, 0x04, 0xf5, 0x5e, 0x71, 0x88, 0xc9, 0x51,
	0xbc, 0xdc, 0x7f, 0xdc, 0xfd, 0xe7, 0xbf, 0x4d, 0x1a, 0x9f, 0xfa, 0xf8, 0x1b, 0x00, 0x00, 0xff,
	0xff, 0x94, 0x0f, 0x6b, 0xc1, 0x8e, 0x01, 0x00, 0x00,
}
