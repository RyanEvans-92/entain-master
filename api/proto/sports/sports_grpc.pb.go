// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sports

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

// EventsClient is the client API for Events service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventsClient interface {
	// ListEvents will return a collection of all Events.
	ListEvents(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ListEventsResponse, error)
}

type eventsClient struct {
	cc grpc.ClientConnInterface
}

func NewEventsClient(cc grpc.ClientConnInterface) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) ListEvents(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ListEventsResponse, error) {
	out := new(ListEventsResponse)
	err := c.cc.Invoke(ctx, "/sports.Events/ListEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventsServer is the server API for Events service.
// All implementations must embed UnimplementedEventsServer
// for forward compatibility
type EventsServer interface {
	// ListEvents will return a collection of all Events.
	ListEvents(context.Context, *ListEventsRequest) (*ListEventsResponse, error)
	mustEmbedUnimplementedEventsServer()
}

// UnimplementedEventsServer must be embedded to have forward compatible implementations.
type UnimplementedEventsServer struct {
}

func (UnimplementedEventsServer) ListEvents(context.Context, *ListEventsRequest) (*ListEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEvents not implemented")
}
func (UnimplementedEventsServer) mustEmbedUnimplementedEventsServer() {}

// UnsafeEventsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventsServer will
// result in compilation errors.
type UnsafeEventsServer interface {
	mustEmbedUnimplementedEventsServer()
}

func RegisterEventsServer(s grpc.ServiceRegistrar, srv EventsServer) {
	s.RegisterService(&Events_ServiceDesc, srv)
}

func _Events_ListEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).ListEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sports.Events/ListEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).ListEvents(ctx, req.(*ListEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Events_ServiceDesc is the grpc.ServiceDesc for Events service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Events_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sports.Events",
	HandlerType: (*EventsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListEvents",
			Handler:    _Events_ListEvents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sports/sports.proto",
}
