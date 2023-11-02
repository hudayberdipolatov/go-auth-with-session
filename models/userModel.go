package models

import (
	"github.com/jinzhu/gorm"
)

type Users struct {
	gorm.Model
	FullName string
	Username string
	Email    string
	Password string
}

func (user Users) GetUser(username string) Users {
	DB.Where("username =?", username).First(&user)
	return user
}

func (user Users) CreateUser() {
	DB.Create(&user)
}
