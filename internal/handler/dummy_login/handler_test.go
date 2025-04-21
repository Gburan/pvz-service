package dummy_login

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwt2 "pvz-service/internal/jwt"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDummyLogin(t *testing.T) {
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

	secret := "test-secret"
	reqDTO := dummyLoginIn{
		Role: "moderator",
	}

	tests := []struct {
		name          string
		reqBody       string
		expectedCode  int
		expected      *dummyLoginOut
		expectedError map[string]string
	}{
		{
			name:         "successful login",
			reqBody:      fmt.Sprintf(`{"role":"%s"}`, reqDTO.Role),
			expectedCode: http.StatusOK,
			expected: &dummyLoginOut{
				Token: "",
			},
		},
		{
			name:         "empty body",
			reqBody:      "",
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "invalid request body",
			},
		},
		{
			name:         "validation failed - empty role",
			reqBody:      `{"role":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name:         "validation failed - invalid role",
			reqBody:      `{"role":"superuser"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(secret, valid)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				"/dummylogin",
				strings.NewReader(tt.reqBody),
			)
			req.Header.Set("Content-Type", "application/json")

			handler.DummyLogin(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expected != nil {
				var response dummyLoginOut
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)

				if tt.expected.Token == "" {
					assert.NotEmpty(t, response.Token)
					_, err := jwt2.ParseToken(response.Token, secret)
					assert.NoError(t, err)
				} else {
					assert.Equal(t, tt.expected.Token, response.Token)
				}
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
