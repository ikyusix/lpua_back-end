package user

import "lpua_back-end/models"

type UserUsecaseInterface interface {
	CheckUserEmail(email string) bool
	GetUserByEmail(email string) (*models.MstUser, error)
	RegisterUser(data *models.MstUser) (map[string]interface{}, error)
	GetAllUser()(map[string]interface{}, error)
	GetUserById(id uint) (map[string]interface{}, error)
	UpdateUser(id uint, req *models.MstUser) (map[string]interface{}, error)
	DeleteUser(id uint) (map[string]interface{}, error)
	UploadImage(id uint, path string) (*models.MstUser, string, error)
}