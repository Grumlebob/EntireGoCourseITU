// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: grpc/interface.proto

package go_grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NodeClient is the client API for Node service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeClient interface {
	// Heartbeat:
	Heartbeat(ctx context.Context, in *HeartbeatTimestamp, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Announce leader:
	Coordinator(ctx context.Context, in *CoordinatorPort, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Replication:
	HandleAgreementAndReplicationFromLeader(ctx context.Context, in *Replicate, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Broadcast
	BroadcastMessage(ctx context.Context, in *MessageString, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Auction
	Bid(ctx context.Context, in *Bid, opts ...grpc.CallOption) (*Acknowledgement, error)
	Result(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Outcome, error)
	// Add
	Add(ctx context.Context, in *DictionaryAdd, opts ...grpc.CallOption) (*AddResult, error)
	// Read
	Read(ctx context.Context, in *ReadWord, opts ...grpc.CallOption) (*ReadResult, error)
}

type nodeClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeClient(cc grpc.ClientConnInterface) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) Heartbeat(ctx context.Context, in *HeartbeatTimestamp, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/grpc.Node/heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Coordinator(ctx context.Context, in *CoordinatorPort, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/grpc.Node/coordinator", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) HandleAgreementAndReplicationFromLeader(ctx context.Context, in *Replicate, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/grpc.Node/handleAgreementAndReplicationFromLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) BroadcastMessage(ctx context.Context, in *MessageString, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/grpc.Node/broadcastMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Bid(ctx context.Context, in *Bid, opts ...grpc.CallOption) (*Acknowledgement, error) {
	out := new(Acknowledgement)
	err := c.cc.Invoke(ctx, "/grpc.Node/bid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Result(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Outcome, error) {
	out := new(Outcome)
	err := c.cc.Invoke(ctx, "/grpc.Node/result", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Add(ctx context.Context, in *DictionaryAdd, opts ...grpc.CallOption) (*AddResult, error) {
	out := new(AddResult)
	err := c.cc.Invoke(ctx, "/grpc.Node/add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Read(ctx context.Context, in *ReadWord, opts ...grpc.CallOption) (*ReadResult, error) {
	out := new(ReadResult)
	err := c.cc.Invoke(ctx, "/grpc.Node/read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServer is the server API for Node service.
// All implementations must embed UnimplementedNodeServer
// for forward compatibility
type NodeServer interface {
	// Heartbeat:
	Heartbeat(context.Context, *HeartbeatTimestamp) (*emptypb.Empty, error)
	// Announce leader:
	Coordinator(context.Context, *CoordinatorPort) (*emptypb.Empty, error)
	// Replication:
	HandleAgreementAndReplicationFromLeader(context.Context, *Replicate) (*emptypb.Empty, error)
	// Broadcast
	BroadcastMessage(context.Context, *MessageString) (*emptypb.Empty, error)
	// Auction
	Bid(context.Context, *Bid) (*Acknowledgement, error)
	Result(context.Context, *emptypb.Empty) (*Outcome, error)
	// Add
	Add(context.Context, *DictionaryAdd) (*AddResult, error)
	// Read
	Read(context.Context, *ReadWord) (*ReadResult, error)
	mustEmbedUnimplementedNodeServer()
}

// UnimplementedNodeServer must be embedded to have forward compatible implementations.
type UnimplementedNodeServer struct {
}

func (UnimplementedNodeServer) Heartbeat(context.Context, *HeartbeatTimestamp) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedNodeServer) Coordinator(context.Context, *CoordinatorPort) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Coordinator not implemented")
}
func (UnimplementedNodeServer) HandleAgreementAndReplicationFromLeader(context.Context, *Replicate) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleAgreementAndReplicationFromLeader not implemented")
}
func (UnimplementedNodeServer) BroadcastMessage(context.Context, *MessageString) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastMessage not implemented")
}
func (UnimplementedNodeServer) Bid(context.Context, *Bid) (*Acknowledgement, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bid not implemented")
}
func (UnimplementedNodeServer) Result(context.Context, *emptypb.Empty) (*Outcome, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Result not implemented")
}
func (UnimplementedNodeServer) Add(context.Context, *DictionaryAdd) (*AddResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedNodeServer) Read(context.Context, *ReadWord) (*ReadResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedNodeServer) mustEmbedUnimplementedNodeServer() {}

// UnsafeNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServer will
// result in compilation errors.
type UnsafeNodeServer interface {
	mustEmbedUnimplementedNodeServer()
}

func RegisterNodeServer(s grpc.ServiceRegistrar, srv NodeServer) {
	s.RegisterService(&Node_ServiceDesc, srv)
}

func _Node_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatTimestamp)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Node/heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Heartbeat(ctx, req.(*HeartbeatTimestamp))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Coordinator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CoordinatorPort)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Coordinator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Node/coordinator",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Coordinator(ctx, req.(*CoordinatorPort))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_HandleAgreementAndReplicationFromLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Replicate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).HandleAgreementAndReplicationFromLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Node/handleAgreementAndReplicationFromLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).HandleAgreementAndReplicationFromLeader(ctx, req.(*Replicate))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_BroadcastMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageString)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).BroadcastMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Node/broadcastMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).BroadcastMessage(ctx, req.(*MessageString))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Bid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Bid)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Bid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Node/bid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Bid(ctx, req.(*Bid))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Result_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Result(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Node/result",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Result(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DictionaryAdd)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Node/add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Add(ctx, req.(*DictionaryAdd))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadWord)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Node/read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Read(ctx, req.(*ReadWord))
	}
	return interceptor(ctx, in, info, handler)
}

// Node_ServiceDesc is the grpc.ServiceDesc for Node service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Node_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "heartbeat",
			Handler:    _Node_Heartbeat_Handler,
		},
		{
			MethodName: "coordinator",
			Handler:    _Node_Coordinator_Handler,
		},
		{
			MethodName: "handleAgreementAndReplicationFromLeader",
			Handler:    _Node_HandleAgreementAndReplicationFromLeader_Handler,
		},
		{
			MethodName: "broadcastMessage",
			Handler:    _Node_BroadcastMessage_Handler,
		},
		{
			MethodName: "bid",
			Handler:    _Node_Bid_Handler,
		},
		{
			MethodName: "result",
			Handler:    _Node_Result_Handler,
		},
		{
			MethodName: "add",
			Handler:    _Node_Add_Handler,
		},
		{
			MethodName: "read",
			Handler:    _Node_Read_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/interface.proto",
}
