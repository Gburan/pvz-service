package delete_product

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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currTime := time.Now()

	reqData := In{
		PVZID: "6d132f66-dcfe-493e-965d-95c99e5f325d",
	}

	retReception := entity.Reception{
		Uuid:     "1becb717-0ace-41e4-a711-37402f10cb51",
		DateTime: currTime,
		PVZID:    reqData.PVZID,
		Status:   "in_progress",
	}
	retProduct := entity.Product{
		Uuid:        "671353c3-d091-4de8-83f9-983fb6e34ecf",
		DateTime:    currTime,
		Type:        "Электроника",
		ReceptionID: retReception.Uuid,
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
		expectedError error
	}{
		{
			name: "successful delete product",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.Reception, error) {
						return &retReception, nil
					})

				mockProduct.EXPECT().
					GetLastProductByReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, recId string) (*entity.Product, error) {
						return &retProduct, nil
					})

				mockProduct.EXPECT().
					DeleteProduct(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) error {
						return nil
					})
			},
		},
		{
			name: "db error happened product",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.Reception, error) {
						return &retReception, nil
					})

				mockProduct.EXPECT().
					GetLastProductByReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, recId string) (*entity.Product, error) {
						return &retProduct, nil
					})

				mockProduct.EXPECT().
					DeleteProduct(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) error {
						return errors.New("some db error")
					})
			},
			expectedError: usecase2.ErrDeleteProduct,
		},
		{
			name: "db error happened product",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.Reception, error) {
						return &retReception, nil
					})

				mockProduct.EXPECT().
					GetLastProductByReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, recId string) (*entity.Product, error) {
						return nil, errors.New("some db error")
					})
			},
			expectedError: usecase2.ErrGetProduct,
		},
		{
			name: "no product in active reception",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.Reception, error) {
						return &retReception, nil
					})

				mockProduct.EXPECT().
					GetLastProductByReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, recId string) (*entity.Product, error) {
						return nil, repository2.ErrProductNotFound
					})
			},
			expectedError: usecase2.ErrNotFoundProduct,
		},
		{
			name: "no receptions ever exist at this pvz",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.Reception, error) {
						return nil, repository2.ErrReceptionNotFound
					})
			},
			expectedError: usecase2.ErrNotFoundReception,
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
					GetPVZByID(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.Reception, error) {
						return nil, errors.New("some db error")
					})
			},
			expectedError: usecase2.ErrGetReception,
		},
		{
			name: "no active reception",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
						return &retPVZ, nil
					})

				mockReception.EXPECT().
					GetLastReceptionPVZ(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.Reception, error) {
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
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockPVZ.EXPECT().
					GetPVZByID(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
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
					GetPVZByID(gomock.Any(), gomock.AssignableToTypeOf(In{}.PVZID)).
					DoAndReturn(func(_ context.Context, pvzId string) (*entity.PVZ, error) {
						return nil, errors.New("some db error")
					})
			},
			expectedError: usecase2.ErrGetPVZByID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoProduct := product.NewMockRepositoryProduct(ctrl)
			mockRepoReception := reception.NewMockRepositoryReception(ctrl)
			mockRepoPVZ := pvz.NewMockRepositoryPVZ(ctrl)
			tt.setupMock(mockRepoProduct, mockRepoReception, mockRepoPVZ)

			u := NewUsecase(mockRepoPVZ, mockRepoReception, mockRepoProduct)
			err := u.Run(context.Background(), tt.req)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
