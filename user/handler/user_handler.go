package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"lpua_back-end/common"
	"lpua_back-end/middleware"
	"lpua_back-end/models"

	"lpua_back-end/user"
)

type UserHandler struct {
	userUsecase user.UserUsecaseInterface
}

func CreateUserHandler(route *mux.Router, userUsecase user.UserUsecaseInterface)  {
	userHandler := UserHandler{userUsecase}
	route.HandleFunc("/users", middleware.TokenVerifyMiddleware(userHandler.FindAllUser)).Methods(http.MethodGet)
	route.HandleFunc("/signin", userHandler.SignIn).Methods(http.MethodPost)
	route.HandleFunc("/signup", userHandler.SignUp).Methods(http.MethodPost)
	sub := route.PathPrefix("/user").Subrouter()
	sub.HandleFunc("/{id}", middleware.TokenVerifyMiddleware(userHandler.FindUserById)).Methods(http.MethodGet)
	sub.HandleFunc("/{id}", middleware.TokenVerifyMiddleware(userHandler.UpdateUser)).Methods(http.MethodPut)
	sub.HandleFunc("/{id}", middleware.TokenVerifyMiddleware(userHandler.DeleteUser)).Methods(http.MethodDelete)
	sub.HandleFunc("/upload/{id}", middleware.TokenVerifyMiddleware(userHandler.UploadImage)).Methods(http.MethodPost)
}

func (h *UserHandler) FindAllUser(writer http.ResponseWriter, request *http.Request) {
	response, err := h.userUsecase.GetAllUser()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, response)
		return
	}

	common.Response(writer, response)
	return
}

func (h *UserHandler) SignIn(writer http.ResponseWriter, request *http.Request) {

	var jwt models.JWT
	var user models.MstUser

	dataLogin, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, common.Message(false, "Login error "+err.Error()))
		logrus.Error(err)
	}

	err = json.Unmarshal(dataLogin, &user)
	if user.Email == ""{
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, common.Message(false, "Email is missing "+err.Error()))
		return
	}
	if user.Password == ""{
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, common.Message(false,  "Password is missing "+ err.Error()))
		return
	}
	password := user.Password
	dataUser, err := h.userUsecase.GetUserByEmail(user.Email)
	if err != nil {
		logrus.Error(err)
		return
	}

	hashedPassword := dataUser.Password
	checkPassword := common.ComparePassword(hashedPassword, []byte(password))
	if checkPassword {
		token, err := common.GenerateToken(user)
		if err != nil {
			logrus.Fatal(err)
		}
		jwt.Token = token
		jwt.ID = dataUser.ID
		//common.ResponseJSON(writer, http.StatusCreated, jwt)
		common.HandleSuccess(writer, jwt, nil)
		return
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
		common.Response(writer, common.Message(false, "Invalid password"))
	}
}

func (h *UserHandler) SignUp(writer http.ResponseWriter, request *http.Request) {
	user := new(models.MstUser)
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, common.Message(false, "Invalid Request "+err.Error()))
		return
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		logrus.Error(err)
		fmt.Errorf("[UserHandler.SignUp] Hash password: %w", err)
	}

	user.Password = string(hashPass)
	response, err := h.userUsecase.RegisterUser(user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, response)
		return
	}
	common.Response(writer, response)
	return
}

func (h *UserHandler) FindUserById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, common.Message(false, "Please provide valid id"))
		return
	}

	response, err := h.userUsecase.GetUserById(uint(id))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, response)
		return
	}

	common.Response(writer, response)
	return
}

func (h *UserHandler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	user := new(models.MstUser)
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, common.Message(false, "Please provide valid id"))
		return
	}
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, common.Message(false, "Invalid request "+err.Error()))
		return
	}
	response, err := h.userUsecase.UpdateUser(uint(id), user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, response)
		return
	}
	common.Response(writer, response)
	return
}

func (h *UserHandler) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, common.Message(false, "Please provide valid id"))
		return
	}

	response, err := h.userUsecase.DeleteUser(uint(id))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, response)
		return
	}
	common.Response(writer, response)
	return
}

func (h *UserHandler) UploadImage(writer http.ResponseWriter, request *http.Request) {
	pathvar := mux.Vars(request)
	id, err := strconv.Atoi(pathvar["id"])
	if err != nil {
		logrus.Error(err)
		fmt.Errorf("[UserHandler.UploadImage] Upload profile: %w", err)
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer,common.Message(false, "Please provide valid id"))
		return
	}

	uploadImage, handler, err := request.FormFile("UrlImage")
	if err != nil {
		logrus.Error(err)
		fmt.Errorf("[UserHandler.UploadImage] Upload profile: %w", err)
	}
	defer uploadImage.Close()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	fileName := handler.Filename
	tempFile := filepath.Join("assets", "user", fileName)
	targetFile, err := os.OpenFile(tempFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		logrus.Error(err)
		fmt.Errorf("[UserHandler.UploadImage]: %w", err)
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadImage); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	user, message, err := h.userUsecase.UploadImage(uint(id), tempFile)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, common.Message(false, "Ooops something wrong"))
		logrus.Error(err)
		fmt.Errorf("[UserHandler.UploadImage]: %w", err)
		return
	}
	if message != "" {
		writer.WriteHeader(http.StatusBadRequest)
		common.Response(writer, common.Message(false, message))
		return
	}
	common.HandleSuccess(writer, http.StatusOK, user)
}

