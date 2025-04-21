package register_user

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	register_user2 "pvz-service/internal/handler/register_user/mocks"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/register_user"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		log.Fatal(err)
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
			ID:    "6d132f66-dcfe-493e-965d-95c99e5f325d",
			Email: validEmail,
			Role:  validRole,
		},
	}

	successResponse := registerUserOut{
		Uuid:  retUser.User.ID,
		Email: retUser.User.Email,
		Role:  retUser.User.Role,
	}

	tests := []struct {
		name          string
		setupMock     func(*register_user2.Mockusecase)
		requestBody   interface{}
		expectedCode  int
		expected      *registerUserOut
		expectedError map[string]string
	}{
		{
			name: "successful registration",
			setupMock: func(mockUsecase *register_user2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(context.TODO(), usecaseIn).
					Return(retUser, nil)
			},
			requestBody: registerUserIn{
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
			requestBody: registerUserIn{
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
			requestBody: registerUserIn{
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
			requestBody: registerUserIn{
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
			requestBody: registerUserIn{
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
					Run(context.TODO(), usecaseIn).
					Return(nil, usecase2.ErrUserAlreadyExist)
			},
			requestBody: registerUserIn{
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
					Run(context.TODO(), usecaseIn).
					Return(nil, usecase2.ErrGetUser)
			},
			requestBody: registerUserIn{
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
					Run(context.TODO(), usecaseIn).
					Return(nil, usecase2.ErrGenHashedPass)
			},
			requestBody: registerUserIn{
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
					Run(context.TODO(), usecaseIn).
					Return(nil, usecase2.ErrAddUser)
			},
			requestBody: registerUserIn{
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
				var response registerUserOut
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
