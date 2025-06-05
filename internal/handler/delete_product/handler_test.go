package delete_product

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	delete_product2 "pvz-service/internal/handler/delete_product/mocks"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/delete_product"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	valid := validator.New(validator.WithRequiredStructEnabled())

	pvzID := uuid.New()
	usecaseIn := delete_product.In{
		PVZID: pvzID,
	}

	successResponse := map[string]string{
		"message": "success delete product",
	}

	tests := []struct {
		name          string
		setupMock     func(*delete_product2.Mockusecase)
		pvzId         uuid.UUID
		expectedCode  int
		expected      map[string]string
		expectedError map[string]string
	}{
		{
			name: "successful delete product",
			setupMock: func(mockUsecase *delete_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusOK,
			expected:     successResponse,
		},
		{
			name:         "validation failed - empty pvzId",
			setupMock:    func(mockUsecase *delete_product2.Mockusecase) {},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name: "usecase error - PVZ not found",
			setupMock: func(mockUsecase *delete_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(usecase2.ErrNotFoundPVZ)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "pvz with such id not exist",
			},
		},
		{
			name: "usecase error - failed to look up for PVZ",
			setupMock: func(mockUsecase *delete_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(usecase2.ErrGetPVZByID)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to look up for pvz",
			},
		},
		{
			name: "usecase error - no receptions at all",
			setupMock: func(mockUsecase *delete_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(usecase2.ErrNotFoundReception)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "there is no receptions at all",
			},
		},
		{
			name: "usecase error - failed to get reception",
			setupMock: func(mockUsecase *delete_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(usecase2.ErrGetReception)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to get reception",
			},
		},
		{
			name: "usecase error - no opened reception",
			setupMock: func(mockUsecase *delete_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(usecase2.ErrNotFoundOpenedReception)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "no opened reception",
			},
		},
		{
			name: "usecase error - no product to delete",
			setupMock: func(mockUsecase *delete_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(usecase2.ErrNotFoundProduct)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "no product to delete",
			},
		},
		{
			name: "usecase error - failed to find product",
			setupMock: func(mockUsecase *delete_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(usecase2.ErrGetProduct)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to find product to delete",
			},
		},
		{
			name: "usecase error - failed to delete product",
			setupMock: func(mockUsecase *delete_product2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(usecase2.ErrDeleteProduct)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to delete product",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := delete_product2.NewMockusecase(ctrl)
			handler := New(mockUsecase, valid)

			tt.setupMock(mockUsecase)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				"/pvz/"+tt.pvzId.String()+"/delete_last_product",
				nil,
			)
			req = mux.SetURLVars(req, map[string]string{"pvzId": tt.pvzId.String()})

			handler.DeleteProduct(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expected != nil {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.expected["message"], response["message"])
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
