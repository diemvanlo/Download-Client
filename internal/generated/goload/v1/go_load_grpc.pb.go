// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: goload/v1/go_load.proto

package go_loadv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	GoLoadService_CreateAccount_FullMethodName       = "/go_load.v1.GoLoadService/CreateAccount"
	GoLoadService_CreateSession_FullMethodName       = "/go_load.v1.GoLoadService/CreateSession"
	GoLoadService_CreateDownloadTask_FullMethodName  = "/go_load.v1.GoLoadService/CreateDownloadTask"
	GoLoadService_GetDownloadTaskList_FullMethodName = "/go_load.v1.GoLoadService/GetDownloadTaskList"
	GoLoadService_UpdateDownloadTask_FullMethodName  = "/go_load.v1.GoLoadService/UpdateDownloadTask"
	GoLoadService_DeleteDownloadTask_FullMethodName  = "/go_load.v1.GoLoadService/DeleteDownloadTask"
	GoLoadService_GetDownloadTaskFile_FullMethodName = "/go_load.v1.GoLoadService/GetDownloadTaskFile"
)

// GoLoadServiceClient is the client API for GoLoadService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GoLoadServiceClient interface {
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error)
	CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*CreateSessionResponse, error)
	CreateDownloadTask(ctx context.Context, in *CreateDownloadTaskRequest, opts ...grpc.CallOption) (*CreateDownloadTaskResponse, error)
	GetDownloadTaskList(ctx context.Context, in *GetDownloadTaskListRequest, opts ...grpc.CallOption) (*GetDownloadTaskListResponse, error)
	UpdateDownloadTask(ctx context.Context, in *UpdateDownloadTaskRequest, opts ...grpc.CallOption) (*UpdateDownloadTaskResponse, error)
	DeleteDownloadTask(ctx context.Context, in *DeleteDownloadTaskRequest, opts ...grpc.CallOption) (*DeleteDownloadTaskResponse, error)
	GetDownloadTaskFile(ctx context.Context, in *GetDownloadTaskFileRequest, opts ...grpc.CallOption) (GoLoadService_GetDownloadTaskFileClient, error)
}

type goLoadServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGoLoadServiceClient(cc grpc.ClientConnInterface) GoLoadServiceClient {
	return &goLoadServiceClient{cc}
}

func (c *goLoadServiceClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAccountResponse)
	err := c.cc.Invoke(ctx, GoLoadService_CreateAccount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goLoadServiceClient) CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*CreateSessionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateSessionResponse)
	err := c.cc.Invoke(ctx, GoLoadService_CreateSession_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goLoadServiceClient) CreateDownloadTask(ctx context.Context, in *CreateDownloadTaskRequest, opts ...grpc.CallOption) (*CreateDownloadTaskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateDownloadTaskResponse)
	err := c.cc.Invoke(ctx, GoLoadService_CreateDownloadTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goLoadServiceClient) GetDownloadTaskList(ctx context.Context, in *GetDownloadTaskListRequest, opts ...grpc.CallOption) (*GetDownloadTaskListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDownloadTaskListResponse)
	err := c.cc.Invoke(ctx, GoLoadService_GetDownloadTaskList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goLoadServiceClient) UpdateDownloadTask(ctx context.Context, in *UpdateDownloadTaskRequest, opts ...grpc.CallOption) (*UpdateDownloadTaskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateDownloadTaskResponse)
	err := c.cc.Invoke(ctx, GoLoadService_UpdateDownloadTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goLoadServiceClient) DeleteDownloadTask(ctx context.Context, in *DeleteDownloadTaskRequest, opts ...grpc.CallOption) (*DeleteDownloadTaskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteDownloadTaskResponse)
	err := c.cc.Invoke(ctx, GoLoadService_DeleteDownloadTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goLoadServiceClient) GetDownloadTaskFile(ctx context.Context, in *GetDownloadTaskFileRequest, opts ...grpc.CallOption) (GoLoadService_GetDownloadTaskFileClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &GoLoadService_ServiceDesc.Streams[0], GoLoadService_GetDownloadTaskFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &goLoadServiceGetDownloadTaskFileClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GoLoadService_GetDownloadTaskFileClient interface {
	Recv() (*GetDownloadTaskFileResponse, error)
	grpc.ClientStream
}

type goLoadServiceGetDownloadTaskFileClient struct {
	grpc.ClientStream
}

func (x *goLoadServiceGetDownloadTaskFileClient) Recv() (*GetDownloadTaskFileResponse, error) {
	m := new(GetDownloadTaskFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GoLoadServiceServer is the server API for GoLoadService service.
// All implementations must embed UnimplementedGoLoadServiceServer
// for forward compatibility
type GoLoadServiceServer interface {
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error)
	CreateSession(context.Context, *CreateSessionRequest) (*CreateSessionResponse, error)
	CreateDownloadTask(context.Context, *CreateDownloadTaskRequest) (*CreateDownloadTaskResponse, error)
	GetDownloadTaskList(context.Context, *GetDownloadTaskListRequest) (*GetDownloadTaskListResponse, error)
	UpdateDownloadTask(context.Context, *UpdateDownloadTaskRequest) (*UpdateDownloadTaskResponse, error)
	DeleteDownloadTask(context.Context, *DeleteDownloadTaskRequest) (*DeleteDownloadTaskResponse, error)
	GetDownloadTaskFile(*GetDownloadTaskFileRequest, GoLoadService_GetDownloadTaskFileServer) error
	mustEmbedUnimplementedGoLoadServiceServer()
}

// UnimplementedGoLoadServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGoLoadServiceServer struct {
}

func (UnimplementedGoLoadServiceServer) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (UnimplementedGoLoadServiceServer) CreateSession(context.Context, *CreateSessionRequest) (*CreateSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}
func (UnimplementedGoLoadServiceServer) CreateDownloadTask(context.Context, *CreateDownloadTaskRequest) (*CreateDownloadTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDownloadTask not implemented")
}
func (UnimplementedGoLoadServiceServer) GetDownloadTaskList(context.Context, *GetDownloadTaskListRequest) (*GetDownloadTaskListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDownloadTaskList not implemented")
}
func (UnimplementedGoLoadServiceServer) UpdateDownloadTask(context.Context, *UpdateDownloadTaskRequest) (*UpdateDownloadTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDownloadTask not implemented")
}
func (UnimplementedGoLoadServiceServer) DeleteDownloadTask(context.Context, *DeleteDownloadTaskRequest) (*DeleteDownloadTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDownloadTask not implemented")
}
func (UnimplementedGoLoadServiceServer) GetDownloadTaskFile(*GetDownloadTaskFileRequest, GoLoadService_GetDownloadTaskFileServer) error {
	return status.Errorf(codes.Unimplemented, "method GetDownloadTaskFile not implemented")
}
func (UnimplementedGoLoadServiceServer) mustEmbedUnimplementedGoLoadServiceServer() {}

// UnsafeGoLoadServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GoLoadServiceServer will
// result in compilation errors.
type UnsafeGoLoadServiceServer interface {
	mustEmbedUnimplementedGoLoadServiceServer()
}

func RegisterGoLoadServiceServer(s grpc.ServiceRegistrar, srv GoLoadServiceServer) {
	s.RegisterService(&GoLoadService_ServiceDesc, srv)
}

func _GoLoadService_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoLoadServiceServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoLoadService_CreateAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoLoadServiceServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoLoadService_CreateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoLoadServiceServer).CreateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoLoadService_CreateSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoLoadServiceServer).CreateSession(ctx, req.(*CreateSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoLoadService_CreateDownloadTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDownloadTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoLoadServiceServer).CreateDownloadTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoLoadService_CreateDownloadTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoLoadServiceServer).CreateDownloadTask(ctx, req.(*CreateDownloadTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoLoadService_GetDownloadTaskList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDownloadTaskListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoLoadServiceServer).GetDownloadTaskList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoLoadService_GetDownloadTaskList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoLoadServiceServer).GetDownloadTaskList(ctx, req.(*GetDownloadTaskListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoLoadService_UpdateDownloadTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDownloadTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoLoadServiceServer).UpdateDownloadTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoLoadService_UpdateDownloadTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoLoadServiceServer).UpdateDownloadTask(ctx, req.(*UpdateDownloadTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoLoadService_DeleteDownloadTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDownloadTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoLoadServiceServer).DeleteDownloadTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoLoadService_DeleteDownloadTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoLoadServiceServer).DeleteDownloadTask(ctx, req.(*DeleteDownloadTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoLoadService_GetDownloadTaskFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetDownloadTaskFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GoLoadServiceServer).GetDownloadTaskFile(m, &goLoadServiceGetDownloadTaskFileServer{ServerStream: stream})
}

type GoLoadService_GetDownloadTaskFileServer interface {
	Send(*GetDownloadTaskFileResponse) error
	grpc.ServerStream
}

type goLoadServiceGetDownloadTaskFileServer struct {
	grpc.ServerStream
}

func (x *goLoadServiceGetDownloadTaskFileServer) Send(m *GetDownloadTaskFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

// GoLoadService_ServiceDesc is the grpc.ServiceDesc for GoLoadService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GoLoadService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "go_load.v1.GoLoadService",
	HandlerType: (*GoLoadServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAccount",
			Handler:    _GoLoadService_CreateAccount_Handler,
		},
		{
			MethodName: "CreateSession",
			Handler:    _GoLoadService_CreateSession_Handler,
		},
		{
			MethodName: "CreateDownloadTask",
			Handler:    _GoLoadService_CreateDownloadTask_Handler,
		},
		{
			MethodName: "GetDownloadTaskList",
			Handler:    _GoLoadService_GetDownloadTaskList_Handler,
		},
		{
			MethodName: "UpdateDownloadTask",
			Handler:    _GoLoadService_UpdateDownloadTask_Handler,
		},
		{
			MethodName: "DeleteDownloadTask",
			Handler:    _GoLoadService_DeleteDownloadTask_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetDownloadTaskFile",
			Handler:       _GoLoadService_GetDownloadTaskFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "goload/v1/go_load.proto",
}
