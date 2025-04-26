package services

import (
	"context"
	pb "github.com/gslnkvmx/gox/proto/gen"
	"github.com/gslnkvmx/gox/server/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type FileService struct {
	pb.UnimplementedFileServiceServer
	minio *storage.MinIOClient
}

func RegisterFileService(server *grpc.Server, minio *storage.MinIOClient) {
	pb.RegisterFileServiceServer(server, &FileService{minio: minio})
}

func (s *FileService) SendFile(ctx context.Context, req *pb.SendFileRequest) (*pb.SendFileResponse, error) {
	// Saves file to MinIO. Saves it to a bucket with a name corresponding to the receiver's name.
	fileID, err := s.minio.SaveFile(ctx, req.ReceiverName, req.FileName, req.FileContent)
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		return nil, status.Error(codes.Internal, "File save failed")
	}

	log.Printf("File %s saved to %s with id: %s", req.FileName, req.ReceiverName, fileID)

	return &pb.SendFileResponse{FileId: fileID, Bucket: req.ReceiverName, Status: "OK"}, nil
}
