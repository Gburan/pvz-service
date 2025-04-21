package pvz

import (
	"context"
	"errors"
	"testing"
	"time"

	repository2 "pvz-service/internal/infrastructure/repository"
	mockdb "pvz-service/internal/infrastructure/repository/mocks"
	"pvz-service/internal/model/entity"

	"github.com/golang/mock/gomock"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSavePVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	city := "Москва"
	currTime := time.Now()

	retPVZ := entity.PVZ{
		Uuid:             "671353c3-d091-4de8-83f9-983fb6e34ecf",
		RegistrationDate: currTime,
		City:             city,
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expected      *entity.PVZ
		expectedError error
	}{
		{
			name: "successful SavePVZ",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						registrationDateColumnName,
						cityColumnName,
					}).
					AddRow(
						retPVZ.Uuid,
						retPVZ.RegistrationDate,
						retPVZ.City,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), city).
					Return(rows, nil)
			},
			expected: &retPVZ,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), city).
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
						registrationDateColumnName,
						cityColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), city).
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
					Query(gomock.Any(), gomock.Any(), city).
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

			result, err := repo.SavePVZ(context.Background(), city)

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

	pvzID := "6d132f66-dcfe-493e-965d-95c99e5f325d"
	currTime := time.Now()

	retPVZ := entity.PVZ{
		Uuid:             pvzID,
		RegistrationDate: currTime,
		City:             "Москва",
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expected      *entity.PVZ
		expectedError error
	}{
		{
			name: "successful GetPVZByID",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						registrationDateColumnName,
						cityColumnName,
					}).
					AddRow(
						retPVZ.Uuid,
						retPVZ.RegistrationDate,
						retPVZ.City,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), pvzID).
					Return(rows, nil)
			},
			expected: &retPVZ,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), pvzID).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "pvz not found",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						registrationDateColumnName,
						cityColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), pvzID).
					Return(rows, nil)
			},
			expectedError: repository2.ErrPVZNotFound,
		},
		{
			name: "pgx.CollectOneRow - got some wrong columns data from db",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{"some_unexpected", "some_unexpected_2"}).
					AddRow("unexp_data", "unexp_data_2").
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), pvzID).
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

			result, err := repo.GetPVZByID(context.Background(), pvzID)

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

	pvzID1 := "6d132f66-dcfe-493e-965d-95c99e5f325d"
	pvzID2 := "6d132f66-dcfe-493e-965d-95c99e5f965d"
	currTime := time.Now()

	retPVZ1 := entity.PVZ{
		Uuid:             pvzID1,
		RegistrationDate: currTime,
		City:             "Москва",
	}

	retPVZ2 := entity.PVZ{
		Uuid:             pvzID2,
		RegistrationDate: currTime,
		City:             "Санкт-Петербург",
	}

	tests := []struct {
		name          string
		inputIDs      []string
		setupMock     func(*mockdb.MockDBContract)
		expected      *[]entity.PVZ
		expectedError error
	}{
		{
			name:     "successful GetPVZsByIDs",
			inputIDs: []string{pvzID1, pvzID2},
			setupMock: func(mockDB *mockdb.MockDBContract) {
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
			expected: &[]entity.PVZ{retPVZ1, retPVZ2},
		},
		{
			name:      "empty input list returns empty result",
			inputIDs:  []string{},
			setupMock: func(mockDB *mockdb.MockDBContract) {},
			expected:  &[]entity.PVZ{},
		},
		{
			name:     "query db error",
			inputIDs: []string{pvzID1},
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name:     "pgx.CollectRows - wrong columns from db",
			inputIDs: []string{pvzID1},
			setupMock: func(mockDB *mockdb.MockDBContract) {
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
			repo := &repository{db: mockDB}

			if tt.setupMock != nil {
				tt.setupMock(mockDB)
			}

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
