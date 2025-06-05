package register_user

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	register_user2 "pvz-service/internal/handler/register_user/mocks"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/register_user"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	valid := validator.New()
	err := valid.RegisterValidation("oneof_user", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		for _, allowed := range []string{"moderator", "employee"} {
			if value == allowed {
				return true
			}
		}
		return false
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	validEmail := "test@example.com"
	validPassword := "valid_password"
	validRole := "employee"
	usecaseIn := register_user.In{
		Email:    validEmail,
		Password: validPassword,
		Role:     validRole,
	}

	retUser := &register_user.Out{
		User: entity.User{
			Uuid:  uuid.New(),
			Email: validEmail,
			Role:  validRole,
		},
	}

	successResponse := dto.RegisterUserOut{
		Uuid:  retUser.User.Uuid,
		Email: retUser.User.Email,
		Role:  retUser.User.Role,
	}

	tests := []struct {
		name          string
		setupMock     func(*register_user2.Mockusecase)
		requestBody   interface{}
		expectedCode  int
		expected      *dto.RegisterUserOut
		expectedError map[string]string
	}{
		{
			name: "successful registration",
			setupMock: func(mockUsecase *register_user2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(retUser, nil)
			},
			requestBody: dto.RegisterUserIn{
				Email:    validEmail,
				Password: validPassword,
				Role:     validRole,
			},
			expectedCode: http.StatusOK,
			expected:     &successResponse,
		},
		{
			name:      "validation failed - empty email",
			setupMock: func(mockUsecase *register_user2.Mockusecase) {},
			requestBody: dto.RegisterUserIn{
				Email:    "",
				Password: validPassword,
				Role:     validRole,
			},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name:      "validation failed - invalid email format",
			setupMock: func(mockUsecase *register_user2.Mockusecase) {},
			requestBody: dto.RegisterUserIn{
				Email:    "invalid_email",
				Password: validPassword,
				Role:     validRole,
			},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name:      "validation failed - empty password",
			setupMock: func(mockUsecase *register_user2.Mockusecase) {},
			requestBody: dto.RegisterUserIn{
				Email:    validEmail,
				Password: "",
				Role:     validRole,
			},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name:      "validation failed - empty role",
			setupMock: func(mockUsecase *register_user2.Mockusecase) {},
			requestBody: dto.RegisterUserIn{
				Email:    validEmail,
				Password: validPassword,
				Role:     "",
			},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name: "usecase error - user already exists",
			setupMock: func(mockUsecase *register_user2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrUserAlreadyExist)
			},
			requestBody: dto.RegisterUserIn{
				Email:    validEmail,
				Password: validPassword,
				Role:     validRole,
			},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "such user already exist",
			},
		},
		{
			name: "usecase error - failed to get user",
			setupMock: func(mockUsecase *register_user2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrGetUser)
			},
			requestBody: dto.RegisterUserIn{
				Email:    validEmail,
				Password: validPassword,
				Role:     validRole,
			},
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to look up for existing user",
			},
		},
		{
			name: "usecase error - failed to generate password hash",
			setupMock: func(mockUsecase *register_user2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrGenHashedPass)
			},
			requestBody: dto.RegisterUserIn{
				Email:    validEmail,
				Password: validPassword,
				Role:     validRole,
			},
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "error while gen has password",
			},
		},
		{
			name: "usecase error - failed to add user",
			setupMock: func(mockUsecase *register_user2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrAddUser)
			},
			requestBody: dto.RegisterUserIn{
				Email:    validEmail,
				Password: validPassword,
				Role:     validRole,
			},
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "error while adding user in db",
			},
		},
		{
			name:         "invalid request body",
			setupMock:    func(mockUsecase *register_user2.Mockusecase) {},
			requestBody:  "invalid_body",
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "failed to decode request",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := register_user2.NewMockusecase(ctrl)
			handler := New(mockUsecase, valid)

			tt.setupMock(mockUsecase)

			w := httptest.NewRecorder()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(
				"POST",
				"/register",
				bytes.NewReader(body),
			)

			handler.RegisterUser(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expected != nil {
				var response dto.RegisterUserOut
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Uuid, response.Uuid)
				assert.Equal(t, tt.expected.Email, response.Email)
				assert.Equal(t, tt.expected.Role, response.Role)
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
