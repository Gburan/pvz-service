package server

import (
	"context"
	"errors"

	pvz_v1 "pvz-service/internal/generated/api/v1/proto"
	usecase2 "pvz-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PVZServerImpl struct {
	pvz_v1.UnimplementedPVZServiceServer
	usecase usecase
}

func New(usecase usecase) *PVZServerImpl {
	return &PVZServerImpl{
		usecase: usecase,
	}
}

func (s *PVZServerImpl) GetPVZList(_ context.Context, _ *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	ctx := context.TODO()

	result, err := s.usecase.Run(ctx)
	if err != nil {
		if errors.Is(err, usecase2.ErrListPVZs) {
			return nil, status.Errorf(codes.Internal, "failed to list PVZs: %v", err)
		}
		return nil, status.Errorf(codes.Unknown, "unexpected error: %v", err)
	}

	out := make([]*pvz_v1.PVZ, 0, len(result.PVZs))
	for _, pvz := range result.PVZs {
		out = append(out, &pvz_v1.PVZ{
			Id:               pvz.Uuid.String(),
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
			City:             pvz.City,
		})
	}

	return &pvz_v1.GetPVZListResponse{Pvzs: out}, nil
}
