package main

import (
	"context"
	"log"
	"time"

	pb "github.com/xiaogao688/zf-project/learn/grpc/go-grpc-middleware/prometheus/pb"
	"google.golang.org/grpc"
)

func main() {
	// 1. 建立到 gRPC 服务器的连接（假设监听在 localhost:50051）
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithInsecure(), // 若服务器未启用 TLS
		grpc.WithBlock(),    // 阻塞直到连接建立
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// 2. 创建客户端 stub
	client := pb.NewRpcServerClient(conn)

	for i := 0; i < 1000; i++ {
		// 3. 构造请求，并设置超时
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		req := &pb.Request{Msg: "Hello from Go client"}

		// 4. 发起 RPC 调用
		resp, err := client.SayHello(ctx, req)
		if err != nil {
			log.Fatalf("SayHello RPC failed: %v", err)
		}

		// 5. 输出结果
		log.Printf("Server replied: %q", resp.Msg)
	}
}
