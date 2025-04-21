package reception

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

func TestStartReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzID := "6d132f66-dcfe-493e-965d-95c99e5f325d"
	currTime := time.Now()

	retReception := entity.Reception{
		Uuid:     "671353c3-d091-4de8-83f9-983fb6e34ecf",
		DateTime: currTime,
		PVZID:    pvzID,
		Status:   "opened",
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expected      *entity.Reception
		expectedError error
	}{
		{
			name: "successful StartReception",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					AddRow(
						retReception.Uuid,
						retReception.DateTime,
						retReception.PVZID,
						retReception.Status,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), pvzID).
					Return(rows, nil)
			},
			expected: &retReception,
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
			name: "no rows returned",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), pvzID).
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

			result, err := repo.StartReception(context.Background(), pvzID)

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

func TestCloseReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	recID := "1becb717-0ace-41e4-a711-37402f10cb51"
	currTime := time.Now()

	retReception := entity.Reception{
		Uuid:     recID,
		DateTime: currTime,
		PVZID:    "6d132f66-dcfe-493e-965d-95c99e5f325d",
		Status:   statusReceptionDone,
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expected      *entity.Reception
		expectedError error
	}{
		{
			name: "successful CloseReception",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					AddRow(
						retReception.Uuid,
						retReception.DateTime,
						retReception.PVZID,
						retReception.Status,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), statusReceptionDone, recID).
					Return(rows, nil)
			},
			expected: &retReception,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), statusReceptionDone, recID).
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
						pvzIDColumnName,
						statusColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), statusReceptionDone, recID).
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
					Query(gomock.Any(), gomock.Any(), statusReceptionDone, recID).
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

			result, err := repo.CloseReception(context.Background(), recID)

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

func TestGetLastReceptionPVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzID := "6d132f66-dcfe-493e-965d-95c99e5f325d"
	currTime := time.Now()

	retReception := entity.Reception{
		Uuid:     "1becb717-0ace-41e4-a711-37402f10cb51",
		DateTime: currTime,
		PVZID:    pvzID,
		Status:   statusReceptionDone,
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expected      *entity.Reception
		expectedError error
	}{
		{
			name: "successful GetLastReceptionPVZ",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					AddRow(
						retReception.Uuid,
						retReception.DateTime,
						retReception.PVZID,
						retReception.Status,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), pvzID).
					Return(rows, nil)
			},
			expected: &retReception,
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
			name: "no rows returned",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), pvzID).
					Return(rows, nil)
			},
			expectedError: repository2.ErrReceptionNotFound,
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

			result, err := repo.GetLastReceptionPVZ(context.Background(), pvzID)

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

func TestGetReceptionsByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	receptionID1 := "11111111-1111-1111-1111-111111111111"
	receptionID2 := "22222222-2222-2222-2222-222222222222"
	currTime := time.Now()

	retReception1 := entity.Reception{
		Uuid:     receptionID1,
		DateTime: currTime,
		PVZID:    "pvz-1",
		Status:   statusReceptionDone,
	}

	retReception2 := entity.Reception{
		Uuid:     receptionID2,
		DateTime: currTime,
		PVZID:    "pvz-2",
		Status:   "in_progress",
	}

	tests := []struct {
		name          string
		inputIDs      []string
		setupMock     func(*mockdb.MockDBContract)
		expected      *[]entity.Reception
		expectedError error
	}{
		{
			name:     "successful GetReceptionsByIDs",
			inputIDs: []string{receptionID1, receptionID2},
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					AddRow(retReception1.Uuid, retReception1.DateTime, retReception1.PVZID, retReception1.Status).
					AddRow(retReception2.Uuid, retReception2.DateTime, retReception2.PVZID, retReception2.Status).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rows, nil)
			},
			expected: &[]entity.Reception{retReception1, retReception2},
		},
		{
			name:      "empty input list returns empty result",
			inputIDs:  []string{},
			setupMock: func(mockDB *mockdb.MockDBContract) {},
			expected:  &[]entity.Reception{},
		},
		{
			name:     "query db error",
			inputIDs: []string{receptionID1},
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name:     "pgx.CollectRows - got some wrong columns data from db",
			inputIDs: []string{receptionID1},
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
		{
			name:     "no rows returned",
			inputIDs: []string{receptionID1},
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rows, nil)
			},
			expectedError: repository2.ErrReceptionsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mockdb.NewMockDBContract(ctrl)
			repo := &repository{db: mockDB}

			if tt.setupMock != nil {
				tt.setupMock(mockDB)
			}

			result, err := repo.GetReceptionsByIDs(context.Background(), tt.inputIDs)

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
