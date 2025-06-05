package product

import (
	"context"
	"errors"
	"testing"
	"time"

	repository2 "pvz-service/internal/infrastructure/repository"
	mockdb "pvz-service/internal/infrastructure/repository/mocks"
	"pvz-service/internal/model/entity"
	nower "pvz-service/internal/usecase/contract/nower/mocks"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prod := entity.Product{
		Uuid:        uuid.New(),
		DateTime:    time.Now(),
		Type:        "Электроника",
		ReceptionID: uuid.New(),
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expectedError error
	}{
		{
			name: "successful delete product",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Exec(gomock.Any(), gomock.Any(), prod.Uuid.String()).
					Return(pgconn.NewCommandTag("DELETE 1"), nil)
			},
			expectedError: nil,
		},
		{
			name: "successful delete product",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Exec(gomock.Any(), gomock.Any(), gomock.Eq(prod.Uuid.String())).
					Return(pgconn.NewCommandTag("DELETE 1"), nil)
			},
			expectedError: nil,
		},
		{
			name: "error executing query",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Exec(gomock.Any(), gomock.Any(), gomock.Eq(prod.Uuid.String())).
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
			mockNow := nower.NewMockNower(ctrl)
			repo := &repository{
				db:    mockDB,
				nower: mockNow,
			}

			tt.setupMock(mockDB)

			err := repo.DeleteProduct(context.Background(), entity.Product{
				Uuid: prod.Uuid,
			})

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestGetLastProductByReceptionPVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prod := entity.Product{
		Uuid:        uuid.New(),
		DateTime:    time.Now(),
		Type:        "Электроника",
		ReceptionID: uuid.New(),
	}

	tests := []struct {
		name          string
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      *entity.Product
		expectedError error
	}{
		{
			name: "successful GetLastProductByReceptionPVZ",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						typeColumnName,
						receptionIdColumnName,
					}).
					AddRow(
						prod.Uuid.String(),
						prod.DateTime,
						prod.Type,
						prod.ReceptionID,
					).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{prod.ReceptionID.String()},
					).
					Return(rows, nil)
			},
			expected: &prod,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{prod.ReceptionID.String()},
					).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "no product found",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						typeColumnName,
						receptionIdColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{prod.ReceptionID.String()},
					).
					Return(rows, nil)
			},
			expectedError: repository2.ErrProductNotFound,
		},
		{
			name: "pgx.CollectOneRow - got some wrong columns data from db",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{"some_unexected", "some_unexected_2"}).
					AddRow("unexp_data", "unexp_data_2").
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{prod.ReceptionID.String()},
					).
					Return(rows, nil)
			},
			expectedError: repository2.ErrScanResult,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mockdb.NewMockDBContract(ctrl)
			mockNow := nower.NewMockNower(ctrl)
			repo := &repository{
				db:    mockDB,
				nower: mockNow,
			}

			tt.setupMock(mockDB, mockNow)

			result, err := repo.GetLastProductByReceptionPVZ(context.Background(), entity.Product{
				ReceptionID: prod.ReceptionID},
			)

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

	prod := entity.Product{
		Uuid:        uuid.New(),
		DateTime:    time.Now(),
		Type:        "Электроника",
		ReceptionID: uuid.New(),
	}

	tests := []struct {
		name          string
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      *entity.Product
		expectedError error
	}{
		{
			name: "successful AddProduct",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				nower.EXPECT().
					Now().
					Return(prod.DateTime)

				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						typeColumnName,
						receptionIdColumnName,
					}).
					AddRow(
						prod.Uuid,
						prod.DateTime,
						prod.Type,
						prod.ReceptionID,
					).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{prod.Uuid.String(), prod.DateTime, prod.Type, prod.ReceptionID.String()},
					).
					Return(rows, nil)
			},
			expected: &prod,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				nower.EXPECT().
					Now().
					Return(prod.DateTime)

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{prod.Uuid.String(), prod.DateTime, prod.Type, prod.ReceptionID.String()},
					).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "no rows returned",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				nower.EXPECT().
					Now().
					Return(prod.DateTime)

				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						typeColumnName,
						receptionIdColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{prod.Uuid.String(), prod.DateTime, prod.Type, prod.ReceptionID.String()},
					).
					Return(rows, nil)
			},
			expectedError: repository2.ErrScanResult,
		},
		{
			name: "pgx.CollectOneRow - got some wrong columns data from db",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				nower.EXPECT().
					Now().
					Return(prod.DateTime)

				rows := pgxmock.
					NewRows([]string{
						"some_unexpected",
						"some_unexpected_2",
						"some_unexpected_3",
						"some_unexpected_4"}).
					AddRow(
						"unexp_data",
						"unexp_data_2",
						"unexp_data_3",
						"unexp_data_4").
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{prod.Uuid.String(), prod.DateTime, prod.Type, prod.ReceptionID.String()},
					).
					Return(rows, nil)
			},
			expectedError: repository2.ErrScanResult,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mockdb.NewMockDBContract(ctrl)
			mockNow := nower.NewMockNower(ctrl)
			repo := &repository{
				db:    mockDB,
				nower: mockNow,
			}

			tt.setupMock(mockDB, mockNow)

			result, err := repo.AddProduct(context.Background(), entity.Product{
				Uuid:        prod.Uuid,
				Type:        prod.Type,
				ReceptionID: prod.ReceptionID,
			})

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
		Uuid:        uuid.New(),
		DateTime:    startDate.Add(2 * time.Hour),
		Type:        "Электроника",
		ReceptionID: uuid.New(),
	}

	product2 := entity.Product{
		Uuid:        uuid.New(),
		DateTime:    startDate.Add(3 * time.Hour),
		Type:        "Одежда",
		ReceptionID: uuid.New(),
	}

	tests := []struct {
		name          string
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      *[]entity.Product
		expectedError error
	}{
		{
			name: "successful GetProductsByTimeRange",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
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
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "pgx.CollectRows returns unexpected scan error",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
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
			mockNow := nower.NewMockNower(ctrl)
			repo := &repository{
				db:    mockDB,
				nower: mockNow,
			}

			tt.setupMock(mockDB, mockNow)

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
