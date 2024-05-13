// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: server.proto

package pb

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

// RsyncClient is the client API for Rsync service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RsyncClient interface {
	Copy(ctx context.Context, in *CopyRequest, opts ...grpc.CallOption) (*CopyReply, error)
}

type rsyncClient struct {
	cc grpc.ClientConnInterface
}

func NewRsyncClient(cc grpc.ClientConnInterface) RsyncClient {
	return &rsyncClient{cc}
}

func (c *rsyncClient) Copy(ctx context.Context, in *CopyRequest, opts ...grpc.CallOption) (*CopyReply, error) {
	out := new(CopyReply)
	err := c.cc.Invoke(ctx, "/Rsync/Copy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RsyncServer is the server API for Rsync service.
// All implementations must embed UnimplementedRsyncServer
// for forward compatibility
type RsyncServer interface {
	Copy(context.Context, *CopyRequest) (*CopyReply, error)
	mustEmbedUnimplementedRsyncServer()
}

// UnimplementedRsyncServer must be embedded to have forward compatible implementations.
type UnimplementedRsyncServer struct {
}

func (UnimplementedRsyncServer) Copy(context.Context, *CopyRequest) (*CopyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Copy not implemented")
}
func (UnimplementedRsyncServer) mustEmbedUnimplementedRsyncServer() {}

// UnsafeRsyncServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RsyncServer will
// result in compilation errors.
type UnsafeRsyncServer interface {
	mustEmbedUnimplementedRsyncServer()
}

func RegisterRsyncServer(s grpc.ServiceRegistrar, srv RsyncServer) {
	s.RegisterService(&Rsync_ServiceDesc, srv)
}

func _Rsync_Copy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CopyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RsyncServer).Copy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Rsync/Copy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RsyncServer).Copy(ctx, req.(*CopyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Rsync_ServiceDesc is the grpc.ServiceDesc for Rsync service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Rsync_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Rsync",
	HandlerType: (*RsyncServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Copy",
			Handler:    _Rsync_Copy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
