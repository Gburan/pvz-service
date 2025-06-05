package login_user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	login_user2 "pvz-service/internal/handler/login_user/mocks"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/login_user"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	valid := validator.New(validator.WithRequiredStructEnabled())
	secret := "test_secret"

	validEmail := "test@example.com"
	validPassword := "valid_password"
	usecaseIn := login_user.In{
		Email:    validEmail,
		Password: validPassword,
	}

	retUser := login_user.Out{
		User: entity.User{
			Uuid:         uuid.New(),
			Email:        validEmail,
			PasswordHash: "asdfasdgdsghsghawefasdgfjgdfhsdfgsdfg",
			Role:         "employee",
		},
	}

	successResponse := dto.LoginUserOut{
		Token: "generated_token",
	}

	tests := []struct {
		name          string
		setupMock     func(*login_user2.Mockusecase)
		requestBody   interface{}
		expectedCode  int
		expected      *dto.LoginUserOut
		expectedError map[string]string
	}{
		{
			name: "successful login",
			setupMock: func(mockUsecase *login_user2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(&retUser, nil)
			},
			requestBody: dto.LoginUserIn{
				Email:    validEmail,
				Password: validPassword,
			},
			expectedCode: http.StatusOK,
			expected:     &successResponse,
		},
		{
			name:      "validation failed - empty email",
			setupMock: func(mockUsecase *login_user2.Mockusecase) {},
			requestBody: dto.LoginUserIn{
				Email:    "",
				Password: validPassword,
			},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name:      "validation failed - invalid email format",
			setupMock: func(mockUsecase *login_user2.Mockusecase) {},
			requestBody: dto.LoginUserIn{
				Email:    "invalid_email",
				Password: validPassword,
			},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name:      "validation failed - empty password",
			setupMock: func(mockUsecase *login_user2.Mockusecase) {},
			requestBody: dto.LoginUserIn{
				Email:    validEmail,
				Password: "",
			},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name: "usecase error - user not found",
			setupMock: func(mockUsecase *login_user2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrNotFoundUser)
			},
			requestBody: dto.LoginUserIn{
				Email:    validEmail,
				Password: validPassword,
			},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "not found such user",
			},
		},
		{
			name: "usecase error - get user failed",
			setupMock: func(mockUsecase *login_user2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrGetUser)
			},
			requestBody: dto.LoginUserIn{
				Email:    validEmail,
				Password: validPassword,
			},
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "error while get user data",
			},
		},
		{
			name: "usecase error - incorrect password",
			setupMock: func(mockUsecase *login_user2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrIncorrectPass)
			},
			requestBody: dto.LoginUserIn{
				Email:    validEmail,
				Password: validPassword,
			},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "incorrect password",
			},
		},
		{
			name:         "invalid request body",
			setupMock:    func(mockUsecase *login_user2.Mockusecase) {},
			requestBody:  "invalid_body",
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "failed to decode request",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := login_user2.NewMockusecase(ctrl)
			handler := New(secret, mockUsecase, valid)

			tt.setupMock(mockUsecase)

			w := httptest.NewRecorder()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(
				"POST",
				"/login",
				bytes.NewReader(body),
			)

			handler.LoginUser(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expected != nil {
				var response dto.LoginUserOut
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				assert.NotEmpty(t, response.Token)
			}

			if tt.expectedError != nil {
				var errorResponse map[string]string
				err := json.NewDecoder(w.Body).Decode(&errorResponse)
				require.NoError(t, err)
				assert.Contains(t, errorResponse["message"], tt.expectedError["message"])
			}
		})
	}
}
