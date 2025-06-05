package pvz_info

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	dto "pvz-service/internal/generated/api/v1/dto/handler"
	pvz_info2 "pvz-service/internal/handler/pvz_info/mocks"
	"pvz-service/internal/model/entity"
	usecase2 "pvz-service/internal/usecase"
	"pvz-service/internal/usecase/pvz_info"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetPVZInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currTime := time.Now()
	valid := validator.New(validator.WithRequiredStructEnabled())

	reqDTO := dto.PvzInfoIn{
		StartDate: currTime.Add(-24 * time.Hour),
		EndDate:   currTime,
		Limit:     1,
		Page:      1,
	}

	retPVZ := entity.PVZ{
		Uuid:             uuid.New(),
		RegistrationDate: currTime.Add(-48 * time.Hour),
		City:             "Санкт-Петербург",
	}
	retReception := entity.Reception{
		Uuid:     uuid.New(),
		DateTime: currTime.Add(-12 * time.Hour),
		PVZID:    retPVZ.Uuid,
		Status:   "completed",
	}
	retProduct := entity.Product{
		Uuid:        uuid.New(),
		DateTime:    currTime.Add(-6 * time.Hour),
		Type:        "электроника",
		ReceptionID: retReception.Uuid,
	}

	usecaseResult := []pvz_info.Out{
		{
			PVZ: retPVZ,
			Receptions: []entity.ReceptionWithProducts{
				{
					Reception: retReception,
					Products:  []entity.Product{retProduct},
				},
			},
		},
	}

	handlerOut := []dto.PvzInfoOut{
		{
			Pvz: dto.PvzInfoPvzOut{
				Uuid:             retPVZ.Uuid,
				RegistrationDate: retPVZ.RegistrationDate,
				City:             retPVZ.City,
			},
			Receptions: []dto.PvzInfoReceptionWithProductsOut{
				{
					Reception: dto.PvzInfoReceptionOut{
						Id:       retReception.Uuid,
						DateTime: retReception.DateTime,
						PvzId:    retReception.PVZID,
						Status:   retReception.Status,
					},
					Products: []dto.PvzInfoProductOut{
						{
							Uuid:        retProduct.Uuid,
							DateTime:    retProduct.DateTime,
							Type:        retProduct.Type,
							ReceptionID: retProduct.ReceptionID,
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		setupMock     func(*pvz_info2.Mockusecase)
		reqBody       string
		expectedCode  int
		expected      []dto.PvzInfoOut
		expectedError map[string]string
	}{
		{
			name: "successful get pvz info",
			setupMock: func(mockUsecase *pvz_info2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), gomock.AssignableToTypeOf(pvz_info.In{})).
					Return(usecaseResult, nil)
			},
			reqBody: fmt.Sprintf(
				`{"startDate":"%s","endDate":"%s","limit":%d,"page":%d}`,
				reqDTO.StartDate.Format(time.RFC3339),
				reqDTO.EndDate.Format(time.RFC3339),
				reqDTO.Limit,
				reqDTO.Page,
			),
			expectedCode: http.StatusOK,
			expected:     handlerOut,
		},
		{
			name:      "incorrect number of pages to list",
			setupMock: func(mockUsecase *pvz_info2.Mockusecase) {},
			reqBody: fmt.Sprintf(
				`{"startDate":"%s","endDate":"%s","limit":%d,"page":%d}`,
				reqDTO.StartDate.Format(time.RFC3339),
				reqDTO.EndDate.Format(time.RFC3339),
				reqDTO.Limit,
				-1,
			),
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "incorrect page request",
			},
		},
		{
			name:      "incorrect number of limit to list",
			setupMock: func(mockUsecase *pvz_info2.Mockusecase) {},
			reqBody: fmt.Sprintf(
				`{"startDate":"%s","endDate":"%s","limit":%d,"page":%d}`,
				reqDTO.StartDate.Format(time.RFC3339),
				reqDTO.EndDate.Format(time.RFC3339),
				-1,
				reqDTO.Page,
			),
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "incorrect limit request",
			},
		},
		{
			name:      "validation failed",
			setupMock: func(mockUsecase *pvz_info2.Mockusecase) {},
			reqBody: fmt.Sprintf(
				`{"startDate":"%s","limit":%d,"page":%d}`,
				reqDTO.StartDate.Format(time.RFC3339),
				reqDTO.Limit,
				reqDTO.Page,
			),
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "validation failed",
			},
		},
		{
			name:         "empty body",
			setupMock:    func(mockUsecase *pvz_info2.Mockusecase) {},
			reqBody:      "",
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "failed to decode request",
			},
		},
		{
			name: "usecase error - failed to get products",
			setupMock: func(mockUsecase *pvz_info2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					Return(nil, usecase2.ErrGetProducts)
			},
			reqBody: fmt.Sprintf(
				`{"startDate":"%s","endDate":"%s","limit":%d,"page":%d}`,
				reqDTO.StartDate.Format(time.RFC3339),
				reqDTO.EndDate.Format(time.RFC3339),
				reqDTO.Limit,
				reqDTO.Page,
			),
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to get products",
			},
		},
		{
			name: "usecase error - failed products not found",
			setupMock: func(mockUsecase *pvz_info2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					Return(nil, usecase2.ErrNotFoundProducts)
			},
			reqBody: fmt.Sprintf(
				`{"startDate":"%s","endDate":"%s","limit":%d,"page":%d}`,
				reqDTO.StartDate.Format(time.RFC3339),
				reqDTO.EndDate.Format(time.RFC3339),
				reqDTO.Limit,
				reqDTO.Page,
			),
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"message": "there is no products at this time interval",
			},
		},
		{
			name: "usecase error - failed to get receptions",
			setupMock: func(mockUsecase *pvz_info2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					Return(nil, usecase2.ErrGetReceptions)
			},
			reqBody: fmt.Sprintf(
				`{"startDate":"%s","endDate":"%s","limit":%d,"page":%d}`,
				reqDTO.StartDate.Format(time.RFC3339),
				reqDTO.EndDate.Format(time.RFC3339),
				reqDTO.Limit,
				reqDTO.Page,
			),
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to get receptions",
			},
		},
		{
			name: "usecase error - failed to get pvz list",
			setupMock: func(mockUsecase *pvz_info2.Mockusecase) {
				mockUsecase.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					Return(nil, usecase2.ErrGetPVZs)
			},
			reqBody: fmt.Sprintf(
				`{"startDate":"%s","endDate":"%s","limit":%d,"page":%d}`,
				reqDTO.StartDate.Format(time.RFC3339),
				reqDTO.EndDate.Format(time.RFC3339),
				reqDTO.Limit,
				reqDTO.Page,
			),
			expectedCode: http.StatusInternalServerError,
			expectedError: map[string]string{
				"message": "failed to get pvz list",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := pvz_info2.NewMockusecase(ctrl)
			handler := New(mockUsecase, valid)

			tt.setupMock(mockUsecase)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				"/pvz",
				strings.NewReader(tt.reqBody),
			)
			req.Header.Set("Content-Type", "application/json")

			handler.GetPVZInfo(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expected != nil {
				var response []dto.PvzInfoOut
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)

				require.Len(t, response, len(tt.expected))
				for i, expectedPVZ := range tt.expected {
					assert.Equal(t, expectedPVZ.Pvz.Uuid, response[i].Pvz.Uuid)
					assert.Equal(t, expectedPVZ.Pvz.City, response[i].Pvz.City)
					assert.True(t, expectedPVZ.Pvz.RegistrationDate.Equal(response[i].Pvz.RegistrationDate),
						"RegistrationDate mismatch: expected %v, got %v",
						expectedPVZ.Pvz.RegistrationDate,
						response[i].Pvz.RegistrationDate,
					)

					require.Len(t, response[i].Receptions, len(expectedPVZ.Receptions))
					for j, expectedReception := range expectedPVZ.Receptions {
						assert.Equal(t, expectedReception.Reception.Id, response[i].Receptions[j].Reception.Id)
						assert.Equal(t, expectedReception.Reception.Status, response[i].Receptions[j].Reception.Status)
						assert.True(t, expectedReception.Reception.DateTime.Equal(response[i].Receptions[j].Reception.DateTime),
							"Reception DateTime mismatch: expected %v, got %v",
							expectedReception.Reception.DateTime,
							response[i].Receptions[j].Reception.DateTime,
						)

						require.Len(t, response[i].Receptions[j].Products, len(expectedReception.Products))
						for k, expectedProduct := range expectedReception.Products {
							assert.Equal(t, expectedProduct.Uuid, response[i].Receptions[j].Products[k].Uuid)
							assert.Equal(t, expectedProduct.Type, response[i].Receptions[j].Products[k].Type)
							assert.True(t, expectedProduct.DateTime.Equal(response[i].Receptions[j].Products[k].DateTime),
								"Product DateTime mismatch: expected %v, got %v",
								expectedProduct.DateTime,
								response[i].Receptions[j].Products[k].DateTime,
							)
						}
					}
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
