package main

import (
	"context"
	"fmt"
	proto "github.com/grpc-ecosystem/go-grpc-prometheus/examples/grpc-server-with-prometheus/protobuf"
	"github.com/xiaogao688/zf-project/learn/ego/grpc/pb"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egrpc"
	"google.golang.org/grpc"
)

type DemoService struct {
	pb.UnimplementedDemoServiceServer
}

func (s *DemoService) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	elog.Info("Received request", elog.String("name", req.Name))
	return &proto.HelloResponse{
		Message: fmt.Sprintf("Hello %s from Ego gRPC Server!", req.Name),
	}, nil
}

func main() {
	if err := ego.New().
		Serve(func() *egrpc.Component {
			// 创建 gRPC 服务器组件
			server := egrpc.Load("server.grpc").Build()

			// 注册服务
			proto.RegisterDemoServiceServer(server.Server, &DemoService{})
			
			// 添加自定义拦截器
			server.Use(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				elog.Info("Custom server interceptor", elog.String("method", info.FullMethod))
				return handler(ctx, req)
			})

			return server
		}()).
		Run(); err != nil {
		elog.Panic("startup failed", elog.FieldErr(err))
	}
}
