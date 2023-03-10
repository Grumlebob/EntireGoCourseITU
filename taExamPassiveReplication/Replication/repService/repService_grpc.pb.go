// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package repService

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

// ReplicationClient is the client API for Replication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReplicationClient interface {
	Put(ctx context.Context, in *KeyVal, opts ...grpc.CallOption) (*Response, error)
	Get(ctx context.Context, in *Int, opts ...grpc.CallOption) (*Int, error)
}

type replicationClient struct {
	cc grpc.ClientConnInterface
}

func NewReplicationClient(cc grpc.ClientConnInterface) ReplicationClient {
	return &replicationClient{cc}
}

func (c *replicationClient) Put(ctx context.Context, in *KeyVal, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/repService.Replication/put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicationClient) Get(ctx context.Context, in *Int, opts ...grpc.CallOption) (*Int, error) {
	out := new(Int)
	err := c.cc.Invoke(ctx, "/repService.Replication/get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicationServer is the server API for Replication service.
// All implementations must embed UnimplementedReplicationServer
// for forward compatibility
type ReplicationServer interface {
	Put(context.Context, *KeyVal) (*Response, error)
	Get(context.Context, *Int) (*Int, error)
	mustEmbedUnimplementedReplicationServer()
}

// UnimplementedReplicationServer must be embedded to have forward compatible implementations.
type UnimplementedReplicationServer struct {
}

func (UnimplementedReplicationServer) Put(context.Context, *KeyVal) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedReplicationServer) Get(context.Context, *Int) (*Int, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedReplicationServer) mustEmbedUnimplementedReplicationServer() {}

// UnsafeReplicationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReplicationServer will
// result in compilation errors.
type UnsafeReplicationServer interface {
	mustEmbedUnimplementedReplicationServer()
}

func RegisterReplicationServer(s grpc.ServiceRegistrar, srv ReplicationServer) {
	s.RegisterService(&Replication_ServiceDesc, srv)
}

func _Replication_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyVal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/repService.Replication/put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServer).Put(ctx, req.(*KeyVal))
	}
	return interceptor(ctx, in, info, handler)
}

func _Replication_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Int)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/repService.Replication/get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServer).Get(ctx, req.(*Int))
	}
	return interceptor(ctx, in, info, handler)
}

// Replication_ServiceDesc is the grpc.ServiceDesc for Replication service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Replication_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "repService.Replication",
	HandlerType: (*ReplicationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "put",
			Handler:    _Replication_Put_Handler,
		},
		{
			MethodName: "get",
			Handler:    _Replication_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "repService/repService.proto",
}
