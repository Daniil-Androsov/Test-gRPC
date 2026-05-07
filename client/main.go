package main

import (
	// импорты

	pb "TestGrpc/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"log"
)

const (
	Address string = "localhost:50051"
)

func main() {
	// 1. подключиться к серверу
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// 2. не забыть закрыть соединение
	defer conn.Close()

	// 3. создать клиента (stub)
	client := pb.NewGreeterClient(conn)

	// 4. вызвать метод
	res, err := client.SayHello(context.Background(), &pb.HelloRequest{Hello: "Hello"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// 5. вывести результат
	log.Println(res)
}
