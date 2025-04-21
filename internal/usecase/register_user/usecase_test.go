package register_user

import (
	"context"
	"errors"
	"testing"

	repository2 "pvz-service/internal/infrastructure/repository"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	user "pvz-service/internal/usecase/contract/repository/user/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqData := In{
		Email:    "test@example.com",
		Password: "securepassword",
		Role:     "user",
	}

	u := entity.User{
		ID:           "1becb717-0ace-41e4-a711-37402f10cb51",
		Email:        reqData.Email,
		PasswordHash: "hashedpassword",
		Role:         reqData.Role,
	}
	retUser := &Out{
		User: u,
	}

	tests := []struct {
		name          string
		req           In
		setupMock     func(*user.MockRepositoryUser)
		expectedError error
		expectedUser  *Out
	}{
		{
			name: "successful user creation",
			req:  reqData,
			setupMock: func(mockUser *user.MockRepositoryUser) {
				mockUser.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.AssignableToTypeOf(reqData.Email)).
					Return(nil, repository2.ErrUserNotFound)

				mockUser.EXPECT().
					AddUser(gomock.Any(), reqData.Email, gomock.Any(), reqData.Role).
					Return(&u, nil)
			},
			expectedUser: retUser,
		},
		{
			name: "user already exists",
			req:  reqData,
			setupMock: func(mockUser *user.MockRepositoryUser) {
				mockUser.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.AssignableToTypeOf(reqData.Email)).
					Return(&u, nil)
			},
			expectedError: usecase2.ErrUserAlreadyExist,
		},
		{
			name: "db error when checking user existence",
			req:  reqData,
			setupMock: func(mockUser *user.MockRepositoryUser) {
				mockUser.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.AssignableToTypeOf(reqData.Email)).
					Return(nil, errors.New("some db error"))
			},
			expectedError: usecase2.ErrGetUser,
		},
		{
			name: "error adding user to db",
			req:  reqData,
			setupMock: func(mockUser *user.MockRepositoryUser) {
				mockUser.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.AssignableToTypeOf(reqData.Email)).
					Return(nil, repository2.ErrUserNotFound)

				mockUser.EXPECT().
					AddUser(gomock.Any(), reqData.Email, gomock.Any(), reqData.Role).
					Return(nil, errors.New("some db error"))
			},
			expectedError: usecase2.ErrAddUser,
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
