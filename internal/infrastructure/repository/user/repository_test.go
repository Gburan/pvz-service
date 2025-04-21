package user

import (
	"context"
	"errors"
	"testing"

	repository2 "pvz-service/internal/infrastructure/repository"
	mockdb "pvz-service/internal/infrastructure/repository/mocks"
	"pvz-service/internal/model/entity"

	"github.com/golang/mock/gomock"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	email := "test@example.com"
	passHash := "hashedpassword"
	role := "employee"
	userID := "123e4567-e89b-12d3-a456-426614174000"

	expectedUser := entity.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passHash,
		Role:         role,
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
						expectedUser.ID,
						expectedUser.Email,
						expectedUser.PasswordHash,
						expectedUser.Role,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rows, nil)
			},
			expected: &expectedUser,
		},
		{
			name: "query execution error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
					Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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

			result, err := repo.AddUser(context.Background(), email, passHash, role)

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

	email := "test@example.com"
	userID := "123e4567-e89b-12d3-a456-426614174000"
	passHash := "hashedpassword"
	role := "employee"

	expectedUser := entity.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passHash,
		Role:         role,
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
						expectedUser.ID,
						expectedUser.Email,
						expectedUser.PasswordHash,
						expectedUser.Role,
					).
					Kind()

				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), email).
					Return(rows, nil)
			},
			expected: &expectedUser,
		},
		{
			name: "query execution error",
			setupMock: func(mockDB *mockdb.MockDBContract) {
				mockDB.EXPECT().
					Query(gomock.Any(), gomock.Any(), email).
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
					Query(gomock.Any(), gomock.Any(), email).
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
					Query(gomock.Any(), gomock.Any(), email).
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

			result, err := repo.GetUserByEmail(context.Background(), email)

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
