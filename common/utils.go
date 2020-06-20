package common

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"lpua_back-end/models"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Response(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(data)
}

func ResponseJSON(w http.ResponseWriter, status int,  data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ComparePassword(hashedPassword string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
	if err != nil {
		return false
	}
	return true
}
func HandleSuccess(resp http.ResponseWriter, data interface{}, user *models.MstUser) {
	returnData:=models.ResponseWrapper{
		Success:true,
		Message: "Success",
		Data:data,
	}

	jsonData, err := json.Marshal(returnData)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Something when wrong"))
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonData)
}

func GenerateToken(user models.MstUser) (string, error) {
	secret := viper.GetString("secret.secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id" : user.ID,
		"email" : user.Email,
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString, nil
}