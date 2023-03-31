package model

import (
	"gorm.io/gorm"
)

type LogRegister struct {
	gorm.Model

	Sub      string `gorm:"uniqueIndex;type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100);default:''"`
	Name     string `gorm:"type:varchar(100);default:''"`
	Img      string `gorm:"type:varchar(100);default:''"`
	Username string `gorm:"type:varchar(100);default:''"`
	Password string `gorm:"type:varchar(100);default:''"`
	Status   string `gorm:"type:varchar(50);default:''"`
}
