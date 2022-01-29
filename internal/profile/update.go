package profile

import (
	"context"
	"log"

	"github.com/firstcontributions/backend/internal/profile/models/mongo"
	"github.com/firstcontributions/backend/internal/profile/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) UpdateProfile(ctx context.Context, req *proto.Profile) (*proto.Profile, error) {
	if err := mongo.UpdateProfile(ctx, s.MongoClient, req); err != nil {
		log.Printf("error on updating profile [%v] %v", req.Uuid, err)
		return nil, status.Errorf(codes.Internal, "error on creating profile  %w", err)
	}
	return req, nil
}
