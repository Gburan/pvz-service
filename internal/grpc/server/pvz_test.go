package server

import (
	"context"
	"errors"
	"testing"
	"time"

	pvz_v1 "pvz-service/internal/generated/api/v1/proto"
	server "pvz-service/internal/grpc/server/mocks"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/list_pvzs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetPVZList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	retPVZ1 := &entity.PVZ{
		Uuid:             uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Москва",
	}
	retPVZ2 := &entity.PVZ{
		Uuid:             uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Санкт-Петербург",
	}

	usecaseOut := &list_pvzs.Out{
		PVZs: []*entity.PVZ{retPVZ1, retPVZ2},
	}

	expectedResp := &pvz_v1.GetPVZListResponse{
		Pvzs: []*pvz_v1.PVZ{
			{
				Id:               retPVZ1.Uuid.String(),
				RegistrationDate: timestamppb.New(retPVZ1.RegistrationDate),
				City:             retPVZ1.City,
			},
			{
				Id:               retPVZ2.Uuid.String(),
				RegistrationDate: timestamppb.New(retPVZ2.RegistrationDate),
				City:             retPVZ2.City,
			},
		},
	}

	tests := []struct {
		name          string
		setupMock     func(*server.Mockusecase)
		expected      *pvz_v1.GetPVZListResponse
		expectedError error
	}{
		{
			name: "successful list pvzs",
			setupMock: func(mockUsecase *server.Mockusecase) {
				mockUsecase.EXPECT().
					Run(context.TODO()).
					Return(usecaseOut, nil)
			},
			expected: expectedResp,
		},
		{
			name: "usecase error - failed to list pvzs",
			setupMock: func(mockUsecase *server.Mockusecase) {
				mockUsecase.EXPECT().
					Run(context.TODO()).
					Return(nil, usecase2.ErrListPVZs)
			},
			expectedError: status.Error(codes.Internal, "failed to list PVZs: "+usecase2.ErrListPVZs.Error()),
		},
		{
			name: "usecase error - unexpected error",
			setupMock: func(mockUsecase *server.Mockusecase) {
				mockUsecase.EXPECT().
					Run(context.TODO()).
					Return(nil, errors.New("some unknown error"))
			},
			expectedError: status.Error(codes.Unknown, "unexpected error: some unknown error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := server.NewMockusecase(ctrl)
			tt.setupMock(mockUsecase)

			server := New(mockUsecase)

			resp, err := server.GetPVZList(context.Background(), &pvz_v1.GetPVZListRequest{})

			if tt.expectedError != nil {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				expectedStatus, _ := status.FromError(tt.expectedError)
				assert.Equal(t, expectedStatus.Code(), st.Code())
				assert.Contains(t, st.Message(), expectedStatus.Message())
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, len(tt.expected.Pvzs), len(resp.Pvzs))

				for i, expectedPVZ := range tt.expected.Pvzs {
					actualPVZ := resp.Pvzs[i]
					assert.Equal(t, expectedPVZ.Id, actualPVZ.Id)
					assert.Equal(t, expectedPVZ.City, actualPVZ.City)
					assert.True(t, proto.Equal(expectedPVZ.RegistrationDate, actualPVZ.RegistrationDate))
				}
			}
		})
	}
}
