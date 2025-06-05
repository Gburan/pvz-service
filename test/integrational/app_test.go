package integrational

// TODO
//func (s *AppTestSuite) TestAppWith50Products() {
//	resp, err := dummyLogin("moderator")
//	s.NoError(err)
//	tokenModerator := assertDummyLogin(s.T(), resp, http.StatusOK, "")
//
//	resp, err = dummyLogin("employee")
//	s.NoError(err)
//	tokenEmployee := assertDummyLogin(s.T(), resp, http.StatusOK, "")
//
//	resp, err = createPVZ(tokenModerator, "Москва")
//	s.NoError(err)
//	pvz := assertCreatePVZ(s.T(), resp, http.StatusOK, "")
//
//	resp, err = startReception(tokenEmployee, pvz.Uuid)
//	s.NoError(err)
//	assertStartReception(s.T(), resp, pvz.Uuid, http.StatusOK, "")
//
//	for i := 0; i < 50; i++ {
//		resp, err = addProduct(tokenEmployee, pvz.Uuid, "электроника")
//		s.NoError(err)
//		assertAddProduct(s.T(), resp, http.StatusOK, "")
//	}
//
//	resp, err = closeReception(tokenEmployee, pvz.Uuid)
//	s.NoError(err)
//	assertCloseReception(s.T(), resp, http.StatusOK, "")
//
//	startDate := time.Now().Add(-time.Minute).UTC()
//	endDate := time.Now().Add(time.Minute).UTC()
//	resp, err = getPVZInfo(tokenModerator, startDate, endDate, 1, 50)
//	s.NoError(err)
//	assertGetPVZInfo(s.T(), resp, http.StatusOK, "")
//}
