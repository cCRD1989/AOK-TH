package model

import (
	"gorm.io/gorm"
)

type LogRegister struct {
	gorm.Model

	Sub      string `gorm:"uniqueIndex;type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100)"`
	Name     string `gorm:"type:varchar(100)"`
	Img      string `gorm:"type:varchar(100)"`
	Username string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100)"`
	Status   string `gorm:"type:varchar(50)"`
}
