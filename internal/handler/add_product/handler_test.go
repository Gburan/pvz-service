package add_product

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	add_product2 "pvz-service/internal/handler/add_product/mocks"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/add_product"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currTime := time.Now()

	valid := validator.New()
	err := valid.RegisterValidation("oneof_category", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		for _, allowed := range []string{"электроника", "одежда", "обувь"} {
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

	reqDTO := dto.AddProductIn{
		Type:  "электроника",
		PVZID: uuid.New(),
	}
	usecaseIn := add_product.In{
		Type:  reqDTO.Type,
		PVZID: reqDTO.PVZID,
	}
	usecaseOut := add_product.Out{
		Product: entity.Product{
			Uuid:        uuid.New(),
			DateTime:    currTime,
			Type:        reqDTO.Type,
			ReceptionID: uuid.New(),
		},
	}
	handlerOut := dto.AddProductOut{
		Uuid:        usecaseOut.Product.Uuid,
		DateTime:    usecaseOut.Product.DateTime.UTC(),
		Type:        usecaseOut.Product.Type,
		ReceptionID: usecaseOut.Product.ReceptionID,
	}

	tests := []struct {
		name          string
		setupMock     func(*add_product2.Mockusecase)
		reqBody       string
		expectedCode  int
		expected      *dto.AddProductOut
		expectedError map[string]string
	}{
		{
			name: "successful add product",
			setupMock: func(mockUsecase *add_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(&usecaseOut, nil)
			},
			reqBody: fmt.Sprintf(
				`{"type":"%s","pvzId":"%s"}`,
				reqDTO.Type,
				reqDTO.PVZID,
			),
			expectedCode: http.StatusOK,
			expected:     &handlerOut,
		},
		{
			name:         "empty body",
			setupMock:    func(mockUsecase *add_product2.Mockusecase) {},
			reqBody:      "",
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "failed to decode request",
			},
		},
		{
			name:         "validation failed - empty type",
			setupMock:    func(mockUsecase *add_product2.Mockusecase) {},
			reqBody:      fmt.Sprintf(`{"type":"%s","pvzId":"%s"}`, "", reqDTO.PVZID),
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name: "usecase error - PVZ not found",
			setupMock: func(mockUsecase *add_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrNotFoundPVZ)
			},
			reqBody: fmt.Sprintf(
				`{"type":"%s","pvzId":"%s"}`,
				reqDTO.Type,
				reqDTO.PVZID,
			),
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "pvz with such id not exist",
			},
		},
		{
			name: "usecase error - failed to look up for PVZ",
			setupMock: func(mockUsecase *add_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrGetPVZByID)
			},
			reqBody: fmt.Sprintf(
				`{"type":"%s","pvzId":"%s"}`,
				reqDTO.Type,
				reqDTO.PVZID,
			),
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to look up for pvz",
			},
		},
		{
			name: "usecase error - no receptions at all",
			setupMock: func(mockUsecase *add_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrNotFoundReception)
			},
			reqBody: fmt.Sprintf(
				`{"type":"%s","pvzId":"%s"}`,
				reqDTO.Type,
				reqDTO.PVZID,
			),
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "there is no receptions at all",
			},
		},
		{
			name: "usecase error - failed to get reception",
			setupMock: func(mockUsecase *add_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrGetReception)
			},
			reqBody: fmt.Sprintf(
				`{"type":"%s","pvzId":"%s"}`,
				reqDTO.Type,
				reqDTO.PVZID,
			),
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to get reception",
			},
		},
		{
			name: "usecase error - no opened reception",
			setupMock: func(mockUsecase *add_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrNotFoundOpenedReception)
			},
			reqBody: fmt.Sprintf(
				`{"type":"%s","pvzId":"%s"}`,
				reqDTO.Type,
				reqDTO.PVZID,
			),
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "no opened reception",
			},
		},
		{
			name: "usecase error - failed to add product",
			setupMock: func(mockUsecase *add_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrAddProduct)
			},
			reqBody: fmt.Sprintf(
				`{"type":"%s","pvzId":"%s"}`,
				reqDTO.Type,
				reqDTO.PVZID,
			),
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to add product",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := add_product2.NewMockusecase(ctrl)
			handler := New(mockUsecase, valid)

			tt.setupMock(mockUsecase)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				"/products",
				strings.NewReader(tt.reqBody),
			)
			req.Header.Set("Content-Type", "application/json")

			handler.AddProduct(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expected != nil {
				var response dto.AddProductOut
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)

				assert.Equal(t, tt.expected.Uuid, response.Uuid)
				assert.Equal(t, tt.expected.Type, response.Type)
				assert.Equal(t, tt.expected.ReceptionID, response.ReceptionID)
				assert.Equal(t, tt.expected.DateTime, response.DateTime)
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
