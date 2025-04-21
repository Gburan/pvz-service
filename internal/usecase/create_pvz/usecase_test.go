package create_pvz

import (
	"context"
	"errors"
	"testing"
	"time"

	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	pvz "pvz-service/internal/usecase/contract/repository/pvz/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currTime := time.Now()

	reqData := In{
		City: "Москва",
	}

	retPVZ := entity.PVZ{
		Uuid:             "1becb717-0ace-41e4-a711-37401becb717",
		RegistrationDate: currTime,
		City:             reqData.City,
	}

	tests := []struct {
		name      string
		req       In
		setupMock func(
			*pvz.MockRepositoryPVZ,
		)
		expected      *entity.PVZ
		expectedError error
	}{
		{
			name: "successful create pvz",
			req:  reqData,
			setupMock: func(
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					SavePVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.City)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
						return &retPVZ, nil
					})
			},
			expected: &retPVZ,
		},
		{
			name: "unexpected error - some kinda 500",
			req:  reqData,
			setupMock: func(
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					SavePVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.City)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
						return nil, errors.New("some db error")
					})
			},
			expectedError: usecase2.ErrAddPVZ,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoPVZ := pvz.NewMockRepositoryPVZ(ctrl)
			tt.setupMock(mockRepoPVZ)

			u := NewUsecase(mockRepoPVZ)
			result, err := u.Run(context.Background(), tt.req)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tt.expected != nil {
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.Uuid, result.PVZ.Uuid)
				assert.Equal(t, tt.expected.RegistrationDate, result.PVZ.RegistrationDate)
				assert.Equal(t, tt.expected.City, result.PVZ.City)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
