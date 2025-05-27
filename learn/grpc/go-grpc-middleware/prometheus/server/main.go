// main.go
package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	prom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/xiaogao688/zf-project/learn/grpc/go-grpc-middleware/prometheus/pb"
)

// 实现 Greeter 服务
type greeterServer struct {
	pb.UnimplementedRpcServerServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *pb.Request) (*pb.Reply, error) {
	return &pb.Reply{Msg: "Hello, " + req.Msg}, nil
}

func main() {
	// 1. 启动 TCP 监听
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. 创建 ServerMetrics，用于拦截器
	srvMetrics := prom.NewServerMetrics(
		// 如果想开启延迟直方图：
		prom.WithServerHandlingTimeHistogram(prom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 1, 5})),
	)

	// 3. 将 metrics 注册到默认的 Prometheus 注册表
	prometheus.MustRegister(srvMetrics) // :contentReference[oaicite:0]{index=0}

	// 4. 创建 gRPC Server，并添加 Unary/Stream 拦截器链
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			srvMetrics.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			srvMetrics.UnaryServerInterceptor(),
		)),
	)

	// 5. 注册业务服务和开启反射（可选）
	pb.RegisterRpcServerServer(grpcServer, &greeterServer{})
	reflection.Register(grpcServer)

	// 6. 在所有服务注册之后，初始化所有方法的 metrics（置零）
	srvMetrics.InitializeMetrics(grpcServer) // :contentReference[oaicite:1]{index=1}

	// 7. 启动一个 HTTP Server 暴露 /metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics listening on :9092/metrics")
		if err := http.ListenAndServe(":9092", nil); err != nil {
			log.Fatalf("failed to start http server: %v", err)
		}
	}()

	// 8. 启动 gRPC 服务
	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
