package list_pvzs

import (
	"context"
	"errors"
	"testing"
	"time"

	"pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	pvz "pvz-service/internal/usecase/contract/repository/pvz/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListPVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currTime := time.Now()

	retPVZ1 := &entity.PVZ{
		Uuid:             uuid.New(),
		RegistrationDate: currTime,
		City:             "Москва",
	}

	retPVZ2 := &entity.PVZ{
		Uuid:             uuid.New(),
		RegistrationDate: currTime,
		City:             "Санкт-Петербург",
	}

	tests := []struct {
		name      string
		setupMock func(
			*pvz.MockRepositoryPVZ,
		)
		expected      []*entity.PVZ
		expectedError error
	}{
		{
			name: "successful list pvzs",
			setupMock: func(
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZList(gomock.Any()).
					Return([]*entity.PVZ{retPVZ1, retPVZ2}, nil)
			},
			expected: []*entity.PVZ{retPVZ1, retPVZ2},
		},
		{
			name: "empty list (ErrPVZNotFound)",
			setupMock: func(
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZList(gomock.Any()).
					Return(nil, repository.ErrPVZNotFound)
			},
		},
		{
			name: "unexpected error - some kinda 500",
			setupMock: func(
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZList(gomock.Any()).
					Return(nil, errors.New("some db error"))
			},
			expectedError: usecase2.ErrListPVZs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoPVZ := pvz.NewMockRepositoryPVZ(ctrl)
			tt.setupMock(mockRepoPVZ)

			u := NewUsecase(mockRepoPVZ)
			result, err := u.Run(context.Background())

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tt.expected != nil {
				require.NotNil(t, result)
				require.NotNil(t, result.PVZs)
				assert.Equal(t, len(tt.expected), len(result.PVZs))
				for i, pvz := range tt.expected {
					assert.Equal(t, pvz.Uuid, result.PVZs[i].Uuid)
					assert.Equal(t, pvz.RegistrationDate, result.PVZs[i].RegistrationDate)
					assert.Equal(t, pvz.City, result.PVZs[i].City)
				}
			}
		})
	}
}
