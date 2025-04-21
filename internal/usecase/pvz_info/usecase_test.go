package pvz_info

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

func TestPVZInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currTime := time.Now()

	reqData := In{
		StartData: currTime.Add(-24 * time.Hour),
		EndDate:   currTime,
		Page:      1,
		Limit:     1,
	}

	retPVZ := entity.PVZ{
		Uuid:             "8a4b3c2d-1e2f-3g4h-5i6j-7k8l9m0n1o2p",
		RegistrationDate: currTime.Add(-48 * time.Hour),
		City:             "Санкт-Петербург",
	}
	retReception := entity.Reception{
		Uuid:     "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
		DateTime: currTime.Add(-12 * time.Hour),
		PVZID:    retPVZ.Uuid,
		Status:   "completed",
	}
	retProducts := []entity.Product{
		{
			Uuid:        "p1a2b3c4-5d6e-7f8g-9h0i-j1k2l3m4n5o",
			DateTime:    currTime.Add(-6 * time.Hour),
			Type:        "одежда",
			ReceptionID: retReception.Uuid,
		},
		{
			Uuid:        "q1w2e3r4-5t6y-7u8i-9o0p-a1s2d3f4g5h",
			DateTime:    currTime.Add(-3 * time.Hour),
			Type:        "электроника",
			ReceptionID: retReception.Uuid,
		},
	}

	expectedSuccess := []Out{
		{
			PVZ: retPVZ,
			Receptions: []entity.ReceptionWithProducts{
				{
					Reception: retReception,
					Products:  retProducts,
				},
			},
		},
	}

	tests := []struct {
		name      string
		req       In
		setupMock func(
			*product.MockRepositoryProduct,
			*reception.MockRepositoryReception,
			*pvz.MockRepositoryPVZ,
		)
		expected      []Out
		expectedError error
	}{
		{
			name: "successful get pvz with receptions and products",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockProduct.EXPECT().
					GetProductsByTimeRange(gomock.Any(), reqData.StartData, reqData.EndDate).
					Return(&retProducts, nil)

				receptionIDs := []string{retReception.Uuid}
				mockReception.EXPECT().
					GetReceptionsByIDs(gomock.Any(), receptionIDs).
					Return(&[]entity.Reception{retReception}, nil)

				pvzIDs := []string{retPVZ.Uuid}
				mockPVZ.EXPECT().
					GetPVZsByIDs(gomock.Any(), pvzIDs).
					Return(&[]entity.PVZ{retPVZ}, nil)
			},
			expected: expectedSuccess,
		},
		{
			name: "page greater that count of pvzs",
			req: In{
				StartData: currTime.Add(-24 * time.Hour),
				EndDate:   currTime,
				Page:      2,
				Limit:     10,
			},
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockProduct.EXPECT().
					GetProductsByTimeRange(gomock.Any(), reqData.StartData, reqData.EndDate).
					Return(&retProducts, nil)

				receptionIDs := []string{retReception.Uuid}
				mockReception.EXPECT().
					GetReceptionsByIDs(gomock.Any(), receptionIDs).
					Return(&[]entity.Reception{retReception}, nil)

				pvzIDs := []string{retPVZ.Uuid}
				mockPVZ.EXPECT().
					GetPVZsByIDs(gomock.Any(), pvzIDs).
					Return(&[]entity.PVZ{retPVZ}, nil)
			},
			expectedError: usecase2.ErrNotPageTooBig,
		},
		{
			name: "limit is greater than all cnt of pvzs in db",
			req: In{
				StartData: currTime.Add(-24 * time.Hour),
				EndDate:   currTime,
				Page:      1,
				Limit:     10,
			},
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockProduct.EXPECT().
					GetProductsByTimeRange(gomock.Any(), reqData.StartData, reqData.EndDate).
					Return(&retProducts, nil)

				receptionIDs := []string{retReception.Uuid}
				mockReception.EXPECT().
					GetReceptionsByIDs(gomock.Any(), receptionIDs).
					Return(&[]entity.Reception{retReception}, nil)

				pvzIDs := []string{retPVZ.Uuid}
				mockPVZ.EXPECT().
					GetPVZsByIDs(gomock.Any(), pvzIDs).
					Return(&[]entity.PVZ{retPVZ}, nil)
			},
			expected: expectedSuccess,
		},
		{
			name: "db error happened products",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockProduct.EXPECT().
					GetProductsByTimeRange(gomock.Any(), reqData.StartData, reqData.EndDate).
					Return(nil, errors.New("some db error"))
			},
			expectedError: usecase2.ErrGetProducts,
		},
		{
			name: "no products found",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockProduct.EXPECT().
					GetProductsByTimeRange(gomock.Any(), reqData.StartData, reqData.EndDate).
					Return(nil, repository2.ErrProductsNotFound)
			},
			expectedError: usecase2.ErrNotFoundProducts,
		},
		{
			name: "error getting receptions",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockProduct.EXPECT().
					GetProductsByTimeRange(gomock.Any(), reqData.StartData, reqData.EndDate).
					Return(&retProducts, nil)

				mockReception.EXPECT().
					GetReceptionsByIDs(gomock.Any(), gomock.Any()).
					Return(nil, repository2.ErrReceptionNotFound)
			},
			expectedError: usecase2.ErrGetReceptions,
		},
		{
			name: "db error happened pvz",
			req:  reqData,
			setupMock: func(
				mockProduct *product.MockRepositoryProduct,
				mockReception *reception.MockRepositoryReception,
				mockPVZ *pvz.MockRepositoryPVZ,
			) {
				mockProduct.EXPECT().
					GetProductsByTimeRange(gomock.Any(), reqData.StartData, reqData.EndDate).
					Return(&retProducts, nil)

				receptionIDs := []string{retReception.Uuid}
				mockReception.EXPECT().
					GetReceptionsByIDs(gomock.Any(), receptionIDs).
					Return(&[]entity.Reception{retReception}, nil)

				pvzIDs := []string{retPVZ.Uuid}
				mockPVZ.EXPECT().
					GetPVZsByIDs(gomock.Any(), pvzIDs).
					Return(nil, errors.New("some db error"))
			},
			expectedError: usecase2.ErrGetPVZs,
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
				assert.ErrorContains(t, err, tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Len(t, result, len(tt.expected))

				for i, expectedPVZ := range tt.expected {
					assert.Equal(t, expectedPVZ.PVZ.Uuid, result[i].PVZ.Uuid)
					assert.Equal(t, expectedPVZ.PVZ.City, result[i].PVZ.City)
					assert.Equal(t, expectedPVZ.PVZ.RegistrationDate, result[i].PVZ.RegistrationDate)

					require.Len(t, result[i].Receptions, len(expectedPVZ.Receptions))
					for j, expectedReception := range expectedPVZ.Receptions {
						assert.Equal(t, expectedReception.Reception.Uuid, result[i].Receptions[j].Reception.Uuid)
						assert.Equal(t, expectedReception.Reception.Status, result[i].Receptions[j].Reception.Status)
						assert.Equal(t, expectedReception.Reception.DateTime, result[i].Receptions[j].Reception.DateTime)

						require.Len(t, result[i].Receptions[j].Products, len(expectedReception.Products))
						for k, expectedProduct := range expectedReception.Products {
							assert.Equal(t, expectedProduct.Uuid, result[i].Receptions[j].Products[k].Uuid)
							assert.Equal(t, expectedProduct.Type, result[i].Receptions[j].Products[k].Type)
							assert.Equal(t, expectedProduct.DateTime, result[i].Receptions[j].Products[k].DateTime)
						}
					}
				}
			}
		})
	}
}
