package integrational

// TODO
//import (
//	"net/http"
//	"strconv"
//	"testing"
//	"time"
//
//	"pvz-service/internal/config"
//	"pvz-service/internal/jwt"
//
//	"github.com/google/uuid"
//)
//
//func (s *AppTestSuite) TestJWTTokenMethods() {
//	invalidToken := "some_invalid_token"
//	pvzID := uuid.New()
//	startDate := time.Now().Add(-time.Minute).UTC().Format(time.RFC3339Nano)
//	endDate := time.Now().Add(time.Minute).UTC().Format(time.RFC3339Nano)
//
//	resp, err := dummyLogin("employee")
//	s.NoError(err)
//	employeeToken := assertDummyLogin(s.T(), resp, http.StatusOK, "")
//
//	resp, err = dummyLogin("moderator")
//	s.NoError(err)
//	moderatorToken := assertDummyLogin(s.T(), resp, http.StatusOK, "")
//
//	cfg := config.MustLoad("../../config/config.yaml")
//	insufficientToken, _ := jwt.GenerateToken(cfg.App.JWTToken, "strange_role", uuid.New(), time.Minute)
//	expiredToken, _ := jwt.GenerateToken(cfg.App.JWTToken, "moderator", uuid.New(), time.Microsecond)
//	time.Sleep(time.Millisecond)
//
//	testCases := []struct {
//		name        string
//		requestFunc func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error)
//		assertFunc  func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string)
//		args        []string
//		errorMsg    string
//		status      int
//		token       string
//	}{
//		{
//			name: "AddProductWithInvalidToken",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				return addProduct(token, pvzID, args[0])
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertAddProduct(t, resp, expectedStatus, errorMsg)
//			},
//			args:     []string{"электроника"},
//			status:   http.StatusUnauthorized,
//			errorMsg: "invalid token",
//			token:    invalidToken,
//		},
//		{
//			name: "AddProductWithInsufficientPermissions",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				return addProduct(token, pvzID, args[0])
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertAddProduct(t, resp, expectedStatus, errorMsg)
//			},
//			args:     []string{"электроника"},
//			status:   http.StatusForbidden,
//			errorMsg: "insufficient permissions",
//			token:    moderatorToken,
//		},
//		{
//			name: "CreatePVZWithInvalidToken",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				return createPVZ(token, args[0])
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertCreatePVZ(t, resp, expectedStatus, errorMsg)
//			},
//			args:     []string{"Москва"},
//			status:   http.StatusUnauthorized,
//			errorMsg: "invalid token",
//			token:    invalidToken,
//		},
//		{
//			name: "CreatePVZWithInsufficientPermissions",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				return createPVZ(token, args[0])
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertCreatePVZ(t, resp, expectedStatus, errorMsg)
//			},
//			args:     []string{"Москва"},
//			status:   http.StatusForbidden,
//			errorMsg: "insufficient permissions",
//			token:    employeeToken,
//		},
//		{
//			name: "StartReceptionWithInvalidToken",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				return startReception(token, pvzID)
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertStartReception(t, resp, pvzID, expectedStatus, errorMsg)
//			},
//			args:     []string{},
//			status:   http.StatusUnauthorized,
//			errorMsg: "invalid token",
//			token:    invalidToken,
//		},
//		{
//			name: "StartReceptionWithInsufficientPermissions",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				return startReception(token, pvzID)
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertStartReception(t, resp, pvzID, expectedStatus, errorMsg)
//			},
//			args:     []string{},
//			status:   http.StatusForbidden,
//			errorMsg: "insufficient permissions",
//			token:    moderatorToken,
//		},
//		{
//			name: "CloseReceptionWithInvalidToken",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				return closeReception(token, pvzID)
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertCloseReception(t, resp, expectedStatus, errorMsg)
//			},
//			args:     []string{},
//			status:   http.StatusUnauthorized,
//			errorMsg: "invalid token",
//			token:    invalidToken,
//		},
//		{
//			name: "CloseReceptionWithInsufficientPermissions",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				return closeReception(token, pvzID)
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertCloseReception(t, resp, expectedStatus, errorMsg)
//			},
//			args:     []string{},
//			status:   http.StatusForbidden,
//			errorMsg: "insufficient permissions",
//			token:    moderatorToken,
//		},
//		{
//			name: "GetPVZInfoWithInvalidToken",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				page, _ := strconv.Atoi(args[2])
//				limit, _ := strconv.Atoi(args[3])
//				start, _ := time.Parse(time.RFC3339Nano, args[0])
//				end, _ := time.Parse(time.RFC3339Nano, args[1])
//				return getPVZInfo(token, start, end, page, limit)
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertGetPVZInfo(t, resp, expectedStatus, errorMsg)
//			},
//			args:     []string{startDate, endDate, "1", "50"},
//			status:   http.StatusUnauthorized,
//			errorMsg: "invalid token",
//			token:    invalidToken,
//		},
//		{
//			name: "GetPVZInfoWithInsufficientPermissions",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				page, _ := strconv.Atoi(args[2])
//				limit, _ := strconv.Atoi(args[3])
//				start, _ := time.Parse(time.RFC3339Nano, args[0])
//				end, _ := time.Parse(time.RFC3339Nano, args[1])
//				return getPVZInfo(token, start, end, page, limit)
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertGetPVZInfo(t, resp, expectedStatus, errorMsg)
//			},
//			args:     []string{startDate, endDate, "1", "50"},
//			status:   http.StatusForbidden,
//			errorMsg: "insufficient permissions",
//			token:    insufficientToken,
//		},
//		{
//			name: "GetPVZInfoWithExpiredToken",
//			requestFunc: func(token string, pvzID uuid.UUID, args ...string) (*http.Response, error) {
//				page, _ := strconv.Atoi(args[2])
//				limit, _ := strconv.Atoi(args[3])
//				start, _ := time.Parse(time.RFC3339Nano, args[0])
//				end, _ := time.Parse(time.RFC3339Nano, args[1])
//				return getPVZInfo(token, start, end, page, limit)
//			},
//			assertFunc: func(t *testing.T, resp *http.Response, expectedStatus int, pvzID uuid.UUID, errorMsg string) {
//				assertGetPVZInfo(t, resp, expectedStatus, errorMsg)
//			},
//			args:     []string{startDate, endDate, "1", "50"},
//			status:   http.StatusUnauthorized,
//			errorMsg: "invalid token",
//			token:    expiredToken,
//		},
//	}
//
//	for _, tc := range testCases {
//		s.T().Run(tc.name, func(t *testing.T) {
//			resp, err = tc.requestFunc(tc.token, pvzID, tc.args...)
//			s.NoError(err)
//			tc.assertFunc(t, resp, tc.status, pvzID, tc.errorMsg)
//		})
//	}
//}
