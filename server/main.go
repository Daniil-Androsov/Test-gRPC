package main

import (
	// импорты
	pb "TestGrpc/proto"
	"context"
	"log"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
)

const (
	Address string = "localhost:50051"
)

// структура которая реализует наш сервис
type Server struct {
	pb.UnimplementedGreeterServer
	// mu sync.RWMutex
}

// реализация метода SayHello из proto
func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// логика ответа
	return &pb.HelloResponse{Msg: req.Hello}, nil
}

func main() {
	// 1. создать listener на порту 50051
	lis, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 2. создать gRPC сервер
	srv := grpc.NewServer()
	// 3. зарегистрировать наш сервис
	pb.RegisterGreeterServer(srv, &Server{})
	// 4. запустить
	if err := srv.Serve(lis); err != nil { // srv.Serve(lis) именно тут мы начинаем слушать порт
		slog.Error("serve", "err", err)
		os.Exit(1)
	}
}
