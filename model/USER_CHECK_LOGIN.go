package model

import "gorm.io/gorm"

type USER_CHECK_LOGIN struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;type:varchar(16);not null"`
	Password string `gorm:"type:varchar(250);not null"`
	Idcode   string `gorm:"type:varchar(13);not null"`
}
