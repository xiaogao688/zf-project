package main

import (
	"context"
	"fmt"
	"github.com/xiaogao688/zf-project/learn/grpc/interceptors/pb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server 需要实现 pb.RpcServerServer 接口
type server struct {
	pb.UnimplementedRpcServerServer
}

// SayHello 是一个简单的 Unary RPC
func (s *server) SayHello(ctx context.Context, req *pb.Result) (*pb.Reply, error) {
	msg := fmt.Sprintf("Hello, you sent url=%q and id=%d", req.GetUrl(), req.GetId())
	return &pb.Reply{Msg: msg}, nil
}

// ClientStream 接收客户端流，处理完再返回单一回复
func (s *server) ClientStream(stream pb.RpcServer_ClientStreamServer) error {
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// 客户端已发送完毕
			reply := &pb.Reply{Msg: fmt.Sprintf("Received %d messages", count)}
			return stream.SendAndClose(reply)
		}
		if err != nil {
			return err
		}
		log.Printf("ClientStream got: url=%q id=%d", req.GetUrl(), req.GetId())
		count++
	}
}

// ServerStream 给客户端返回一个流
func (s *server) ServerStream(req *pb.Result, stream pb.RpcServer_ServerStreamServer) error {
	// 简单地根据 id 连续回送多条
	for i := 0; i < int(req.GetId()); i++ {
		msg := fmt.Sprintf("Stream #%d for url=%q", i, req.GetUrl())
		if err := stream.Send(&pb.Reply{Msg: msg}); err != nil {
			return err
		}
	}
	return nil
}

// BidirectionalStream 双向流示例（对应 “双向流” 注释）
// proto 文件中若有 rpc Chat(stream Result) returns (stream Reply);
func (s *server) BidiStream(stream pb.RpcServer_BidiStreamServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// 回显：将收到的每个 Result 转为一条 Reply
		resp := &pb.Reply{
			Msg: fmt.Sprintf("Echo: url=%q, id=%d", in.GetUrl(), in.GetId()),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(orderUnaryServerInterceptor))
	pb.RegisterRpcServerServer(grpcServer, &server{})
	log.Printf("gRPC server listening on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func orderUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Pre-processing logic
	s := time.Now()

	// Invoking the handler to complete the normal execution of a unary RPC.
	m, err := handler(ctx, req)

	// Post processing logic
	log.Printf("Method: %s, req: %s, resp: %s, latency: %s\n",
		info.FullMethod, req, m, time.Now().Sub(s))

	return m, err
}

func StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Printf("Stream请求拦截器：Method: %s, IsServerStream: %t, IsClientStream: %t",
		info.FullMethod, info.IsServerStream, info.IsClientStream)

	// 可以在此处实现认证等逻辑
	err := handler(srv, ss)

	log.Printf("Stream请求结束")
	return err
}
