package usecase

import (
	"lpua_back-end/common"
	"lpua_back-end/models"
	"lpua_back-end/user"
)

type UserUsecaseImpl struct {
	UserRepository user.UserRepositoryInterface
}

func CreateUserUsecaseImpl(userRepo user.UserRepositoryInterface) user.UserUsecaseInterface {
	return &UserUsecaseImpl{UserRepository:userRepo}
}

func (u *UserUsecaseImpl) CheckUserEmail(email string) bool {
	return u.UserRepository.CheckUserEmail(email)
}

func (u *UserUsecaseImpl) GetUserByEmail(email string) (*models.MstUser, error) {
	return u.UserRepository.GetUserByEmail(email)
}

func (u *UserUsecaseImpl) RegisterUser(data *models.MstUser) (map[string]interface{}, error) {
	response, err := u.UserRepository.RegisterUser(data)
	if err != nil {
		return common.Message(false, err.Error()), err
	}
	mapResponse := common.Message(true, "Create user data: success")
	mapResponse["response"] = response
	return mapResponse, nil
}

func (u *UserUsecaseImpl) GetAllUser() (map[string]interface{}, error) {
	response, err := u.UserRepository.GetAllUser()
	if err != nil {
		return common.Message(false, err.Error()), err
	}

	mapResponse := common.Message(true, "read all user data: success")
	mapResponse["response"] = response
	return mapResponse, nil
}

func (u *UserUsecaseImpl) GetUserById(id uint) (map[string]interface{}, error) {
	response, err := u.UserRepository.GetUserById(id)
	if err != nil {
		return common.Message(false, err.Error()), err
	}

	mapResponse := common.Message(true, "Read user data: success")
	mapResponse["response"] = response
	return mapResponse, nil
}

func (u *UserUsecaseImpl) UpdateUser(id uint, req *models.MstUser) (map[string]interface{}, error) {
	response, err := u.UserRepository.UpdateUser(id, req)
	if err != nil {
		return common.Message(false, err.Error()), err
	}

	mapResponse := common.Message(true, "Update user data: success")
	mapResponse["response"] = response
	return mapResponse, nil
}

func (u *UserUsecaseImpl) DeleteUser(id uint) (map[string]interface{}, error) {
	_, err := u.UserRepository.DeleteUser(id)
	if err != nil {
		return common.Message(false, err.Error()), err
	}

	mapResponse := common.Message(true, "Delete user data: success")
	return mapResponse, nil
}

func (u *UserUsecaseImpl) UploadImage(id uint, path string) (*models.MstUser, string, error) {
	valUser, err := u.UserRepository.GetUserById(id)
	if valUser == nil {
		return nil, "Id user not found", nil
	}
	user, err := u.UserRepository.UploadImage(id, path)
	return user, "", err
}