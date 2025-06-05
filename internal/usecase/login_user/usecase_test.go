package login_user

import (
	"context"
	"errors"
	"testing"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	user "pvz-service/internal/usecase/contract/repository/user/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserAuthUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	correctPassword := "correctpassword"
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate password hash: %v", err)
	}

	reqData := In{
		Email:    "test@example.com",
		Password: correctPassword,
	}

	retUser := &entity.User{
		Uuid:         uuid.New(),
		Email:        reqData.Email,
		PasswordHash: string(hashedPass),
		Role:         "user",
	}
	out := &Out{User: *retUser}

	tests := []struct {
		name          string
		req           In
		setupMock     func(*user.MockRepositoryUser)
		expectedError error
		expectedUser  *Out
	}{
		{
			name: "successful user authentication",
			req:  reqData,
			setupMock: func(mockUser *user.MockRepositoryUser) {
				mockUser.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Return(retUser, nil)
			},
			expectedUser: out,
		},
		{
			name: "user not found",
			req:  reqData,
			setupMock: func(mockUser *user.MockRepositoryUser) {
				mockUser.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Return(nil, repository2.ErrUserNotFound)
			},
			expectedError: usecase2.ErrNotFoundUser,
		},
		{
			name: "db error when getting user",
			req:  reqData,
			setupMock: func(mockUser *user.MockRepositoryUser) {
				mockUser.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("some db error"))
			},
			expectedError: usecase2.ErrGetUser,
		},
		{
			name: "incorrect password",
			req: In{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(mockUser *user.MockRepositoryUser) {
				mockUser.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Return(retUser, nil)
			},
			expectedError: usecase2.ErrIncorrectPass,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoUser := user.NewMockRepositoryUser(ctrl)
			tt.setupMock(mockRepoUser)

			u := NewUsecase(mockRepoUser)
			result, err := u.Run(context.Background(), tt.req)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedUser, result)
			}
		})
	}
}
