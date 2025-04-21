package create_pvz

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	create_pvz2 "pvz-service/internal/handler/create_pvz/mocks"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/create_pvz"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currTime := time.Now()

	valid := validator.New()
	err := valid.RegisterValidation("oneof_city", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		for _, allowed := range []string{"Москва", "Санкт-Петербург", "Казань"} {
			if value == allowed {
				return true
			}
		}
		return false
	})
	if err != nil {
		log.Fatal(err)
	}

	reqDTO := createPVZIn{
		City: "Москва",
	}
	usecaseIn := create_pvz.In{
		City: reqDTO.City,
	}
	usecaseOut := create_pvz.Out{
		PVZ: entity.PVZ{
			Uuid:             "671353c3-d091-4de8-83f9-983fb6e34ecf",
			RegistrationDate: currTime,
			City:             reqDTO.City,
		},
	}
	handlerOut := createPVZOut{
		Uuid:             usecaseOut.PVZ.Uuid,
		RegistrationDate: usecaseOut.PVZ.RegistrationDate.UTC().Format(time.RFC3339Nano),
		City:             usecaseOut.PVZ.City,
	}

	tests := []struct {
		name          string
		setupMock     func(*create_pvz2.Mockusecase)
		reqBody       string
		expectedCode  int
		expected      *createPVZOut
		expectedError map[string]string
	}{
		{
			name: "successful create PVZ",
			setupMock: func(mockUsecase *create_pvz2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(context.TODO(), usecaseIn).
					Return(&usecaseOut, nil)
			},
			reqBody: fmt.Sprintf(
				`{"city":"%s"}`,
				reqDTO.City,
			),
			expectedCode: http.StatusOK,
			expected:     &handlerOut,
		},
		{
			name:         "empty body",
			setupMock:    func(mockUsecase *create_pvz2.Mockusecase) {},
			reqBody:      "",
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "failed to decode request",
			},
		},
		{
			name:         "validation failed - empty city",
			setupMock:    func(mockUsecase *create_pvz2.Mockusecase) {},
			reqBody:      `{"city":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name:         "validation failed - invalid city",
			setupMock:    func(mockUsecase *create_pvz2.Mockusecase) {},
			reqBody:      `{"city":"Нью-Йорк"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name: "usecase error - failed to add PVZ",
			setupMock: func(mockUsecase *create_pvz2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(context.TODO(), usecaseIn).
					Return(nil, usecase2.ErrAddPVZ)
			},
			reqBody: fmt.Sprintf(
				`{"city":"%s"}`,
				reqDTO.City,
			),
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to add pvz",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := create_pvz2.NewMockusecase(ctrl)
			handler := New(mockUsecase, valid)

			tt.setupMock(mockUsecase)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				"/pvz",
				strings.NewReader(tt.reqBody),
			)
			req.Header.Set("Content-Type", "application/json")

			handler.CreatePVZ(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expected != nil {
				var response createPVZOut
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)

				assert.Equal(t, tt.expected.Uuid, response.Uuid)
				assert.Equal(t, tt.expected.City, response.City)
				assert.Equal(t, tt.expected.RegistrationDate, response.RegistrationDate)
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
