package integrational

// TODO
//func (s *AppTestSuite) TestDeletePruductWithNonExistPVZ() {
//	resp, err := dummyLogin("employee")
//	s.NoError(err)
//	token := assertDummyLogin(s.T(), resp, http.StatusOK, "")
//
//	nonExistentPVZID := uuid.New()
//
//	resp, err = deleteProduct(token, nonExistentPVZID)
//	s.NoError(err)
//	assertDeleteProduct(s.T(), resp, http.StatusBadRequest, "pvz with such id not exist")
//}
//
//func (s *AppTestSuite) TestDeletePruductWithPVZNoReceptions() {
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
//	resp, err = deleteProduct(tokenEmployee, pvz.Uuid)
//	s.NoError(err)
//	assertDeleteProduct(s.T(), resp, http.StatusBadRequest, "there is no receptions at all")
//}
//
//func (s *AppTestSuite) TestDeleteProductWithNoOpenedReception() {
//	resp, err := dummyLogin("moderator")
//	s.NoError(err)
//	tokenModerator := assertDummyLogin(s.T(), resp, http.StatusOK, "")
//
//	resp, err = createPVZ(tokenModerator, "Москва")
//	s.NoError(err)
//	pvz := assertCreatePVZ(s.T(), resp, http.StatusOK, "")
//
//	resp, err = dummyLogin("employee")
//	s.NoError(err)
//	tokenEmployee := assertDummyLogin(s.T(), resp, http.StatusOK, "")
//
//	resp, err = startReception(tokenEmployee, pvz.Uuid)
//	s.NoError(err)
//	assertStartReception(s.T(), resp, pvz.Uuid, http.StatusOK, "")
//
//	resp, err = closeReception(tokenEmployee, pvz.Uuid)
//	s.NoError(err)
//	assertCloseReception(s.T(), resp, http.StatusOK, "")
//
//	resp, err = deleteProduct(tokenEmployee, pvz.Uuid)
//	s.NoError(err)
//	assertDeleteProduct(s.T(), resp, http.StatusBadRequest, "no opened reception")
//}
//
//func (s *AppTestSuite) TestDeleteProductWithNoProductInReception() {
//	resp, err := dummyLogin("moderator")
//	s.NoError(err)
//	tokenModerator := assertDummyLogin(s.T(), resp, http.StatusOK, "")
//
//	resp, err = createPVZ(tokenModerator, "Москва")
//	s.NoError(err)
//	pvz := assertCreatePVZ(s.T(), resp, http.StatusOK, "")
//
//	resp, err = dummyLogin("employee")
//	s.NoError(err)
//	tokenEmployee := assertDummyLogin(s.T(), resp, http.StatusOK, "")
//
//	resp, err = startReception(tokenEmployee, pvz.Uuid)
//	s.NoError(err)
//	assertStartReception(s.T(), resp, pvz.Uuid, http.StatusOK, "")
//
//	resp, err = deleteProduct(tokenEmployee, pvz.Uuid)
//	s.NoError(err)
//	assertDeleteProduct(s.T(), resp, http.StatusBadRequest, "no product to delete")
//}
