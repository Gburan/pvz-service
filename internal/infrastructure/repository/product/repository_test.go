package product

import (
	"context"
	"errors"
	"testing"
	"time"

	repository2 "pvz-service/internal/infrastructure/repository"
	mockdb "pvz-service/internal/infrastructure/repository/mocks"
	"pvz-service/internal/model/entity"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodId := "1becb717-0ace-41e4-a711-37402f10cb51"

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expectedError error
	}{
		{
			name: "successful delete product",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Exec(gomock.Any(), gomock.Any(), prodId).
					Return(pgconn.NewCommandTag("DELETE 1"), nil)
			},
			expectedError: nil,
		},
		{
			name: "successful delete product",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Exec(gomock.Any(), gomock.Any(), gomock.Eq(prodId)).
					Return(pgconn.NewCommandTag("DELETE 1"), nil)
			},
			expectedError: nil,
		},
		{
			name: "error executing query",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Exec(gomock.Any(), gomock.Any(), gomock.Eq(prodId)).
					Return(pgconn.NewCommandTag(""), &pgconn.PgError{
						Code: "2281337",
					})
			},
			expectedError: repository2.ErrExecuteQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mockdb.NewMockDBContract(ctrl)
			repo := &repository{db: mockDB}

			tt.setupMock(mockDB)

			err := repo.DeleteProduct(context.Background(), prodId)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestGetLastProductByReceptionPVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	recId := "1becb717-0ace-41e4-a711-37402f10cb51"
	currTime := time.Now()

	retProduct := entity.Product{
		Uuid:        "671353c3-d091-4de8-83f9-983fb6e34ecf",
		DateTime:    currTime,
		Type:        "Электроника",
		ReceptionID: recId,
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expected      *entity.Product
		expectedError error
	}{
		{
			name: "successful GetLastProductByReceptionPVZ",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						typeColumnName,
						receptionIdColumnName,
					}).
					AddRow(
						retProduct.Uuid,
						retProduct.DateTime,
						retProduct.Type,
						retProduct.ReceptionID,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), recId).
					Return(rows, nil)
			},
			expected: &retProduct,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), recId).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "no product found",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						typeColumnName,
						receptionIdColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), recId).
					Return(rows, nil)
			},
			expectedError: repository2.ErrProductNotFound,
		},
		{
			name: "pgx.CollectOneRow - got some wrong columns data from db",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{"some_unexected", "some_unexected_2"}).
					AddRow("unexp_data", "unexp_data_2").
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), recId).
					Return(rows, nil)
			},
			expectedError: repository2.ErrScanResult,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mockdb.NewMockDBContract(ctrl)
			repo := &repository{db: mockDB}

			tt.setupMock(mockDB)

			result, err := repo.GetLastProductByReceptionPVZ(context.Background(), recId)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expected != nil {
				require.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	receptionID := "1becb717-0ace-41e4-a711-37402f10cb51"
	productType := "Электроника"
	currTime := time.Now()

	retProduct := entity.Product{
		Uuid:        "671353c3-d091-4de8-83f9-983fb6e34ecf",
		DateTime:    currTime,
		Type:        productType,
		ReceptionID: receptionID,
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expected      *entity.Product
		expectedError error
	}{
		{
			name: "successful AddProduct",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						typeColumnName,
						receptionIdColumnName,
					}).
					AddRow(
						retProduct.Uuid,
						retProduct.DateTime,
						retProduct.Type,
						retProduct.ReceptionID,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), productType, receptionID).
					Return(rows, nil)
			},
			expected: &retProduct,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), productType, receptionID).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "no rows returned",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						typeColumnName,
						receptionIdColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), productType, receptionID).
					Return(rows, nil)
			},
			expectedError: repository2.ErrScanResult,
		},
		{
			name: "pgx.CollectOneRow - got some wrong columns data from db",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{"some_unexpected", "some_unexpected_2"}).
					AddRow("unexp_data", "unexp_data_2").
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), productType, receptionID).
					Return(rows, nil)
			},
			expectedError: repository2.ErrScanResult,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mockdb.NewMockDBContract(ctrl)
			repo := &repository{db: mockDB}

			tt.setupMock(mockDB)

			result, err := repo.AddProduct(context.Background(), receptionID, productType)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expected != nil {
				require.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestGetProductsByTimeRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	product1 := entity.Product{
		Uuid:        "prod-1",
		DateTime:    startDate.Add(2 * time.Hour),
		Type:        "Электроника",
		ReceptionID: "rec-1",
	}

	product2 := entity.Product{
		Uuid:        "prod-2",
		DateTime:    startDate.Add(3 * time.Hour),
		Type:        "Одежда",
		ReceptionID: "rec-2",
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expected      *[]entity.Product
		expectedError error
	}{
		{
			name: "successful GetProductsByTimeRange",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						typeColumnName,
						receptionIdColumnName,
					}).
					AddRow(product1.Uuid, product1.DateTime, product1.Type, product1.ReceptionID).
					AddRow(product2.Uuid, product2.DateTime, product2.Type, product2.ReceptionID).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rows, nil)
			},
			expected: &[]entity.Product{product1, product2},
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "pgx.CollectRows returns unexpected scan error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{"some", "wrong"}).
					AddRow("unexpected", "columns").
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rows, nil)
			},
			expectedError: repository2.ErrScanResult,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mockdb.NewMockDBContract(ctrl)
			repo := &repository{db: mockDB}

			tt.setupMock(mockDB)

			result, err := repo.GetProductsByTimeRange(context.Background(), startDate, endDate)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expected != nil {
				require.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
