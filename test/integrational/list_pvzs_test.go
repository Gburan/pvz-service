package integrational

// TODO
//func (s *AppTestSuite) TestListPVZs() {
//	expectedPVZs := []*pvz_v1.PVZ{
//		{
//			Id:               "671353c3-d091-4de8-83f9-983fb6e34ecf",
//			City:             "Москва",
//			RegistrationDate: timestamppb.New(time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)),
//		},
//		{
//			Id:               "6d132f66-dcfe-493e-965d-95c99e5f325d",
//			City:             "Санкт-Петербург",
//			RegistrationDate: timestamppb.New(time.Date(2025, 5, 2, 13, 0, 0, 0, time.UTC)),
//		},
//	}
//
//	client, err := grpc.NewClient("http://localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
//	s.NoError(err)
//	grpcClient := pvz_v1.NewPVZServiceClient(client)
//
//	listResp, err := listPVZs(context.Background(), grpcClient)
//	s.NoError(err)
//	s.NotNil(listResp)
//	s.Len(listResp.Pvzs, len(expectedPVZs))
//
//	for i, pvz := range listResp.Pvzs {
//		s.Equal(expectedPVZs[i].Id, pvz.Id)
//		s.Equal(expectedPVZs[i].City, pvz.City)
//		s.True(expectedPVZs[i].RegistrationDate.AsTime().Equal(pvz.RegistrationDate.AsTime()))
//	}
//}
