package integrational

// TODO
//func (s *AppTestSuite) TestStartReceptionWithNonExistentPVZ() {
//	resp, err := dummyLogin("employee")
//	s.NoError(err)
//	token := assertDummyLogin(s.T(), resp, http.StatusOK, "")
//
//	nonExistentPVZID := uuid.New()
//
//	resp, err = startReception(token, nonExistentPVZID)
//	s.NoError(err)
//	assertStartReception(s.T(), resp, uuid.New(), http.StatusBadRequest, "pvz with such id not exist")
//}
//
//func (s *AppTestSuite) TestStartReceptionWithReceptionAlreadyInProgress() {
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
//	resp, err = startReception(tokenEmployee, pvz.Uuid)
//	s.NoError(err)
//	assertStartReception(s.T(), resp, pvz.Uuid, http.StatusBadRequest, "opened reception already exist")
//}
