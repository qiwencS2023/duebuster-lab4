// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: coordinator.proto

package main

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CoordinatorServiceClient is the client API for CoordinatorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CoordinatorServiceClient interface {
	CreateTable(ctx context.Context, in *Table, opts ...grpc.CallOption) (*CoordinatorResponse, error)
	DeleteTable(ctx context.Context, in *Table, opts ...grpc.CallOption) (*CoordinatorResponse, error)
	InsertLine(ctx context.Context, in *Line, opts ...grpc.CallOption) (*CoordinatorResponse, error)
	DeleteLine(ctx context.Context, in *Line, opts ...grpc.CallOption) (*CoordinatorResponse, error)
	GetLine(ctx context.Context, in *Line, opts ...grpc.CallOption) (*Line, error)
	UpdateLine(ctx context.Context, in *Line, opts ...grpc.CallOption) (*CoordinatorResponse, error)
}

type coordinatorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCoordinatorServiceClient(cc grpc.ClientConnInterface) CoordinatorServiceClient {
	return &coordinatorServiceClient{cc}
}

func (c *coordinatorServiceClient) CreateTable(ctx context.Context, in *Table, opts ...grpc.CallOption) (*CoordinatorResponse, error) {
	out := new(CoordinatorResponse)
	err := c.cc.Invoke(ctx, "/CoordinatorService/CreateTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coordinatorServiceClient) DeleteTable(ctx context.Context, in *Table, opts ...grpc.CallOption) (*CoordinatorResponse, error) {
	out := new(CoordinatorResponse)
	err := c.cc.Invoke(ctx, "/CoordinatorService/DeleteTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coordinatorServiceClient) InsertLine(ctx context.Context, in *Line, opts ...grpc.CallOption) (*CoordinatorResponse, error) {
	out := new(CoordinatorResponse)
	err := c.cc.Invoke(ctx, "/CoordinatorService/InsertLine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coordinatorServiceClient) DeleteLine(ctx context.Context, in *Line, opts ...grpc.CallOption) (*CoordinatorResponse, error) {
	out := new(CoordinatorResponse)
	err := c.cc.Invoke(ctx, "/CoordinatorService/DeleteLine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coordinatorServiceClient) GetLine(ctx context.Context, in *Line, opts ...grpc.CallOption) (*Line, error) {
	out := new(Line)
	err := c.cc.Invoke(ctx, "/CoordinatorService/GetLine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coordinatorServiceClient) UpdateLine(ctx context.Context, in *Line, opts ...grpc.CallOption) (*CoordinatorResponse, error) {
	out := new(CoordinatorResponse)
	err := c.cc.Invoke(ctx, "/CoordinatorService/UpdateLine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CoordinatorServiceServer is the server API for CoordinatorService service.
// All implementations must embed UnimplementedCoordinatorServiceServer
// for forward compatibility
type CoordinatorServiceServer interface {
	CreateTable(context.Context, *Table) (*CoordinatorResponse, error)
	DeleteTable(context.Context, *Table) (*CoordinatorResponse, error)
	InsertLine(context.Context, *Line) (*CoordinatorResponse, error)
	DeleteLine(context.Context, *Line) (*CoordinatorResponse, error)
	GetLine(context.Context, *Line) (*Line, error)
	UpdateLine(context.Context, *Line) (*CoordinatorResponse, error)
	mustEmbedUnimplementedCoordinatorServiceServer()
}

// UnimplementedCoordinatorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCoordinatorServiceServer struct {
}

func (UnimplementedCoordinatorServiceServer) CreateTable(context.Context, *Table) (*CoordinatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTable not implemented")
}
func (UnimplementedCoordinatorServiceServer) DeleteTable(context.Context, *Table) (*CoordinatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTable not implemented")
}
func (UnimplementedCoordinatorServiceServer) InsertLine(context.Context, *Line) (*CoordinatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertLine not implemented")
}
func (UnimplementedCoordinatorServiceServer) DeleteLine(context.Context, *Line) (*CoordinatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLine not implemented")
}
func (UnimplementedCoordinatorServiceServer) GetLine(context.Context, *Line) (*Line, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLine not implemented")
}
func (UnimplementedCoordinatorServiceServer) UpdateLine(context.Context, *Line) (*CoordinatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLine not implemented")
}
func (UnimplementedCoordinatorServiceServer) mustEmbedUnimplementedCoordinatorServiceServer() {}

// UnsafeCoordinatorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CoordinatorServiceServer will
// result in compilation errors.
type UnsafeCoordinatorServiceServer interface {
	mustEmbedUnimplementedCoordinatorServiceServer()
}

func RegisterCoordinatorServiceServer(s grpc.ServiceRegistrar, srv CoordinatorServiceServer) {
	s.RegisterService(&CoordinatorService_ServiceDesc, srv)
}

func _CoordinatorService_CreateTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Table)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinatorServiceServer).CreateTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CoordinatorService/CreateTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinatorServiceServer).CreateTable(ctx, req.(*Table))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoordinatorService_DeleteTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Table)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinatorServiceServer).DeleteTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CoordinatorService/DeleteTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinatorServiceServer).DeleteTable(ctx, req.(*Table))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoordinatorService_InsertLine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Line)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinatorServiceServer).InsertLine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CoordinatorService/InsertLine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinatorServiceServer).InsertLine(ctx, req.(*Line))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoordinatorService_DeleteLine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Line)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinatorServiceServer).DeleteLine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CoordinatorService/DeleteLine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinatorServiceServer).DeleteLine(ctx, req.(*Line))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoordinatorService_GetLine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Line)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinatorServiceServer).GetLine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CoordinatorService/GetLine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinatorServiceServer).GetLine(ctx, req.(*Line))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoordinatorService_UpdateLine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Line)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoordinatorServiceServer).UpdateLine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CoordinatorService/UpdateLine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoordinatorServiceServer).UpdateLine(ctx, req.(*Line))
	}
	return interceptor(ctx, in, info, handler)
}

// CoordinatorService_ServiceDesc is the grpc.ServiceDesc for CoordinatorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CoordinatorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "CoordinatorService",
	HandlerType: (*CoordinatorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTable",
			Handler:    _CoordinatorService_CreateTable_Handler,
		},
		{
			MethodName: "DeleteTable",
			Handler:    _CoordinatorService_DeleteTable_Handler,
		},
		{
			MethodName: "InsertLine",
			Handler:    _CoordinatorService_InsertLine_Handler,
		},
		{
			MethodName: "DeleteLine",
			Handler:    _CoordinatorService_DeleteLine_Handler,
		},
		{
			MethodName: "GetLine",
			Handler:    _CoordinatorService_GetLine_Handler,
		},
		{
			MethodName: "UpdateLine",
			Handler:    _CoordinatorService_UpdateLine_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "coordinator.proto",
}
