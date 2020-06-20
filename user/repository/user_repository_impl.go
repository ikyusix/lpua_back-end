package repository

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"lpua_back-end/models"
	"lpua_back-end/user"
)

type UserRepositoryImpl struct {
	Conn *gorm.DB
}

func CreateUserRepositoryImpl(DB *gorm.DB) user.UserRepositoryInterface {
	return &UserRepositoryImpl{DB}
}

func (u *UserRepositoryImpl) CheckUserEmail(email string) bool {
	var total int
	u.Conn.Table("mstUser").Where("user_email = ?", email).Count(&total)
	if total > 0 {
		return true
	}
	return false
}

func (u *UserRepositoryImpl) GetUserByEmail(email string) (*models.MstUser, error) {
	dataUser := new(models.MstUser)

	if err := u.Conn.Table("mstUser").Where("user_email = ?", email).First(&dataUser).Error; err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("[UserRepositoryImpl.GetUserByEmail] Get user data by email: %w", err)
	}
	return dataUser, nil
}

func (u *UserRepositoryImpl) RegisterUser(data *models.MstUser) (*models.MstUser, error) {
	if err := u.Conn.Table("mstUser").Save(&data).Error; err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("[UserRepositoryImpl.RegisterUser] Add user data: %w", err)
	}
	return data, nil
}

func (u *UserRepositoryImpl) GetAllUser() ([]*models.MstUser, error) {
	userList := make([]*models.MstUser, 0)
	u.Conn.Find(&userList)
	if err := u.Conn.Table("mstUser").Find(&userList).Error; err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("[UserRepositoryImpl.GetAllUser] Get user list data: %w", err)
	}
	return userList, nil
}

func (u *UserRepositoryImpl) GetUserById(id uint) (*models.MstUser, error) {
	user := new(models.MstUser)
	if err := u.Conn.Table("mstUser").Where("id = ?", id).First(&user).Error; err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("[UserRepositoryImpl.GetUserById] Get user data: %w", err)
	}
	return user, nil
}

func (u *UserRepositoryImpl) UpdateUser(id uint, req *models.MstUser) (*models.MstUser, error) {
	user := new(models.MstUser)
	if err := u.Conn.Table("mstUser").Where("id = ?", id).First(&user).Update(&req).Error; err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("[UserRepositoryImpl.UpdateUser] Update user data: %w", err)
	}
	return user, nil
}

func (u *UserRepositoryImpl) DeleteUser(id uint) (*models.MstUser, error) {
	if err := u.Conn.Table("mstUser").Where("id = ?", id).Delete(&models.MstUser{}).Error; err != nil {
		logrus.Error(err)
		return nil, errors.New("delete user data: error")
	}

	return nil, nil
}

func (u *UserRepositoryImpl) UploadImage(id uint, path string) (*models.MstUser, error) {
	dataUser := models.MstUser{}
	err := u.Conn.Model(&dataUser).Where("id = ?", id).Update("url_image", path).Error
	if err != nil {
		return nil, fmt.Errorf("[UserRepositoryImpl.UploadImage] Error update model data user: %w", err)
	}
	return &dataUser, nil
}