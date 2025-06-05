package user

import (
	"context"
	"errors"
	"testing"

	repository2 "pvz-service/internal/infrastructure/repository"
	mockdb "pvz-service/internal/infrastructure/repository/mocks"
	"pvz-service/internal/model/entity"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.User{
		Uuid:         uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "employee",
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expected      *entity.User
		expectedError error
	}{
		{
			name: "successful AddUser",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						emailColumnName,
						passHashColumnName,
						roleColumnName,
					}).
					AddRow(
						user.Uuid,
						user.Email,
						user.PasswordHash,
						user.Role,
					).
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{user.Uuid, user.Email, user.PasswordHash, user.Role},
					).
					Return(rows, nil)
			},
			expected: &user,
		},
		{
			name: "query execution error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{user.Uuid, user.Email, user.PasswordHash, user.Role},
					).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "scan result error - invalid columns",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{"unexpected_col"}).
					AddRow("unexpected_value").
					Kind()

				mockDB.EXPECT().
					Query(
						gomock.Any(),
						gomock.Any(),
						[]interface{}{user.Uuid, user.Email, user.PasswordHash, user.Role},
					).
					Return(rows, nil)
			},
			expectedError: repository2.ErrScanResult,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mockdb.NewMockDBContract(ctrl)
			repo := &repository{
				db: mockDB,
			}

			tt.setupMock(mockDB)

			result, err := repo.AddUser(context.Background(), entity.User{
				Uuid:         user.Uuid,
				Email:        user.Email,
				PasswordHash: user.PasswordHash,
				Role:         user.Role,
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

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.User{
		Uuid:         uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "employee",
	}

	tests := []struct {
		name          string
		setupMock     func(*mockdb.MockDBContract)
		expected      *entity.User
		expectedError error
	}{
		{
			name: "successful GetUserByEmail",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						emailColumnName,
						passHashColumnName,
						roleColumnName,
					}).
					AddRow(
						user.Uuid,
						user.Email,
						user.PasswordHash,
						user.Role,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), user.Email).
					Return(rows, nil)
			},
			expected: &user,
		},
		{
			name: "query execution error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), user.Email).
					Return(nil, errors.New("query error"))
			},
			expectedError: repository2.ErrExecuteQuery,
		},
		{
			name: "no user found (pgx.ErrNoRows)",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{
						idColumnName,
						emailColumnName,
						passHashColumnName,
						roleColumnName,
					}).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), user.Email).
					Return(rows, nil)
			},
			expectedError: repository2.ErrUserNotFound,
		},
		{
			name: "scan result error (wrong columns)",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				rows := pgxmock.
					NewRows([]string{"unexpected_col"}).
					AddRow("unexpected_val").
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), user.Email).
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

			result, err := repo.GetUserByEmail(context.Background(), entity.User{
				Email: user.Email,
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
