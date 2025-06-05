package start_reception

import (
	"context"
	"errors"
	"testing"
	"time"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	product "pvz-service/internal/usecase/contract/repository/product/mocks"
	pvz "pvz-service/internal/usecase/contract/repository/pvz/mocks"
	reception "pvz-service/internal/usecase/contract/repository/reception/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStartReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currTime := time.Now()

	reqData := In{
		PVZID: uuid.New(),
	}

	retReception := entity.Reception{
		Uuid:     uuid.New(),
		DateTime: currTime,
		PVZID:    reqData.PVZID,
		Status:   "close",
	}
	retPVZ := entity.PVZ{
		Uuid:             reqData.PVZID,
		RegistrationDate: currTime,
		City:             "Москва",
	}

	tests := []struct {
		name      string
		req       In
		setupMock func(
			*product.MockRepositoryProduct,
			*reception.MockRepositoryReception,
			*pvz.MockRepositoryPVZ,
		)
		expected      *entity.Reception
		expectedError error
	}{
		{
			name: "successful start reception",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pvz entity.PVZ) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, rec entity.Reception) (*entity.Reception, error) {
						return &retReception, nil
					})

				mockReception.EXPECT().
					StartReception(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, rec entity.Reception) (*entity.Reception, error) {
						return &retReception, nil
					})
			},
			expected: &retReception,
		},
		{
			name: "requested pvz to update not exist",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pvz entity.PVZ) (*entity.PVZ, error) {
						return nil, repository2.ErrPVZNotFound
					})
			},
			expectedError: usecase2.ErrNotFoundPVZ,
		},
		{
			name: "db error happened pvz",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pvz entity.PVZ) (*entity.PVZ, error) {
						return nil, errors.New("some db error")
					})
			},
			expectedError: usecase2.ErrGetPVZByID,
		},
		{
			name: "db error happened reception",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pvz entity.PVZ) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, rec entity.Reception) (*entity.Reception, error) {
						return nil, errors.New("some db error")
					})
			},
			expectedError: usecase2.ErrNotFoundReception,
		},
		{
			name: "active reception already exist",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pvz entity.PVZ) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, rec entity.Reception) (*entity.Reception, error) {
						return &entity.Reception{
							Status: "in_progress",
						}, nil
					})
			},
			expectedError: usecase2.ErrFoundOpenedReception,
		},
		{
			name: "db error happened reception",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pvz entity.PVZ) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, rec entity.Reception) (*entity.Reception, error) {
						return &retReception, nil
					})

				mockReception.EXPECT().
					StartReception(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, rec entity.Reception) (*entity.Reception, error) {
						return nil, errors.New("some db error")
					})
			},
			expectedError: usecase2.ErrStartReception,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoProduct := product.NewMockRepositoryProduct(ctrl)
			mockRepoReception := reception.NewMockRepositoryReception(ctrl)
			mockRepoPVZ := pvz.NewMockRepositoryPVZ(ctrl)
			tt.setupMock(mockRepoProduct, mockRepoReception, mockRepoPVZ)

			u := NewUsecase(mockRepoPVZ, mockRepoReception, mockRepoProduct)
			result, err := u.Run(context.Background(), tt.req)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tt.expected != nil {
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.Uuid, result.Reception.Uuid)
				assert.Equal(t, tt.expected.DateTime, result.Reception.DateTime)
				assert.Equal(t, tt.expected.PVZID, result.Reception.PVZID)
				assert.Equal(t, tt.expected.Status, result.Reception.Status)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
