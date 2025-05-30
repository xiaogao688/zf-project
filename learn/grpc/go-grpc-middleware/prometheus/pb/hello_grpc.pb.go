// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.7
// source: pb/hello.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RpcServer_SayHello_FullMethodName = "/pb.RpcServer/SayHello"
)

// RpcServerClient is the client API for RpcServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RpcServerClient interface {
	SayHello(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error)
}

type rpcServerClient struct {
	cc grpc.ClientConnInterface
}

func NewRpcServerClient(cc grpc.ClientConnInterface) RpcServerClient {
	return &rpcServerClient{cc}
}

func (c *rpcServerClient) SayHello(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Reply)
	err := c.cc.Invoke(ctx, RpcServer_SayHello_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RpcServerServer is the server API for RpcServer service.
// All implementations must embed UnimplementedRpcServerServer
// for forward compatibility.
type RpcServerServer interface {
	SayHello(context.Context, *Request) (*Reply, error)
	mustEmbedUnimplementedRpcServerServer()
}

// UnimplementedRpcServerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRpcServerServer struct{}

func (UnimplementedRpcServerServer) SayHello(context.Context, *Request) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedRpcServerServer) mustEmbedUnimplementedRpcServerServer() {}
func (UnimplementedRpcServerServer) testEmbeddedByValue()                   {}

// UnsafeRpcServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RpcServerServer will
// result in compilation errors.
type UnsafeRpcServerServer interface {
	mustEmbedUnimplementedRpcServerServer()
}

func RegisterRpcServerServer(s grpc.ServiceRegistrar, srv RpcServerServer) {
	// If the following call pancis, it indicates UnimplementedRpcServerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RpcServer_ServiceDesc, srv)
}

func _RpcServer_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServerServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RpcServer_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServerServer).SayHello(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// RpcServer_ServiceDesc is the grpc.ServiceDesc for RpcServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RpcServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.RpcServer",
	HandlerType: (*RpcServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _RpcServer_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/hello.proto",
}
