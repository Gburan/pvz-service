package start_reception

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	start_reception2 "pvz-service/internal/handler/start_reception/mocks"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/start_reception"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStartReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currTime := time.Now()
	valid := validator.New(validator.WithRequiredStructEnabled())

	pvzID := uuid.New()
	reqDTO := dto.StartReceptionIn{
		PVZID: pvzID,
	}
	usecaseIn := start_reception.In{
		PVZID: reqDTO.PVZID,
	}
	usecaseOut := start_reception.Out{
		Reception: entity.Reception{
			Uuid:     uuid.New(),
			DateTime: currTime,
			PVZID:    pvzID,
			Status:   "opened",
		},
	}
	handlerOut := dto.StartReceptionOut{
		Id:       usecaseOut.Reception.Uuid,
		DateTime: usecaseOut.Reception.DateTime.UTC(),
		PvzId:    usecaseOut.Reception.PVZID,
		Status:   usecaseOut.Reception.Status,
	}

	tests := []struct {
		name          string
		setupMock     func(*start_reception2.Mockusecase)
		reqBody       string
		expectedCode  int
		expected      *dto.StartReceptionOut
		expectedError map[string]string
	}{
		{
			name: "successful start reception",
			setupMock: func(mockUsecase *start_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(&usecaseOut, nil)
			},
			reqBody: fmt.Sprintf(
				`{"pvzId":"%s"}`,
				pvzID,
			),
			expectedCode: http.StatusOK,
			expected:     &handlerOut,
		},
		{
			name:         "empty body",
			setupMock:    func(mockUsecase *start_reception2.Mockusecase) {},
			reqBody:      "",
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "failed to decode request",
			},
		},
		{
			name: "usecase error - PVZ not found",
			setupMock: func(mockUsecase *start_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrNotFoundPVZ)
			},
			reqBody: fmt.Sprintf(
				`{"pvzId":"%s"}`,
				pvzID,
			),
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "pvz with such id not exist",
			},
		},
		{
			name: "usecase error - failed to look up for PVZ",
			setupMock: func(mockUsecase *start_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrGetPVZByID)
			},
			reqBody: fmt.Sprintf(
				`{"pvzId":"%s"}`,
				pvzID,
			),
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to look up for pvz",
			},
		},
		{
			name: "usecase error - failed to get reception",
			setupMock: func(mockUsecase *start_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrGetReception)
			},
			reqBody: fmt.Sprintf(
				`{"pvzId":"%s"}`,
				pvzID,
			),
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to get reception",
			},
		},
		{
			name: "usecase error - opened reception already exist",
			setupMock: func(mockUsecase *start_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrFoundOpenedReception)
			},
			reqBody: fmt.Sprintf(
				`{"pvzId":"%s"}`,
				pvzID,
			),
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "opened reception already exist",
			},
		},
		{
			name: "usecase error - failed to start reception",
			setupMock: func(mockUsecase *start_reception2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), usecaseIn).
					Return(nil, usecase2.ErrStartReception)
			},
			reqBody: fmt.Sprintf(
				`{"pvzId":"%s"}`,
				pvzID,
			),
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to start reception",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := start_reception2.NewMockusecase(ctrl)
			handler := New(mockUsecase, valid)

			tt.setupMock(mockUsecase)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				"/receptions",
				strings.NewReader(tt.reqBody),
			)
			req.Header.Set("Content-Type", "application/json")

			handler.StartReception(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expected != nil {
				var response dto.StartReceptionOut
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)

				assert.Equal(t, tt.expected.Id, response.Id)
				assert.Equal(t, tt.expected.PvzId, response.PvzId)
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
