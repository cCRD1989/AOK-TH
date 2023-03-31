package model

import (
	"gorm.io/gorm"
)

type LogRegister struct {
	gorm.Model

	Sub      string
	Email    string
	Name     string
	Img      string
	Username string
	Password string
	Status   string
}
