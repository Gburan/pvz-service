package integrational

import (
	"context"
	"time"

	"pvz-service/internal/config"
	nower2 "pvz-service/internal/infrastructure/nower"
	"pvz-service/internal/infrastructure/repository/pvz"
	"pvz-service/internal/infrastructure/repository/reception"
	"pvz-service/internal/jwt"
	"pvz-service/internal/model/entity"

	"github.com/google/uuid"
)

func (s *AppTestSuite) TestCloseReceptionSuccessful() {
	nower := nower2.Nower{}
	repPVZ := pvz.NewRepository(s.app.Pool, nower)
	repReception := reception.NewRepository(s.app.Pool, nower)

	pvzIn := entity.PVZ{
		Uuid:             uuid.New(),
		RegistrationDate: nower.Now(),
		City:             "Санкт-Петербург",
	}
	pvzOut, err := repPVZ.SavePVZ(context.Background(), pvzIn)
	s.Equal(pvzIn, *pvzOut)
	s.NoError(err)

	recIn := entity.Reception{
		Uuid:     uuid.New(),
		DateTime: nower.Now(),
		PVZID:    pvzIn.Uuid,
	}
	recOut, err := repReception.StartReception(context.Background(), recIn)
	recIn.Status = "in_progress"
	s.Equal(recIn, *recOut)
	s.NoError(err)

	cfg := config.MustLoad("../../config/config.yaml")
	token, _ := jwt.GenerateToken(cfg.App.JWTToken, "employee", uuid.New(), time.Minute)
	err = closeReception(token, recOut.PVZID)
	s.NoError(err)
}
