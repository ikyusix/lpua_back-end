package user

import "lpua_back-end/models"

type UserRepositoryInterface interface {
	CheckUserEmail(email string) bool
	GetUserByEmail(email string) (*models.MstUser, error)
	RegisterUser(data *models.MstUser) (*models.MstUser, error)
	GetAllUser()([]*models.MstUser, error)
	GetUserById(id uint) (*models.MstUser, error)
	UpdateUser(id uint, req *models.MstUser) (*models.MstUser, error)
	DeleteUser(id uint) (*models.MstUser, error)
	UploadImage(id uint, path string) (*models.MstUser, error)
}
