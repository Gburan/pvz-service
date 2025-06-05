package close_reception

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	close_reception2 "pvz-service/internal/handler/close_reception/mocks"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/close_reception"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCloseReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currTime := time.Now()

	valid := validator.New(validator.WithRequiredStructEnabled())

	pvzID := uuid.New()

	usecaseIn := close_reception.In{
		PVZID: pvzID,
	}
	usecaseOut := close_reception.Out{
		Reception: entity.Reception{
			Uuid:     uuid.New(),
			DateTime: currTime,
			PVZID:    pvzID,
			Status:   "closed",
		},
	}
	handlerOut := dto.CloseReceptionOut{
		Uuid:     usecaseOut.Reception.Uuid,
		DateTime: usecaseOut.Reception.DateTime.UTC(),
		PVZID:    usecaseOut.Reception.PVZID,
		Status:   usecaseOut.Reception.Status,
	}

	tests := []struct {
		name          string
		setupMock     func(*close_reception2.Mockusecase)
		pvzId         uuid.UUID
		expectedCode  int
		expected      *dto.CloseReceptionOut
		expectedError map[string]string
	}{
		{
			name: "successful close reception",
			setupMock: func(mockUsecase *close_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(&usecaseOut, nil)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusOK,
			expected:     &handlerOut,
		},
		{
			name:         "validation failed - empty pvzId",
			setupMock:    func(mockUsecase *close_reception2.Mockusecase) {},
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name: "usecase error - PVZ not found",
			setupMock: func(mockUsecase *close_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrNotFoundPVZ)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "pvz with such id not exist",
			},
		},
		{
			name: "usecase error - failed to look up for PVZ",
			setupMock: func(mockUsecase *close_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrGetPVZByID)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to look up for pvz",
			},
		},
		{
			name: "usecase error - no receptions at all",
			setupMock: func(mockUsecase *close_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrNotFoundReception)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "there is no receptions at all",
			},
		},
		{
			name: "usecase error - failed to get reception",
			setupMock: func(mockUsecase *close_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrGetReception)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to get reception",
			},
		},
		{
			name: "usecase error - no opened reception",
			setupMock: func(mockUsecase *close_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrNotFoundOpenedReception)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "no opened reception",
			},
		},
		{
			name: "usecase error - failed to close reception",
			setupMock: func(mockUsecase *close_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrCloseReception)
			},
			pvzId:        pvzID,
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to close reception",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := close_reception2.NewMockusecase(ctrl)
			handler := New(mockUsecase, valid)

			tt.setupMock(mockUsecase)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				"/pvz/"+tt.pvzId.String()+"/close_last_reception",
				nil,
			)
			req = mux.SetURLVars(req, map[string]string{"pvzId": tt.pvzId.String()})

			handler.CloseReception(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expected != nil {
				var response dto.CloseReceptionOut
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)

				assert.Equal(t, tt.expected.Uuid, response.Uuid)
				assert.Equal(t, tt.expected.PVZID, response.PVZID)
				assert.Equal(t, tt.expected.Status, response.Status)
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
