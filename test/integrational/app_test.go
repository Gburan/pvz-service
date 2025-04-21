package integrational

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"pvz-service/internal/app"
	"pvz-service/internal/config"
	"pvz-service/internal/model/entity"

	"github.com/stretchr/testify/assert"
)

var testServerAddr = "http://localhost:8080"

func TestApp(t *testing.T) {
	ctx := context.Background()

	cfg := config.MustLoad("../../config/config.yaml")
	cfg.DB.MigrationsDir = "../../migrations"

	testDB := SetupTestDatabase()
	defer testDB.TearDown()

	cfg.DB.Conn = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, testDB.DbAddress, DbName)

	go func() {
		a, err := app.NewApp(ctx, cfg)
		if err != nil {
			log.Fatal(err)
		}
		err = a.Run()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	time.Sleep(2 * time.Second)

	tokenModerator := doDummyLogin(t, "moderator")
	tokenEmployee := doDummyLogin(t, "employee")

	pvz := createPVZ(t, tokenModerator)

	doStartReception(t, tokenEmployee, pvz.Uuid)

	for i := 0; i < 50; i++ {
		addProduct(t, tokenEmployee, pvz.Uuid)
	}

	closeReception(t, tokenEmployee, pvz.Uuid)

	getPVZInfo(
		t,
		tokenModerator,
		time.Now().Add(-time.Minute).UTC().Format(time.RFC3339Nano),
		time.Now().Add(time.Minute).UTC().Format(time.RFC3339Nano),
		1,
		50,
	)
}

func doDummyLogin(t *testing.T, role string) string {
	resp := post(t, "/dummyLogin", fmt.Sprintf(`{"role":"%s"}`, role))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var out struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out.Token
}

func createPVZ(t *testing.T, token string) entity.PVZ {
	body := `{"city": "Москва"}`
	resp := postAuth(t, "/pvz", token, body)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var pvz_ struct {
		Uuid             string    `json:"id"`
		RegistrationDate time.Time `json:"registrationDate"`
		City             string    `json:"city"`
	}
	err := json.NewDecoder(resp.Body).Decode(&pvz_)
	assert.NoError(t, err)

	assert.NotEmpty(t, pvz_.Uuid)
	assert.Equal(t, "Москва", pvz_.City)
	assert.False(t, pvz_.RegistrationDate.IsZero())

	return entity.PVZ{
		Uuid:             pvz_.Uuid,
		RegistrationDate: pvz_.RegistrationDate,
		City:             pvz_.City,
	}
}

func doStartReception(t *testing.T, token, pvzId string) entity.Reception {
	body := fmt.Sprintf(`{"pvzId": "%s"}`, pvzId)
	resp := postAuth(t, "/receptions", token, body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var rec struct {
		Uuid             string    `json:"id"`
		RegistrationDate time.Time `json:"registrationDate"`
		PVZID            string    `json:"pvzId"`
		Status           string    `json:"status"`
	}
	err := json.NewDecoder(resp.Body).Decode(&rec)
	assert.NoError(t, err)

	assert.NotEmpty(t, rec.Uuid)
	assert.Equal(t, pvzId, rec.PVZID)
	assert.Equal(t, "in_progress", rec.Status)

	return entity.Reception{
		Uuid:     rec.Uuid,
		DateTime: rec.RegistrationDate,
		PVZID:    rec.PVZID,
		Status:   rec.Status,
	}
}

func addProduct(t *testing.T, token string, pvzID string) entity.Product {
	body := fmt.Sprintf(`{"type":"электроника", "pvzId":"%s"}`, pvzID)
	resp := postAuth(t, "/products", token, body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var out struct {
		Uuid        string `json:"id"`
		Type        string `json:"type"`
		ReceptionID string `json:"receptionId"`
	}
	err := json.NewDecoder(resp.Body).Decode(&out)
	assert.NoError(t, err)

	assert.NotEmpty(t, out.Uuid)
	assert.Equal(t, "электроника", out.Type)
	assert.NotEmpty(t, out.ReceptionID)

	return entity.Product{
		Uuid:        out.Uuid,
		Type:        out.Type,
		ReceptionID: out.ReceptionID,
	}
}

func closeReception(t *testing.T, token string, pvzID string) {
	url := fmt.Sprintf("/pvz/%s/close_last_reception", pvzID)
	resp := postAuth(t, url, token, `{}`)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func getPVZInfo(t *testing.T, token string, startDate, endDate string, page, limit int) {
	body := fmt.Sprintf(`{
		"startDate": "%s",
		"endDate": "%s",
		"page": %d,
		"limit": %d
	}`, startDate, endDate, page, limit)

	req, err := http.NewRequest("GET", testServerAddr+"/pvz", strings.NewReader(body))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	type pvzOut struct {
		Uuid             string `json:"id"`
		RegistrationDate string `json:"registrationDate"`
		City             string `json:"city"`
	}
	type receptionOut struct {
		Id       string `json:"id"`
		DateTime string `json:"dateTime"`
		PvzId    string `json:"pvzId"`
		Status   string `json:"status"`
	}
	type productOut struct {
		Uuid        string `json:"id"`
		DateTime    string `json:"dateTime"`
		Type        string `json:"type"`
		ReceptionID string `json:"receptionId"`
	}
	type receptionWithProductsOut struct {
		Reception receptionOut `json:"reception"`
		Products  []productOut `json:"products"`
	}

	var pvzInfoOut []struct {
		PVZ        pvzOut                     `json:"pvz"`
		Receptions []receptionWithProductsOut `json:"receptions"`
	}

	err = json.NewDecoder(resp.Body).Decode(&pvzInfoOut)
	assert.NoError(t, err)

	for _, pvzInfo := range pvzInfoOut {
		assert.NotEmpty(t, pvzInfo.PVZ.Uuid)
		assert.NotEmpty(t, pvzInfo.PVZ.City)
		assert.False(t, pvzInfo.PVZ.RegistrationDate == "")

		registrDate, err := time.Parse(time.RFC3339, pvzInfo.PVZ.RegistrationDate)
		assert.NoError(t, err)
		assert.False(t, registrDate.IsZero())

		for _, receptionWithProducts := range pvzInfo.Receptions {
			assert.NotEmpty(t, receptionWithProducts.Reception.Id)
			assert.False(t, receptionWithProducts.Reception.DateTime == "")
			assert.Equal(t, pvzInfo.PVZ.Uuid, receptionWithProducts.Reception.PvzId)
			assert.NotEmpty(t, receptionWithProducts.Reception.Status)

			receptionDate, err := time.Parse(time.RFC3339, receptionWithProducts.Reception.DateTime)
			assert.NoError(t, err)
			assert.False(t, receptionDate.IsZero())

			for _, product := range receptionWithProducts.Products {
				assert.NotEmpty(t, product.Uuid)
				assert.False(t, product.DateTime == "")
				assert.NotEmpty(t, product.Type)

				productDate, err := time.Parse(time.RFC3339, product.DateTime)
				assert.NoError(t, err)
				assert.False(t, productDate.IsZero())

				assert.Equal(t, receptionWithProducts.Reception.Id, product.ReceptionID)
			}
		}
	}
}

func post(t *testing.T, path, body string) *http.Response {
	resp, err := http.Post(testServerAddr+path, "application/json", bytes.NewBuffer([]byte(body)))
	assert.NoError(t, err)

	return resp
}

func postAuth(t *testing.T, path, token, body string) *http.Response {
	req, err := http.NewRequest("POST", testServerAddr+path, bytes.NewBuffer([]byte(body)))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	return resp
}
