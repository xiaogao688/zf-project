package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"

	"github.com/xiaogao688/zf-project/learn/grpc/interceptors/pb"
	"google.golang.org/grpc"
)

func main() {
	// 建立连接
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// … 其它 DialOption
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewRpcServerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Unary 调用 SayHello
	unaryReq := &pb.Result{Url: "https://example.com", Id: 1}
	unaryResp, err := client.SayHello(ctx, unaryReq)
	if err != nil {
		log.Fatalf("SayHello failed: %v", err)
	}
	log.Printf("Unary response: %s", unaryResp.GetMsg())

	// 客户端流 ClientStream
	cs, err := client.ClientStream(ctx)
	if err != nil {
		log.Fatalf("ClientStream error: %v", err)
	}
	for i := 0; i < 5; i++ {
		if err := cs.Send(&pb.Result{Url: "client_stream", Id: int32(i)}); err != nil {
			log.Fatalf("send error: %v", err)
		}
	}
	csReply, err := cs.CloseAndRecv()
	if err != nil {
		log.Fatalf("ClientStream CloseAndRecv error: %v", err)
	}
	log.Printf("ClientStream reply: %s", csReply.GetMsg())

	// 服务端流 ServerStream
	ss, err := client.ServerStream(ctx, &pb.Result{Url: "server_stream", Id: 3})
	if err != nil {
		log.Fatalf("ServerStream error: %v", err)
	}
	for {
		r, err := ss.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("recv error: %v", err)
		}
		log.Printf("ServerStream got: %s", r.GetMsg())
	}

	// 双向流 Chat（如 proto 中定义）
	bidi, err := client.BidiStream(ctx)
	if err != nil {
		log.Fatalf("Chat error: %v", err)
	}
	// 并发发送与接收
	waitChen := make(chan struct{})
	go func() {
		for i := 0; i < 3; i++ {
			if err := bidi.Send(&pb.Result{Url: "bidi", Id: int32(i)}); err != nil {
				log.Fatalf("Chat send error: %v", err)
			}
		}
		bidi.CloseSend()
	}()
	go func() {
		for {
			in, err := bidi.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Chat recv error: %v", err)
			}
			log.Printf("Chat got: %s", in.GetMsg())
		}
		close(waitChen)
	}()
	<-waitChen
}
