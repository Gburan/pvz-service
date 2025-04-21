package handler

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, status int, errorMsg string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	var err_ error
	if err != nil {
		err_ = json.NewEncoder(w).Encode(map[string]string{
			"message": errorMsg,
			"details": err.Error(),
		})
	} else {
		err_ = json.NewEncoder(w).Encode(map[string]string{
			"message": errorMsg,
		})
	}
	if err_ != nil {
		return
	}
}
