package main

import (
	"fmt"
	"net/http"
	"net/url"

	userHandler "lpua_back-end/user/handler"
	_userRepository "lpua_back-end/user/repository"
	_userUsecase "lpua_back-end/user/usecase"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"lpua_back-end/models"
)

func init() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func main() {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	port := viper.GetString("port.port")
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("sslmode", "disable")
	connStr := fmt.Sprintf("%s?%s", connection, val.Encode())

	dbConn, err := gorm.Open("postgres", connStr)
	if err != nil {
		panic("Failed to connect database")
	}

	err = dbConn.DB().Ping()
	if err != nil {
		logrus.Error(err)
		return
	}
	fmt.Println("Success connect DB")

	defer func() {
		err = dbConn.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	dbConn.Debug().AutoMigrate(&models.MstUser{})

	route := mux.NewRouter().StrictSlash(true)

	userRepository := _userRepository.CreateUserRepositoryImpl(dbConn)
	userService := _userUsecase.CreateUserUsecaseImpl(userRepository)
	userHandler.CreateUserHandler(route, userService)

	fmt.Println("Starting web server at port: "+port)
	err = http.ListenAndServe(": "+port, route)
	if err != nil {
		logrus.Fatal(err)
	}

}
