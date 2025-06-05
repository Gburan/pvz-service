package integrational

// TODO
//func (s *AppTestSuite) TestAddProduct() {
//	ctrl := gomock.NewController(s.T())
//	defer ctrl.Finish()
//
//	nower := nower2.Nower{}
//	repPVZ := pvz.NewRepository(s.app.Pool, nower)
//	repReception := reception.NewRepository(s.app.Pool, nower)
//
//	pvzIn := entity.PVZ{
//		Uuid:             uuid.New(),
//		RegistrationDate: nower.Now(),
//		City:             "Санкт-Петербург",
//	}
//	pvzOut, err := repPVZ.SavePVZ(context.Background(), pvzIn)
//	s.Equal(pvzIn, *pvzOut)
//	s.NoError(err)
//
//	recIn := entity.Reception{
//		Uuid:     uuid.New(),
//		DateTime: nower.Now(),
//		PVZID:    pvzIn.Uuid,
//	}
//	recOut, err := repReception.StartReception(context.Background(), recIn)
//	recIn.Status = "in_progress"
//	s.Equal(recIn, *recOut)
//	s.NoError(err)
//
//	cfg := config.MustLoad("../../config/config.yaml")
//	token, _ := jwt.GenerateToken(cfg.App.JWTToken, "employee", uuid.New(), time.Minute)
//	prod := entity.Product{
//		Uuid:        uuid.New(),
//		DateTime:    nower.Now(),
//		Type:        "электроника",
//		ReceptionID: recOut.Uuid,
//	}
//	mock := add_product.NewMockusecase(ctrl)
//
//	mock.EXPECT().
//		Run(context.Background(), gomock.Any()).
//		SetArg(0, prod.DateTime)
//
//	prodOut, err := addProduct(token, recOut.PVZID, prod.Type)
//	s.Equal(dto.AddProductOut{
//		DateTime:    prod.DateTime,
//		Uuid:        prod.Uuid,
//		ReceptionID: prod.ReceptionID,
//		Type:        prod.Type,
//	}, prodOut)
//	s.NoError(err)
//}
