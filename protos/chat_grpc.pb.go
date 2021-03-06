// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package chat

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

// GatewayConnectorClient is the client API for GatewayConnector service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GatewayConnectorClient interface {
	SendChat(ctx context.Context, in *ChatMessage, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetChat(ctx context.Context, in *RetrieveChatMessages, opts ...grpc.CallOption) (*Chats, error)
}

type gatewayConnectorClient struct {
	cc grpc.ClientConnInterface
}

func NewGatewayConnectorClient(cc grpc.ClientConnInterface) GatewayConnectorClient {
	return &gatewayConnectorClient{cc}
}

func (c *gatewayConnectorClient) SendChat(ctx context.Context, in *ChatMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/chat.GatewayConnector/SendChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayConnectorClient) GetChat(ctx context.Context, in *RetrieveChatMessages, opts ...grpc.CallOption) (*Chats, error) {
	out := new(Chats)
	err := c.cc.Invoke(ctx, "/chat.GatewayConnector/GetChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GatewayConnectorServer is the server API for GatewayConnector service.
// All implementations should embed UnimplementedGatewayConnectorServer
// for forward compatibility
type GatewayConnectorServer interface {
	SendChat(context.Context, *ChatMessage) (*emptypb.Empty, error)
	GetChat(context.Context, *RetrieveChatMessages) (*Chats, error)
}

// UnimplementedGatewayConnectorServer should be embedded to have forward compatible implementations.
type UnimplementedGatewayConnectorServer struct {
}

func (UnimplementedGatewayConnectorServer) SendChat(context.Context, *ChatMessage) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendChat not implemented")
}
func (UnimplementedGatewayConnectorServer) GetChat(context.Context, *RetrieveChatMessages) (*Chats, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChat not implemented")
}

// UnsafeGatewayConnectorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GatewayConnectorServer will
// result in compilation errors.
type UnsafeGatewayConnectorServer interface {
	mustEmbedUnimplementedGatewayConnectorServer()
}

func RegisterGatewayConnectorServer(s grpc.ServiceRegistrar, srv GatewayConnectorServer) {
	s.RegisterService(&GatewayConnector_ServiceDesc, srv)
}

func _GatewayConnector_SendChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayConnectorServer).SendChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.GatewayConnector/SendChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayConnectorServer).SendChat(ctx, req.(*ChatMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayConnector_GetChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveChatMessages)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayConnectorServer).GetChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.GatewayConnector/GetChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayConnectorServer).GetChat(ctx, req.(*RetrieveChatMessages))
	}
	return interceptor(ctx, in, info, handler)
}

// GatewayConnector_ServiceDesc is the grpc.ServiceDesc for GatewayConnector service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GatewayConnector_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.GatewayConnector",
	HandlerType: (*GatewayConnectorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendChat",
			Handler:    _GatewayConnector_SendChat_Handler,
		},
		{
			MethodName: "GetChat",
			Handler:    _GatewayConnector_GetChat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat.proto",
}
