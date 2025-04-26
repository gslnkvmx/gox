package main

import (
	"google.golang.org/grpc"
	"log"
	"net"

	"github.com/gslnkvmx/gox/server/services"
	"github.com/gslnkvmx/gox/server/storage"
)

func main() {
	// Инициализация MinIO
	minioClient, err := storage.NewMinIOClient(
		"localhost:9000",
		"3y9OcJ4CSAdsjiXpYZni",
		"HuWfrpEkWESSJAStkwCBDTU9utkf7iqg1VpT397U",
		false, // useSSL
	)
	if err != nil {
		log.Fatalf("Failed to init MinIO: %v", err)
	}

	gRPCServer := grpc.NewServer()
	services.RegisterFileService(gRPCServer, minioClient)

	// Запуск gRPC-сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Server started on :50051")
	err = gRPCServer.Serve(lis)
	if err != nil {
		return
	}
}
