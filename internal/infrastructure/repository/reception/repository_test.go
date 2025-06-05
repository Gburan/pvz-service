package reception

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

func TestStartReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reception := entity.Reception{
		Uuid:     uuid.New(),
		DateTime: time.Now(),
		PVZID:    uuid.New(),
		Status:   statusReceptionProgress,
	}

	tests := []struct {
		name          string
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      *entity.Reception
		expectedError error
	}{
		{
			name: "successful StartReception",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				nower.EXPECT().
					Now().
					Return(reception.DateTime)

				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					AddRow(
						reception.Uuid,
						reception.DateTime,
						reception.PVZID,
						reception.Status,
					).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{reception.Uuid, reception.DateTime, reception.PVZID, reception.Status},
					).
					Return(rows, nil)
			},
			expected: &reception,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				nower.EXPECT().
					Now().
					Return(reception.DateTime)

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{reception.Uuid, reception.DateTime, reception.PVZID, reception.Status},
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
					Return(reception.DateTime)

				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{reception.Uuid, reception.DateTime, reception.PVZID, reception.Status},
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
					Return(reception.DateTime)

				rows := pgxmock.
					NewRows([]string{"some_unexpected", "some_unexpected_2"}).
					AddRow("unexp_data", "unexp_data_2").
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{reception.Uuid, reception.DateTime, reception.PVZID, reception.Status},
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

			result, err := repo.StartReception(context.Background(), entity.Reception{
				PVZID: reception.PVZID,
				Uuid:  reception.Uuid,
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

func TestCloseReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reception := entity.Reception{
		Uuid:     uuid.New(),
		DateTime: time.Now(),
		PVZID:    uuid.New(),
		Status:   statusReceptionDone,
	}

	tests := []struct {
		name          string
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      *entity.Reception
		expectedError error
	}{
		{
			name: "successful CloseReception",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					AddRow(
						reception.Uuid,
						reception.DateTime,
						reception.PVZID,
						reception.Status,
					).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						statusReceptionDone,
						[]interface{}{reception.Uuid.String()}).
					Return(rows, nil)
			},
			expected: &reception,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						statusReceptionDone,
						[]interface{}{reception.Uuid.String()}).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "no rows returned",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						statusReceptionDone,
						[]interface{}{reception.Uuid.String()}).
					Return(rows, nil)
			},
			expectedError: repository2.ErrScanResult,
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
						statusReceptionDone,
						[]interface{}{reception.Uuid.String()}).
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

			result, err := repo.CloseReception(context.Background(), entity.Reception{
				Uuid: reception.Uuid,
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

func TestGetLastReceptionPVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reception := entity.Reception{
		Uuid:     uuid.New(),
		DateTime: time.Now(),
		PVZID:    uuid.New(),
		Status:   statusReceptionDone,
	}

	tests := []struct {
		name          string
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      *entity.Reception
		expectedError error
	}{
		{
			name: "successful GetLastReceptionPVZ",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					AddRow(
						reception.Uuid,
						reception.DateTime,
						reception.PVZID,
						reception.Status,
					).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{reception.PVZID.String()}).
					Return(rows, nil)
			},
			expected: &reception,
		},
		{
			name: "query db error",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{reception.PVZID.String()}).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "no rows returned",
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{reception.PVZID.String()}).
					Return(rows, nil)
			},
			expectedError: repository2.ErrReceptionNotFound,
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
						[]interface{}{reception.PVZID.String()}).
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

			result, err := repo.GetLastReceptionPVZ(context.Background(), entity.Reception{
				PVZID: reception.PVZID,
				Uuid:  reception.Uuid,
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

func TestGetReceptionsByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reception1 := entity.Reception{
		Uuid:     uuid.New(),
		DateTime: time.Now(),
		PVZID:    uuid.New(),
		Status:   statusReceptionDone,
	}

	reception2 := entity.Reception{
		Uuid:     uuid.New(),
		DateTime: time.Now(),
		PVZID:    uuid.New(),
		Status:   "in_progress",
	}

	tests := []struct {
		name          string
		inputIDs      []uuid.UUID
		setupMock     func(db *mockdb.MockDBContract, nower *nower.MockNower)
		expected      *[]entity.Reception
		expectedError error
	}{
		{
			name:     "successful GetReceptionsByIDs",
			inputIDs: []uuid.UUID{reception1.Uuid, reception2.Uuid},
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						dateTimeColumnName,
						pvzIDColumnName,
						statusColumnName,
					}).
					AddRow(reception1.Uuid, reception1.DateTime, reception1.PVZID, reception1.Status).
					AddRow(reception2.Uuid, reception2.DateTime, reception2.PVZID, reception2.Status).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rows, nil)
			},
			expected: &[]entity.Reception{reception1, reception2},
		},
		{
			name:      "empty input list returns empty result",
			inputIDs:  []uuid.UUID{},
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {},
			expected:  &[]entity.Reception{},
		},
		{
			name:     "query db error",
			inputIDs: []uuid.UUID{reception1.Uuid},
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name:     "pgx.CollectRows - got some wrong columns data from db",
			inputIDs: []uuid.UUID{reception1.Uuid},
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
		{
			name:     "no rows returned",
			inputIDs: []uuid.UUID{reception1.Uuid},
			setupMock: func(mockDB *mockdb.MockDBContract, nower *nower.MockNower) {
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
			mockNow := nower.NewMockNower(ctrl)
			repo := &repository{
				db:    mockDB,
				nower: mockNow,
			}

			tt.setupMock(mockDB, mockNow)

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
