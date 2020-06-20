package common

import (
	"net/http"

	"JWT/models"
	"JWT/utils"
)

type Controller struct {
}
func (c Controller) ProtectedEndPoint(w http.ResponseWriter, req *http.Request) {
	user := models.User{
		ID:       0,
		Email:    "email",
		Password: "true",
	}

	authHeader := req.Header.Get("Authorization")
	println(authHeader)

	utils.ResponseJSON(w, http.StatusOK, user)
}
