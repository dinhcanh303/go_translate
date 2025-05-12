package grpc_client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/dinhcanh303/go_translate/example/proto"
)

type LanguageDetectionService interface {
	DetectLanguage(ctx context.Context, text string) (*pb.DetectLanguageResponse, error)
}

// groupGRPCClient implements LanguageDetectionService
type languageDetectionGRPCClient struct {
	client pb.LanguageDetectionServiceClient
}

// DetectLanguage calls gRPC server
func (g *languageDetectionGRPCClient) DetectLanguage(ctx context.Context, text string) (*pb.DetectLanguageResponse, error) {
	return g.client.DetectLanguage(ctx, &pb.DetectLanguageRequest{
		Text: text,
	})
}

// NewGRPCLanguageDetectionClient creates a new gRPC client connection
func NewGRPCLanguageDetectionClient(grpcURL string) (LanguageDetectionService, error) {
	conn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewLanguageDetectionServiceClient(conn)

	return &languageDetectionGRPCClient{
		client: client,
	}, nil
}
