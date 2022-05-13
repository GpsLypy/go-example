// Code generated by protoc-gen-go.
// source: google/bigtable/v1/bigtable_service.proto
// DO NOT EDIT!

package bigtable

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf2 "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for BigtableService service

type BigtableServiceClient interface {
	// Streams back the contents of all requested rows, optionally applying
	// the same Reader filter to each. Depending on their size, rows may be
	// broken up across multiple responses, but atomicity of each row will still
	// be preserved.
	ReadRows(ctx context.Context, in *ReadRowsRequest, opts ...grpc.CallOption) (BigtableService_ReadRowsClient, error)
	// Returns a sample of row keys in the table. The returned row keys will
	// delimit contiguous sections of the table of approximately equal size,
	// which can be used to break up the data for distributed tasks like
	// mapreduces.
	SampleRowKeys(ctx context.Context, in *SampleRowKeysRequest, opts ...grpc.CallOption) (BigtableService_SampleRowKeysClient, error)
	// Mutates a row atomically. Cells already present in the row are left
	// unchanged unless explicitly changed by 'mutation'.
	MutateRow(ctx context.Context, in *MutateRowRequest, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
	// Mutates multiple rows in a batch. Each individual row is mutated
	// atomically as in MutateRow, but the entire batch is not executed
	// atomically.
	MutateRows(ctx context.Context, in *MutateRowsRequest, opts ...grpc.CallOption) (*MutateRowsResponse, error)
	// Mutates a row atomically based on the output of a predicate Reader filter.
	CheckAndMutateRow(ctx context.Context, in *CheckAndMutateRowRequest, opts ...grpc.CallOption) (*CheckAndMutateRowResponse, error)
	// Modifies a row atomically, reading the latest existing timestamp/value from
	// the specified columns and writing a new value at
	// max(existing timestamp, current server time) based on pre-defined
	// read/modify/write rules. Returns the new contents of all modified cells.
	ReadModifyWriteRow(ctx context.Context, in *ReadModifyWriteRowRequest, opts ...grpc.CallOption) (*Row, error)
}

type bigtableServiceClient struct {
	cc *grpc.ClientConn
}

func NewBigtableServiceClient(cc *grpc.ClientConn) BigtableServiceClient {
	return &bigtableServiceClient{cc}
}

func (c *bigtableServiceClient) ReadRows(ctx context.Context, in *ReadRowsRequest, opts ...grpc.CallOption) (BigtableService_ReadRowsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_BigtableService_serviceDesc.Streams[0], c.cc, "/google.bigtable.v1.BigtableService/ReadRows", opts...)
	if err != nil {
		return nil, err
	}
	x := &bigtableServiceReadRowsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BigtableService_ReadRowsClient interface {
	Recv() (*ReadRowsResponse, error)
	grpc.ClientStream
}

type bigtableServiceReadRowsClient struct {
	grpc.ClientStream
}

func (x *bigtableServiceReadRowsClient) Recv() (*ReadRowsResponse, error) {
	m := new(ReadRowsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *bigtableServiceClient) SampleRowKeys(ctx context.Context, in *SampleRowKeysRequest, opts ...grpc.CallOption) (BigtableService_SampleRowKeysClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_BigtableService_serviceDesc.Streams[1], c.cc, "/google.bigtable.v1.BigtableService/SampleRowKeys", opts...)
	if err != nil {
		return nil, err
	}
	x := &bigtableServiceSampleRowKeysClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BigtableService_SampleRowKeysClient interface {
	Recv() (*SampleRowKeysResponse, error)
	grpc.ClientStream
}

type bigtableServiceSampleRowKeysClient struct {
	grpc.ClientStream
}

func (x *bigtableServiceSampleRowKeysClient) Recv() (*SampleRowKeysResponse, error) {
	m := new(SampleRowKeysResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *bigtableServiceClient) MutateRow(ctx context.Context, in *MutateRowRequest, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/google.bigtable.v1.BigtableService/MutateRow", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableServiceClient) MutateRows(ctx context.Context, in *MutateRowsRequest, opts ...grpc.CallOption) (*MutateRowsResponse, error) {
	out := new(MutateRowsResponse)
	err := grpc.Invoke(ctx, "/google.bigtable.v1.BigtableService/MutateRows", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableServiceClient) CheckAndMutateRow(ctx context.Context, in *CheckAndMutateRowRequest, opts ...grpc.CallOption) (*CheckAndMutateRowResponse, error) {
	out := new(CheckAndMutateRowResponse)
	err := grpc.Invoke(ctx, "/google.bigtable.v1.BigtableService/CheckAndMutateRow", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableServiceClient) ReadModifyWriteRow(ctx context.Context, in *ReadModifyWriteRowRequest, opts ...grpc.CallOption) (*Row, error) {
	out := new(Row)
	err := grpc.Invoke(ctx, "/google.bigtable.v1.BigtableService/ReadModifyWriteRow", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for BigtableService service

type BigtableServiceServer interface {
	// Streams back the contents of all requested rows, optionally applying
	// the same Reader filter to each. Depending on their size, rows may be
	// broken up across multiple responses, but atomicity of each row will still
	// be preserved.
	ReadRows(*ReadRowsRequest, BigtableService_ReadRowsServer) error
	// Returns a sample of row keys in the table. The returned row keys will
	// delimit contiguous sections of the table of approximately equal size,
	// which can be used to break up the data for distributed tasks like
	// mapreduces.
	SampleRowKeys(*SampleRowKeysRequest, BigtableService_SampleRowKeysServer) error
	// Mutates a row atomically. Cells already present in the row are left
	// unchanged unless explicitly changed by 'mutation'.
	MutateRow(context.Context, *MutateRowRequest) (*google_protobuf2.Empty, error)
	// Mutates multiple rows in a batch. Each individual row is mutated
	// atomically as in MutateRow, but the entire batch is not executed
	// atomically.
	MutateRows(context.Context, *MutateRowsRequest) (*MutateRowsResponse, error)
	// Mutates a row atomically based on the output of a predicate Reader filter.
	CheckAndMutateRow(context.Context, *CheckAndMutateRowRequest) (*CheckAndMutateRowResponse, error)
	// Modifies a row atomically, reading the latest existing timestamp/value from
	// the specified columns and writing a new value at
	// max(existing timestamp, current server time) based on pre-defined
	// read/modify/write rules. Returns the new contents of all modified cells.
	ReadModifyWriteRow(context.Context, *ReadModifyWriteRowRequest) (*Row, error)
}

func RegisterBigtableServiceServer(s *grpc.Server, srv BigtableServiceServer) {
	s.RegisterService(&_BigtableService_serviceDesc, srv)
}

func _BigtableService_ReadRows_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ReadRowsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BigtableServiceServer).ReadRows(m, &bigtableServiceReadRowsServer{stream})
}

type BigtableService_ReadRowsServer interface {
	Send(*ReadRowsResponse) error
	grpc.ServerStream
}

type bigtableServiceReadRowsServer struct {
	grpc.ServerStream
}

func (x *bigtableServiceReadRowsServer) Send(m *ReadRowsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _BigtableService_SampleRowKeys_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SampleRowKeysRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BigtableServiceServer).SampleRowKeys(m, &bigtableServiceSampleRowKeysServer{stream})
}

type BigtableService_SampleRowKeysServer interface {
	Send(*SampleRowKeysResponse) error
	grpc.ServerStream
}

type bigtableServiceSampleRowKeysServer struct {
	grpc.ServerStream
}

func (x *bigtableServiceSampleRowKeysServer) Send(m *SampleRowKeysResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _BigtableService_MutateRow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutateRowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableServiceServer).MutateRow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.v1.BigtableService/MutateRow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableServiceServer).MutateRow(ctx, req.(*MutateRowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableService_MutateRows_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutateRowsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableServiceServer).MutateRows(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.v1.BigtableService/MutateRows",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableServiceServer).MutateRows(ctx, req.(*MutateRowsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableService_CheckAndMutateRow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckAndMutateRowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableServiceServer).CheckAndMutateRow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.v1.BigtableService/CheckAndMutateRow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableServiceServer).CheckAndMutateRow(ctx, req.(*CheckAndMutateRowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableService_ReadModifyWriteRow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadModifyWriteRowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableServiceServer).ReadModifyWriteRow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.v1.BigtableService/ReadModifyWriteRow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableServiceServer).ReadModifyWriteRow(ctx, req.(*ReadModifyWriteRowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BigtableService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.bigtable.v1.BigtableService",
	HandlerType: (*BigtableServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MutateRow",
			Handler:    _BigtableService_MutateRow_Handler,
		},
		{
			MethodName: "MutateRows",
			Handler:    _BigtableService_MutateRows_Handler,
		},
		{
			MethodName: "CheckAndMutateRow",
			Handler:    _BigtableService_CheckAndMutateRow_Handler,
		},
		{
			MethodName: "ReadModifyWriteRow",
			Handler:    _BigtableService_ReadModifyWriteRow_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReadRows",
			Handler:       _BigtableService_ReadRows_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SampleRowKeys",
			Handler:       _BigtableService_SampleRowKeys_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "google/bigtable/v1/bigtable_service.proto",
}

func init() { proto.RegisterFile("google/bigtable/v1/bigtable_service.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 521 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xcd, 0x6e, 0xd4, 0x30,
	0x10, 0xc7, 0x65, 0x0e, 0xa8, 0x58, 0x42, 0x08, 0x4b, 0x14, 0x69, 0xe1, 0x14, 0xa0, 0xa2, 0x11,
	0x8d, 0xdb, 0x72, 0x0b, 0xe2, 0xd0, 0x45, 0x50, 0x21, 0x58, 0x51, 0xa5, 0xe2, 0x43, 0xe5, 0xb0,
	0x78, 0x93, 0x69, 0x08, 0x4d, 0xe2, 0x60, 0x7b, 0x37, 0x5a, 0xaa, 0x5e, 0x38, 0x71, 0xe7, 0x11,
	0x10, 0x17, 0x5e, 0x80, 0x23, 0xef, 0x00, 0x67, 0x6e, 0x3c, 0x08, 0xb2, 0x63, 0x2f, 0x2c, 0x0d,
	0xcb, 0x8a, 0xee, 0x29, 0x8e, 0xe6, 0x3f, 0x33, 0xbf, 0xff, 0xf8, 0x03, 0xaf, 0xa6, 0x9c, 0xa7,
	0x39, 0xd0, 0x41, 0x96, 0x2a, 0x36, 0xc8, 0x81, 0x8e, 0x36, 0x26, 0xeb, 0xbe, 0x04, 0x31, 0xca,
	0x62, 0x08, 0x2a, 0xc1, 0x15, 0x27, 0xa4, 0x91, 0x06, 0x2e, 0x1c, 0x8c, 0x36, 0x3a, 0x97, 0x6d,
	0x3a, 0xab, 0x32, 0xca, 0xca, 0x92, 0x2b, 0xa6, 0x32, 0x5e, 0xca, 0x26, 0xa3, 0xb3, 0x32, 0xab,
	0x78, 0xc2, 0x14, 0xb3, 0xba, 0xcd, 0x39, 0x20, 0xfa, 0x05, 0x48, 0xc9, 0x52, 0x70, 0xb5, 0x2f,
	0xd9, 0x1c, 0xf3, 0x37, 0x18, 0xee, 0x53, 0x28, 0x2a, 0x35, 0x6e, 0x82, 0x9b, 0xdf, 0x97, 0xf0,
	0xb9, 0xae, 0x2d, 0xb0, 0xdb, 0xe4, 0x93, 0x8f, 0x08, 0x2f, 0x45, 0xc0, 0x92, 0x88, 0xd7, 0x92,
	0x5c, 0x09, 0x8e, 0x9b, 0x09, 0x5c, 0x34, 0x82, 0xd7, 0x43, 0x90, 0xaa, 0x73, 0x75, 0xb6, 0x48,
	0x56, 0xbc, 0x94, 0xe0, 0x3d, 0x7c, 0xfb, 0xed, 0xc7, 0xfb, 0x53, 0xf7, 0xbc, 0x2d, 0x4d, 0x7d,
	0xd8, 0x30, 0x97, 0xac, 0x80, 0xdb, 0x95, 0xe0, 0xaf, 0x20, 0x56, 0x92, 0xfa, 0xf4, 0x0d, 0x2f,
	0x41, 0x7f, 0xe3, 0x7c, 0x28, 0x15, 0x08, 0xbd, 0x34, 0x42, 0x49, 0xfd, 0x23, 0x2a, 0x78, 0x2d,
	0x43, 0x01, 0x2c, 0x09, 0x91, 0xbf, 0x8e, 0xc8, 0x67, 0x84, 0xcf, 0xee, 0xb2, 0xa2, 0xca, 0x21,
	0xe2, 0xf5, 0x03, 0x18, 0x4b, 0x72, 0xbd, 0x8d, 0x63, 0x4a, 0xe2, 0x88, 0x57, 0xe7, 0x50, 0x5a,
	0xec, 0x47, 0x06, 0xfb, 0x3e, 0xd9, 0x3e, 0x11, 0xb6, 0x34, 0xb5, 0x75, 0xe1, 0x75, 0x44, 0x3e,
	0x20, 0x7c, 0xa6, 0x37, 0x54, 0x4c, 0xe9, 0x66, 0xa4, 0x75, 0x7a, 0x93, 0xb0, 0x23, 0x5e, 0x76,
	0x2a, 0xb7, 0x8f, 0xc1, 0x5d, 0xbd, 0x8f, 0xde, 0x33, 0x83, 0x17, 0x79, 0xbd, 0x93, 0xe0, 0xd1,
	0x43, 0xc1, 0xeb, 0xfe, 0x01, 0x8c, 0x8f, 0xc2, 0xc2, 0x34, 0x0e, 0x91, 0x4f, 0x3e, 0x21, 0x8c,
	0x27, 0x18, 0x92, 0x5c, 0x9b, 0x89, 0x39, 0x99, 0xec, 0xca, 0xbf, 0x64, 0x76, 0xac, 0x3d, 0xc3,
	0xbd, 0xed, 0x75, 0xff, 0x93, 0xdb, 0x82, 0xea, 0x9a, 0x1a, 0xf6, 0x2b, 0xc2, 0xe7, 0xef, 0xbc,
	0x84, 0xf8, 0x60, 0xab, 0x4c, 0x7e, 0x8d, 0xf6, 0x46, 0x1b, 0xcc, 0x31, 0x99, 0x43, 0x5f, 0x9b,
	0x53, 0x6d, 0x1d, 0xbc, 0x30, 0x0e, 0xf6, 0xbc, 0xc7, 0x0b, 0x9a, 0x7c, 0x3c, 0xd5, 0x49, 0x9b,
	0xfa, 0x82, 0x30, 0xd1, 0xd7, 0xa8, 0xc7, 0x93, 0x6c, 0x7f, 0xfc, 0x54, 0x64, 0x8d, 0xab, 0xb5,
	0xbf, 0x5d, 0xb7, 0x69, 0x9d, 0xb3, 0x75, 0xb1, 0x55, 0xce, 0x6b, 0x8f, 0x19, 0x03, 0xcf, 0xbd,
	0x27, 0x0b, 0x32, 0x20, 0xa6, 0x11, 0x42, 0xe4, 0x77, 0x2b, 0xbc, 0x1c, 0xf3, 0xa2, 0x05, 0xa0,
	0x7b, 0xe1, 0x8f, 0x67, 0x47, 0xee, 0xe8, 0x73, 0xbd, 0x83, 0xf6, 0x42, 0x2b, 0x4e, 0x79, 0xce,
	0xca, 0x34, 0xe0, 0x22, 0xa5, 0x29, 0x94, 0xe6, 0xd4, 0xd3, 0x26, 0xc4, 0xaa, 0x4c, 0xfe, 0xfe,
	0x04, 0xde, 0x72, 0xeb, 0x77, 0x08, 0x0d, 0x4e, 0x1b, 0xe5, 0xcd, 0x9f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x4c, 0x27, 0x6e, 0x9a, 0xb0, 0x05, 0x00, 0x00,
}
