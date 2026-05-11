package main

import (
	pb "TestGrpc/proto"
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	HttpAdress string = "localhost:8080"
	GrpcAdress string = "localhost:50051"
)

func sumHandler(client2 pb.CalculatorServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. декодируем JSON body → {a, b}
		var body pb.AddRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body.A == 0 && body.B == 0 {
			http.Error(w, "a and b are required", http.StatusBadRequest)
			return
		}
		// 2. вызываем client2.Sum(...)
		res, err := client2.Sum(r.Context(), &body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 3. возвращаем JSON с результатом
		json.NewEncoder(w).Encode(res)
	}
}

func helloHandler(client pb.GreeterClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. декодируем JSON body → {a, b}
		var body pb.HelloRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body.Hello == "" {
			http.Error(w, "hello required", http.StatusBadRequest)
			return
		}
		// 2. вызываем client2.Sum(...)
		res, err := client.SayHello(r.Context(), &body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 3. возвращаем JSON с результатом
		json.NewEncoder(w).Encode(res)
	}
}

func main() {
	// 1. подключиться к серверу
	conn, err := grpc.Dial(GrpcAdress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// 2. не забыть закрыть соединение
	defer conn.Close()

	// Создание двух клиентов
	client := pb.NewGreeterClient(conn)
	client2 := pb.NewCalculatorServiceClient(conn)

	//  создать роутер
	mux := http.NewServeMux()

	//  зарегистрировать маршруты
	mux.HandleFunc("POST /sum", sumHandler(client2))
	mux.HandleFunc("POST /hello", helloHandler(client))

	//  запустить HTTP сервер
	http.ListenAndServe(HttpAdress, mux)
}
