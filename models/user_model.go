package models

import (
	"github.com/jinzhu/gorm"
)



type MstUser struct {
	gorm.Model
	FullName string `gorm:"column:user_fullName";json:"full_name"`
	Gender   string `gorm:"column:user_gender";json:"gender"`
	Email    string `gorm:"column:user_email";json:"email";type:varchar(100);unique_index;`
	Password string `gorm:"column:user_password";json:"password"`
	UrlImage string `gorm:"column:user_urlImage";json:"url_image"`
}

func (h MstUser) TableName() string {
	return "mstUser"
}
