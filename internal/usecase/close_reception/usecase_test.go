package close_reception

import (
	"context"
	"errors"
	"testing"
	"time"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	pvz "pvz-service/internal/usecase/contract/repository/pvz/mocks"
	reception "pvz-service/internal/usecase/contract/repository/reception/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCloseReception(t *testing.T) {
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
		Status:   "in_progress",
	}
	retReceptionClosed := entity.Reception{
		Uuid:     uuid.New(),
		DateTime: currTime,
		PVZID:    reqData.PVZID,
		Status:   "in_progress",
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
			*reception.MockRepositoryReception,
			*pvz.MockRepositoryPVZ,
		)
		expected      *entity.Reception
		expectedError error
	}{
		{
			name: "successful close reception",
			req:  reqData,
			setupMock: func(
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
					CloseReception(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, rec entity.Reception) (*entity.Reception, error) {
						return &retReceptionClosed, nil
					})
			},
			expected: &retReceptionClosed,
		},
		{
			name: "db error happened reception",
			req:  reqData,
			setupMock: func(
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
					CloseReception(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, rec entity.Reception) (*entity.Reception, error) {
						return nil, errors.New("some db error")
					})
			},
			expectedError: usecase2.ErrCloseReception,
		},
		{
			name: "db error happened reception",
			req:  reqData,
			setupMock: func(
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
					CloseReception(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, rec entity.Reception) (*entity.Reception, error) {
						return nil, errors.New("some db error")
					})
			},
			expectedError: usecase2.ErrCloseReception,
		},
		{
			name: "no receptions ever exist at this pvz",
			req:  reqData,
			setupMock: func(
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
						return nil, repository2.ErrReceptionNotFound
					})
			},
			expectedError: usecase2.ErrNotFoundReception,
		},
		{
			name: "db error happened reception",
			req:  reqData,
			setupMock: func(
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
			expectedError: usecase2.ErrGetReception,
		},
		{
			name: "no active reception",
			req:  reqData,
			setupMock: func(
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
							Status: "close",
						}, nil
					})
			},
			expectedError: usecase2.ErrNotFoundOpenedReception,
		},
		{
			name: "requested pvz to update not exist",
			req:  reqData,
			setupMock: func(
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoReception := reception.NewMockRepositoryReception(ctrl)
			mockRepoPVZ := pvz.NewMockRepositoryPVZ(ctrl)
			tt.setupMock(mockRepoReception, mockRepoPVZ)

			u := NewUsecase(mockRepoPVZ, mockRepoReception)
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
				assert.Equal(t, tt.expected.PVZID, result.Reception.PVZID)
				assert.Equal(t, tt.expected.Status, result.Reception.Status)
				assert.Equal(t, tt.expected.DateTime, result.Reception.DateTime)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
