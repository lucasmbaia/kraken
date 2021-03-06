// Code generated by protoc-gen-go.
// source: orchestration.proto
// DO NOT EDIT!

/*
Package orchestrator is a generated protocol buffer package.

It is generated from these files:
	orchestration.proto

It has these top-level messages:
	Task
	Status
	Response
*/
package orchestrator

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StatusCode int32

const (
	StatusCode_Unknown    StatusCode = 0
	StatusCode_Ok         StatusCode = 1
	StatusCode_Failed     StatusCode = 2
	StatusCode_InProgress StatusCode = 3
)

var StatusCode_name = map[int32]string{
	0: "Unknown",
	1: "Ok",
	2: "Failed",
	3: "InProgress",
}
var StatusCode_value = map[string]int32{
	"Unknown":    0,
	"Ok":         1,
	"Failed":     2,
	"InProgress": 3,
}

func (x StatusCode) String() string {
	return proto.EnumName(StatusCode_name, int32(x))
}
func (StatusCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Task struct {
	Name       string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Parameters []byte `protobuf:"bytes,2,opt,name=parameters,proto3" json:"parameters,omitempty"`
	Version    int32  `protobuf:"varint,3,opt,name=version" json:"version,omitempty"`
}

func (m *Task) Reset()                    { *m = Task{} }
func (m *Task) String() string            { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()               {}
func (*Task) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Task) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Task) GetParameters() []byte {
	if m != nil {
		return m.Parameters
	}
	return nil
}

func (m *Task) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

type Status struct {
	Message  string            `protobuf:"bytes,1,opt,name=Message,json=message" json:"Message,omitempty"`
	Code     StatusCode        `protobuf:"varint,2,opt,name=Code,json=code,enum=orchestrator.StatusCode" json:"Code,omitempty"`
	Response map[string][]byte `protobuf:"bytes,3,rep,name=response" json:"response,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *Status) Reset()                    { *m = Status{} }
func (m *Status) String() string            { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()               {}
func (*Status) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Status) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Status) GetCode() StatusCode {
	if m != nil {
		return m.Code
	}
	return StatusCode_Unknown
}

func (m *Status) GetResponse() map[string][]byte {
	if m != nil {
		return m.Response
	}
	return nil
}

type Response struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Response) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*Task)(nil), "orchestrator.Task")
	proto.RegisterType((*Status)(nil), "orchestrator.Status")
	proto.RegisterType((*Response)(nil), "orchestrator.Response")
	proto.RegisterEnum("orchestrator.StatusCode", StatusCode_name, StatusCode_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for OrchestratorService service

type OrchestratorServiceClient interface {
	Workflow(ctx context.Context, in *Task, opts ...grpc.CallOption) (*Response, error)
}

type orchestratorServiceClient struct {
	cc *grpc.ClientConn
}

func NewOrchestratorServiceClient(cc *grpc.ClientConn) OrchestratorServiceClient {
	return &orchestratorServiceClient{cc}
}

func (c *orchestratorServiceClient) Workflow(ctx context.Context, in *Task, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/orchestrator.OrchestratorService/Workflow", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OrchestratorService service

type OrchestratorServiceServer interface {
	Workflow(context.Context, *Task) (*Response, error)
}

func RegisterOrchestratorServiceServer(s *grpc.Server, srv OrchestratorServiceServer) {
	s.RegisterService(&_OrchestratorService_serviceDesc, srv)
}

func _OrchestratorService_Workflow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServiceServer).Workflow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/orchestrator.OrchestratorService/Workflow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServiceServer).Workflow(ctx, req.(*Task))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrchestratorService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "orchestrator.OrchestratorService",
	HandlerType: (*OrchestratorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Workflow",
			Handler:    _OrchestratorService_Workflow_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "orchestration.proto",
}

func init() { proto.RegisterFile("orchestration.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 331 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x91, 0x41, 0x6b, 0xfa, 0x40,
	0x10, 0xc5, 0xdd, 0x24, 0x46, 0xff, 0xa3, 0x7f, 0x09, 0x63, 0x29, 0xc1, 0x43, 0x09, 0x39, 0x85,
	0x52, 0x72, 0xb0, 0x17, 0x69, 0x69, 0x2f, 0xa5, 0x85, 0x1e, 0x8a, 0x25, 0x5a, 0x7a, 0xde, 0x9a,
	0xa9, 0x0d, 0xd1, 0x5d, 0xd9, 0x5d, 0x15, 0x3f, 0x64, 0xbf, 0x53, 0x31, 0x26, 0xd5, 0x40, 0x6f,
	0x33, 0x6f, 0xe6, 0x37, 0xbc, 0xc7, 0x40, 0x5f, 0xaa, 0xd9, 0x17, 0x69, 0xa3, 0xb8, 0xc9, 0xa4,
	0x88, 0x57, 0x4a, 0x1a, 0x89, 0xdd, 0xa3, 0x28, 0x55, 0x38, 0x05, 0x67, 0xca, 0x75, 0x8e, 0x08,
	0x8e, 0xe0, 0x4b, 0xf2, 0x59, 0xc0, 0xa2, 0x7f, 0x49, 0x51, 0xe3, 0x05, 0xc0, 0x8a, 0x2b, 0xbe,
	0x24, 0x43, 0x4a, 0xfb, 0x56, 0xc0, 0xa2, 0x6e, 0x72, 0xa2, 0xa0, 0x0f, 0xad, 0x0d, 0x29, 0x9d,
	0x49, 0xe1, 0xdb, 0x01, 0x8b, 0x9a, 0x49, 0xd5, 0x86, 0xdf, 0x0c, 0xdc, 0x89, 0xe1, 0x66, 0x5d,
	0x2c, 0xbd, 0x90, 0xd6, 0x7c, 0x5e, 0xdd, 0x6e, 0x2d, 0x0f, 0x2d, 0x5e, 0x81, 0xf3, 0x20, 0x53,
	0x2a, 0x0e, 0xf7, 0x86, 0x7e, 0x7c, 0xea, 0x2b, 0x3e, 0xd0, 0xfb, 0x79, 0xe2, 0xcc, 0x64, 0x4a,
	0x78, 0x0f, 0x6d, 0x45, 0x7a, 0x25, 0x85, 0x26, 0xdf, 0x0e, 0xec, 0xa8, 0x33, 0x0c, 0xff, 0x22,
	0xe2, 0xa4, 0x5c, 0x7a, 0x14, 0x46, 0xed, 0x92, 0x5f, 0x66, 0x70, 0x0b, 0xff, 0x6b, 0x23, 0xf4,
	0xc0, 0xce, 0x69, 0x57, 0x9a, 0xda, 0x97, 0x78, 0x06, 0xcd, 0x0d, 0x5f, 0xac, 0xa9, 0x8c, 0x7a,
	0x68, 0x6e, 0xac, 0x11, 0x0b, 0x07, 0xd0, 0xae, 0x60, 0xec, 0x81, 0x95, 0xa5, 0x25, 0x66, 0x65,
	0xe9, 0xe5, 0x1d, 0xc0, 0xd1, 0x2c, 0x76, 0xa0, 0xf5, 0x26, 0x72, 0x21, 0xb7, 0xc2, 0x6b, 0xa0,
	0x0b, 0xd6, 0x38, 0xf7, 0x18, 0x02, 0xb8, 0x4f, 0x3c, 0x5b, 0x50, 0xea, 0x59, 0xd8, 0x03, 0x78,
	0x16, 0xaf, 0x4a, 0xce, 0x15, 0x69, 0xed, 0xd9, 0xc3, 0x31, 0xf4, 0xc7, 0x27, 0x31, 0x26, 0xa4,
	0x36, 0xd9, 0x8c, 0x70, 0x04, 0xed, 0x77, 0xa9, 0xf2, 0xcf, 0x85, 0xdc, 0x22, 0xd6, 0x83, 0xee,
	0xff, 0x35, 0x38, 0xaf, 0x6b, 0x95, 0xbb, 0xb0, 0xf1, 0xe1, 0x16, 0x6f, 0xbe, 0xfe, 0x09, 0x00,
	0x00, 0xff, 0xff, 0xbe, 0xdd, 0xae, 0x21, 0xfd, 0x01, 0x00, 0x00,
}
