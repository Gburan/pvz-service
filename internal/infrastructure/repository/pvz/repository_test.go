package pvz

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
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSavePVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvz := entity.PVZ{
		Uuid:             uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Москва",
	}

	tests := []struct {
		name          string
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      *entity.PVZ
		expectedError error
	}{
		{
			name: "successful SavePVZ",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				nower.EXPECT().
					Now().
					Return(pvz.RegistrationDate)

				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						registrationDateColumnName,
						cityColumnName,
					}).
					AddRow(
						pvz.Uuid,
						pvz.RegistrationDate,
						pvz.City,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rows, nil)
			},
			expected: &pvz,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				nower.EXPECT().
					Now().
					Return(pvz.RegistrationDate)

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "no rows returned",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				nower.EXPECT().
					Now().
					Return(pvz.RegistrationDate)

				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						registrationDateColumnName,
						cityColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rows, nil)
			},
			expectedError: repository2.ErrScanResult,
		},
		{
			name: "pgx.CollectOneRow - got some wrong columns data from db",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				nower.EXPECT().
					Now().
					Return(pvz.RegistrationDate)

				rows := pgxmock.
					NewRows([]string{"some_unexpected", "some_unexpected_2"}).
					AddRow("unexp_data", "unexp_data_2").
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
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

			result, err := repo.SavePVZ(context.Background(), entity.PVZ{
				City: pvz.City,
				Uuid: pvz.Uuid,
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

func TestGetPVZByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvz := entity.PVZ{
		Uuid:             uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Москва",
	}

	tests := []struct {
		name          string
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      *entity.PVZ
		expectedError error
	}{
		{
			name: "successful GetPVZByID",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						registrationDateColumnName,
						cityColumnName,
					}).
					AddRow(
						pvz.Uuid,
						pvz.RegistrationDate,
						pvz.City,
					).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{pvz.Uuid.String()},
					).
					Return(rows, nil)
			},
			expected: &pvz,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{pvz.Uuid.String()},
					).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "pvz not found",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						registrationDateColumnName,
						cityColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{pvz.Uuid.String()},
					).
					Return(rows, nil)
			},
			expectedError: repository2.ErrPVZNotFound,
		},
		{
			name: "pgx.CollectOneRow - got some wrong columns data from db",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{"some_unexpected", "some_unexpected_2"}).
					AddRow("unexp_data", "unexp_data_2").
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{pvz.Uuid.String()},
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

			result, err := repo.GetPVZByID(context.Background(), entity.PVZ{
				Uuid: pvz.Uuid,
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

func TestGetPVZsByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvz1 := entity.PVZ{
		Uuid:             uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Москва",
	}

	pvz2 := entity.PVZ{
		Uuid:             uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Санкт-Петербург",
	}

	tests := []struct {
		name          string
		inputIDs      []uuid.UUID
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      *[]entity.PVZ
		expectedError error
	}{
		{
			name:     "successful GetPVZsByIDs",
			inputIDs: []uuid.UUID{pvz1.Uuid, pvz2.Uuid},
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						registrationDateColumnName,
						cityColumnName,
					}).
					AddRow(pvz1.Uuid, pvz1.RegistrationDate, pvz1.City).
					AddRow(pvz2.Uuid, pvz2.RegistrationDate, pvz2.City).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rows, nil)
			},
			expected: &[]entity.PVZ{pvz1, pvz2},
		},
		{
			name:      "empty input list returns empty result",
			inputIDs:  []uuid.UUID{},
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {},
			expected:  &[]entity.PVZ{},
		},
		{
			name:     "query db error",
			inputIDs: []uuid.UUID{pvz1.Uuid},
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name:     "pgx.CollectRows - wrong columns from db",
			inputIDs: []uuid.UUID{pvz1.Uuid},
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{"wrong_col1", "wrong_col2"}).
					AddRow("wrong1", "wrong2").
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
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

			result, err := repo.GetPVZsByIDs(context.Background(), tt.inputIDs)

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

	tests := []struct {
		name          string
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      []*entity.PVZ
		expectedError error
	}{
		{
			name: "successful GetPVZList",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						registrationDateColumnName,
						cityColumnName,
					}).
					AddRow(retPVZ1.Uuid, retPVZ1.RegistrationDate, retPVZ1.City).
					AddRow(retPVZ2.Uuid, retPVZ2.RegistrationDate, retPVZ2.City).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rows, nil)
			},
			expected: []*entity.PVZ{retPVZ1, retPVZ2},
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "pgx.CollectRows - wrong columns from db",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{"wrong_col1", "wrong_col2"}).
					AddRow("wrong1", "wrong2").
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
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

			result, err := repo.GetPVZList(context.Background())

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expected != nil {
				require.NotNil(t, result)
				assert.Equal(t, tt.expected, result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
