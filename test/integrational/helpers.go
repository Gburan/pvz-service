package integrational

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	dto "pvz-service/internal/generated/api/v1/dto/handler"

	"github.com/google/uuid"
)

const (
	testServerAddr = "http://localhost:8080"

	methodPost = "POST"
	methodGet  = "GET"

	apiV1 = "/api/v1"
)

func apiPathV1(path string) string {
	return fmt.Sprintf("%s%s", apiV1, path)
}

// nolint:unused
func dummyLogin(role string) (dto.DummyLoginOut, error) {
	in := dto.DummyLoginIn{Role: role}
	return doRequest[dto.DummyLoginOut](methodPost, apiPathV1("/dummyLogin"), "", in)
}

// nolint:unused
func createPVZ(token, city string) (dto.CreatePVZOut, error) {
	in := dto.CreatePVZIn{City: city}
	return doRequest[dto.CreatePVZOut](methodPost, apiPathV1("/pvz"), token, in)
}

// nolint:unused
func startReception(token string, pvzId uuid.UUID) (dto.StartReceptionOut, error) {
	in := dto.StartReceptionIn{PVZID: pvzId}
	return doRequest[dto.StartReceptionOut](methodPost, apiPathV1("/receptions"), token, in)
}

// nolint:unused
func addProduct(token string, pvzId uuid.UUID, prodType string) (dto.AddProductOut, error) {
	in := dto.AddProductIn{
		PVZID: pvzId,
		Type:  prodType,
	}
	return doRequest[dto.AddProductOut](methodPost, apiPathV1("/products"), token, in)
}

// nolint:unused
func deleteProduct(token string, pvzId uuid.UUID) error {
	url := fmt.Sprintf("/pvz/%s/delete_last_product", pvzId.String())
	_, err := doRequest[struct{}](methodPost, apiPathV1(url), token)
	return err
}

func closeReception(token string, pvzId uuid.UUID) error {
	url := fmt.Sprintf("/pvz/%s/close_last_reception", pvzId.String())
	_, err := doRequest[struct{}](methodPost, apiPathV1(url), token)
	return err
}

// nolint:unused
func getPVZInfo(token string, startDate, endDate time.Time, page, limit int) (dto.PvzInfoOut, error) {
	in := dto.PvzInfoIn{
		EndDate:   endDate,
		Limit:     limit,
		Page:      page,
		StartDate: startDate,
	}
	return doRequest[dto.PvzInfoOut](methodGet, apiPathV1("/pvz"), token, in)
}

// TODO
//
//	func listPVZs(ctx context.Context, client pvz_v1.PVZServiceClient) (*pvz_v1.GetPVZListResponse, error) {
//		req := &pvz_v1.GetPVZListRequest{}
//		return client.GetPVZList(ctx, req)
//	}

func doRequest[T any](method, path, token string, inData ...interface{}) (T, error) {
	var result T
	var body io.Reader

	if len(inData) > 0 {
		jsonBody, err := json.Marshal(inData[0])
		if err != nil {
			return result, err
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", testServerAddr, path), body)
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}
