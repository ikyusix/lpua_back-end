package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"

	"lpua_back-end/common"
)

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 {
			authToken := bearerToken[1]
			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte(viper.GetString("secret.secret")), nil
			})
			if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				common.Response(writer, common.Message(false, "error "+ err.Error()))
				return
			}
			//spew.Dump(token)
			if token.Valid {
				next.ServeHTTP(writer, req)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
				common.Response(writer, common.Message(false, "error "+ err.Error()))
			}
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			common.Response(writer, common.Message(true, "Invalid Token"))
			return
		}
	})
}